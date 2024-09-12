package validate

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type ObjectType string

const (
	TriggerTime ObjectType = "trigger_time"
)

var (
	ErrUnknownObjectType = errors.New("unknown object type")

	ErrInvalidHours   = errors.New("invalid Hours value. Value should be between 0 and 23")
	ErrInvalidMinutes = errors.New("invalid Minutes value. Value should be between 0 and 59")
	ErrInvalidSeconds = errors.New("invalid Seconds value. Value should be between 0 and 59")
)

type Validator interface {
	Validate(object string, objectType ObjectType) error
	ValidateTriggerTime(object string) error
	ValidatePeriod(object string) error
	ValidateLink(object string) error
	ValidateTonWallet(object string) error
}

type validator struct {
}

func New() Validator {
	return &validator{}
}

func (v *validator) Validate(object string, objectType ObjectType) error {
	switch objectType {
	case TriggerTime:
		return v.ValidateTriggerTime(object)
	default:
		return ErrUnknownObjectType
	}
}

func (v *validator) ValidateTriggerTime(object string) error {
	arr := strings.Split(object, ":")

	if len(arr) != 2 {
		return errors.New("please provide 2 values separated by colon")
	}

	hoursStr := arr[0]
	minutesStr := arr[1]

	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		return errors.New("invalid hours value")
	}

	if hours < 0 || hours > 23 {
		return ErrInvalidHours
	}

	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		return errors.New("invalid minutes value")
	}

	if minutes < 0 || minutes > 59 {
		return ErrInvalidMinutes
	}

	return nil
}

func (v *validator) ValidatePeriod(object string) error {
	arr := strings.Split(object, ":")

	if len(arr) < 2 {
		return errors.New("please provide 2 or 3 values separated by colon")
	}

	hoursStr := arr[0]
	minutesStr := arr[1]

	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		return errors.New("invalid hours value")
	}

	if hours < 0 || hours > 23 {
		return ErrInvalidHours
	}

	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		return errors.New("invalid minutes value")
	}

	if minutes < 0 || minutes > 59 {
		return ErrInvalidMinutes
	}

	var seconds int
	if len(arr) == 3 {
		secondsStr := arr[2]

		seconds, err = strconv.Atoi(secondsStr)
		if err != nil {
			return errors.New("invalid seconds value")
		}

		if seconds < 0 || seconds > 59 {
			return ErrInvalidSeconds
		}
	}

	return nil
}

func (v *validator) ValidateLink(object string) error {
	switch {
	case strings.HasPrefix(object, "@"), strings.HasPrefix(object, "t.me"):
		return nil
	default:
		if isUrl(object) {
			return nil
		}

		return errors.New("invalid link")
	}
}

func (v *validator) ValidateTonWallet(object string) error {
	if len(object) != 48 {
		return errors.New("invalid wallet length")
	}

	if !strings.HasPrefix(object, "UQ") {
		return errors.New("invalid wallet prefix")
	}

	return nil
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
