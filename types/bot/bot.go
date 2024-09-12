package bot

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ReferralBonus int64 = 500

	Admins = map[int64]struct{}{
		128883002: {},
	}
)

const (
	ChannelID = -1002074185522
	ChatRuID  = -1002145320176
	ChatEnID  = -1002080041294
)

type User struct {
	ID               uuid.UUID       `db:"id"`
	TgID             int64           `db:"tg_id"`
	CreatedAt        time.Time       `db:"created_at"`
	Username         string          `db:"username"`
	TimerLimit       int             `db:"timer_limit"`
	ReferrerID       sql.NullString  `db:"referrer_id"`
	Balance          decimal.Decimal `db:"balance"`
	ReferralsAmount  int             `db:"referrals"`
	Timezone         string          `db:"timezone"`
	Locale           sql.NullString  `db:"locale"`
	RefBalance       decimal.Decimal `db:"ref_balance"`
	Score            decimal.Decimal `db:"score"`
	Wallet           sql.NullString  `db:"wallet"`
	IsWalletVerified bool            `db:"is_wallet_verified"`
	IsNewUser        bool
}

type Task struct {
	ID          int        `db:"id"`
	Type        TaskType   `db:"task_type"`
	Reward      int        `db:"reward"`
	RewardType  RewardType `db:"reward_type"`
	Link        string     `db:"link"`
	CreatedAt   time.Time  `db:"created_at"`
	ChatID      int64      `db:"chat_id"`
	Description string     `db:"description"`
	Status      string     `db:"status"`
}

type File struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	FileID string `db:"file_id"`
}
