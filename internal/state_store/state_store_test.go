package state_store

import (
	"testing"

	state_types "github.com/errorboi/boo/types/user_state"
)

func TestInmemStateStore_Set(t *testing.T) {
	type args struct {
		userID int64
		state  *state_types.State
	}

	stepStore := NewInmemStateStore("test")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "SetState user_state",
			args: args{
				userID: 1,
				state:  state_types.NewState(state_types.WaitForName),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := stepStore.SetState(tt.args.userID, tt.args.state); (err != nil) != tt.wantErr {
				t.Errorf("SetState() error = %v, wantErr %v", err, tt.wantErr)
			}

			res, _ := stepStore.Get(tt.args.userID)
			if res.Step != state_types.WaitForName {
				t.Errorf("SetState() error = %v, wantErr %v", res.Step, state_types.WaitForName)
			}
		})
	}
}
