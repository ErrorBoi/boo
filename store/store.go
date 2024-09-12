package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/errorboi/boo/internal/locale"
	bot_types "github.com/errorboi/boo/types/bot"
	"github.com/errorboi/boo/types/timer"
)

type UserStore interface {
	CreateIfNotExist(ctx context.Context, userID, referrerID sql.NullInt64, username string, lcl locale.Locale) (bool, error)
	Get(userID sql.NullInt64) (*bot_types.User, error)
	ListUsers() ([]bot_types.User, error)
	UpdateUser(user *bot_types.User) error

	GetBalance(userID sql.NullInt64) (int64, error)
	IncreaseBalance(userID sql.NullInt64, increase int64) (int64, error)
	SaveTx(userID sql.NullInt64, amount int64, description string) error

	IncrementReferralsAmount(userID sql.NullInt64) error

	CompleteTask(userID sql.NullInt64, taskID int) error
	IsTaskCompleted(userID sql.NullInt64, taskID int) (bool, error)
	GetTask(taskID int) (*bot_types.Task, error)
	GetTasks(tgID int) ([]bot_types.Task, error)
	GetAllTasks(tgID int) ([]bot_types.Task, error)

	CreateTimerByName(userID int64, name string) (*timer.Timer, error)
	CreateTimer(timer *timer.Timer) error
	ListTimers(userID int64) ([]timer.Timer, error)
	GetTimer(timerID int64) (*timer.Timer, error)
	UpdateTimer(timer *timer.Timer) error
	DeleteTimer(timerID int64) error
	GetReadyTimers() ([]*timer.Timer, error)
	ListPresetTimers() ([]*timer.Timer, error)
	GetPresetTimer(timerID int64) (*timer.Timer, error)

	GetLeaderboard() ([]bot_types.User, error)
	GetLeaderboardPosition(tgID int64) (int, float64, error)

	MigrateUp() error
}

type BonusStore interface {
	CreateWithTTL(ctx context.Context, key string, val interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type FileStore interface {
	LoadFile(name string) (string, error)
	SaveFile(name, fileID string) error
}
