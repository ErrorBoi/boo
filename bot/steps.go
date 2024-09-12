package bot

import (
	"database/sql"
	"strings"
	"time"

	"github.com/errorboi/boo/internal/locale"
	"github.com/errorboi/boo/internal/utils"
	"github.com/errorboi/boo/keyboard"
	"github.com/errorboi/boo/text"
	state_types "github.com/errorboi/boo/types/user_state"
	tgbotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) StepInit(m *tgbotAPI.Message, lcl locale.Locale) {
	msg := tgbotAPI.NewMessage(m.Chat.ID, text.NewTimerText[lcl])

	_, err := b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)
		return
	}

	state := state_types.NewState(state_types.Name)

	err = b.stateStore.SetState(m.From.ID, state)
	if err != nil {
		b.l.Errorf("SetState user_state error: %w", err)
		return
	}

	b.ExecuteStep(m, state, lcl)
}

func (b *Bot) StepName(m *tgbotAPI.Message, lcl locale.Locale) {
	msg := tgbotAPI.NewMessage(m.Chat.ID, text.TimerNameText[lcl])
	// msg.ReplyMarkup = bot.NewKeyboard(bot.NewTimerKeyboard)

	_, err := b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)
		return
	}

	state := state_types.NewState(state_types.WaitForName)

	err = b.stateStore.SetState(m.From.ID, state)
	if err != nil {
		b.l.Errorf("SetState user_state error: %w", err)
		return
	}
}

func (b *Bot) StepWaitForName(m *tgbotAPI.Message, lcl locale.Locale) {
	if len(m.Text) < 3 || len(m.Text) > 100 {
		msg := tgbotAPI.NewMessage(m.Chat.ID, text.InvalidTimerNameText[lcl])
		_, err := b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}

		return
	}

	timer, err := b.store.CreateTimerByName(m.From.ID, m.Text)
	if err != nil {
		b.l.Errorf("CreateWithTTL timer error: %w", err)
		return
	}

	successMsg := tgbotAPI.NewMessage(m.Chat.ID, text.TimerSetupSuccessText[lcl])
	_, err = b.BotAPI.Send(successMsg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)
	}

	msg := tgbotAPI.NewMessage(m.Chat.ID, getTimerText(timer, b.BotAPI.Self.UserName, lcl))
	msg.ParseMode = tgbotAPI.ModeMarkdown
	replyMarkup := keyboard.TimerInlineKeyboard(timer, lcl)
	msg.ReplyMarkup = &replyMarkup
	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)
	}

	err = b.stateStore.Del(m.From.ID)
	if err != nil {
		b.l.Errorf("Del user_state error: %w", err)
	}
}

func (b *Bot) StepEditName(m *tgbotAPI.Message, state *state_types.State, lcl locale.Locale) {
	if len(m.Text) < 3 || len(m.Text) > 100 {
		msg := tgbotAPI.NewMessage(m.Chat.ID, text.InvalidTimerNameText[lcl])
		_, err := b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}

		return
	}

	timer, err := b.store.GetTimer(*state.TimerID)
	if err != nil {
		b.l.Errorf("Get timer error: %w", err)
		return
	}

	name := strings.TrimSpace(m.Text)
	name = strings.ReplaceAll(name, "\n", "")

	name = utils.EscapeMarkdown(name)

	timer.Name = name

	err = b.store.UpdateTimer(timer)
	if err != nil {
		b.l.Errorf("Update timer error: %w", err)
		return
	}

	b.TimerUpdateSuccessMessage(m, timer, lcl)
}

func (b *Bot) StepEditDescription(m *tgbotAPI.Message, state *state_types.State, lcl locale.Locale) {
	if len(m.Text) < 3 || len(m.Text) > 1000 {
		msg := tgbotAPI.NewMessage(m.Chat.ID, text.InvalidTimerDescriptionText[lcl])
		_, err := b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}

		return
	}

	timer, err := b.store.GetTimer(*state.TimerID)
	if err != nil {
		b.l.Errorf("Get timer error: %w", err)
		return
	}

	description := strings.TrimSpace(m.Text)
	description = strings.ReplaceAll(description, "\n", "")

	timer.Description = sql.NullString{
		String: m.Text,
		Valid:  m.Text != "",
	}

	err = b.store.UpdateTimer(timer)
	if err != nil {
		b.l.Errorf("Update timer error: %w", err)
		return
	}

	b.TimerUpdateSuccessMessage(m, timer, lcl)
}

