package validate

import "testing"

func TestValidator_ValidateLink(t *testing.T) {
	type args struct {
		object string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "plain text",
			args: args{
				object: "lol",
			},
			wantErr: true,
		},
		{
			name: "valid https link",
			args: args{
				object: "https://google.com",
			},
			wantErr: false,
		},
		{
			name: "telegram link",
			args: args{
				object: "t.me/BlumCryptoBot/app?startapp=ref_BHtV2K2haY",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			if err := v.ValidateLink(tt.args.object); (err != nil) != tt.wantErr {
				t.Errorf("ValidateLink() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
