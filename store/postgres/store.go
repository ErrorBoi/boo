package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/errorboi/boo/internal/locale"
	bot_types "github.com/errorboi/boo/types/bot"
	"github.com/errorboi/boo/types/timer"
	"github.com/errorboi/boo/types/user"
	"github.com/lib/pq"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	// "github.com/golang-migrate/migrate/v4"
	// "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	ErrNotFound            = errors.New("entity not found")
	ErrWalletAlreadyExists = errors.New("wallet is taken by another user")
)

type Params struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type DB struct {
	client *sqlx.DB
}

func New(p *Params) (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.DBName)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{client: db}, nil
}

func (db *DB) Close() error {
	return db.client.Close()
}

func (db *DB) Extensions() error {
	_, err := db.client.Exec(`CREATE EXTENSION "pgcrypto";`)
	return err
}

func (db *DB) NewAirdropDB() error {
	statement := fmt.Sprintf(`SELECT 'CREATE DATABASE %s'
		WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')`, "airdrop_bot", "airdrop_bot")
	_, err := db.client.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) NewUsersTable() error {
	_, err := db.client.Exec(`CREATE TABLE IF NOT EXISTS users(
		id UUID NOT NULL DEFAULT gen_random_uuid(),
		tg_id BIGINT UNIQUE NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    	CONSTRAINT pkey PRIMARY KEY ( id )
	 );`)
	return err
}

func (db *DB) NewBalanceTable() error {
	_, err := db.client.Exec(`CREATE TABLE IF NOT EXISTS balance(
		tg_id BIGINT UNIQUE NOT NULL,
		amount DECIMAL NOT NULL DEFAULT 0
	 );`)
	return err
}

func (db *DB) NewReferralTable() error {
	_, err := db.client.Exec(`CREATE TABLE IF NOT EXISTS referrals(
		tg_id BIGINT UNIQUE NOT NULL REFERENCES users(tg_id) ON DELETE CASCADE,
		referrer_id INTEGER,
		referrals INTEGER NOT NULL DEFAULT 0
	 );`)

	return err
}

func (db *DB) NewTaskTable() error {
	_, err := db.client.Exec(`CREATE TABLE IF NOT EXISTS tasks(
		id BIGSERIAL PRIMARY KEY,
		reward BIGINT NOT NULL,
		link TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    	chat_id BIGINT NOT NULL
	 );`)

	return err
}

func (db *DB) NewTaskUserTable() error {
	_, err := db.client.Exec(`CREATE TABLE IF NOT EXISTS tasks_users(
		task_id BIGINT REFERENCES tasks(id) ON DELETE CASCADE,
		tg_id BIGINT,
		created_at VARCHAR (50) NOT NULL
	 );`)

	return err
}

func (db *DB) InitialTasks() error {
	_, err := db.client.Exec(`INSERT INTO tasks
	(reward, link, chat_id)
	VALUES (1200, 'https://t.me/tokerkarlsen', -1002074185522),
	(1000, 'https://t.me/tokerkarlsen_chat', -1002145320176),
	(1000, 'https://t.me/tokerkarlsen_en_chat', -1002080041294);`)

	return err

}

func (db *DB) MigrateUp() error {
	driver, err := postgres.WithInstance(db.client.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err

	}
	return m.Up()
}

func (db *DB) NewUsersIndex() error {
	_, err := db.client.Exec(`CREATE INDEX IF NOT EXISTS users_tg_id_idx ON users(tg_id);`)

	return err
}

func (db *DB) NewBalanceIndex() error {
	_, err := db.client.Exec(`CREATE INDEX IF NOT EXISTS balance_tg_id_idx ON balance(tg_id);`)

	return err
}

func (db *DB) NewReferralIndex() error {
	_, err := db.client.Exec(`CREATE INDEX IF NOT EXISTS referrals_tg_id_idx ON referrals(tg_id);`)

	return err
}

