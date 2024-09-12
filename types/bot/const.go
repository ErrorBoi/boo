package bot

type TaskType string

const (
	SubscribeTaskType     TaskType = "subscribe"
	AddToNicknameTaskType TaskType = "add_to_nickname"
)

type RewardType string

const (
	EbooRewardType  RewardType = "eboo"
	TimerRewardType RewardType = "timer"
)

func (rt RewardType) ToText() string {
	switch rt {
	case EbooRewardType:
		return "$eBOO"
	case TimerRewardType:
		return "Timer"
	default:
		return ""
	}
}
