package keyboard

import (
	"fmt"

	"github.com/errorboi/boo/internal/locale"
	"github.com/errorboi/boo/internal/utils"
	"github.com/errorboi/boo/text"
	"github.com/errorboi/boo/types/bot"
	timer_types "github.com/errorboi/boo/types/timer"
	tgbotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetStartInlineKeyboard(lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	return tgbotAPI.NewInlineKeyboardMarkup(
		tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonData(text.AddTimersButtonText[lcl], "list_preset_timers"),
			tgbotAPI.NewInlineKeyboardButtonData(text.NewTimerButtonText[lcl], "new_timer"),
		),
		tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonData(text.BoobleJumpButtonText[lcl], "booble_jump"),
			tgbotAPI.NewInlineKeyboardButtonData(text.LeaderboardButtonText[lcl], "leaderboard"),
		),
		tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonData(text.MyTimersButtonText[lcl], "timers"),
			tgbotAPI.NewInlineKeyboardButtonData(text.TasksButtonText[lcl], "tasks"),
		),
		tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonData(text.ProfileButtonText[lcl], "profile"),
			tgbotAPI.NewInlineKeyboardButtonData(text.InviteButtonText[lcl], "invite"),
		),
	)
}

func GetMainReplyKeyboard(lcl locale.Locale) tgbotAPI.ReplyKeyboardMarkup {
	return tgbotAPI.NewReplyKeyboard(
		tgbotAPI.NewKeyboardButtonRow(
			tgbotAPI.NewKeyboardButton(text.AddTimersButtonText[lcl]),
		),
		tgbotAPI.NewKeyboardButtonRow(
			tgbotAPI.NewKeyboardButton(text.MyTimersButtonText[lcl]),
			tgbotAPI.NewKeyboardButton(text.TasksButtonText[lcl]),
		),
		tgbotAPI.NewKeyboardButtonRow(
			tgbotAPI.NewKeyboardButton(text.ProfileButtonText[lcl]),
			tgbotAPI.NewKeyboardButton(text.InviteButtonText[lcl]),
		),
	)
}

func GetTaskCenterInlineKeyboard(tgID int64, tasks []bot.Task, lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	for _, task := range tasks {
		rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonData(fmt.Sprintf("[%d %s] %s", task.Reward, task.RewardType.ToText(), task.Description), fmt.Sprintf("task-%d-%d", tgID, task.ID)),
		))
	}

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.HomepageButtonText[lcl], "homepage"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func GetTimersInlineKeyboard(timers []timer_types.Timer, lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	for _, t := range timers {
		status := "🔴"
		if t.Status == timer_types.Active {
			status = "🟢"
		}

		rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonData(fmt.Sprintf("[%s] %s", status, utils.EscapeMarkdown(t.Name)), fmt.Sprintf("timer-%d", t.ID)),
		))
	}

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.HomepageButtonText[lcl], "homepage"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func TimerInlineKeyboard(timer *timer_types.Timer, lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	if timer.Link.Valid {
		rows = append(rows,
			tgbotAPI.NewInlineKeyboardRow(
				tgbotAPI.NewInlineKeyboardButtonURL(fmt.Sprintf(text.GoToButtonText[lcl], utils.EscapeMarkdown(timer.Name)), timer.Link.String),
			),
		)
	}

	buttons := make([]tgbotAPI.InlineKeyboardButton, 0)

	if timer.Status == timer_types.Active {
		buttons = append(buttons, tgbotAPI.NewInlineKeyboardButtonData(text.StopButtonText[lcl], fmt.Sprintf("stop-%d", timer.ID)))
	} else {
		buttons = append(buttons, tgbotAPI.NewInlineKeyboardButtonData(text.StartButtonText[lcl], fmt.Sprintf("start-%d", timer.ID)))
	}

	if timer.Type != nil && *timer.Type == timer_types.Periodical {
		buttons = append(buttons, tgbotAPI.NewInlineKeyboardButtonData(text.RebootTimerButtonText[lcl], fmt.Sprintf("reboot-%d", timer.ID)))
	}

	rows = append(rows, buttons)

	rows = append(rows,
		tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonData(text.EditTimerButtonText[lcl], fmt.Sprintf("edit_timer-%d", timer.ID)),
			tgbotAPI.NewInlineKeyboardButtonData(text.DeleteTimerButtonText[lcl], fmt.Sprintf("delete-%d", timer.ID)),
		),
		tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonData(text.BackToTimersListButtonText[lcl], "timers"),
		),
	)

	return tgbotAPI.NewInlineKeyboardMarkup(rows...)
}

