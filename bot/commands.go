package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/errorboi/boo/internal/locale"
	"github.com/errorboi/boo/internal/utils"
	"github.com/errorboi/boo/keyboard"
	"github.com/errorboi/boo/store/postgres"
	"github.com/errorboi/boo/text"
	"github.com/errorboi/boo/types/bot"
	timer_types "github.com/errorboi/boo/types/timer"
	user_types "github.com/errorboi/boo/types/user"
	state_types "github.com/errorboi/boo/types/user_state"
	tgbotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) with(m *tgbotAPI.Message, lcl locale.Locale, fn func(*tgbotAPI.Message, locale.Locale), mws ...func(*tgbotAPI.Message, locale.Locale) bool) {
	for _, mw := range mws {
		if !mw(m, lcl) {
			return
		}
	}

	fn(m, lcl)
}

func (b *Bot) stepWith(m *tgbotAPI.Message, state *state_types.State, lcl locale.Locale, fn func(*tgbotAPI.Message, *state_types.State, locale.Locale), mws ...func(message *tgbotAPI.Message, lcl locale.Locale) bool) {
	for _, mw := range mws {
		if !mw(m, lcl) {
			return
		}
	}

	fn(m, state, lcl)
}

func (b *Bot) prestart(m *tgbotAPI.Message, lcl locale.Locale) bool {
	user, err := b.store.Get(sql.NullInt64{Int64: m.From.ID, Valid: true})
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			b.l.Errorf("Get user error: %w", err)

			return false
		}

		b.start(m, lcl)

		user, err = b.store.Get(sql.NullInt64{Int64: m.From.ID, Valid: true})
		if err != nil {
			b.l.Errorf("Get new user error: %w", err)

			return false
		}
	}

	if !user.Locale.Valid {
		msg := tgbotAPI.NewMessage(m.Chat.ID, text.SelectLanguageText[lcl])
		msg.ReplyMarkup = keyboard.SelectLanguageKeyboard(nil)

		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)

			return false
		}

		return false
	}

	tgID := m.From.ID

	completed, err := b.store.IsTaskCompleted(sql.NullInt64{Int64: tgID, Valid: tgID != 0}, 1)
	if err != nil {
		b.l.Errorf("Is task completed error: %w", err)

		return false
	}

	if completed {
		return true
	}

	msg := tgbotAPI.NewAnimation(m.Chat.ID, tgbotAPI.FilePath("./media/vectorboo.mp4"))
	msg.Caption = text.PrestartText[lcl]
	msg.ParseMode = tgbotAPI.ModeMarkdown
	msg.ReplyMarkup = tgbotAPI.NewInlineKeyboardMarkup(
		tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonURL(text.SubscribeButtonText[lcl], "https://t.me/BooTimer"),
			tgbotAPI.NewInlineKeyboardButtonData(text.CheckButtonText[lcl], fmt.Sprintf("checkSubscribe-%d-%d", tgID, 1)),
		),
	)

	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)

		return false
	}

	return false
}

func (b *Bot) log(m *tgbotAPI.Message, lcl locale.Locale) bool {
	if !b.cfg.DebugMode {
		return true
	}

	if m.Chat.IsPrivate() {
		b.l.Infof("Message from %s: %s", m.From.UserName, m.Text)
	}

	return true
}

func (b *Bot) deleteCommandsInChat(m *tgbotAPI.Message, lcl locale.Locale) bool {
	if !m.Chat.IsPrivate() {
		msg := tgbotAPI.NewDeleteMessage(m.Chat.ID, m.MessageID)
		b.BotAPI.Send(msg)

		return false
	}

	return true
}

func (b *Bot) onlyPrivate(m *tgbotAPI.Message, lcl locale.Locale) bool {
	return m.Chat.IsPrivate()
}

func (b *Bot) onlyAdmin(m *tgbotAPI.Message, lcl locale.Locale) bool {
	_, ok := bot.Admins[m.From.ID]

	if !ok {
		b.BotAPI.Send(tgbotAPI.NewDeleteMessage(m.Chat.ID, m.MessageID))
	}

	return ok
}

func (b *Bot) updateUsername(m *tgbotAPI.Message, lcl locale.Locale) bool {
	user, err := b.store.Get(sql.NullInt64{Int64: m.From.ID, Valid: true})
	if err != nil {
		b.l.Errorf("Get user error: %w", err)
		return true
	}

	if user.Username == m.From.UserName {
		return true
	}

	user.Username = m.From.UserName

	err = b.store.UpdateUser(user)
	if err != nil {
		b.l.Errorf("Update user error: %w", err)
		return true
	}

	return true
}