func (b *Bot) StepEditTriggerTime(m *tgbotAPI.Message, state *state_types.State, lcl locale.Locale) {
	timer, err := b.store.GetTimer(*state.TimerID)
	if err != nil {
		b.l.Errorf("Get timer error: %w", err)
		return
	}

	msgText := strings.TrimSpace(m.Text)

	err = b.validator.ValidateTriggerTime(msgText)
	if err != nil {
		msg := tgbotAPI.NewMessage(m.Chat.ID, err.Error())
		_, err := b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}

		return
	}

	timer.TriggerTime = sql.NullString{String: msgText, Valid: msgText != ""}

	triggerTime, err := time.Parse("15:04", timer.TriggerTime.String)
	if err != nil {
		b.l.Errorf("Parse trigger time error: %w", err)
		return
	}

	now := time.Now().UTC()
	nextTrigger := time.Date(now.Year(), now.Month(), now.Day(), triggerTime.Hour(), triggerTime.Minute(), triggerTime.Second(), 0, time.UTC)

	if nextTrigger.Before(now) {
		nextTrigger = nextTrigger.Add(24 * time.Hour)
	}

	timer.NextTrigger = nextTrigger

	err = b.store.UpdateTimer(timer)
	if err != nil {
		b.l.Errorf("Update timer error: %w", err)
		return
	}

	b.TimerUpdateSuccessMessage(m, timer, lcl)
}

func (b *Bot) StepEditPeriod(m *tgbotAPI.Message, state *state_types.State, lcl locale.Locale) {
	timer, err := b.store.GetTimer(*state.TimerID)
	if err != nil {
		b.l.Errorf("Get timer error: %w", err)
		return
	}

	msgText := strings.TrimSpace(m.Text)

	err = b.validator.ValidatePeriod(msgText)
	if err != nil {
		msg := tgbotAPI.NewMessage(m.Chat.ID, err.Error())
		_, err := b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}

		return
	}

	layout := "15:04"
	if arr := strings.Split(msgText, ":"); len(arr) == 3 {
		layout = "15:04:05"
	}

	period, err := time.Parse(layout, msgText)
	if err != nil {
		b.l.Errorf("Parse period time error: %w", err)

		msg := tgbotAPI.NewMessage(m.Chat.ID, err.Error())
		_, err = b.BotAPI.Send(msg)

		return
	}

	periodSeconds := periodToSeconds(period)

	if periodSeconds < 5*60 {
		msg := tgbotAPI.NewMessage(m.Chat.ID, text.PeriodTooSmallText[lcl])
		_, err := b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}

		return
	}

	timer.PeriodSeconds = sql.NullInt64{Int64: periodSeconds, Valid: !period.IsZero()}

	timer.NextTrigger = timer.LastTrigger.Add(time.Duration(timer.PeriodSeconds.Int64) * time.Second)

	err = b.store.UpdateTimer(timer)
	if err != nil {
		b.l.Errorf("Update timer error: %w", err)
		return
	}

	b.TimerUpdateSuccessMessage(m, timer, lcl)
}

func (b *Bot) StepEditLink(m *tgbotAPI.Message, state *state_types.State, lcl locale.Locale) {
	timer, err := b.store.GetTimer(*state.TimerID)
	if err != nil {
		b.l.Errorf("Get timer error: %w", err)
		return
	}

	msgText := strings.TrimSpace(m.Text)

	err = b.validator.ValidateLink(msgText)
	if err != nil {
		msg := tgbotAPI.NewMessage(m.Chat.ID, err.Error())
		_, err := b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}

		return
	}

	if strings.HasPrefix(msgText, "@") {
		msgText = "t.me/" + strings.TrimPrefix(msgText, "@")

	}

	timer.Link = sql.NullString{String: msgText, Valid: msgText != ""}

	err = b.store.UpdateTimer(timer)
	if err != nil {
		b.l.Errorf("Update timer error: %w", err)
		return
	}

	b.TimerUpdateSuccessMessage(m, timer, lcl)
}

func (b *Bot) StepEditUserWallet(m *tgbotAPI.Message, state *state_types.State, lcl locale.Locale) {
	user, err := b.store.Get(sql.NullInt64{Int64: *state.UserID, Valid: state.UserID != nil && *state.UserID != 0})
	if err != nil {
		b.l.Errorf("Get user error: %w", err)
		return
	}

	msgText := strings.TrimSpace(m.Text)

	err = b.validator.ValidateTonWallet(msgText)
	if err != nil {
		msg := tgbotAPI.NewMessage(m.Chat.ID, err.Error())
		_, err := b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}

		return
	}

	user.Wallet = sql.NullString{String: msgText, Valid: msgText != ""}

	err = b.store.UpdateUser(user)
	if err != nil {
		msg := tgbotAPI.NewMessage(m.Chat.ID, err.Error())
		b.BotAPI.Send(msg)

		b.l.Errorf("Update user error: %w", err)
		return
	}

	b.UserUpdateSuccessMessage(m, user, lcl)
}
