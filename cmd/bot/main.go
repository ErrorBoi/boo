package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/errorboi/boo/bot"
	"github.com/errorboi/boo/internal/config"
	"github.com/errorboi/boo/internal/notify"
	"github.com/errorboi/boo/internal/random"
	"github.com/errorboi/boo/internal/redis"
	schedule "github.com/errorboi/boo/internal/schedule"
	"github.com/errorboi/boo/internal/state_store"
	"github.com/errorboi/boo/internal/validate"
	postgres_store "github.com/errorboi/boo/store/postgres"
	tgbotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-migrate/migrate"
	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

// MainHandler responds to http request
func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm Bot!"))
}

func main() {
	atom := zap.NewAtomicLevel()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	logger := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.Lock(os.Stdout), atom))

	defer logger.Sync()

	atom.SetLevel(zap.DebugLevel)

	sugar := logger.Sugar()

	cfg := config.New()

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		sugar.Fatalf("Parse port error: %w", err)
	}

	store, err := postgres_store.New(&postgres_store.Params{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	})
	if err != nil {
		sugar.Fatalf("Database connection error: %s", err.Error())
	}

	defer store.Close()

	err = store.MigrateUp()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		sugar.Fatalf("Migrate up error: %w", err)
	}

	stepStore := state_store.NewInmemStateStore("timer")

	botAPI, err := tgbotAPI.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		sugar.Fatalf("Bot init error: %w", err)
	}

	botAPI.Debug = false

	redisStore := redis.New(cfg.RedisCfg)

	randomizer := random.NewRandomizer()

	notifier := notify.New(botAPI, store, sugar, redisStore, randomizer)

	go notifier.Run(context.TODO())

	scheduler := schedule.NewScheduler(botAPI, store, sugar, cfg.DebugMode, redisStore)

	go scheduler.Run()

	validator := validate.New()

	b, err := bot.New(botAPI, store, sugar, stepStore, validator, cfg, redisStore, notifier, store)
	if err != nil {
		sugar.Fatalf("Bot init error: %w", err)
	}

	http.HandleFunc("/", MainHandler)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	go b.InitUpdates()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	notifier.Close()
	scheduler.Close()

	log.Println("Graceful shutdown complete.")
}