func (b *Bot) start(m *tgbotAPI.Message, lcl locale.Locale) {
	userID := sql.NullInt64{
		Int64: m.From.ID,
		Valid: true,
	}

	var referrerID sql.NullInt64
	if m.Text != "" {
		arr := strings.Split(m.Text, " ")
		if len(arr) > 1 {
			metadata := arr[1]
			metadataArr := strings.Split(metadata, "-")

			for _, value := range metadataArr {
				switch {
				case strings.HasPrefix(value, "r"):
					value = strings.TrimPrefix(value, "r")
					id, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						b.l.Errorf("ParseInt error: %w", err)
					}

					referrerID = sql.NullInt64{
						Int64: id,
						Valid: id != 0,
					}
				case strings.HasPrefix(value, "t"):
					value = strings.TrimPrefix(value, "t")
					timerID, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						b.l.Errorf("ParseInt error: %w", err)
						return
					}

					timer, err := b.store.GetTimer(timerID)
					if err != nil {
						b.l.Errorf("Get timer error: %w", err)
						return
					}

					timer.TgID = m.From.ID
					timer.Status = timer_types.Active

					timers, err := b.store.ListTimers(m.From.ID)
					if err != nil {
						b.l.Errorf("List timers error: %w", err)
						return
					}

					user, err := b.store.Get(sql.NullInt64{Int64: m.From.ID, Valid: true})
					if err != nil && !errors.Is(err, postgres.ErrNotFound) {
						b.l.Errorf("Get user error: %w", err)
						return
					}

					timerLimit := user_types.DefaultTimersLimit
					if user != nil {
						timerLimit = user.TimerLimit
					}

					if len(timers) >= timerLimit {
						msg := tgbotAPI.NewMessage(m.Chat.ID, text.TimersLimitText[lcl])
						msg.ReplyMarkup = keyboard.BackToTimersListKeyboard(lcl)

						_, err = b.BotAPI.Send(msg)
						if err != nil {
							b.l.Errorf("Send message error: %w", err)
						}
						return
					}

					var alreadyAdded bool
					var addedTimerID int64
					for _, t := range timers {
						if t.Name == timer.Name {
							alreadyAdded = true
							addedTimerID = t.ID
							break
						}
					}

					if alreadyAdded {
						addedTimer, err := b.store.GetTimer(addedTimerID)
						if err != nil {
							b.l.Errorf("Get timer error: %w", err)
							return
						}

						msgText := getTimerText(addedTimer, b.BotAPI.Self.UserName, lcl)

						msg := tgbotAPI.NewMessage(m.Chat.ID, msgText)
						msg.ParseMode = tgbotAPI.ModeMarkdown
						replyMarkup := keyboard.TimerInlineKeyboard(addedTimer, lcl)
						msg.ReplyMarkup = &replyMarkup

						_, err = b.BotAPI.Send(msg)
						if err != nil {
							b.l.Errorf("Send message error: %s", err.Error())
						}

						b.BotAPI.Send(tgbotAPI.NewDeleteMessage(m.Chat.ID, m.MessageID))

						return
					}

					err = b.store.CreateTimer(timer)
					if err != nil {
						b.l.Errorf("Create timer error: %w", err)
						return
					}

					msg := tgbotAPI.NewMessage(m.Chat.ID, getTimerText(timer, b.BotAPI.Self.UserName, lcl))
					msg.ParseMode = tgbotAPI.ModeMarkdown
					replyMarkup := keyboard.TimerInlineKeyboard(timer, lcl)
					msg.ReplyMarkup = &replyMarkup
					_, err = b.BotAPI.Send(msg)
					if err != nil {
						b.l.Errorf("Send message error: %w", err)
					}
				}
			}
		}
	}

	username := m.From.UserName
	if username == "" {
		username = m.From.FirstName + " " + m.From.LastName
	}

	isNewUser, err := b.store.CreateIfNotExist(context.Background(), userID, referrerID, username, lcl)
	if err != nil {
		b.l.Errorf("Create user error: %w", err)
	}

	if isNewUser {
		// referral bonus
		if referrerID.Valid {
			referrerBalance, err := b.store.IncreaseBalance(referrerID, bot.ReferralBonus)
			if err != nil && err != postgres.ErrNotFound {
				b.l.Errorf("Increase balance error: %w", err)
			}

			err = b.store.IncrementReferralsAmount(referrerID)
			if err != nil {
				b.l.Errorf("Increment referrals amount error: %w", err)
			}

			msg := tgbotAPI.NewMessage(referrerID.Int64, fmt.Sprintf(text.ReferralAcceptedText[lcl], utils.GetBalanceText(referrerBalance, lcl)))
			msg.ParseMode = tgbotAPI.ModeMarkdown
			_, err = b.BotAPI.Send(msg)
			if err != nil {
				b.l.Errorf("Send message error: %w", err)
			}
		}
	}

	var filePath tgbotAPI.RequestFileData
	var needUpload bool
	fileID, err := b.fileStore.LoadFile("vectorboo.mp4")
	if err != nil {
		b.l.Errorf("Load file error: %w", err)

		filePath = tgbotAPI.FilePath("./media/vectorboo.mp4")

		needUpload = true
	} else {
		filePath = tgbotAPI.FileID(fileID)
	}

	msg := tgbotAPI.NewAnimation(m.Chat.ID, filePath)
	msg.Caption = text.StartText[lcl]
	msg.ParseMode = tgbotAPI.ModeHTML
	if m.Chat.IsPrivate() {
		msg.ReplyMarkup = keyboard.GetStartInlineKeyboard(lcl)
	} else {
		msg.ReplyMarkup = tgbotAPI.NewRemoveKeyboard(true)
	}
	res, err := b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send photo error: %w", err)

		needUpload = true
	}

	if needUpload && res.Animation != nil {
		err = b.fileStore.SaveFile("vectorboo.mp4", res.Animation.FileID)
		if err != nil {
			b.l.Errorf("Save file error: %w", err)
		}
	}
}

