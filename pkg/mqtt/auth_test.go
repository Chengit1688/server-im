package mqtt

import "testing"

func TestDeleteAuthUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				username: "im_test_3061040876",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteAuthUsername(tt.args.username); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAuthUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
