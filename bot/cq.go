package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/errorboi/boo/internal/locale"
	"github.com/errorboi/boo/internal/utils"
	"github.com/errorboi/boo/keyboard"
	"github.com/errorboi/boo/store"
	"github.com/errorboi/boo/store/postgres"
	"github.com/errorboi/boo/text"
	"github.com/errorboi/boo/types/bot"
	timer_types "github.com/errorboi/boo/types/timer"
	"github.com/errorboi/boo/types/user_state"
	tgbotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) ExecuteCallbackQuery(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	switch {
	case cq.Data == "homepage":
		b.executeCqHomepage(cq, lcl)
	case cq.Data == "timers":
		b.executeCqTimers(cq, lcl)
	case cq.Data == "list_preset_timers":
		b.executeCqPresetTimers(cq, lcl)
	case cq.Data == "language":
		b.executeCqLanguage(cq, lcl)
	case cq.Data == "tasks":
		b.executeCqTasks(cq, lcl)
	case cq.Data == "profile":
		b.executeCqProfile(cq, lcl)
	case cq.Data == "invite":
		b.executeCqInvite(cq, lcl)
	case cq.Data == "booble_jump":
		b.executeCqBoobleJump(cq, lcl)
	case cq.Data == "leaderboard":
		b.executeCqLeaderboard(cq, lcl)
	case cq.Data == "new_timer":
		b.executeCqNewTimer(cq, lcl)
	case strings.HasPrefix(cq.Data, "timer"):
		b.executeCqTimer(cq, lcl)
	case strings.HasPrefix(cq.Data, "start"):
		b.executeCqStart(cq, lcl)
	case strings.HasPrefix(cq.Data, "stop"):
		b.executeCqStop(cq, lcl)
	case strings.HasPrefix(cq.Data, "delete"):
		b.executeCqDelete(cq, lcl)
	case strings.HasPrefix(cq.Data, "confirm_delete"):
		b.executeCqConfirmDelete(cq, lcl)
	case strings.HasPrefix(cq.Data, "edit_timer"):
		b.executeCqEditTimer(cq, lcl)
	case strings.HasPrefix(cq.Data, "editTimer"):
		b.executeCqTimerFieldEdit(cq, lcl)
	case strings.HasPrefix(cq.Data, "editUser"):
		b.executeCqUserEdit(cq, lcl)
	case strings.HasPrefix(cq.Data, "edit_type"):
		b.executeCqEditType(cq, lcl)
	case strings.HasPrefix(cq.Data, "cancel"):
		b.executeCqCancel(cq, lcl)
	case strings.HasPrefix(cq.Data, "checkSubscribe"):
		b.executeCqCheckSubscribe(cq, lcl)
	case strings.HasPrefix(cq.Data, "checkAddToNickname"):
		b.executeCqCheckAddToNickname(cq, lcl)
	case strings.HasPrefix(cq.Data, "show_preset"):
		b.executeCqShowPreset(cq, lcl)
	case strings.HasPrefix(cq.Data, "reboot"):
		b.executeCqReboot(cq, lcl)
	case strings.HasPrefix(cq.Data, "add_preset"):
		b.executeCqAddPreset(cq, lcl)
	case strings.HasPrefix(cq.Data, "lang"):
		b.executeCqLang(cq, lcl)
	case strings.HasPrefix(cq.Data, "task"):
		b.executeCqTask(cq, lcl)
	case strings.HasPrefix(cq.Data, "claimbonus"):
		b.executeCqClaimBonus(cq, lcl)
	}
}

