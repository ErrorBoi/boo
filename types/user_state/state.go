package user_state

import "errors"

var (
	ErrStepNotFound = errors.New("user_state not found")
)

type Step string

const (
	Init        Step = "init"
	Name        Step = "name"
	WaitForName Step = "wait_for_name"

	NameEdit        Step = "editTimer_name"
	DescriptionEdit Step = "editTimer_description"
	TypeEdit        Step = "editTimer_type"
	TriggerTimeEdit Step = "editTimer_triggerTime"
	PeriodEdit      Step = "editTimer_period"
	RepeatTypeEdit  Step = "editTimer_repeatType"
	LinkEdit        Step = "editTimer_link"

	UserWalletEdit Step = "editUser_wallet"

	None Step = "none"
)

type State struct {
	Step    Step
	TimerID *int64
	UserID  *int64
}

func NewState(step Step) *State {
	return &State{Step: step}
}
