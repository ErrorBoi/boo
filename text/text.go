package text

import "github.com/errorboi/boo/internal/locale"

var (
	PrestartText = map[locale.Locale]string{
		locale.English: PrestartTextEN,
		locale.Russian: PrestartTextRU,
	}
	StartText = map[locale.Locale]string{
		locale.English: StartTextEN,
		locale.Russian: StartTextRU,
	}
	ProfileText = map[locale.Locale]string{
		locale.English: ProfileTextEN,
		locale.Russian: ProfileTextRU,
	}
	ReferralAcceptedText = map[locale.Locale]string{
		locale.English: ReferralAcceptedTextEN,
		locale.Russian: ReferralAcceptedTextRU,
	}
	InviteLinkText = map[locale.Locale]string{
		locale.English: InviteLinkTextEN,
		locale.Russian: InviteLinkTextRU,
		locale.Uzbek:   InviteLinkTextUZ,
	}
	NewTimerText = map[locale.Locale]string{
		locale.English: NewTimerTextEN,
		locale.Russian: NewTimerTextRU,
	}
	TimerNameText = map[locale.Locale]string{
		locale.English: TimerNameTextEN,
		locale.Russian: TimerNameTextRU,
	}
	InvalidTimerNameText = map[locale.Locale]string{
		locale.English: InvalidTimerNameTextEN,
		locale.Russian: InvalidTimerNameTextRU,
	}
	InvalidTimerDescriptionText = map[locale.Locale]string{
		locale.English: InvalidTimerDescriptionTextEN,
		locale.Russian: InvalidTimerDescriptionTextRU,
	}
	TimerSetupSuccessText = map[locale.Locale]string{
		locale.English: TimerSetupSuccessTextEN,
		locale.Russian: TimerSetupSuccessTextRU,
	}
	TimerUpdateSuccessText = map[locale.Locale]string{
		locale.English: TimerUpdateSuccessTextEN,
		locale.Russian: TimerUpdateSuccessTextRU,
	}
	UserUpdateSuccessText = map[locale.Locale]string{
		locale.English: UserUpdateSuccessTextEN,
		locale.Russian: UserUpdateSuccessTextRU,
	}
	NoTimersText = map[locale.Locale]string{
		locale.English: NoTimersTextEN,
		locale.Russian: NoTimersTextRU,
	}
	YourTimersText = map[locale.Locale]string{
		locale.English: YourTimersTextEN,
		locale.Russian: YourTimersTextRU,
	}
	TimersLimitText = map[locale.Locale]string{
		locale.English: TimersLimitTextEN,
		locale.Russian: TimersLimitTextRU,
	}
	TimerAlreadyAddedText = map[locale.Locale]string{
		locale.English: TimerAlreadyAddedTextEN,
		locale.Russian: TimerAlreadyAddedTextRU,
	}
	PresetTimersText = map[locale.Locale]string{
		locale.English: PresetTimersTextEN,
		locale.Russian: PresetTimersTextRU,
	}
	EditNameText = map[locale.Locale]string{
		locale.English: EditNameTextEN,
		locale.Russian: EditNameTextRU,
	}
	EditDescriptionText = map[locale.Locale]string{
		locale.English: EditDescriptionTextEN,
		locale.Russian: EditDescriptionTextRU,
	}
	EditTypeText = map[locale.Locale]string{
		locale.English: EditTypeTextEN,
		locale.Russian: EditTypeTextRU,
	}
	EditTriggerTimeText = map[locale.Locale]string{
		locale.English: EditTriggerTimeTextEN,
		locale.Russian: EditTriggerTimeTextRU,
	}
	EditPeriodText = map[locale.Locale]string{
		locale.English: EditPeriodTextEN,
		locale.Russian: EditPeriodTextRU,
	}
	EditRepeatTypeText = map[locale.Locale]string{
		locale.English: EditRepeatTypeTextEN,
		locale.Russian: EditRepeatTypeTextRU,
	}
	EditLinkText = map[locale.Locale]string{
		locale.English: EditLinkTextEN,
		locale.Russian: EditLinkTextRU,
	}
	EditUserWalletText = map[locale.Locale]string{
		locale.English: EditUserWalletTextEN,
		locale.Russian: EditUserWalletTextRU,
	}
	SelectLanguageText = map[locale.Locale]string{
		locale.English: SelectLanguageTextEN,
		locale.Russian: SelectLanguageTextRU,
		locale.Uzbek:   SelectLanguageTextUZ,
	}
	AllTasksFinishedText = map[locale.Locale]string{
		locale.English: AllTasksFinishedTextEN,
		locale.Russian: AllTasksFinishedTextRU,
	}
	TaskCenterText = map[locale.Locale]string{
		locale.English: TaskCenterTextEN,
		locale.Russian: TaskCenterTextRU,
	}
	SubscribeTaskMessage = map[locale.Locale]string{
		locale.English: SubscribeTaskMessageEN,
		locale.Russian: SubscribeTaskMessageRU,
	}
	AddToNicknameTaskMessage = map[locale.Locale]string{
		locale.English: AddToNicknameTaskMessageEN,
		locale.Russian: AddToNicknameTaskMessageRU,
	}
	PeriodTooSmallText = map[locale.Locale]string{
		locale.English: PeriodTooSmallTextEN,
		locale.Russian: PeriodTooSmallTextRU,
	}
	BonusIsNoLongerActiveText = map[locale.Locale]string{
		locale.English: BonusIsNoLongerActiveTextEN,
		locale.Russian: BonusIsNoLongerActiveTextRU,
	}
	AlertBonusClaimedText = map[locale.Locale]string{
		locale.English: AlertBonusClaimedTextEN,
		locale.Russian: AlertBonusClaimedTextRU,
	}
	MintReceivedText = map[locale.Locale]string{
		locale.English: MintReceivedTextEN,
		locale.Russian: MintReceivedTextRU,
	}
	YourPositionText = map[locale.Locale]string{
		locale.English: YourPositionTextEN,
		locale.Russian: YourPositionTextRU,
	}
	LeaderboardGoogleText = map[locale.Locale]string{
		locale.English: LeaderboardGoogleTextEN,
		locale.Russian: LeaderboardGoogleTextRU,
	}
	BoobleJumpText = map[locale.Locale]string{
		locale.English: BoobleJumpTextEN,
		locale.Russian: BoobleJumpTextRU,
	}
	TapToSeeReadyTimersText = map[locale.Locale]string{
		locale.English: TapToSeeReadyTimersTextEN,
		locale.Russian: TapToSeeReadyTimersTextRU,
	}
	YouCannotPerformThisText = map[locale.Locale]string{
		locale.English: YouCannotPerformThisTextEN,
		locale.Russian: YouCannotPerformThisTextRU,
	}
	TimerTypeNotSetText = map[locale.Locale]string{
		locale.English: TimerTypeNotSetTextEN,
		locale.Russian: TimerTypeNotSetTextRU,
	}
	TimerAlertTimeNotSetText = map[locale.Locale]string{
		locale.English: TimerAlertTimeNotSetTextEN,
		locale.Russian: TimerAlertTimeNotSetTextRU,
	}
	TimerPeriodIsNotSetText = map[locale.Locale]string{
		locale.English: TimerPeriodIsNotSetTextEN,
		locale.Russian: TimerPeriodIsNotSetTextRU,
	}
	ConfirmDeleteText = map[locale.Locale]string{
		locale.English: ConfirmDeleteTextEN,
		locale.Russian: ConfirmDeleteTextRU,
	}
	TimerDeletedText = map[locale.Locale]string{
		locale.English: TimerDeletedTextEN,
		locale.Russian: TimerDeletedTextRU,
	}
	TimerLimitReachedText = map[locale.Locale]string{
		locale.English: TimerLimitReachedTextEN,
		locale.Russian: TimerLimitReachedTextRU,
	}
	TimerAddedText = map[locale.Locale]string{
		locale.English: TimerAddedTextEN,
		locale.Russian: TimerAddedTextRU,
	}
	GiveawayText = map[locale.Locale]string{
		locale.English: GiveawayTextEN,
		locale.Russian: GiveawayTextRU,
		locale.Uzbek:   GiveawayTextUZ,
	}
	StopButtonText = map[locale.Locale]string{
		locale.English: StopButtonTextEN,
		locale.Russian: StopButtonTextRU,
	}
	StartButtonText = map[locale.Locale]string{
		locale.English: StartButtonTextEN,
		locale.Russian: StartButtonTextRU,
	}
	EditTimerButtonText = map[locale.Locale]string{
		locale.English: EditTimerButtonTextEN,
		locale.Russian: EditTimerButtonTextRU,
	}
	DeleteTimerButtonText = map[locale.Locale]string{
		locale.English: DeleteTimerButtonTextEN,
		locale.Russian: DeleteTimerButtonTextRU,
	}
	BackToTimersListButtonText = map[locale.Locale]string{
		locale.English: BackToTimersListButtonTextEN,
		locale.Russian: BackToTimersListButtonTextRU,
	}
	EditNameButtonText = map[locale.Locale]string{
		locale.English: EditNameButtonTextEN,
		locale.Russian: EditNameButtonTextRU,
	}
	EditDescriptionButtonText = map[locale.Locale]string{
		locale.English: EditDescriptionButtonTextEN,
		locale.Russian: EditDescriptionButtonTextRU,
	}
	EditTypeButtonText = map[locale.Locale]string{
		locale.English: EditTypeButtonTextEN,
		locale.Russian: EditTypeButtonTextRU,
	}
	EditTriggerTimeButtonText = map[locale.Locale]string{
		locale.English: EditTriggerTimeButtonTextEN,
		locale.Russian: EditTriggerTimeButtonTextRU,
	}
	EditPeriodButtonText = map[locale.Locale]string{
		locale.English: EditPeriodButtonTextEN,
		locale.Russian: EditPeriodButtonTextRU,
	}
	EditLinkButtonText = map[locale.Locale]string{
		locale.English: EditLinkButtonTextEN,
		locale.Russian: EditLinkButtonTextRU,
	}
	BackToTimerButtonText = map[locale.Locale]string{
		locale.English: BackToTimerButtonTextEN,
		locale.Russian: BackToTimerButtonTextRU,
	}
	BackToProfileButtonText = map[locale.Locale]string{
		locale.English: BackToProfileButtonTextEN,
		locale.Russian: BackToProfileButtonTextRU,
	}
	ConfirmDeleteButtonText = map[locale.Locale]string{
		locale.English: ConfirmDeleteButtonTextEN,
		locale.Russian: ConfirmDeleteButtonTextRU,
	}
	CancelButtonText = map[locale.Locale]string{
		locale.English: CancelButtonTextEN,
		locale.Russian: CancelButtonTextRU,
	}
	BackToTimerEditButtonText = map[locale.Locale]string{
		locale.English: BackToTimerEditButtonTextEN,
		locale.Russian: BackToTimerEditButtonTextRU,
	}
	GoToButtonText = map[locale.Locale]string{
		locale.English: GoToButtonTextEN,
		locale.Russian: GoToButtonTextRU,
	}
	ManageTimerButtonText = map[locale.Locale]string{
		locale.English: ManageTimerButtonTextEN,
		locale.Russian: ManageTimerButtonTextRU,
	}
	RebootTimerButtonText = map[locale.Locale]string{
		locale.English: RebootTimerButtonTextEN,
		locale.Russian: RebootTimerButtonTextRU,
	}
	ClaimBonusButtonText = map[locale.Locale]string{
		locale.English: ClaimBonusButtonTextEN,
		locale.Russian: ClaimBonusButtonTextRU,
	}
	AddToMyTimersButtonText = map[locale.Locale]string{
		locale.English: AddToMyTimersButtonTextEN,
		locale.Russian: AddToMyTimersButtonTextRU,
	}
	BackToPresetTimersListButtonText = map[locale.Locale]string{
		locale.English: BackToPresetTImersListButtonTextEN,
		locale.Russian: BackToPresetTImersListButtonTextRU,
	}
	EditLanguageButtonText = map[locale.Locale]string{
		locale.English: EditLanguageButtonTextEN,
		locale.Russian: EditLanguageButtonTextRU,
	}
	EditWalletButtonText = map[locale.Locale]string{
		locale.English: EditWalletButtonTextEN,
		locale.Russian: EditWalletButtonTextRU,
	}
	ProfileSlug10Text = map[locale.Locale]string{
		locale.English: ProfileSlug10TextEN,
		locale.Russian: ProfileSlug10TextRU,
	}
	ProfileSlug50Text = map[locale.Locale]string{
		locale.English: ProfileSlug50TextEN,
		locale.Russian: ProfileSlug50TextRU,
	}
	ProfileSlug100Text = map[locale.Locale]string{
		locale.English: ProfileSlug100TextEN,
		locale.Russian: ProfileSlug100TextRU,
	}
	ProfileSlug500Text = map[locale.Locale]string{
		locale.English: ProfileSlug500TextEN,
		locale.Russian: ProfileSlug500TextRU,
	}
	ProfileSlug1000Text = map[locale.Locale]string{
		locale.English: ProfileSlug1000TextEN,
		locale.Russian: ProfileSlug1000TextRU,
	}
	ProfileSlug1000PlusText = map[locale.Locale]string{
		locale.English: ProfileSlug1000PlusTextEN,
		locale.Russian: ProfileSlug1000PlusTextRU,
	}
	YourStatusText = map[locale.Locale]string{
		locale.English: YourStatusTextEN,
		locale.Russian: YourStatusTextRU,
	}
	BalanceText = map[locale.Locale]string{
		locale.English: BalanceTextEN,
		locale.Russian: BalanceTextRU,
	}
	TimerText = map[locale.Locale]string{
		locale.English: TimerTextEN,
		locale.Russian: TimerTextRU,
	}
	ReadyTimerText = map[locale.Locale]string{
		locale.English: ReadyTimerTextEN,
		locale.Russian: ReadyTimerTextRU,
	}
	TypeText = map[locale.Locale]string{
		locale.English: TypeTextEN,
		locale.Russian: TypeTextRU,
	}
	StatusText = map[locale.Locale]string{
		locale.English: StatusTextEN,
		locale.Russian: StatusTextRU,
	}
	DescriptionText = map[locale.Locale]string{
		locale.English: DescriptionTextEN,
		locale.Russian: DescriptionTextRU,
	}
	TriggerTimeText = map[locale.Locale]string{
		locale.English: TriggerTimeTextEN,
		locale.Russian: TriggerTimeTextRU,
	}
	PeriodText = map[locale.Locale]string{
		locale.English: PeriodTextEN,
		locale.Russian: PeriodTextRU,
	}
	LastAlertText = map[locale.Locale]string{
		locale.English: LastAlertTextEN,
		locale.Russian: LastAlertTextRU,
	}
	NextAlertText = map[locale.Locale]string{
		locale.English: NextAlertTextEN,
		locale.Russian: NextAlertTextRU,
	}
	LinkText = map[locale.Locale]string{
		locale.English: LinkTextEN,
		locale.Russian: LinkTextRU,
	}
	ShareThisTimerText = map[locale.Locale]string{
		locale.English: ShareThisTimerTextEN,
		locale.Russian: ShareThisTimerTextRU,
	}
	TimersLimitFieldText = map[locale.Locale]string{
		locale.English: TimersLimitFieldTextEN,
		locale.Russian: TimersLimitFieldTextRU,
	}
	WalletFieldText = map[locale.Locale]string{
		locale.English: WalletFieldTextEN,
		locale.Russian: WalletFieldTextRU,
	}
	TimerRingingText = map[locale.Locale]string{
		locale.English: TimerRingingTextEN,
		locale.Russian: TimerRingingTextRU,
		locale.Uzbek:   TimerRingingTextUZ,
	}
)