func EditTimerInlineKeyboard(timer *timer_types.Timer, lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.EditNameButtonText[lcl], fmt.Sprintf("editTimer_name-%d", timer.ID)),
		tgbotAPI.NewInlineKeyboardButtonData(text.EditDescriptionButtonText[lcl], fmt.Sprintf("editTimer_description-%d", timer.ID)),
	))

	typeRowButtons := make([]tgbotAPI.InlineKeyboardButton, 0)

	typeRowButtons = append(typeRowButtons,
		tgbotAPI.NewInlineKeyboardButtonData(text.EditTypeButtonText[lcl], fmt.Sprintf("editTimer_type-%d", timer.ID)))

	if timer.Type != nil {
		if *timer.Type == timer_types.Daily {
			typeRowButtons = append(typeRowButtons,
				tgbotAPI.NewInlineKeyboardButtonData(text.EditTriggerTimeButtonText[lcl], fmt.Sprintf("editTimer_triggerTime-%d", timer.ID)),
			)
		} else {
			typeRowButtons = append(typeRowButtons,
				tgbotAPI.NewInlineKeyboardButtonData(text.EditPeriodButtonText[lcl], fmt.Sprintf("editTimer_period-%d", timer.ID)),
			)
		}
	}

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(typeRowButtons...))

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		// tgbotAPI.NewInlineKeyboardButtonData("Edit Repeat Type", fmt.Sprintf("editTimer_repeatType-%d", timer.ID)),
		tgbotAPI.NewInlineKeyboardButtonData(text.EditLinkButtonText[lcl], fmt.Sprintf("editTimer_link-%d", timer.ID)),
	),
	)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.BackToTimerButtonText[lcl], fmt.Sprintf("timer-%d", timer.ID)),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func DeleteTimerInlineKeyboard(timerID int64, lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.ConfirmDeleteButtonText[lcl], fmt.Sprintf("confirm_delete-%d", timerID)),
		tgbotAPI.NewInlineKeyboardButtonData(text.CancelButtonText[lcl], fmt.Sprintf("timer-%d", timerID)),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func EditTimerTypeInlineKeyboard(timer *timer_types.Timer, lcl locale.Locale) *tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	daily := "Daily"
	periodical := "Periodical"
	if timer.Type != nil {
		if *timer.Type == timer_types.Daily {
			daily = fmt.Sprintf("✅ %s", daily)
		} else {
			periodical = fmt.Sprintf("✅ %s", periodical)
		}
	}

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(daily, fmt.Sprintf("edit_type-%d-daily", timer.ID)),
		tgbotAPI.NewInlineKeyboardButtonData(periodical, fmt.Sprintf("edit_type-%d-periodical", timer.ID)),
	))

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.BackToTimerEditButtonText[lcl], fmt.Sprintf("edit_timer-%d", timer.ID)),
	))

	res := tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)

	return &res
}

func SuccessMessageKeyboard(timer *timer_types.Timer, lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.BackToTimerButtonText[lcl], fmt.Sprintf("timer-%d", timer.ID)),
		tgbotAPI.NewInlineKeyboardButtonData(text.BackToTimersListButtonText[lcl], "timers"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func SuccessUserMessageKeyboard(lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.BackToProfileButtonText[lcl], "profile"),
		tgbotAPI.NewInlineKeyboardButtonData(text.HomepageButtonText[lcl], "homepage"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)

}