func (b *Bot) getInviteLink(m *tgbotAPI.Message, lcl locale.Locale) {
	link := getInviteLink(b.BotAPI.Self.UserName, m.From.ID, nil)
	photoCfg := tgbotAPI.NewPhoto(m.Chat.ID, tgbotAPI.FilePath("./media/crowd.jpeg"))
	photoCfg.Caption = fmt.Sprintf(text.InviteLinkText[lcl], link)
	photoCfg.ParseMode = tgbotAPI.ModeMarkdown
	_, err := b.BotAPI.Send(photoCfg)
	if err != nil {
		b.l.Errorf("Send video error: %w", err)
	}
}

func (b *Bot) profile(m *tgbotAPI.Message, lcl locale.Locale) {
	user, err := b.store.Get(sql.NullInt64{Int64: m.From.ID, Valid: true})
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			b.l.Errorf("Get user error: %w", err)

			return
		}

		b.start(m, lcl)

		user, err = b.store.Get(sql.NullInt64{Int64: m.From.ID, Valid: true})
		if err != nil {
			b.l.Errorf("Get new user error: %w", err)

			return
		}
	}

	balance, err := b.store.GetBalance(sql.NullInt64{Int64: m.From.ID, Valid: true})
	if err != nil {
		b.l.Errorf("Get balance error: %w", err)
	}

	// none by default
	wallet := "🚫 None"
	if user.Wallet.Valid {
		wallet = user.Wallet.String

	}

	msgText := fmt.Sprintf(text.ProfileText[lcl],
		utils.EscapeMarkdown(m.From.UserName),
		utils.GetBalanceText(balance, lcl),
		getProfileSlug(user.ReferralsAmount, lcl),
		fmt.Sprintf("*%s*: %d", text.TimersLimitFieldText[lcl], user.TimerLimit),
		fmt.Sprintf("*%s*: `%s`", text.WalletFieldText[lcl], wallet),
	)

	msg := tgbotAPI.NewPhoto(m.Chat.ID, tgbotAPI.FilePath("./media/yinyang.jpeg"))
	msg.Caption = msgText
	msg.ReplyMarkup = keyboard.ProfileInlineKeyboard(user, lcl)
	msg.ParseMode = tgbotAPI.ModeMarkdown
	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send image error: %w", err)
	}
}

func (b *Bot) newTimer(m *tgbotAPI.Message, lcl locale.Locale) {
	timers, err := b.store.ListTimers(m.From.ID)
	if err != nil {
		b.l.Errorf("List timers error: %w", err)
		return
	}

	user, err := b.store.Get(sql.NullInt64{Int64: m.From.ID, Valid: true})
	if err != nil {
		b.l.Errorf("Get user error: %w", err)
		return
	}

	if len(timers) >= user.TimerLimit {
		msg := tgbotAPI.NewMessage(m.Chat.ID, text.TimersLimitText[lcl])
		msg.ReplyMarkup = keyboard.BackToTimersListKeyboard(lcl)
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}
		return
	}

	step, err := b.stateStore.Init(m.From.ID)
	if err != nil {
		b.l.Errorf("Init user_state error: %w", err)
		return
	}

	b.ExecuteStep(m, step, lcl)
}

