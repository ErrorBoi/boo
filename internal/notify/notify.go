package notify

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/errorboi/boo/internal/locale"
	"github.com/errorboi/boo/internal/random"
	"github.com/errorboi/boo/internal/utils"
	"github.com/errorboi/boo/keyboard"
	"github.com/errorboi/boo/store"
	"github.com/errorboi/boo/text"
	"github.com/errorboi/boo/types/timer"
	tgbotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Notifier interface {
	Notify(timer *timer.Timer, key string, lcl locale.Locale) error
	BroadcastMessage(message string, fileID *string) error
	Run(ctx context.Context)
	Close()
}

type notifier struct {
	botAPI     *tgbotAPI.BotAPI
	userStore  store.UserStore
	l          *zap.SugaredLogger
	quit       chan struct{}
	bonusStore store.BonusStore
	rnd        *random.Randomizer
}

func New(
	botAPI *tgbotAPI.BotAPI,
	userStore store.UserStore,
	l *zap.SugaredLogger,
	bonusStore store.BonusStore,
	rnd *random.Randomizer,
) Notifier {
	return &notifier{
		botAPI:     botAPI,
		userStore:  userStore,
		l:          l,
		quit:       make(chan struct{}),
		bonusStore: bonusStore,
		rnd:        rnd,
	}
}

func (n *notifier) Notify(timer *timer.Timer, key string, lcl locale.Locale) error {
	msgText := fmt.Sprintf(text.TimerRingingText[lcl], utils.EscapeMarkdown(timer.Name))

	msg := tgbotAPI.NewMessage(timer.TgID, msgText)
	msg.ParseMode = tgbotAPI.ModeMarkdown

	msg.ReplyMarkup = keyboard.NotifyClaimbonusKeyboard(key, lcl)

	res, err := n.botAPI.Send(msg)
	if err != nil {
		return err
	}

	if timer.NotifyID.Valid {
		n.botAPI.Send(tgbotAPI.DeleteMessageConfig{
			ChatID:    timer.TgID,
			MessageID: int(timer.NotifyID.Int64),
		})
	}

	timer.NotifyID = sql.NullInt64{Int64: int64(res.MessageID), Valid: true}

	err = n.userStore.UpdateTimer(timer)
	if err != nil {
		n.l.Errorf("Update timer error: %w", err)
	}

	return nil
}

func (n *notifier) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ticker.C:
			go n.triggerTimers(ctx)
		case <-n.quit:
			ticker.Stop()
			return
		}
	}
}

func (n *notifier) Close() {
	close(n.quit)
}

func (n *notifier) triggerTimers(ctx context.Context) {
	timers, err := n.userStore.GetReadyTimers()
	if err != nil {
		n.l.Errorf("Get ready timers error: %w", err)
		return
	}
	n.l.Infoln("Timers to notify: ", len(timers))

	failedNotifies := 0
	for _, t := range timers {
		period := t.NextTrigger.Sub(t.LastTrigger)
		bonusPeriod := 2 * time.Hour

		t.LastTrigger = time.Now().UTC()

		switch *t.Type {
		case timer.Daily:
			if !t.TriggerTime.Valid || t.TriggerTime.String == "" {
				continue
			}
			triggerTime, err := time.Parse("15:04", t.TriggerTime.String)
			if err != nil {
				n.l.Errorf("Parse trigger time error: %w", err)
				continue
			}
			t.NextTrigger = time.Date(t.LastTrigger.Year(), t.LastTrigger.Month(), t.LastTrigger.Day(), triggerTime.Hour(), triggerTime.Minute(), triggerTime.Second(), 0, time.UTC).Add(24 * time.Hour)
		case timer.Periodical:
			if !t.PeriodSeconds.Valid || t.PeriodSeconds.Int64 == 0 {
				continue
			}

			timerPeriod := time.Duration(t.PeriodSeconds.Int64) * time.Second
			if bonusPeriod > timerPeriod {
				bonusPeriod = timerPeriod
			}
			t.NextTrigger = t.LastTrigger.Add(timerPeriod)
		}

		user, err := n.userStore.Get(sql.NullInt64{Int64: t.TgID, Valid: true})
		if err != nil {
			n.l.Errorf("Get user error: %w", err)
			continue
		}

		err = n.userStore.UpdateTimer(t)
		if err != nil {
			n.l.Errorf("Update timer error: %w", err)
		}

		key := fmt.Sprintf("%d-%d-%s", t.TgID, t.ID, n.rnd.String(10))

		err = n.Notify(t, key, locale.LangToLocale(user.Locale.String))
		if err != nil {
			failedNotifies++
		}

		// 12 tokens per hour
		amount := period.Minutes() / 5

		if amount < 1 {
			amount = 1
		}

		err = n.bonusStore.CreateWithTTL(ctx, key, int64(amount), bonusPeriod)
		if err != nil {
			n.l.Errorf("CreateWithTTL bonus error: %s", err.Error())
		}
	}

	n.l.Infof("%d timers notified, %d failed", len(timers)-failedNotifies, failedNotifies)
}

func (n *notifier) BroadcastMessage(message string, fileID *string) error {
	users, err := n.userStore.ListUsers()
	if err != nil {
		n.l.Errorf("Get all users error: %w", err)
		return err
	}

	lenUsers := len(users)

	var failedNotifies int
	for i, u := range users {
		var msg tgbotAPI.Chattable
		if fileID != nil {
			photoMsg := tgbotAPI.NewPhoto(u.TgID, tgbotAPI.FileID(*fileID))
			photoMsg.Caption = utils.EscapeMarkdown(message)
			photoMsg.ParseMode = tgbotAPI.ModeMarkdown

			msg = photoMsg
		} else {
			textMsg := tgbotAPI.NewMessage(u.TgID, utils.EscapeMarkdown(message))
			textMsg.ParseMode = tgbotAPI.ModeMarkdown

			msg = textMsg
		}

		_, err = n.botAPI.Send(msg)
		if err != nil {
			n.l.Errorf("[%d/%d] Send message error: %w", i+1, lenUsers, err)
			failedNotifies++
		}

		time.Sleep(50 * time.Millisecond)
	}

	n.l.Infof("Sending %d messages. Success: %d, Failed: %d", len(users), len(users)-failedNotifies, failedNotifies)

	return nil
}