func BackToTimersListKeyboard(lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.BackToTimersListButtonText[lcl], "timers"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func CancelInlineKeyboard(timerID int64, lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.CancelButtonText[lcl], fmt.Sprintf("cancel-%d", timerID)),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func CancelProfileInlineKeyboard(lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.CancelButtonText[lcl], "profile"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)

}

func NotifyInlineKeyboard(timer *timer_types.Timer, lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	if timer.Link.Valid {
		rows = append(rows,
			tgbotAPI.NewInlineKeyboardRow(
				tgbotAPI.NewInlineKeyboardButtonURL(fmt.Sprintf(text.GoToButtonText[lcl], utils.EscapeMarkdown(timer.Name)), timer.Link.String),
			))
	}

	rows = append(rows,
		tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonData(text.ManageTimerButtonText[lcl], fmt.Sprintf("timer-%d", timer.ID)),
		),
	)

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func NotifyClaimbonusKeyboard(key string, lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.ClaimBonusButtonText[lcl], fmt.Sprintf("claimbonus-%s", key)),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func PresetTimersKeyboard(timers []*timer_types.Timer, lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	for _, timer := range timers {
		rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonData(timer.Name, fmt.Sprintf("show_preset-%d", timer.ID)),
		))
	}

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.HomepageButtonText[lcl], "homepage"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func PresetTimerInlineKeyboard(timer *timer_types.Timer, lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.AddToMyTimersButtonText[lcl], fmt.Sprintf("add_preset-%d", timer.ID)),
	))

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.BackToPresetTimersListButtonText[lcl], "list_preset_timers"),
	))

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.HomepageButtonText[lcl], "homepage"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func AddPresetTimerSuccessInlineKeyboard(lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.BackToPresetTimersListButtonText[lcl], "list_preset_timers"),
	))

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.HomepageButtonText[lcl], "homepage"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func StartInlineKeyboard(lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.AddTimersButtonText[lcl], "list_preset_timers"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func SelectLanguageKeyboard(lang *string) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	langEngText := text.EnglishButtonText
	langRusText := text.RussianButtonText
	langUzbText := text.UzbekButtonText

	if lang != nil {
		switch *lang {
		case "eng":
			langEngText = fmt.Sprintf("✅ %s", langEngText)
		case "rus":
			langRusText = fmt.Sprintf("✅ %s", langRusText)
		case "uzb":
			langUzbText = fmt.Sprintf("✅ %s", langUzbText)
		}
	}

	langEng := tgbotAPI.NewInlineKeyboardButtonData(langEngText, "lang-eng")
	langRus := tgbotAPI.NewInlineKeyboardButtonData(langRusText, "lang-rus")

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		langRus,
	), tgbotAPI.NewInlineKeyboardRow(
		langEng,
	))

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.HomepageButtonText[locale.English], "homepage"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func ProfileInlineKeyboard(user *bot.User, lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows,
		tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonData(text.EditLanguageButtonText[lcl], "language"),
		),
		tgbotAPI.NewInlineKeyboardRow(
			tgbotAPI.NewInlineKeyboardButtonData(text.EditWalletButtonText[lcl], fmt.Sprintf("editUser_wallet-%d", user.TgID)),
		))

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.HomepageButtonText[lcl], "homepage"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}

func DefaultInlineKeyboard(lcl locale.Locale) tgbotAPI.InlineKeyboardMarkup {
	rows := make([][]tgbotAPI.InlineKeyboardButton, 0)

	rows = append(rows, tgbotAPI.NewInlineKeyboardRow(
		tgbotAPI.NewInlineKeyboardButtonData(text.HomepageButtonText[lcl], "homepage"),
	))

	return tgbotAPI.NewInlineKeyboardMarkup(
		rows...,
	)
}