func (db *DB) CreateIfNotExist(
	ctx context.Context,
	userID,
	referrerID sql.NullInt64,
	username string,
	lcl locale.Locale,
) (bool, error) {
	tx, err := db.client.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return false, err
	}

	res, err := tx.Exec(`INSERT INTO users
    	(tg_id, created_at, username, timer_limit, locale)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (tg_id) DO NOTHING;`, userID, time.Now().UTC(), username, user.DefaultTimersLimit, lcl.Shortcode())
	if err != nil {
		_ = tx.Rollback()

		return false, err
	}

	rowsAffected, err := res.RowsAffected()

	isNewUser := rowsAffected > 0

	_, err = tx.Exec(`INSERT INTO balance
		(tg_id, amount)
		VALUES ($1, $2)
		ON CONFLICT (tg_id) DO NOTHING;`, userID, 0)
	if err != nil {
		_ = tx.Rollback()

		return false, err
	}

	_, err = tx.Exec(`INSERT INTO referrals
    	(tg_id, referrer_id, referrals)
		VALUES ($1, $2, $3)
		ON CONFLICT (tg_id) DO NOTHING;`, userID, referrerID, 0)
	if err != nil {
		_ = tx.Rollback()

		return false, err
	}

	if err = tx.Commit(); err != nil {
		return false, err
	}

	return isNewUser, err
}

