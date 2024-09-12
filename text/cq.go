package text

import "github.com/errorboi/boo/internal/locale"

var (
	TaskAlreadyCompletedText = map[locale.Locale]string{
		locale.English: TaskAlreadyCompletedTextEN,
		locale.Russian: TaskAlreadyCompletedTextRU,
		locale.Uzbek:   TaskAlreadyCompletedTextUZ,
	}
	NeedSubscribeText = map[locale.Locale]string{
		locale.English: NeedSubscribeTextEN,
		locale.Russian: NeedSubscribeTextRU,
		locale.Uzbek:   NeedSubscribeTextUZ,
	}
	NeedBooNicknameText = map[locale.Locale]string{
		locale.English: NeedBooNicknameTextEN,
		locale.Russian: NeedBooNicknameTextRU,
		locale.Uzbek:   NeedBooNicknameTextUZ,
	}
	TaskCompleteSuccessText = map[locale.Locale]string{
		locale.English: TaskCompleteSuccessTextEN,
		locale.Russian: TaskCompleteSuccessTextRU,
		locale.Uzbek:   TaskCompleteSuccessTextUZ,
	}
)

const (
	TaskAlreadyCompletedTextEN = "You completed that task already (¬‿¬)"
	TaskAlreadyCompletedTextRU = "Ты уже выполнил это задание (¬‿¬)"
	TaskAlreadyCompletedTextUZ = "Siz ushbu vazifani avval bajarib bo'ldingiz (¬‿¬)"
	NeedSubscribeTextEN        = "You need to subscribe to the channel to complete the task"
	NeedSubscribeTextRU        = "Нужно подписаться на канал, чтобы выполнить задание"
	NeedSubscribeTextUZ        = "Siz vazifani bajarish uchun kanalga obuna bo'lishingiz kerak"
	NeedBooNicknameTextEN      = "You need to add 👻 to your nickname to complete the task"
	NeedBooNicknameTextRU      = "Нужно добавить 👻 в свой ник, чтобы выполнить задание"
	NeedBooNicknameTextUZ      = "Vazifani bajarish uchun 👻 ni nikizga qo'shishingiz kerak"
	TaskCompleteSuccessTextEN  = "You've earned %d %s for completing the task! 👻"
	TaskCompleteSuccessTextRU  = "Ты заработал %d %s за выполнение задания! 👻"
	TaskCompleteSuccessTextUZ  = "Siz vazifani bajarish uchun %d %s ga ega bo'ldingiz! 👻"
)
