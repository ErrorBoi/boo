package utils

import (
	"fmt"
	"time"

	"github.com/errorboi/boo/internal/locale"
	"github.com/errorboi/boo/text"
	"github.com/errorboi/boo/types/bot"
	tgbotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetBalanceText(balance int64, lcl locale.Locale) string {
	symbols := []string{"", "k", "M", "B"}

	var i int
	for i = 0; i < len(symbols) && balance >= 1_000_000; i++ {
		balance /= 1000
	}

	return fmt.Sprintf(text.BalanceText[lcl], balance, symbols[i])
}

func GetBalance(balance int64, lcl locale.Locale) string {
	symbols := []string{"", "k", "M", "B"}

	var i int
	for i = 0; i < len(symbols) && balance >= 1_000_000; i++ {
		balance /= 1000
	}

	return fmt.Sprintf("%d%s", balance, symbols[i])
}

func FormatLeaderboard(users []bot.User) string {
	places := []string{"🥇", "🥈", "🥉", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"}

	res := "🏆 <b>Leaderboard</b>\n\n"

	for i, user := range users {
		res += fmt.Sprintf(
			"%s %s - %s\n",
			places[i],
			user.Username,
			GetBalance(user.Score.IntPart(), locale.English),
		)
	}

	res += fmt.Sprintf("\n\n<b>Score = Balance + 0.2 * Total Ref. Balance</b>\n\n<b>Updated at</b>: %s (UTC Time)\n\nClaim your bonus 🎁  at <b>@BooTimerBot</b>", time.Now().UTC().Format("2006-01-02 15:04:05"))

	return res
}

func EscapeMarkdown(text string) string {
	return tgbotAPI.EscapeText(tgbotAPI.ModeMarkdown, text)
}
