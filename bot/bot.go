package bot

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/errorboi/boo/internal/config"
	"github.com/errorboi/boo/internal/locale"
	"github.com/errorboi/boo/internal/notify"
	"github.com/errorboi/boo/internal/state_store"
	"github.com/errorboi/boo/internal/validate"
	"github.com/errorboi/boo/keyboard"
	"github.com/errorboi/boo/store"
	"github.com/errorboi/boo/store/postgres"
	"github.com/errorboi/boo/text"
	"github.com/errorboi/boo/types/bot"
	"github.com/errorboi/boo/types/timer"
	state_types "github.com/errorboi/boo/types/user_state"
	tgbotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Bot struct {
	BotAPI     *tgbotAPI.BotAPI
	store      store.UserStore
	l          *zap.SugaredLogger
	cfg        config.Config
	stateStore state_store.StateStore
	validator  validate.Validator
	bonusStore store.BonusStore
	notifier   notify.Notifier
	fileStore  store.FileStore
}

func New(
	botAPI *tgbotAPI.BotAPI,
	store store.UserStore,
	l *zap.SugaredLogger,
	stateStore state_store.StateStore,
	validator validate.Validator,
	cfg config.Config,
	bonusStore store.BonusStore,
	notifier notify.Notifier,
	fileStore store.FileStore,
) (*Bot, error) {
	return &Bot{
		BotAPI:     botAPI,
		store:      store,
		l:          l,
		cfg:        cfg,
		stateStore: stateStore,
		validator:  validator,
		bonusStore: bonusStore,
		notifier:   notifier,
		fileStore:  fileStore,
	}, nil
}

func (b *Bot) InitUpdates() {
	ucfg := tgbotAPI.NewUpdate(0)
	ucfg.Timeout = 60

	updates := b.BotAPI.GetUpdatesChan(ucfg)

	// updates := b.BotAPI.ListenForWebhook("/" + BotToken)
	log.Printf("Authorized on account %s", b.BotAPI.Self.UserName)

	b.processUpdates(updates)
}

func (b *Bot) processUpdates(updates tgbotAPI.UpdatesChannel) {
	for update := range updates {
		lcl := locale.English
		if update.Message == nil {
			if update.CallbackQuery != nil && update.CallbackQuery.Message.Chat.IsPrivate() {
				user, err := b.store.Get(sql.NullInt64{Int64: update.CallbackQuery.From.ID, Valid: update.CallbackQuery.From.ID != 0})
				if err != nil && !errors.Is(err, postgres.ErrNotFound) {
					b.l.Errorf("Get user error: %s", err.Error())
				}

				if user != nil && user.Locale.String != "" && user.Locale.Valid {
					lcl = locale.LangToLocale(user.Locale.String)
				}

				b.ExecuteCallbackQuery(update.CallbackQuery, lcl)
			}
		} else {
			if update.Message.From.LanguageCode == "ru" {
				lcl = locale.Russian
			}

			user, err := b.store.Get(sql.NullInt64{Int64: update.Message.From.ID, Valid: update.Message.From.ID != 0})
			if err != nil && !errors.Is(err, postgres.ErrNotFound) {
				b.l.Errorf("Get user error: %s", err.Error())
			}

			if user != nil && user.Locale.String != "" && user.Locale.Valid {
				lcl = locale.LangToLocale(user.Locale.String)
			}

			switch {
			case update.Message.IsCommand():
				b.ExecuteCommand(update.Message, lcl)
			default:
				if update.Message.Chat.IsPrivate() {
					b.ExecuteText(update.Message, lcl)
				}
			}
		}
	}
}

