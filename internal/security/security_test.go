package security

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCryptoKey(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Positive case",
			args: args{
				path: "./cert2.pem",
			},
			wantErr: false,
			want:    "[]byte",
		},
		{
			name: "Negative case",
			args: args{
				path: "../cert2.pem",
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCryptoKey(tt.args.path)
			if (err == nil) == tt.wantErr {
				t.Errorf("GetCryptoKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got).String() != "[]uint8" {
				t.Errorf("GetCryptoKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPrivateKey(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Positive case",
			args: args{
				path: "./key.pem",
			},
			wantErr: false,
		},
		{
			name: "Negative case (wrong path)",
			args: args{
				path: "../key.pem",
			},
			wantErr: true,
		},
		{
			name: "Negative case (wrong file)",
			args: args{
				path: "../cert2.pem",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPrivateKey(tt.args.path)
			fmt.Println(err != nil, tt.wantErr)
			if tt.wantErr {
				require.Nil(t, got)
				return
			} else {
				require.NotNil(t, got)
			}
		})
	}
}

func TestGetCertificate(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *x509.Certificate
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Positive case",
			args: args{
				path: "./cert2.pem",
			},
			wantErr: false,
		},
		{
			name: "Negative case (wrong path)",
			args: args{
				path: "../cert2.pem",
			},
			wantErr: true,
		},
		{
			name: "Negative case (wrong file)",
			args: args{
				path: "../key.pem",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCertificate(tt.args.path)
			fmt.Println(err != nil, tt.wantErr)
			if tt.wantErr {
				require.Nil(t, got)
				return
			} else {
				require.NotNil(t, got)
			}
		})
	}
}

func TestEncodeMsg(t *testing.T) {
	cert, _ := GetCertificate("./cert2.pem")
	type args struct {
		payload     []byte
		certificate *x509.Certificate
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Positive case",
			args: args{
				payload:     []byte(""),
				certificate: cert,
			},
			want:    []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeMsg(tt.args.payload, tt.args.certificate)
			if tt.wantErr {
				require.NotNil(t, err)
				return
			}
			require.NotNil(t, got)
		})
	}
}

func TestDecodeMsg(t *testing.T) {
	key, _ := GetPrivateKey("./key.pem")
	cert, _ := GetCertificate("./cert2.pem")
	encodedMsg, _ := EncodeMsg([]byte("sdf"), cert)
	type args struct {
		msg        []byte
		privateKey *rsa.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Positive case",
			args: args{
				msg:        encodedMsg,
				privateKey: key,
			},
			want:    []byte("sdf"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeMsg(tt.args.msg, tt.args.privateKey)
			if tt.wantErr {
				require.NotNil(t, err)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
