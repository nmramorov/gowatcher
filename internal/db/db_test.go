package db

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestNewCursor(t *testing.T) {
	parent := context.Background()
	type args struct {
		parent  context.Context
		link    string
		adaptor string
	}
	type test struct {
		name    string
		args    args
		wantErr bool
	}
	tests := []test{
		{
			name: "Positive Cursor creation",
			args: args{
				parent:  parent,
				link:    "",
				adaptor: "pgx",
			},
			wantErr: false,
		},
		{
			name: "Negative Cursor creation",
			args: args{
				parent:  parent,
				link:    "",
				adaptor: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCursor(tt.args.parent, tt.args.link, tt.args.adaptor)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCursor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCursor_Close(t *testing.T) {
	parent := context.Background()

	type args struct {
		parent context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Close Positive case",
			args: args{
				parent: parent,
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCursor(tt.args.parent, "", "pgx")
			if err != nil {
				t.Errorf("Error creating Cursror %v", err)
			}
			if err := c.Close(tt.args.parent); (err != nil) != tt.wantErr {
				t.Errorf("Cursor.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCursor_Ping(t *testing.T) {
	parent := context.Background()

	cursor, _ := NewCursor(parent, "", "pgx")

	type args struct {
		parent context.Context
	}
	tests := []struct {
		name    string
		c       *Cursor
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Ping Negative",
			c: cursor,
			args: args{
				parent: parent,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Ping(tt.args.parent); (err != nil) != tt.wantErr {
				t.Errorf("Cursor.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