func (b *Bot) myTimers(m *tgbotAPI.Message, lcl locale.Locale) {
	timers, err := b.store.ListTimers(m.From.ID)
	if err != nil {
		b.l.Errorf("List timers error: %w", err)
		return
	}

	if len(timers) == 0 {
		msg := tgbotAPI.NewMessage(m.Chat.ID, text.NoTimersText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}
		return
	}

	msg := tgbotAPI.NewMessage(m.Chat.ID, text.YourTimersText[lcl])
	msg.ReplyMarkup = keyboard.GetTimersInlineKeyboard(timers, lcl)

	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)
	}
}

func (b *Bot) presetTimers(m *tgbotAPI.Message, lcl locale.Locale) {
	timers, err := b.store.ListPresetTimers()
	if err != nil {
		b.l.Errorf("List preset timers error: %w", err)
		return
	}

	msg := tgbotAPI.NewMessage(m.Chat.ID, text.PresetTimersText[lcl])
	msg.ReplyMarkup = keyboard.PresetTimersKeyboard(timers, lcl)

	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)
	}
}

func (b *Bot) taskCenter(m *tgbotAPI.Message, lcl locale.Locale) {
	var (
		tasks []bot.Task
		err   error
	)

	if _, ok := bot.Admins[m.From.ID]; ok {
		tasks, err = b.store.GetAllTasks(int(m.From.ID))
	} else {
		tasks, err = b.store.GetTasks(int(m.From.ID))
	}
	if err != nil {
		b.l.Errorf("Get tasks error: %w", err)
		return
	}

	if len(tasks) == 0 {
		msg := tgbotAPI.NewMessage(m.Chat.ID, text.AllTasksFinishedText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}
		return

	}

	msg := tgbotAPI.NewMessage(m.Chat.ID, text.TaskCenterText[lcl])

	msg.ReplyMarkup = keyboard.GetTaskCenterInlineKeyboard(m.From.ID, tasks, lcl)
	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send msg error: %w", err)
	}
}

func (b *Bot) broadcastMessage(m *tgbotAPI.Message, lcl locale.Locale) {
	if _, ok := bot.Admins[m.From.ID]; !ok {
		return
	}

	arr := strings.Split(m.Text, " ")
	if len(arr) < 2 {
		return
	}

	textMsg := strings.Join(arr[1:], " ")

	var fileID *string
	if m.Photo != nil && len(m.Photo) > 0 {
		fileID = &m.Photo[0].FileID
	}

	err := b.notifier.BroadcastMessage(textMsg, fileID)
	if err != nil {
		b.l.Errorf("Broadcast message error: %w", err)
	}
}

func (b *Bot) mint(m *tgbotAPI.Message, lcl locale.Locale) {
	if m.ReplyToMessage == nil || m.ReplyToMessage.From == nil {
		b.l.Error("Reply to message is nil")

		b.BotAPI.Send(tgbotAPI.NewDeleteMessage(m.Chat.ID, m.MessageID))

		return
	}

	var (
		amount      int64
		description = "increase balance"
		err         error
	)
	if m.Text != "" {
		arr := strings.Split(m.Text, " ")
		if len(arr) < 2 {
			b.l.Error("amount is nil")

			b.BotAPI.Send(tgbotAPI.NewDeleteMessage(m.Chat.ID, m.MessageID))

			return
		}

		amount, err = strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			b.l.Errorf("ParseInt error: %w", err)

			b.BotAPI.Send(tgbotAPI.NewDeleteMessage(m.Chat.ID, m.MessageID))

			return
		}

		if len(arr) > 2 {
			description = strings.Join(arr[2:], " ")
		}
	}

	if amount == 0 {
		b.l.Error("amount is nil")

		b.BotAPI.Send(tgbotAPI.NewDeleteMessage(m.Chat.ID, m.MessageID))

		return
	}

	userID := sql.NullInt64{
		Int64: m.ReplyToMessage.From.ID,
		Valid: m.ReplyToMessage.From.ID != 0,
	}

	newBalance, err := b.store.IncreaseBalance(userID, amount)
	if err != nil {
		b.l.Errorf("Increase balance error: %w", err)

		b.BotAPI.Send(tgbotAPI.NewDeleteMessage(m.Chat.ID, m.MessageID))

		return
	}

	err = b.store.SaveTx(userID, amount, description)
	if err != nil {
		b.l.Errorf("Save tx error: %w", err)

		b.BotAPI.Send(tgbotAPI.NewDeleteMessage(m.Chat.ID, m.MessageID))

		return
	}

	msg := tgbotAPI.NewMessage(m.Chat.ID, fmt.Sprintf(text.MintReceivedText[lcl], m.ReplyToMessage.From.UserName, amount, m.From.UserName, description, newBalance))

	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)
	}
}
