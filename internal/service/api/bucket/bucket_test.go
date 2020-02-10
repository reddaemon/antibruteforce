package bucket

import (
	"golang.org/x/net/context"
	"testing"
	"time"

	"go.uber.org/zap"
)

var memRepo = NewMemRepo(zap.NewNop())
var rate = 300 * time.Millisecond

func TestMemRepo(t *testing.T) {
	type args struct {
		ctx      context.Context
		key      string
		capacity uint
		rate     time.Duration
	}

	tests := []struct {
		sleep   time.Duration
		name    string
		args    args
		wantErr bool
	}{
		{
			sleep: 0,
			name:  "testing leaky bucket, first iteration",
			args: args{
				ctx:      context.Background(),
				key:      "test",
				capacity: 2,
				rate:     rate,
			},
			wantErr: false,
		},
		{
			sleep: 0,
			name:  "testing leaky bucket, second iteration",
			args: args{
				key:      "test",
				capacity: 2,
				rate:     0,
			},
			wantErr: false,
		},
		{
			sleep: 0,
			name:  "testing leaky bucket, another key",
			args: args{
				ctx:      context.Background(),
				key:      "test-another-key",
				capacity: 2,
				rate:     0,
			},
			wantErr: false,
		},
		{
			sleep: 0,
			name:  "testing leaky bucket, third iteration",
			args: args{
				ctx:      context.Background(),
				key:      "test",
				capacity: 2,
				rate:     0,
			},
			wantErr: true,
		},
		{
			sleep: rate,
			name:  "testing leaky bucket, after time drop",
			args: args{
				ctx:      context.Background(),
				key:      "test",
				capacity: 2,
				rate:     0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		time.Sleep(tt.sleep)
		t.Run(tt.name, func(t *testing.T) {
			if err := memRepo.Add(tt.args.ctx, tt.args.key,
				tt.args.capacity, tt.args.rate); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemRepo_Drop(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "drop",
			args: args{
				ctx:  context.Background(),
				keys: []string{"test", "undefined"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := memRepo.Drop(tt.args.ctx, tt.args.keys); (err != nil) != tt.wantErr {
				t.Errorf("Drop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
