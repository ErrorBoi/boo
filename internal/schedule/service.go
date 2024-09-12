package schedule

import (
	"context"
	"strings"
	"time"

	"github.com/errorboi/boo/internal/utils"
	"github.com/errorboi/boo/store"
	tgbotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Scheduler interface {
	UpdateLeaderBoard()
	Run()
	Close()
}

type scheduler struct {
	botAPI     *tgbotAPI.BotAPI
	userStore  store.UserStore
	l          *zap.SugaredLogger
	bonusStore store.BonusStore
	quit       chan struct{}
	debugMode  bool
}

func NewScheduler(
	botAPI *tgbotAPI.BotAPI,
	userStore store.UserStore,
	l *zap.SugaredLogger,
	debugMode bool,
	bonusStore store.BonusStore,
) Scheduler {
	return &scheduler{
		botAPI:     botAPI,
		userStore:  userStore,
		l:          l,
		quit:       make(chan struct{}),
		debugMode:  debugMode,
		bonusStore: bonusStore,
	}
}

func (s *scheduler) Run() {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ticker.C:
			go s.UpdateLeaderBoard()
		case <-s.quit:
			ticker.Stop()
			return
		}
	}
}

func (s *scheduler) Close() {
	close(s.quit)
}

func (s *scheduler) UpdateLeaderBoard() {
	users, err := s.userStore.GetLeaderboard()
	if err != nil {
		s.l.Errorf("get leaderboard error: %w", err)
		return
	}

	res := utils.FormatLeaderboard(users)

	err = s.bonusStore.CreateWithTTL(context.Background(), "leaderboard", res, 24*time.Hour)
	if err != nil {
		s.l.Errorf("cache leaderboard error: %w", err)
		return
	}

	if s.debugMode {
		// beta boo
		msg := tgbotAPI.NewEditMessageText(-1002135004893, 6, res)
		msg.ParseMode = tgbotAPI.ModeHTML

		_, err = s.botAPI.Send(msg)
		if err != nil {
			s.l.Errorf("send leaderboard error: %w", err)
			return
		}
	} else {
		// @BooTimer
		// msg := tgbotAPI.NewEditMessageText(-1002097798490, 14, res)
		// msg.ParseMode = tgbotAPI.ModeHTML
		//
		// _, err = s.botAPI.Send(msg)
		// if err != nil {
		// 	s.l.Errorf("send leaderboard error: %w", err)
		// 	return
		// }

		// @BooDrops
		// msg = tgbotAPI.NewEditMessageText(-1002038494379, 29, res)
		// msg.ParseMode = tgbotAPI.ModeHTML
		//
		// _, err = s.botAPI.Send(msg)
		// if err != nil {
		// 	s.l.Errorf("send leaderboard error: %w", err)
		// 	return
		// }
	}
}

func addSpaces(s string, n int) string {
	if len([]rune(s)) > n {
		return s[:n-3] + "..."
	}

	repeatCount := n - len([]rune(s))
	if repeatCount < 0 {
		repeatCount = 0
	}

	return s + strings.Repeat(" ", repeatCount)
}

func containsEmoji(s string) bool {
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r >= 128 {
			return true
		}
	}

	return false
}