const (
	PrestartTextEN = `BOO!

HI! I'm Boo - your friendly ghost

I will help you not to forget about any degen stuff you have to do. Just set a timer and I will remind you about it.`
	PrestartTextRU = `BOO!
Привет! Я Boo - дружелюбное привидение

Я помогу вам не забыть о любых ваших крипто фармилках. Просто установите таймер, а я напомню`

	StartTextEN = `BOO!

HI! I'm Boo - your friendly ghost

I will help you not to forget about any degen stuff you have to do. Just set a timer and I will remind you about it.

📢 BOO Channel @BooTimer
💬 BOO Chat @BooTimerChat`
	StartTextRU = `BOO!
Привет! Я Boo - дружелюбное привидение

Я помогу вам не забыть о любых ваших крипто фармилках. Просто установите таймер, а я напомню

📢 BOO Channel @BooTimer
💬 BOO Chat @BooTimerChat`

	ProfileTextEN = `Your profile, @%s

%s

%s

%s`

	ProfileTextRU = `Ваш профиль, @%s

%s

%s

%s`

	ReferralAcceptedTextEN = `Someone has chosen to join the afterlife using your referral link. Thanks from Charon and me! 🚣‍♂️

%s`

	ReferralAcceptedTextRU = `Кто-то решил присоединиться к загробному миру, используя вашу реферальную ссылку. Спасибо от Харона и меня! 🚣‍♂️

%s`

	InviteLinkTextEN = `Your referral link: %s`

	InviteLinkTextRU = `Ваша реферальная ссылка: %s`

	InviteLinkTextUZ = `Sizning taklifnomangiz: %s`

	NewTimerTextEN = `BOO! 👻 Let's create a new timer. Please provide the following details:`

	NewTimerTextRU = `BOO! 👻 Давайте создадим новый таймер. Пожалуйста, укажите следующие детали:`

	TimerNameTextEN = `What would you like to name your timer?`

	TimerNameTextRU = `Как вы хотите назвать свой таймер?`

	InvalidTimerNameTextEN        = `Timer Name's length should contain 3-100 symbols. Please provide a name for your timer`
	InvalidTimerNameTextRU        = `Длина имени таймера должна быть от 3 до 100 символов. Пожалуйста, укажите имя для вашего таймера`
	InvalidTimerDescriptionTextEN = `Timer Description's length should contain 3-1000 letters. Please provide a description for your timer`
	InvalidTimerDescriptionTextRU = `Длина описания таймера должна быть от 3 до 1000 символов. Пожалуйста, укажите описание для вашего таймера`

	TimerSetupSuccessTextEN = `Timer successfully set up! To make it running, open timer settings and set up Type, and also Period or Trigger Time. And don't forget to press Start!'`
	TimerSetupSuccessTextRU = `Таймер успешно создан! Чтобы запустить его, откройте настройки таймера и установите Тип, а также Период или Время срабатывания. И не забудьте нажать Старт!`

	TimerUpdateSuccessTextEN = `Timer successfully updated!`
	TimerUpdateSuccessTextRU = `Таймер успешно обновлен!`

	UserUpdateSuccessTextEN = `User settings are saved!`
	UserUpdateSuccessTextRU = `Настройки пользователя сохранены!`

	NoTimersTextEN          = `You have no timers yet`
	NoTimersTextRU          = `У вас пока нет таймеров`
	YourTimersTextEN        = `Choose a timer from the list below`
	YourTimersTextRU        = `Выберите таймер из списка ниже`
	TimersLimitTextEN       = `You have reached the limit of timers you can create. To create a new timer, you need to delete one of the existing ones.`
	TimersLimitTextRU       = `Вы достигли лимита таймеров, которые вы можете создать. Чтобы создать новый таймер, вам нужно удалить один из существующих.`
	TimerAlreadyAddedTextEN = `❌ You have this timer already. Delete or rename old timer to add a new one.`
	TimerAlreadyAddedTextRU = `❌ У вас уже есть этот таймер. Удалите или переименуйте старый таймер, чтобы добавить новый.`

	PresetTimersTextEN = `Choose a preset timer from the list below:`
	PresetTimersTextRU = `Выберите готовый таймер из списка ниже:`

	EditNameTextEN        = `Enter a new name for the timer:`
	EditNameTextRU        = `Введите новое имя для таймера:`
	EditDescriptionTextEN = `Enter a new description for the timer:`
	EditDescriptionTextRU = `Введите новое описание для таймера:`
	EditTypeTextEN        = `Choose a new type for the timer.
*Daily* will notify you at a specific time every day.
*Periodical* will notify you every time the specified period passes.`
	EditTypeTextRU = `Выберите новый тип для таймера.
*Daily* будет уведомлять вас в определенное время каждый день.
*Periodical* будет уведомлять вас каждый раз, когда проходит указанный период.`
	EditTriggerTimeTextEN = `Enter a new trigger time for the timer. Boo 👻 will notify you at this time daily.
Format: 'HH:MM'. 24-hour format.
Hours - \[00:23], Minutes - \[00:59]
Example: 17:45, which is equal to 5:45 PM
All times are in UTC.`
	EditTriggerTimeTextRU = `Введите новое время срабатывания для таймера. Boo 👻 будет уведомлять вас в это время каждый день.
Формат: 'HH:MM'. 24-часовой формат.
Часы - \[00:23], Минуты - \[00:59]
Пример: 17:45, что равно 5:45 PM
Все времена в UTC.`
	EditPeriodTextEN = `Enter a new period for the timer. Boo 👻 will notify you every time this period passes.
Format: HH:MM(:SS)
Examples: 
1) 01:30:15, which is 1 hour, 30 minutes and 15 seconds
2) 00:05, which is 5 minutes`
	EditPeriodTextRU = `Введите новый период для таймера. Boo 👻 будет уведомлять вас каждый раз, когда проходит этот период.
Формат: HH:MM(:SS)
Примеры:
1) 01:30:15, что равно 1 часу, 30 минутам и 15 секундам
2) 00:05, что равно 5 минутам`
	EditRepeatTypeTextEN = `Choose a new repeat type for the timer.
*Once missed*: Boo 👻 will not send you a notification if you didn't acknowledged the previous one.
*Always*: You'll receive a notification every time the timer triggers, no matter if last one was acknowledged or not.`
	EditRepeatTypeTextRU = `Выберите новый тип повторения для таймера.
*Once missed*: Boo 👻 не будет отправлять вам уведомление, если вы не подтвердили предыдущее.
*Always*: Вы будете получать уведомление каждый раз, когда таймер срабатывает, независимо от того, было ли подтверждено последнее уведомление или нет.`
	EditLinkTextEN = `Enter a new link for the timer. Boo 👻 will send you this link with the notification.
Given value must be a valid URL. Examples:
1) https://google.com
2) t.me/BlumCryptoBot/app?startapp=ref_BHtV2K2haY`
	EditLinkTextRU = `Введите новую ссылку для таймера. Boo 👻 отправит вам эту ссылку с уведомлением.
Указанное значение должно быть действительным URL-адресом. Примеры:
1) https://google.com
2) t.me/BlumCryptoBot/app?startapp=ref_BHtV2K2haY`

	EditUserWalletTextRU = `Введите адрес вашего TON кошелька.

**Важно**: используйте только Tonkeeper, MyTonWallet или Tonhub. 

❌ @wallet и биржевые адреса не подходят`
	EditUserWalletTextEN = `Enter your TON wallet address.

**Notice**: use only Tonkeeper, MyTonWallet or Tonhub. 

❌ @wallet and exchange addresses are not suitable`

	SelectLanguageTextEN = `Select your language:`
	SelectLanguageTextRU = `Выберите ваш язык:`
	SelectLanguageTextUZ = `Tilni tanlang:`

	AllTasksFinishedTextEN = "All tasks finished. Hell yeah! 😈"
	AllTasksFinishedTextRU = "Все задания выполнены. Чертовски круто! 😈"

	TaskCenterTextEN = `Choose a task from the list below:`
	TaskCenterTextRU = `Выберите задание из списка ниже:`

	SubscribeTaskMessageEN = `📢 Complete task to get %d %s!
	👉 [press here](%s)`
	SubscribeTaskMessageRU = `📢 Выполните задание, чтобы получить %d %s!
	👉 [жмите сюда](%s)`

	AddToNicknameTaskMessageEN = `Add 👻 - our mascot - to your nickname and get %d %s!`
	AddToNicknameTaskMessageRU = `Добавьте 👻 - нашего маскота - в свой ник и получите %d %s!`

	PeriodTooSmallTextEN = `Period is too small. Please provide a period of at least 5 minute`
	PeriodTooSmallTextRU = `Период слишком маленький. Пожалуйста, укажите период не менее 5 минут`

	BonusIsNoLongerActiveTextEN = `Bonus is no longer active. Please claim the bonus within 2 hours. The next notification will be sent at: %s (UTC)`
	BonusIsNoLongerActiveTextRU = `Бонус больше не активен. Пожалуйста, заберите бонус в течение 2 часов. Следующее уведомление будет отправлено в: %s (UTC)`

	AlertBonusClaimedTextEN = `Bonus %d $eBOO received 🎁 

Next alert (UTC): %s

%s`
	AlertBonusClaimedTextRU = `Бонус %d $eBOO получен 🎁

Следующее уведомление (UTC): %s

%s`
	MintReceivedTextEN = `@%s received %d $eBOO from @%s with coment: %s
new balance: %d $eBOO`
	MintReceivedTextRU = `@%s получил %d $eBOO от @%s с комментарием: %s
новый баланс: %d $eBOO`
	YourPositionTextEN      = `Your position: %d. Score: %.1f`
	YourPositionTextRU      = `Ваша позиция: %d. Очки: %.1f`
	LeaderboardGoogleTextEN = `The contest is over\. Find results [here](https://docs.google.com/spreadsheets/d/1qdQdX5FlJew4yOFNZcu1qB8ETJbCfFVoZhAnUdxpKSs)`
	LeaderboardGoogleTextRU = `Конкурс закончен\. Результаты доступны [тут](https://docs.google.com/spreadsheets/d/1qdQdX5FlJew4yOFNZcu1qB8ETJbCfFVoZhAnUdxpKSs)`

	BoobleJumpTextEN = `SOON!`
	BoobleJumpTextRU = `СКОРО!`

	TapToSeeReadyTimersTextEN = "Tap button to see ready-to-use timers 👻"
	TapToSeeReadyTimersTextRU = "Нажмите на кнопку, чтобы увидеть готовые таймеры 👻"

	YouCannotPerformThisTextEN = "You cannot perform this action"
	YouCannotPerformThisTextRU = "Вы не можете совершать это действие"
	TimerTypeNotSetTextEN      = "Timer type is not set"
	TimerTypeNotSetTextRU      = "Тип таймера не установлен"
	TimerAlertTimeNotSetTextEN = "Timer alert time is not set"
	TimerAlertTimeNotSetTextRU = "Время оповещения таймера не установлено"
	TimerPeriodIsNotSetTextEN  = "Timer period is not set"
	TimerPeriodIsNotSetTextRU  = "Период таймера не установлен"
	ConfirmDeleteTextEN        = "Are you sure you want to delete this timer?"
	ConfirmDeleteTextRU        = "Вы уверены, что хотите удалить этот таймер?"

	TimerDeletedTextEN                 = "Timer '%s' deleted"
	TimerDeletedTextRU                 = "Таймер '%s' удален"
	TimerLimitReachedTextEN            = "You have reached the limit of timers. Delete some timers to add new ones"
	TimerLimitReachedTextRU            = "Вы достигли лимита таймеров. Удалите старые таймеры, чтобы создать новый"
	TimerAddedTextEN                   = "Timer added"
	TimerAddedTextRU                   = "Таймер добавлен"
	GiveawayTextEN                     = "6.66 TON giveaway - [in our channel](https://t.me/boodrops/30) 🎁"
	GiveawayTextRU                     = "Розыгрыш 6.66 TON - [в нашем канале](https://t.me/boodrops/30) 🎁"
	GiveawayTextUZ                     = "6.66 TON olish uchun - [kanalimizda](https://t.me/boodrops/30) 🎁"
	StopButtonTextEN                   = "🔴 STOP"
	StopButtonTextRU                   = "🔴 СТОП"
	StartButtonTextEN                  = "🟢 START"
	StartButtonTextRU                  = "🟢 СТАРТ"
	EditTimerButtonTextEN              = "⏰ Edit Timer"
	EditTimerButtonTextRU              = "⏰ Изменить таймер"
	DeleteTimerButtonTextEN            = "⚠️ Delete Timer"
	DeleteTimerButtonTextRU            = "⚠️ Удалить таймер"
	BackToTimersListButtonTextEN       = "📝 Back to Timers List"
	BackToTimersListButtonTextRU       = "📝 К списку таймеров"
	EditNameButtonTextEN               = "📝 Edit Name"
	EditNameButtonTextRU               = "📝 Изменить имя"
	EditDescriptionButtonTextEN        = "📝 Edit Description"
	EditDescriptionButtonTextRU        = "📝 Изменить описание"
	EditTypeButtonTextEN               = "📝 Edit Type"
	EditTypeButtonTextRU               = "📝 Изменить тип"
	EditTriggerTimeButtonTextEN        = "📝 Edit Trigger Time"
	EditTriggerTimeButtonTextRU        = "📝 Изменить время срабатывания"
	EditPeriodButtonTextEN             = "📝 Edit Period"
	EditPeriodButtonTextRU             = "📝 Изменить период"
	EditLinkButtonTextEN               = "📝 Edit Link"
	EditLinkButtonTextRU               = "📝 Изменить ссылку"
	BackToTimerButtonTextEN            = "◀️ Back to Timer"
	BackToTimerButtonTextRU            = "◀️ Назад к таймеру"
	BackToProfileButtonTextEN          = "📝 Back to Profile"
	BackToProfileButtonTextRU          = "📝 Назад к профилю"
	ConfirmDeleteButtonTextEN          = "Confirm delete"
	ConfirmDeleteButtonTextRU          = "Подтвердить удаление"
	CancelButtonTextEN                 = "❌ CANCEL"
	CancelButtonTextRU                 = "❌ ОТМЕНА"
	BackToTimerEditButtonTextEN        = "◀️ Back to Timer Edit"
	BackToTimerEditButtonTextRU        = "◀️ Назад к редактированию таймера"
	GoToButtonTextEN                   = "🔗 Go to %s"
	GoToButtonTextRU                   = "🔗 Перейти в %s"
	ManageTimerButtonTextEN            = "🕰 Manage Timer"
	ManageTimerButtonTextRU            = "🕰 Управление таймером"
	RebootTimerButtonTextEN            = "🔄 ReBOOt Timer"
	RebootTimerButtonTextRU            = "🔄 Сбросить таймер"
	ClaimBonusButtonTextEN             = "🎁 Claim Bonus"
	ClaimBonusButtonTextRU             = "🎁 Получить бонус"
	AddToMyTimersButtonTextEN          = "📝 Add to My Timers"
	AddToMyTimersButtonTextRU          = "📝 Добавить в Мои таймеры"
	BackToPresetTImersListButtonTextEN = "📝 Back to Preset Timers List"
	BackToPresetTImersListButtonTextRU = "📝 К списку готовых таймеров"
	EnglishButtonText                  = "🇬🇧 English"
	RussianButtonText                  = "🇷🇺 Русский"
	UzbekButtonText                    = "🇺🇿 O'zbekcha"
	EditLanguageButtonTextEN           = "🌐 Edit Language"
	EditLanguageButtonTextRU           = "🌐 Изменить язык"
	EditWalletButtonTextEN             = "💼 Manage Wallet"
	EditWalletButtonTextRU             = "💼 Управление кошельком"

	ProfileSlug10TextEN       = "*Granny Nancy's urn* ⚱️\nYou invited %d people to join the afterlife."
	ProfileSlug10TextRU       = "*Урна бабушки Нэнси* ⚱️\nВы пригласили %d человек в загробный мир."
	ProfileSlug50TextEN       = "*Not bad, not dead* 🪦\nYou invited %d people to join the afterlife."
	ProfileSlug50TextRU       = "*Умереть не встать* 🪦\nВы пригласили %d человек в загробный мир."
	ProfileSlug100TextEN      = "You're a *Casual 5/2 Grim Reaper* 🧟\nYou invited %d people to join the afterlife."
	ProfileSlug100TextRU      = "Вы *5/2 Смерть с косой* 🧟\nВы пригласили %d человек в загробный мир."
	ProfileSlug500TextEN      = "You're a *DraCOOLa* 🧛\nYou invited %d people to join the afterlife."
	ProfileSlug500TextRU      = "Вы *Дра-КУЛ-а* 🧛\nВы пригласили %d человек в загробный мир."
	ProfileSlug1000TextEN     = "*Satan's right hand* 😈\nYou invited %d people to join the afterlife."
	ProfileSlug1000TextRU     = "*Правая рука Сатаны* 😈\nВы пригласили %d человек в загробный мир."
	ProfileSlug1000PlusTextEN = "You're the *Satan himself* 👹\nYou invited %d people to join the afterlife."
	ProfileSlug1000PlusTextRU = "Вы *Сам Сатана* 👹\nВы пригласили %d человек в загробный мир."
	YourStatusTextEN          = "Your status: "
	YourStatusTextRU          = "Ваш статус: "
	BalanceTextEN             = "💰 *Balance*: %d%s $eBOO"
	BalanceTextRU             = "💰 *Баланс*: %d%s $eBOO"

	TimerTextEN            = "Timer"
	TimerTextRU            = "Таймер"
	ReadyTimerTextEN       = "Ready Timer"
	ReadyTimerTextRU       = "Готовый таймер"
	TypeTextEN             = "Type"
	TypeTextRU             = "Тип"
	StatusTextEN           = "Status"
	StatusTextRU           = "Статус"
	DescriptionTextEN      = "Description"
	DescriptionTextRU      = "Описание"
	TriggerTimeTextEN      = "Trigger Time (UTC)"
	TriggerTimeTextRU      = "Время срабатывания (UTC)"
	PeriodTextEN           = "Period"
	PeriodTextRU           = "Период"
	LastAlertTextEN        = "Last alert (UTC)"
	LastAlertTextRU        = "Последнее уведомление (UTC)"
	NextAlertTextEN        = "Next alert (UTC)"
	NextAlertTextRU        = "Следующее уведомление (UTC)"
	LinkTextEN             = "Link"
	LinkTextRU             = "Ссылка"
	ShareThisTimerTextEN   = "Share this timer"
	ShareThisTimerTextRU   = "Поделиться этим таймером"
	TimersLimitFieldTextEN = "Timers limit"
	TimersLimitFieldTextRU = "Лимит таймеров"
	WalletFieldTextEN      = "Wallet"
	WalletFieldTextRU      = "Кошелек"

	TimerRingingTextEN = "*%s* is ringing!"
	TimerRingingTextRU = "*%s* звенит!"
	TimerRingingTextUZ = "*%s* qo'ng'iroqlar!"
)