func (db *DB) Get(userID sql.NullInt64) (*bot_types.User, error) {
	var u bot_types.User
	err := db.client.Get(&u, `SELECT u.*, r.referrer_id, r.referrals
FROM users u
INNER JOIN referrals r
ON r.tg_id = u.tg_id
where u.tg_id = $1;`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &u, nil
}

func (db *DB) ListUsers() ([]bot_types.User, error) {
	var users []bot_types.User

	err := db.client.Select(&users, `SELECT * from users;`)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (db *DB) UpdateUser(user *bot_types.User) error {
	_, err := db.client.Exec(`UPDATE users
	SET username = $1, timer_limit = $2, locale = $3, wallet = $5, is_wallet_verified = $6
	WHERE id = $4;`, user.Username, user.TimerLimit, user.Locale, user.ID, user.Wallet, user.IsWalletVerified)
	if err != nil {
		if IsPgViolationError(err, ErrPgUniqueViolationError) {
			return ErrWalletAlreadyExists
		}
		return err
	}

	return nil
}

func (db *DB) GetBalance(userID sql.NullInt64) (int64, error) {
	var amount int64
	err := db.client.Get(&amount, `SELECT amount FROM balance
	WHERE tg_id = $1;`, userID)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

func (db *DB) IncreaseBalance(userID sql.NullInt64, increase int64) (int64, error) {
	row := db.client.QueryRow(`INSERT INTO balance (tg_id, amount)
	VALUES ($1, $2)
	ON CONFLICT(tg_id)
	DO UPDATE SET
	amount = balance.amount + EXCLUDED.amount
	RETURNING amount;`, userID, increase)

	var amountInt int64
	err := row.Scan(&amountInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrNotFound
		}

		return 0, err
	}

	return amountInt, nil
}

func (db *DB) SaveTx(userID sql.NullInt64, amount int64, description string) error {
	_, err := db.client.Exec(`INSERT INTO txs
	(tg_id, amount, description, created_at)
	VALUES ($1, $2, $3, $4);`, userID, amount, description, time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) IncrementReferralsAmount(userID sql.NullInt64) error {
	_, err := db.client.Exec(`UPDATE referrals
	SET referrals = referrals + 1
	WHERE tg_id = $1;`, userID)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CompleteTask(userID sql.NullInt64, taskID int) error {
	_, err := db.client.Exec(`INSERT INTO tasks_users
	(task_id, tg_id, created_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (task_id, tg_id)
	DO NOTHING;`, taskID, userID, time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) IsTaskCompleted(userID sql.NullInt64, taskID int) (bool, error) {
	var count int
	err := db.client.Get(&count, `SELECT COUNT(*) FROM tasks_users
	WHERE task_id = $1 AND tg_id = $2;`, taskID, userID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (db *DB) GetTask(taskID int) (*bot_types.Task, error) {
	var task bot_types.Task
	err := db.client.Get(&task, `SELECT * FROM tasks where id = $1;`, taskID)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (db *DB) GetTasks(tgID int) ([]bot_types.Task, error) {
	var tasks []bot_types.Task
	err := db.client.Select(&tasks, `
SELECT * FROM tasks
WHERE status = 'active' AND 
      id NOT IN (SELECT task_id FROM tasks_users WHERE tg_id = $1);`, tgID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (db *DB) GetAllTasks(tgID int) ([]bot_types.Task, error) {
	var tasks []bot_types.Task
	err := db.client.Select(&tasks, `
SELECT * FROM tasks
WHERE id NOT IN (SELECT task_id FROM tasks_users WHERE tg_id = $1);`, tgID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (db *DB) CreateTimerByName(userID int64, name string) (*timer.Timer, error) {
	var newTimer timer.Timer
	err := db.client.Get(&newTimer, `INSERT INTO timers
	(tg_id, name, created_at)
	VALUES ($1, $2, $3)
	RETURNING *;`, userID, name, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	return &newTimer, nil
}

func (db *DB) ListTimers(userID int64) ([]timer.Timer, error) {
	var timers []timer.Timer
	err := db.client.Select(&timers, `SELECT * FROM timers
	WHERE tg_id = $1;`, userID)
	if err != nil {
		return nil, err
	}

	return timers, nil
}

func (db *DB) GetTimer(timerID int64) (*timer.Timer, error) {
	var t timer.Timer
	err := db.client.Get(&t, `SELECT * FROM timers
	WHERE id = $1;`, timerID)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (db *DB) CreateTimer(timer *timer.Timer) error {
	_, err := db.client.Exec(`INSERT INTO timers
	(tg_id, name, timer_type, repeat_type, description, period_seconds, trigger_time, status, last_trigger, last_ack, next_trigger, link, img_url)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`, timer.TgID, timer.Name, timer.Type, timer.RepeatType, timer.Description, timer.PeriodSeconds, timer.TriggerTime, timer.Status, timer.LastTrigger, timer.LastAck, timer.NextTrigger, timer.Link, timer.ImgUrl)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateTimer(timer *timer.Timer) error {
	_, err := db.client.Exec(`UPDATE timers
	SET name = $1, timer_type = $2, repeat_type = $3, description = $4, period_seconds = $5, trigger_time = $6, status = $7, last_trigger = $8, last_ack = $9, next_trigger = $10, link = $11, notify_id = $12
	WHERE id = $13;`, timer.Name, timer.Type, timer.RepeatType, timer.Description, timer.PeriodSeconds, timer.TriggerTime, timer.Status, timer.LastTrigger, timer.LastAck, timer.NextTrigger, timer.Link, timer.NotifyID, timer.ID)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteTimer(timerID int64) error {
	_, err := db.client.Exec(`DELETE FROM timers
	WHERE id = $1;`, timerID)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetReadyTimers() ([]*timer.Timer, error) {
	var timers []*timer.Timer
	err := db.client.Select(&timers, `SELECT * FROM timers
	WHERE status = 'active'
	AND next_trigger <= $1
	;`, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	return timers, nil
}

func (db *DB) ListPresetTimers() ([]*timer.Timer, error) {
	var timers []*timer.Timer
	err := db.client.Select(&timers, `SELECT * FROM preset_timers WHERE status = 'active';`)
	if err != nil {
		return nil, err
	}

	return timers, nil
}

func (db *DB) GetPresetTimer(timerID int64) (*timer.Timer, error) {
	var t timer.Timer
	err := db.client.Get(&t, `SELECT * FROM preset_timers
	WHERE id = $1;`, timerID)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (db *DB) GetLeaderboard() ([]bot_types.User, error) {
	var users []bot_types.User
	err := db.client.Select(&users, `SELECT 
    username, 
    SS.ref_balance, 
    balance.amount AS balance, 
    (0.2 * SS.ref_balance + balance.amount) AS score 
FROM 
    users u
INNER JOIN
    (
        SELECT 
            r.referrer_id, 
            SUM(b.amount) AS ref_balance 
        FROM 
            referrals r
        INNER JOIN 
            balance b ON r.tg_id = b.tg_id
        INNER JOIN 
            users u ON r.tg_id = u.tg_id
        WHERE 
            referrer_id IS NOT NULL
--             AND u.created_at > '2024-06-17 08:00' 
            AND u.created_at < '2024-07-01 21:00'
            AND r.referrer_id NOT IN (321649773, 128883002)
        GROUP BY 
            r.referrer_id
        ORDER BY ref_balance DESC
    ) SS ON u.tg_id = SS.referrer_id
INNER JOIN 
    balance ON balance.tg_id = u.tg_id
ORDER BY 
    score DESC
LIMIT 10;
`)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (db *DB) GetLeaderboardPosition(tgID int64) (int, float64, error) {
	var res struct {
		Score    float64 `db:"score"`
		Position int     `db:"position"`
	}
	err := db.client.Get(&res, `SELECT score, position
FROM (
    SELECT
    DENSE_RANK() OVER (ORDER BY score DESC) AS position,
    tg_id,
    score
FROM
    (
    SELECT
    u.tg_id,
    username,
    SS.ref_balance,
    balance.amount AS amount,
    (0.2 * SS.ref_balance + balance.amount) AS score
FROM
    users u
INNER JOIN
    (
        SELECT
            r.referrer_id,
            SUM(b.amount) AS ref_balance
        FROM
            referrals r
        INNER JOIN
            balance b ON r.tg_id = b.tg_id
        INNER JOIN
            users u ON r.tg_id = u.tg_id
        WHERE
            referrer_id IS NOT NULL
--             AND u.created_at > '2024-06-17 08:00'
            AND u.created_at < '2024-07-01 09:00'
            AND r.referrer_id NOT IN (321649773, 128883002)
        GROUP BY
            r.referrer_id
        ORDER BY ref_balance DESC
    ) SS ON u.tg_id = SS.referrer_id
INNER JOIN
    balance ON balance.tg_id = u.tg_id
ORDER BY
    score DESC
    ) as balances
WHERE tg_id NOT IN (321649773, 128883002)
) AS leaderboard
where tg_id = $1;`, tgID)
	if err != nil {
		return 0, 0, err
	}

	return res.Position, res.Score, nil
}

func (db *DB) LoadFile(name string) (string, error) {
	var file bot_types.File
	err := db.client.Get(&file, `SELECT * FROM files
	WHERE name = $1;`, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrNotFound
		}
		return "", err
	}

	return file.FileID, nil
}

func (db *DB) SaveFile(name, fileID string) error {
	_, err := db.client.Exec(`INSERT INTO files
	(name, file_id)
	VALUES ($1, $2)
	ON CONFLICT (name) DO UPDATE SET file_id = $2;`, name, fileID)
	if err != nil {
		return err
	}

	return nil
}

func IsPgViolationError(err error, filter ...pq.ErrorCode) bool {
	pgerr, ok := err.(*pq.Error)
	if !ok {
		return false
	}

	if filter == nil {
		filter = []pq.ErrorCode{
			ErrPgUniqueViolationError,
			ErrPgForeignKeyViolation,
			ErrPgCheckViolation,
			ErrPgNotNullViolation,
		}
	}

	for _, typ := range filter {
		if pgerr.Code == typ {
			return true
		}
	}

	return false
}

const (
	ErrPgUniqueViolationError = "23505"
	ErrPgForeignKeyViolation  = "23503"
	ErrPgCheckViolation       = "23514"
	ErrPgNotNullViolation     = "23502"
)