func (b *Bot) executeCqHomepage(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
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

	msg := tgbotAPI.NewAnimation(cq.Message.Chat.ID, filePath)
	msg.Caption = text.StartText[lcl]
	msg.ReplyMarkup = keyboard.GetStartInlineKeyboard(lcl)

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

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqTimers(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	timers, err := b.store.ListTimers(cq.From.ID)
	if err != nil {
		b.l.Errorf("List timers error: %w", err)
		return
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

	msg := tgbotAPI.NewAnimation(cq.Message.Chat.ID, filePath)
	msg.Caption = text.YourTimersText[lcl]
	if len(timers) == 0 {
		msg.Caption = text.NoTimersText[lcl]
	}

	replyMarkup := keyboard.GetTimersInlineKeyboard(timers, lcl)
	msg.ReplyMarkup = &replyMarkup

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

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqPresetTimers(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	timers, err := b.store.ListPresetTimers()
	if err != nil {
		b.l.Errorf("List preset timers error: %w", err)
		return
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

	msg := tgbotAPI.NewAnimation(cq.Message.Chat.ID, filePath)

	msg.Caption = text.PresetTimersText[lcl]
	replyMarkup := keyboard.PresetTimersKeyboard(timers, lcl)
	msg.ReplyMarkup = &replyMarkup

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

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqTasks(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	var (
		tasks []bot.Task
		err   error
	)

	if _, ok := bot.Admins[cq.From.ID]; ok {
		tasks, err = b.store.GetAllTasks(int(cq.From.ID))
	} else {
		tasks, err = b.store.GetTasks(int(cq.From.ID))
	}
	if err != nil {
		b.l.Errorf("Get tasks error: %w", err)
		return
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

	msg := tgbotAPI.NewAnimation(cq.Message.Chat.ID, filePath)
	msg.Caption = text.TaskCenterText[lcl]
	if len(tasks) == 0 {
		msg.Caption = text.AllTasksFinishedText[lcl]
	}

	msg.ReplyMarkup = keyboard.GetTaskCenterInlineKeyboard(cq.From.ID, tasks, lcl)

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

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqProfile(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	user, err := b.store.Get(sql.NullInt64{Int64: cq.From.ID, Valid: true})
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			b.l.Errorf("Get user error: %w", err)

			return
		}
	}

	balance, err := b.store.GetBalance(sql.NullInt64{Int64: cq.From.ID, Valid: true})
	if err != nil {
		b.l.Errorf("Get balance error: %w", err)
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

	// none by default
	wallet := "🚫 None"
	if user.Wallet.Valid {
		wallet = user.Wallet.String

	}

	msgText := fmt.Sprintf(text.ProfileText[lcl],
		utils.EscapeMarkdown(cq.From.UserName),
		utils.GetBalanceText(balance, lcl),
		fmt.Sprintf("*%s*: %d", text.TimersLimitFieldText[lcl], user.TimerLimit),
		fmt.Sprintf("*%s*: `%s`", text.WalletFieldText[lcl], utils.EscapeMarkdown(wallet)),
	)

	msg := tgbotAPI.NewAnimation(cq.Message.Chat.ID, filePath)
	msg.Caption = msgText
	msg.ReplyMarkup = keyboard.ProfileInlineKeyboard(user, lcl)
	msg.ParseMode = tgbotAPI.ModeMarkdown

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

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqInvite(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	link := getInviteLink(b.BotAPI.Self.UserName, cq.From.ID, nil)

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

	msg := tgbotAPI.NewAnimation(cq.Message.Chat.ID, filePath)

	msg.Caption = fmt.Sprintf(text.InviteLinkText[lcl], link)
	msg.ParseMode = tgbotAPI.ModeMarkdown
	msg.ReplyMarkup = keyboard.DefaultInlineKeyboard(lcl)

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

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqBoobleJump(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	var filePath tgbotAPI.RequestFileData
	var needUpload bool
	fileID, err := b.fileStore.LoadFile("BOOARDING.mp4")
	if err != nil {
		filePath = tgbotAPI.FilePath("./media/BOOARDING.mp4")

		needUpload = true
	} else {
		filePath = tgbotAPI.FileID(fileID)
	}

	msg := tgbotAPI.NewAnimation(cq.Message.Chat.ID, filePath)
	msg.Caption = text.BoobleJumpText[lcl]
	msg.ReplyMarkup = keyboard.DefaultInlineKeyboard(lcl)

	res, err := b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send animation error: %w", err)

		return
	}

	if needUpload && res.Video != nil {
		err = b.fileStore.SaveFile("BOOARDING.mp4", res.Video.FileID)
		if err != nil {
			b.l.Errorf("Save file error: %w", err)
		}
	}

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqLeaderboard(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
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

	caption := text.LeaderboardGoogleText[lcl]

	msg := tgbotAPI.NewAnimation(cq.Message.Chat.ID, filePath)
	msg.Caption = caption
	msg.ParseMode = tgbotAPI.ModeMarkdownV2
	msg.ReplyMarkup = keyboard.DefaultInlineKeyboard(lcl)

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

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqNewTimer(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	timers, err := b.store.ListTimers(cq.From.ID)
	if err != nil {
		b.l.Errorf("List timers error: %w", err)
		return
	}

	user, err := b.store.Get(sql.NullInt64{Int64: cq.From.ID, Valid: true})
	if err != nil {
		b.l.Errorf("Get user error: %w", err)
		return
	}

	if len(timers) >= user.TimerLimit {
		msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, text.TimersLimitText[lcl])
		msg.ReplyMarkup = keyboard.BackToTimersListKeyboard(lcl)
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}
		return
	}

	step, err := b.stateStore.Init(cq.From.ID)
	if err != nil {
		b.l.Errorf("Init user_state error: %w", err)
		return
	}

	cq.Message.From = cq.From

	b.ExecuteStep(cq.Message, step, lcl)
}

func (b *Bot) executeCqLanguage(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, text.SelectLanguageText[lcl])
	lang := lcl.Shortcode()
	replyMarkup := keyboard.SelectLanguageKeyboard(&lang)
	msg.ReplyMarkup = &replyMarkup

	_, err := b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)

		return
	}

	b.BotAPI.Send(tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID))

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqTimer(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	info := strings.Split(cq.Data, "-")

	timerIdStr := info[1]
	timerID, err := strconv.Atoi(timerIdStr)
	if err != nil {
		b.l.Errorf("timerID string to int convertation error: %s", err.Error())
		return
	}

	timer, err := b.store.GetTimer(int64(timerID))
	if err != nil {
		b.l.Errorf("GetTimer error: %s", err.Error())
		return
	}

	if !timer.IsOwnedBy(cq.From.ID) {
		msg := tgbotAPI.NewCallback(cq.ID, text.YouCannotPerformThisText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
			return
		}
	}

	msgText := getTimerText(timer, b.BotAPI.Self.UserName, lcl)

	msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, msgText)
	msg.ParseMode = tgbotAPI.ModeMarkdown
	replyMarkup := keyboard.TimerInlineKeyboard(timer, lcl)
	msg.ReplyMarkup = &replyMarkup

	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %s", err.Error())
		return
	}

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqStart(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	var validateTimer = func(timer *timer_types.Timer) string {
		if !timer.IsOwnedBy(cq.From.ID) {
			return text.YouCannotPerformThisText[lcl]
		}

		if timer.Type == nil {
			return text.TimerTypeNotSetText[lcl]
		}

		switch *timer.Type {
		case timer_types.Daily:
			if !timer.TriggerTime.Valid {
				return text.TimerAlertTimeNotSetText[lcl]
			}
		case timer_types.Periodical:
			if !timer.PeriodSeconds.Valid {
				return text.TimerPeriodIsNotSetText[lcl]
			}
		}

		return ""
	}

	info := strings.Split(cq.Data, "-")

	timerIdStr := info[1]
	timerID, err := strconv.Atoi(timerIdStr)
	if err != nil {
		b.l.Errorf("timerID string to int convertation error: %s", err.Error())
		return
	}

	timer, err := b.store.GetTimer(int64(timerID))
	if err != nil {
		b.l.Errorf("GetTimer error: %s", err.Error())
		return
	}

	if errSlug := validateTimer(timer); errSlug != "" {
		msg := tgbotAPI.NewCallbackWithAlert(cq.ID, errSlug)
		b.BotAPI.Send(msg)

		return
	}

	timer.Status = timer_types.Active

	err = b.store.UpdateTimer(timer)
	if err != nil {
		b.l.Errorf("UpdateTimer error: %s", err.Error())
		return
	}

	msgText := getTimerText(timer, b.BotAPI.Self.UserName, lcl)
	msg := tgbotAPI.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, msgText)
	msg.ParseMode = tgbotAPI.ModeMarkdown
	replyMarkup := keyboard.TimerInlineKeyboard(timer, lcl)
	msg.ReplyMarkup = &replyMarkup

	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %s", err.Error())
		return
	}

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqReboot(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	info := strings.Split(cq.Data, "-")

	timerIdStr := info[1]
	timerID, err := strconv.Atoi(timerIdStr)
	if err != nil {
		b.l.Errorf("timerID string to int convertation error: %s", err.Error())
		return
	}

	timer, err := b.store.GetTimer(int64(timerID))
	if err != nil {
		b.l.Errorf("GetTimer error: %s", err.Error())
		return
	}

	if !timer.IsOwnedBy(cq.From.ID) {
		msg := tgbotAPI.NewCallback(cq.ID, text.YouCannotPerformThisText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
			return
		}
	}

	if timer.Type != nil && *timer.Type == timer_types.Periodical {
		timer.LastTrigger = time.Now().UTC()
		timer.NextTrigger = time.Now().UTC().Add(time.Duration(timer.PeriodSeconds.Int64) * time.Second)
	}

	err = b.store.UpdateTimer(timer)
	if err != nil {
		b.l.Errorf("UpdateTimer error: %s", err.Error())
		return
	}

	msgText := getTimerText(timer, b.BotAPI.Self.UserName, lcl)
	msg := tgbotAPI.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, msgText)
	msg.ParseMode = tgbotAPI.ModeMarkdown
	replyMarkup := keyboard.TimerInlineKeyboard(timer, lcl)
	msg.ReplyMarkup = &replyMarkup

	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %s", err.Error())
		return
	}

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqStop(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	info := strings.Split(cq.Data, "-")

	timerIdStr := info[1]
	timerID, err := strconv.Atoi(timerIdStr)
	if err != nil {
		b.l.Errorf("timerID string to int convertation error: %s", err.Error())
		return
	}

	timer, err := b.store.GetTimer(int64(timerID))
	if err != nil {
		b.l.Errorf("GetTimer error: %s", err.Error())
		return
	}

	if !timer.IsOwnedBy(cq.From.ID) {

		msg := tgbotAPI.NewCallback(cq.ID, text.YouCannotPerformThisText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
			return
		}
	}

	timer.Status = timer_types.Inactive

	err = b.store.UpdateTimer(timer)
	if err != nil {
		b.l.Errorf("UpdateTimer error: %s", err.Error())
		return
	}

	msgText := getTimerText(timer, b.BotAPI.Self.UserName, lcl)
	msg := tgbotAPI.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, msgText)
	msg.ParseMode = tgbotAPI.ModeMarkdown
	replyMarkup := keyboard.TimerInlineKeyboard(timer, lcl)
	msg.ReplyMarkup = &replyMarkup

	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %s", err.Error())
		return
	}

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqDelete(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	info := strings.Split(cq.Data, "-")

	timerIdStr := info[1]
	timerID, err := strconv.Atoi(timerIdStr)
	if err != nil {
		b.l.Errorf("timerID string to int convertation error: %s", err.Error())
		return
	}

	timer, err := b.store.GetTimer(int64(timerID))
	if err != nil {
		b.l.Errorf("GetTimer error: %s", err.Error())
		return
	}

	if !timer.IsOwnedBy(cq.From.ID) {

		msg := tgbotAPI.NewCallback(cq.ID, text.YouCannotPerformThisText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
			return
		}
	}

	msg := tgbotAPI.NewEditMessageText(cq.From.ID, cq.Message.MessageID, text.ConfirmDeleteText[lcl])
	replyMarkup := keyboard.DeleteTimerInlineKeyboard(timer.ID, lcl)
	msg.ReplyMarkup = &replyMarkup

	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %s", err.Error())
		return
	}

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqConfirmDelete(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	info := strings.Split(cq.Data, "-")

	timerIdStr := info[1]
	timerID, err := strconv.Atoi(timerIdStr)
	if err != nil {
		b.l.Errorf("timerID string to int convertation error: %s", err.Error())
		return
	}

	timer, err := b.store.GetTimer(int64(timerID))
	if err != nil {
		b.l.Errorf("GetTimer error: %s", err.Error())
		return
	}

	if !timer.IsOwnedBy(cq.From.ID) {

		msg := tgbotAPI.NewCallback(cq.ID, text.YouCannotPerformThisText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
			return
		}
	}

	err = b.store.DeleteTimer(timer.ID)
	if err != nil {
		b.l.Errorf("DeleteTimer error: %s", err.Error())
		return
	}

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, fmt.Sprintf(text.TimerDeletedText[lcl], timer.Name))
	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %s", err.Error())
		return
	}

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqEditTimer(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	info := strings.Split(cq.Data, "-")

	timerIdStr := info[1]
	timerID, err := strconv.Atoi(timerIdStr)
	if err != nil {
		b.l.Errorf("timerID string to int convertation error: %s", err.Error())
		return
	}

	timer, err := b.store.GetTimer(int64(timerID))
	if err != nil {
		b.l.Errorf("GetTimer error: %s", err.Error())
		return
	}

	if !timer.IsOwnedBy(cq.From.ID) {

		msg := tgbotAPI.NewCallback(cq.ID, text.YouCannotPerformThisText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
			return
		}
	}

	msg := tgbotAPI.NewEditMessageTextAndMarkup(cq.Message.Chat.ID, cq.Message.MessageID, getTimerText(timer, b.BotAPI.Self.UserName, lcl), keyboard.EditTimerInlineKeyboard(timer, lcl))
	msg.ParseMode = tgbotAPI.ModeMarkdown

	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %s", err.Error())
		return
	}

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqTimerFieldEdit(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	info := strings.Split(cq.Data, "-")

	timerIdStr := info[1]
	timerID, err := strconv.Atoi(timerIdStr)
	if err != nil {
		b.l.Errorf("timerID string to int convertation error: %s", err.Error())
		return
	}

	timer, err := b.store.GetTimer(int64(timerID))
	if err != nil {
		b.l.Errorf("GetTimer error: %s", err.Error())
		return
	}

	if !timer.IsOwnedBy(cq.From.ID) {
		msg := tgbotAPI.NewCallback(cq.ID, text.YouCannotPerformThisText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
			return
		}
	}

	var (
		msgText     string
		step        user_state.Step
		replyMarkup *tgbotAPI.InlineKeyboardMarkup
	)

	parsedStep := strings.TrimSpace(cq.Data)

	switch {
	case strings.HasPrefix(parsedStep, string(user_state.NameEdit)):
		msgText = text.EditNameText[lcl]
		step = user_state.NameEdit
	case strings.HasPrefix(parsedStep, string(user_state.DescriptionEdit)):
		msgText = text.EditDescriptionText[lcl]
		step = user_state.DescriptionEdit
	case strings.HasPrefix(parsedStep, string(user_state.TypeEdit)):
		msgText = text.EditTypeText[lcl]
		step = user_state.None
		replyMarkup = keyboard.EditTimerTypeInlineKeyboard(timer, lcl)

		msg := tgbotAPI.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, msgText)
		msg.ParseMode = tgbotAPI.ModeMarkdown
		if replyMarkup != nil {
			msg.ReplyMarkup = replyMarkup
		}
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}
		return
	case strings.HasPrefix(parsedStep, string(user_state.TriggerTimeEdit)):
		msgText = text.EditTriggerTimeText[lcl]
		step = user_state.TriggerTimeEdit
	case strings.HasPrefix(parsedStep, string(user_state.PeriodEdit)):
		msgText = text.EditPeriodText[lcl]
		step = user_state.PeriodEdit
	case strings.HasPrefix(parsedStep, string(user_state.RepeatTypeEdit)):
		msgText = text.EditRepeatTypeText[lcl]
		step = user_state.RepeatTypeEdit
	case strings.HasPrefix(parsedStep, string(user_state.LinkEdit)):
		msgText = utils.EscapeMarkdown(text.EditLinkText[lcl])
		step = user_state.LinkEdit
	}

	msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, msgText)
	msg.ParseMode = tgbotAPI.ModeMarkdown
	msg.ReplyMarkup = keyboard.CancelInlineKeyboard(timer.ID, lcl)
	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)
		return
	}

	state := user_state.NewState(step)
	state.TimerID = &timer.ID

	err = b.stateStore.SetState(cq.From.ID, state)
	if err != nil {
		b.l.Errorf("SetState user_state error: %w", err)
		return
	}

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqUserEdit(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	info := strings.Split(cq.Data, "-")

	userIdStr := info[1]
	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		b.l.Errorf("userID string to int convertation error: %s", err.Error())
		return
	}

	user, err := b.store.Get(sql.NullInt64{Int64: int64(userID), Valid: true})
	if err != nil {
		b.l.Errorf("GetTimer error: %s", err.Error())
		return
	}

	if user.TgID != cq.From.ID {
		msg := tgbotAPI.NewCallback(cq.ID, text.YouCannotPerformThisText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
			return
		}
	}

	var (
		msgText string
		step    user_state.Step
	)

	parsedStep := strings.TrimSpace(cq.Data)

	switch {
	case strings.HasPrefix(parsedStep, string(user_state.UserWalletEdit)):
		msgText = text.EditUserWalletText[lcl]
		step = user_state.UserWalletEdit
	}

	msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, msgText)
	msg.ParseMode = tgbotAPI.ModeMarkdown
	msg.ReplyMarkup = keyboard.CancelProfileInlineKeyboard(lcl)
	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)
		return
	}

	state := user_state.NewState(step)
	state.UserID = &user.TgID

	err = b.stateStore.SetState(cq.From.ID, state)
	if err != nil {
		b.l.Errorf("SetState user_state error: %w", err)
		return
	}

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqEditType(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	info := strings.Split(cq.Data, "-")

	if len(info) < 2 {
		b.l.Error("EditType callback data is invalid: %s", cq.Data)
		return
	}

	timerIdStr := info[1]
	timerID, err := strconv.Atoi(timerIdStr)
	if err != nil {
		b.l.Errorf("timerID string to int convertation error: %s", err.Error())
		return
	}

	timer, err := b.store.GetTimer(int64(timerID))
	if err != nil {
		b.l.Errorf("GetTimer error: %s", err.Error())
		return
	}

	if !timer.IsOwnedBy(cq.From.ID) {
		msg := tgbotAPI.NewCallback(cq.ID, text.YouCannotPerformThisText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
			return
		}
	}

	timerType := timer_types.TimerType(info[2])
	if timerType != timer_types.Daily {
		timerType = timer_types.Periodical
	}

	timer.Type = &timerType

	err = b.store.UpdateTimer(timer)
	if err != nil {
		b.l.Errorf("UpdateTimer error: %w", err)
		return
	}

	msg := tgbotAPI.NewEditMessageReplyMarkup(cq.Message.Chat.ID, cq.Message.MessageID, *keyboard.EditTimerTypeInlineKeyboard(timer, lcl))

	b.BotAPI.Send(msg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqCancel(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	state, err := b.stateStore.Get(cq.From.ID)
	if err != nil {
		b.l.Errorf("Get user_state error: %w", err)
		return
	}

	if state == nil {
		return
	}

	err = b.stateStore.Del(cq.From.ID)
	if err != nil {
		b.l.Errorf("Del user_state error: %w", err)
		return
	}

	b.BotAPI.Send(tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID))

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqCheckSubscribe(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	data := strings.Split(cq.Data, "-")

	tgIDStr := data[1]
	tgID, err := strconv.Atoi(tgIDStr)
	if err != nil {
		b.l.Errorf("tgID string to int convertation error: %s", err.Error())
		return
	}

	if tgID != int(cq.From.ID) {
		return
	}

	taskIDStr := data[2]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		b.l.Errorf("taskID string to int convertation error: %s", err.Error())
		return
	}

	completed, err := b.store.IsTaskCompleted(sql.NullInt64{Int64: int64(tgID), Valid: tgID != 0}, taskID)
	if err != nil {
		b.l.Errorf("IsTaskCompleted error: %s", err.Error())
		return
	}

	if completed {
		b.BotAPI.Send(tgbotAPI.NewCallback(cq.ID, text.TaskAlreadyCompletedText[lcl]))

		b.BotAPI.Send(tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID))

		return
	}

	task, err := b.store.GetTask(taskID)
	if err != nil {
		b.l.Errorf("GetTask error: %s", err.Error())
		return
	}

	chatMember, err := b.BotAPI.GetChatMember(tgbotAPI.GetChatMemberConfig{
		ChatConfigWithUser: tgbotAPI.ChatConfigWithUser{
			UserID: int64(tgID),
			ChatID: task.ChatID,
		},
	})
	if err != nil {
		b.l.Errorf("GetChatMember error: %s", err.Error())
		return
	}

	if chatMember.Status == "left" || chatMember.Status == "kicked" {
		msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, text.NeedSubscribeText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
		}
		return
	}

	err = b.store.CompleteTask(sql.NullInt64{Int64: int64(tgID), Valid: tgID != 0}, taskID)
	if err != nil {
		b.l.Errorf("CompleteTask error: %s", err.Error())
		return
	}

	user, err := b.store.Get(sql.NullInt64{Int64: int64(tgID), Valid: tgID != 0})
	if err != nil {
		if !errors.Is(err, postgres.ErrNotFound) {
			b.l.Errorf("Get user error: %s", err.Error())
			return
		}

		_, err = b.store.CreateIfNotExist(context.Background(), sql.NullInt64{Int64: int64(tgID), Valid: tgID != 0}, sql.NullInt64{}, cq.From.UserName, lcl)
		if err != nil {
			b.l.Errorf("CreateIfNotExist error: %w", err)
			return
		}

		user, err = b.store.Get(sql.NullInt64{Int64: int64(tgID), Valid: tgID != 0})
		if err != nil {
			b.l.Errorf("Get user error: %w", err)
			return
		}
	}

	err = b.reward(user, task)
	if err != nil {
		b.l.Errorf("Reward error: %s", err.Error())
		return
	}

	msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, fmt.Sprintf(text.TaskCompleteSuccessText[lcl], task.Reward, task.RewardType.ToText()))
	msg.ReplyMarkup = keyboard.DefaultInlineKeyboard(lcl)

	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %s", err.Error())
	}

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqCheckAddToNickname(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	data := strings.Split(cq.Data, "-")

	tgIDStr := data[1]
	tgID, err := strconv.Atoi(tgIDStr)
	if err != nil {
		b.l.Errorf("tgID string to int convertation error: %s", err.Error())
		return
	}

	if tgID != int(cq.From.ID) {
		return
	}

	taskIDStr := data[2]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		b.l.Errorf("taskID string to int convertation error: %s", err.Error())
		return
	}

	completed, err := b.store.IsTaskCompleted(sql.NullInt64{Int64: int64(tgID), Valid: tgID != 0}, taskID)
	if err != nil {
		b.l.Errorf("IsTaskCompleted error: %s", err.Error())
		return
	}

	if completed {
		b.BotAPI.Send(tgbotAPI.NewCallback(cq.ID, text.TaskAlreadyCompletedText[lcl]))

		b.BotAPI.Send(tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID))

		return
	}

	task, err := b.store.GetTask(taskID)
	if err != nil {
		b.l.Errorf("GetTask error: %s", err.Error())
		return
	}

	fullname := cq.From.FirstName + " " + cq.From.LastName
	if !strings.Contains(fullname, "👻") {
		b.BotAPI.Send(tgbotAPI.NewCallback(cq.ID, text.NeedBooNicknameText[lcl]))

		return
	}

	err = b.store.CompleteTask(sql.NullInt64{Int64: int64(tgID), Valid: tgID != 0}, taskID)
	if err != nil {
		b.l.Errorf("CompleteTask error: %s", err.Error())
		return
	}

	user, err := b.store.Get(sql.NullInt64{Int64: int64(tgID), Valid: tgID != 0})
	if err != nil {
		b.l.Errorf("Get user error: %s", err.Error())
		return
	}

	err = b.reward(user, task)
	if err != nil {
		b.l.Errorf("Reward error: %s", err.Error())
		return
	}

	msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, fmt.Sprintf(text.TaskCompleteSuccessText[lcl], task.Reward, task.RewardType.ToText()))
	msg.ReplyMarkup = keyboard.DefaultInlineKeyboard(lcl)

	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %s", err.Error())
	}

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqShowPreset(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	info := strings.Split(cq.Data, "-")

	timerIdStr := info[1]
	timerID, err := strconv.Atoi(timerIdStr)
	if err != nil {
		b.l.Errorf("timerID string to int convertation error: %s", err.Error())
		return
	}

	timer, err := b.store.GetPresetTimer(int64(timerID))
	if err != nil {
		b.l.Errorf("GetPresetTimer error: %s", err.Error())
		return
	}

	var filePath tgbotAPI.RequestFileData
	var needUpload bool
	fileID, err := b.fileStore.LoadFile(timer.ImgUrl.String)
	if err != nil {
		b.l.Errorf("Load file error: %w", err)

		filePath = tgbotAPI.FilePath(timer.ImgUrl.String)

		needUpload = true
	} else {
		filePath = tgbotAPI.FileID(fileID)
	}

	msg := tgbotAPI.NewPhoto(cq.Message.Chat.ID, filePath)
	msg.Caption = getPresetTimerText(timer, lcl)
	msg.ReplyMarkup = keyboard.PresetTimerInlineKeyboard(timer, lcl)
	msg.ParseMode = tgbotAPI.ModeMarkdown

	res, err := b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %s", err.Error())
		return
	}

	if needUpload && len(res.Photo) > 0 {
		err = b.fileStore.SaveFile(timer.ImgUrl.String, res.Photo[0].FileID)
		if err != nil {
			b.l.Errorf("Save file error: %w", err)
		}
	}

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqAddPreset(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	user, err := b.store.Get(sql.NullInt64{Int64: cq.From.ID, Valid: cq.From.ID != 0})
	if err != nil {
		b.l.Errorf("Get user error: %w", err)
		return
	}

	timers, err := b.store.ListTimers(cq.From.ID)
	if err != nil {
		b.l.Errorf("List timers error: %w", err)
		return
	}

	if len(timers) >= user.TimerLimit {
		msg := tgbotAPI.NewCallbackWithAlert(cq.ID, text.TimerLimitReachedText[lcl])
		b.BotAPI.Send(msg)
		return
	}

	info := strings.Split(cq.Data, "-")

	timerIdStr := info[1]
	timerID, err := strconv.Atoi(timerIdStr)
	if err != nil {
		b.l.Errorf("timerID string to int convertation error: %s", err.Error())
		return
	}

	oldTimer, err := b.store.GetPresetTimer(int64(timerID))
	if err != nil {
		b.l.Errorf("GetPresetTimer error: %s", err.Error())
		return
	}

	var alreadyAdded bool
	var addedTimerID int64
	for _, t := range timers {
		if t.Name == oldTimer.Name {
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

		msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, msgText)
		msg.ParseMode = tgbotAPI.ModeMarkdown
		replyMarkup := keyboard.TimerInlineKeyboard(addedTimer, lcl)
		msg.ReplyMarkup = &replyMarkup

		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
		}

		b.BotAPI.Send(tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID))

		return
	}

	now := time.Now().UTC()

	var nextTrigger time.Time
	switch *oldTimer.Type {
	case timer_types.Daily:
		triggerTime, err := time.Parse("15:04", oldTimer.TriggerTime.String)
		if err != nil {
			b.l.Errorf("Parse trigger time error: %w", err)
			return
		}

		nextTrigger = time.Date(now.Year(), now.Month(), now.Day(), triggerTime.Hour(), triggerTime.Minute(), triggerTime.Second(), 0, time.UTC)
		if nextTrigger.Before(now) {
			nextTrigger = nextTrigger.Add(24 * time.Hour)
		}
	case timer_types.Periodical:
		nextTrigger = now.Add(time.Duration(oldTimer.PeriodSeconds.Int64) * time.Second)
	}

	newTimer := &timer_types.Timer{
		TgID:          cq.From.ID,
		Name:          oldTimer.Name,
		Type:          oldTimer.Type,
		CreatedAt:     now,
		Description:   oldTimer.Description,
		PeriodSeconds: oldTimer.PeriodSeconds,
		TriggerTime:   oldTimer.TriggerTime,
		LastTrigger:   now,
		Status:        timer_types.Active,
		NextTrigger:   nextTrigger,
		Link:          oldTimer.Link,
		ImgUrl:        oldTimer.ImgUrl,
	}

	err = b.store.CreateTimer(newTimer)
	if err != nil {
		b.l.Errorf("CreateTimer error: %s", err.Error())
		return
	}

	msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, text.TimerAddedText[lcl])
	msg.ReplyMarkup = keyboard.AddPresetTimerSuccessInlineKeyboard(lcl)
	b.BotAPI.Send(msg)

	deleteMsg := tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.BotAPI.Send(deleteMsg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqLang(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	info := strings.Split(cq.Data, "-")

	lang := info[1]

	l := locale.LangToLocale(lang)

	if l == lcl {
		return
	}

	user, err := b.store.Get(sql.NullInt64{Int64: cq.From.ID, Valid: cq.From.ID != 0})
	if err != nil {
		b.l.Errorf("Get user error: %w", err)
		return
	}

	user.Locale = sql.NullString{String: l.Shortcode(), Valid: true}

	err = b.store.UpdateUser(user)
	if err != nil {
		b.l.Errorf("Update user error: %w", err)
		return
	}

	msg := tgbotAPI.NewEditMessageTextAndMarkup(cq.Message.Chat.ID, cq.Message.MessageID, text.SelectLanguageText[l], keyboard.SelectLanguageKeyboard(&lang))
	b.BotAPI.Send(msg)

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqTask(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	taskInfo := strings.Split(cq.Data, "-")

	tgIDStr := taskInfo[1]
	tgID, err := strconv.Atoi(tgIDStr)
	if err != nil {
		b.l.Errorf("tgID string to int convertation error: %s", err.Error())
		return
	}

	if tgID != int(cq.From.ID) {
		return
	}

	taskIDStr := taskInfo[2]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		b.l.Errorf("taskID string to int convertation error: %s", err.Error())
		return
	}

	completed, err := b.store.IsTaskCompleted(sql.NullInt64{Int64: int64(tgID), Valid: tgID != 0}, taskID)
	if err != nil {
		b.l.Errorf("IsTaskCompleted error: %s", err.Error())
		return
	}

	if completed {
		msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, text.TaskAlreadyCompletedText[lcl])
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
		}
		return
	}

	b.BotAPI.Send(tgbotAPI.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID))

	task, err := b.store.GetTask(taskID)
	if err != nil {
		b.l.Errorf("GetTask error: %s", err.Error())
		return
	}

	switch task.Type {
	case bot.SubscribeTaskType:
		msgText := fmt.Sprintf(text.SubscribeTaskMessage[lcl], task.Reward, task.RewardType.ToText(), task.Link)

		msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, msgText)
		msg.ParseMode = tgbotAPI.ModeMarkdown
		msg.ReplyMarkup = tgbotAPI.NewInlineKeyboardMarkup(
			tgbotAPI.NewInlineKeyboardRow(
				tgbotAPI.NewInlineKeyboardButtonURL(text.SubscribeButtonText[lcl], task.Link),
				tgbotAPI.NewInlineKeyboardButtonData(text.CheckButtonText[lcl], fmt.Sprintf("checkSubscribe-%d-%d", tgID, taskID)),
			),
		)

		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
		}

		return
	case bot.AddToNicknameTaskType:
		msgText := fmt.Sprintf(text.AddToNicknameTaskMessage[lcl], task.Reward, task.RewardType.ToText())

		msg := tgbotAPI.NewMessage(cq.Message.Chat.ID, msgText)
		msg.ParseMode = tgbotAPI.ModeMarkdown
		msg.ReplyMarkup = tgbotAPI.NewInlineKeyboardMarkup(
			tgbotAPI.NewInlineKeyboardRow(
				tgbotAPI.NewInlineKeyboardButtonData(text.CheckButtonText[lcl], fmt.Sprintf("checkAddToNickname-%d-%d", tgID, taskID)),
			),
		)

		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %s", err.Error())
		}

		return
	}

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}