func (b *Bot) ExecuteCommand(m *tgbotAPI.Message, lcl locale.Locale) {
	command := strings.ToLower(m.Command())

	switch command {
	case "start":
		go b.with(m, lcl, b.start, b.log, b.deleteCommandsInChat, b.onlyPrivate, b.updateUsername)
	// case "invite":
	// 	go b.with(m, lcl, b.getInviteLink, b.prestart, b.log, b.deleteCommandsInChat, b.onlyPrivate, b.updateUsername)
	// case "profile":
	// 	go b.with(m, lcl, b.profile, b.prestart, b.log, b.onlyPrivate, b.updateUsername)
	// case "new_timer":
	// 	go b.with(m, lcl, b.newTimer, b.prestart, b.log, b.deleteCommandsInChat, b.onlyPrivate, b.updateUsername)
	// case "my_timers":
	// 	go b.with(m, lcl, b.myTimers, b.prestart, b.log, b.deleteCommandsInChat, b.onlyPrivate, b.updateUsername)
	// case "presets":
	// 	go b.with(m, lcl, b.presetTimers, b.prestart, b.log, b.deleteCommandsInChat, b.onlyPrivate, b.updateUsername)
	// case "task":
	// 	go b.with(m, lcl, b.taskCenter, b.prestart, b.log, b.deleteCommandsInChat, b.onlyPrivate, b.updateUsername)
	case "broadcast":
		go b.with(m, lcl, b.broadcastMessage, b.onlyAdmin, b.onlyPrivate, b.log, b.deleteCommandsInChat)
	case "mint":
		go b.with(m, lcl, b.mint, b.onlyAdmin, b.log)
	}
}

func (b *Bot) ExecuteText(m *tgbotAPI.Message, lcl locale.Locale) {
	switch m.Text {
	default:
		state, err := b.stateStore.Get(m.From.ID)
		if err != nil {
			b.l.Errorf("Get state error: %s", err.Error())
			return
		}

		b.ExecuteStep(m, state, lcl)
	}
}

func (b *Bot) ExecuteStep(m *tgbotAPI.Message, state *state_types.State, lcl locale.Locale) {
	switch state.Step {
	case state_types.Init:
		go b.with(m, lcl, b.StepInit, b.log, b.deleteCommandsInChat, b.onlyPrivate)
	case state_types.Name:
		go b.with(m, lcl, b.StepName, b.log, b.deleteCommandsInChat, b.onlyPrivate)
	case state_types.WaitForName:
		go b.with(m, lcl, b.StepWaitForName, b.log, b.deleteCommandsInChat, b.onlyPrivate)
	case state_types.NameEdit:
		go b.stepWith(m, state, lcl, b.StepEditName, b.log, b.deleteCommandsInChat, b.onlyPrivate)
	case state_types.DescriptionEdit:
		go b.stepWith(m, state, lcl, b.StepEditDescription, b.log, b.deleteCommandsInChat, b.onlyPrivate)
	case state_types.TriggerTimeEdit:
		go b.stepWith(m, state, lcl, b.StepEditTriggerTime, b.log, b.deleteCommandsInChat, b.onlyPrivate)
	case state_types.PeriodEdit:
		go b.stepWith(m, state, lcl, b.StepEditPeriod, b.log, b.deleteCommandsInChat, b.onlyPrivate)
	case state_types.UserWalletEdit:
		go b.stepWith(m, state, lcl, b.StepEditUserWallet, b.log, b.deleteCommandsInChat, b.onlyPrivate)
	case state_types.LinkEdit:
		go b.stepWith(m, state, lcl, b.StepEditLink, b.log, b.deleteCommandsInChat, b.onlyPrivate)
	}
}

func (b *Bot) TimerUpdateSuccessMessage(m *tgbotAPI.Message, timer *timer.Timer, lcl locale.Locale) {
	successMsg := tgbotAPI.NewMessage(m.Chat.ID, text.TimerUpdateSuccessText[lcl])
	successMsg.ReplyMarkup = keyboard.SuccessMessageKeyboard(timer, lcl)
	_, err := b.BotAPI.Send(successMsg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)
		return
	}
}

func (b *Bot) UserUpdateSuccessMessage(m *tgbotAPI.Message, user *bot.User, lcl locale.Locale) {
	successMsg := tgbotAPI.NewMessage(m.Chat.ID, text.UserUpdateSuccessText[lcl])
	successMsg.ReplyMarkup = keyboard.SuccessUserMessageKeyboard(lcl)
	_, err := b.BotAPI.Send(successMsg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)
		return
	}
}

func (b *Bot) reward(user *bot.User, task *bot.Task) error {
	switch task.RewardType {
	case bot.TimerRewardType:
		user.TimerLimit += task.Reward
		err := b.store.UpdateUser(user)
		if err != nil {
			b.l.Errorf("Update user error: %s", err.Error())
		}
		return err
	case bot.EbooRewardType:
		_, err := b.store.IncreaseBalance(sql.NullInt64{Int64: user.TgID, Valid: user.TgID != 0}, int64(task.Reward))
		if err != nil {
			b.l.Errorf("IncreaseBalance error: %w", err)
		}
		return err
	}

	return nil
}
