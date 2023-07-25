package jsonparser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadJSONConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *ServerJSONConfig
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Positive case",
			args: args{
				path: "./server_config.json",
			},
			want: &ServerJSONConfig{
				Address:        "localhost:8080",
				Restore:        true,
				StoreInterval:  "1s",
				StoreFile:      "/path/to/file.db",
				Key:            "",
				Database:       "",
				PrivateKeyPath: "/path/to/key.pem",
			},
			wantErr: false,
		},
		{
			name: "Negative case",
			args: args{
				path: "../agent_config.json",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Negative case (unmarshalling)",
			args: args{
				path: "./invalid_server_config.json",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadJSONConfig[ServerJSONConfig](tt.args.path)
			if err != nil && tt.wantErr {
				require.Nil(t, got)
				return
			}
			require.EqualValues(t, tt.want, got)
		})
	}
}
