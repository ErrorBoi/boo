package bot

import (
	"fmt"
	"time"

	"github.com/errorboi/boo/internal/locale"
	"github.com/errorboi/boo/internal/utils"
	"github.com/errorboi/boo/text"
	timer_types "github.com/errorboi/boo/types/timer"
)

func getProfileSlug(referrals int, lcl locale.Locale) string {
	var slug string

	switch {
	case referrals < 10:
		slug = text.ProfileSlug10Text[lcl]
	case referrals < 50:
		slug = text.ProfileSlug50Text[lcl]
	case referrals < 100:
		slug = text.ProfileSlug100Text[lcl]
	case referrals < 500:
		slug = text.ProfileSlug500Text[lcl]
	case referrals < 1000:
		slug = text.ProfileSlug1000Text[lcl]
	case referrals > 1000:
		slug = text.ProfileSlug1000PlusText[lcl]
	}

	slug = text.YourStatusText[lcl] + slug

	return fmt.Sprintf(slug, referrals)
}

func getTimerText(timer *timer_types.Timer, botName string, lcl locale.Locale) string {
	msgText := fmt.Sprintf("*%s*\n\n", utils.EscapeMarkdown(timer.Name))

	timerType := "🚫"
	if timer.Type != nil {
		timerType = string(*timer.Type)
	}
	msgText += fmt.Sprintf("*%s*: %s\n", text.TypeText[lcl], utils.EscapeMarkdown(timerType))

	switch timerType {
	case timer_types.Daily.String():
		triggerTime := "🚫"
		if timer.TriggerTime.Valid {
			triggerTime = timer.TriggerTime.String
		}

		msgText += fmt.Sprintf("*%s*: %s\n\n", text.TriggerTimeText[lcl], utils.EscapeMarkdown(triggerTime))
	case timer_types.Periodical.String():
		period := "🚫"
		if timer.PeriodSeconds.Valid {
			period = secondsToPeriod(timer.PeriodSeconds.Int64).Format("15:04:05")
		}

		msgText += fmt.Sprintf("*%s*: %s\n\n", text.PeriodText[lcl], period)
	}

	status := "🔴 OFF"
	if timer.Status == timer_types.Active {
		status = "🟢 ON"
	}
	msgText += fmt.Sprintf("*%s*: %s\n\n", text.StatusText[lcl], utils.EscapeMarkdown(status))

	description := "🚫"
	if timer.Description.Valid {
		description = timer.Description.String
	}
	msgText += fmt.Sprintf("*%s*: %s\n\n", text.DescriptionText[lcl], utils.EscapeMarkdown(description))

	nextTrigger := "🚫"
	if !timer.NextTrigger.IsZero() {
		nextTrigger = timer.NextTrigger.Format(time.DateTime)
	}
	msgText += fmt.Sprintf("*%s*: %s\n\n", text.NextAlertText[lcl], utils.EscapeMarkdown(nextTrigger))

	msgText += fmt.Sprintf("*%s*: %s", text.ShareThisTimerText[lcl], getInviteLink(botName, timer.TgID, &timer.ID))

	return msgText
}

func getPresetTimerText(timer *timer_types.Timer, lcl locale.Locale) string {
	msgText := fmt.Sprintf("*%s*\n\n", utils.EscapeMarkdown(timer.Name))

	timerType := "🚫"
	if timer.Type != nil {
		timerType = string(*timer.Type)
	}
	msgText += fmt.Sprintf("*%s*: %s\n\n", text.TypeText[lcl], utils.EscapeMarkdown(timerType))

	description := "🚫"
	if timer.Description.Valid {
		description = timer.Description.String
	}
	msgText += fmt.Sprintf("*%s*: %s\n\n", text.DescriptionText[lcl], utils.EscapeMarkdown(description))

	switch timerType {
	case timer_types.Daily.String():
		triggerTime := "🚫"
		if timer.TriggerTime.Valid {
			triggerTime = timer.TriggerTime.String
		}

		msgText += fmt.Sprintf("*%s*: %s\n", text.TriggerTimeText[lcl], utils.EscapeMarkdown(triggerTime))
	case timer_types.Periodical.String():
		period := "🚫"
		if timer.PeriodSeconds.Valid {
			period = secondsToPeriod(timer.PeriodSeconds.Int64).Format("15:04:05")
		}

		msgText += fmt.Sprintf("*%s*: %s\n", text.PeriodText[lcl], period)
	}

	return msgText
}

func periodToSeconds(period time.Time) int64 {
	return int64(period.Hour()*3600 + period.Minute()*60 + period.Second())
}

func secondsToPeriod(seconds int64) time.Time {
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60

	return time.Date(0, 0, 0, int(hours), int(minutes), int(seconds), 0, time.UTC)
}

func getInviteLink(botName string, tgID int64, timerID *int64) string {
	link := fmt.Sprintf("https://t.me/%s?start=r%d", botName, tgID)

	if timerID != nil {
		link += fmt.Sprintf("-t%d", *timerID)
	}

	return fmt.Sprintf("`%s`", link)
}