func (b *Bot) executeCqClaimBonus(cq *tgbotAPI.CallbackQuery, lcl locale.Locale) {
	info := strings.Split(cq.Data, "-")

	tgIDStr := info[1]
	tgID, err := strconv.Atoi(tgIDStr)
	if err != nil {
		b.l.Errorf("tgID string to int convertation error: %s", err.Error())
		return
	}

	if tgID != int(cq.From.ID) {
		return
	}

	timerIdStr := info[2]
	timerID, err := strconv.Atoi(timerIdStr)
	if err != nil {
		b.l.Errorf("timerID string to int convertation error: %s", err.Error())
		return
	}

	timer, err := b.store.GetTimer(int64(timerID))
	if err != nil {
		b.l.Errorf("GetTimer error: %w", err)
		return
	}

	if len(info) < 4 {
		b.l.Infoln("info len < 4")
		if timer.Type != nil && *timer.Type == timer_types.Periodical {
			timer.LastTrigger = time.Now().UTC()
			timer.NextTrigger = time.Now().UTC().Add(time.Duration(timer.PeriodSeconds.Int64) * time.Second)
			err = b.store.UpdateTimer(timer)
			if err != nil {
				b.l.Errorf("UpdateTimer error: %w", err)
				return
			}
	
		}
		nextNotificationTime := timer.NextTrigger.Format(time.DateTime)
		msg := tgbotAPI.NewEditMessageTextAndMarkup(cq.Message.Chat.ID, cq.Message.MessageID, fmt.Sprintf(text.BonusIsNoLongerActiveText[lcl], nextNotificationTime), keyboard.NotifyInlineKeyboard(timer, lcl))
		_, err = b.BotAPI.Send(msg)
		if err != nil {
			b.l.Errorf("Send message error: %w", err)
		}
		return
	}

	hash := info[3]

	key := fmt.Sprintf("%d-%d-%s", tgID, timerID, hash)

	amount, err := b.bonusStore.Get(context.Background(), key)
	if err != nil {
		if errors.Is(err, store.ErrEntityNotFound) {
			b.l.Infof("bonus %s is not found in bonus store", key)
			msg := tgbotAPI.NewEditMessageTextAndMarkup(cq.Message.Chat.ID, cq.Message.MessageID, text.BonusIsNoLongerActiveText[lcl], keyboard.NotifyInlineKeyboard(timer, lcl))
			_, err = b.BotAPI.Send(msg)
			if err != nil {
				b.l.Errorf("Send message error: %w", err)
			}
			return
		}

		b.l.Errorf("Get error: %w", err)
		return
	}

	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		b.l.Errorf("amount string to int convertation error: %s", err.Error())
		return
	}

	newBalance, err := b.store.IncreaseBalance(sql.NullInt64{Int64: int64(tgID), Valid: tgID != 0}, int64(amountInt))
	if err != nil {
		b.l.Errorf("IncreaseBalance error: %s", err.Error())
		return
	}

	err = b.bonusStore.Delete(context.Background(), key)
	if err != nil {
		b.l.Errorf("Delete error: %w", err)
		return
	}

	if timer.Type != nil && *timer.Type == timer_types.Periodical {
		timer.LastTrigger = time.Now().UTC()
		timer.NextTrigger = time.Now().UTC().Add(time.Duration(timer.PeriodSeconds.Int64) * time.Second)
		err = b.store.UpdateTimer(timer)
		if err != nil {
			b.l.Errorf("UpdateTimer error: %w", err)
			return
		}

	}

	msg := tgbotAPI.NewEditMessageTextAndMarkup(cq.Message.Chat.ID, cq.Message.MessageID, fmt.Sprintf(text.AlertBonusClaimedText[lcl], amountInt, timer.NextTrigger.Format(time.DateTime), utils.GetBalanceText(newBalance, lcl)), keyboard.NotifyInlineKeyboard(timer, lcl))
	msg.ParseMode = tgbotAPI.ModeMarkdown
	_, err = b.BotAPI.Send(msg)
	if err != nil {
		b.l.Errorf("Send message error: %w", err)
	}

	callback := tgbotAPI.NewCallback(cq.ID, "")
	b.BotAPI.Send(callback)
}
