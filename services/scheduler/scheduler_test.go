package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_schedulerService_AddTask(t *testing.T) {
	var runCount int
	type args struct {
		fn                TaskFn
		delay             time.Duration
		maxRetryFromPanic int
	}
	tests := []struct {
		name             string
		args             args
		expectedRunCount int
	}{
		{
			name: "no error, panic, exptected 2",
			args: args{
				fn: func(ctx context.Context) error {
					if runCount < 2 {
						runCount++
					}

					return nil
				},
				delay:             time.Microsecond * 10,
				maxRetryFromPanic: 0,
			},
			expectedRunCount: 2,
		},
		{
			name: "error, no panic, exptected 3",
			args: args{
				fn: func(ctx context.Context) error {
					if runCount < 3 {
						runCount++
					}

					return fmt.Errorf("error")
				},
				delay:             time.Microsecond * 100,
				maxRetryFromPanic: 0,
			},
			expectedRunCount: 3,
		},
		{
			name: "panic exptected 1",
			args: args{
				fn: func(ctx context.Context) error {
					runCount = 1
					panic("only run once")
				},
				delay:             time.Microsecond * 10,
				maxRetryFromPanic: 0,
			},
			expectedRunCount: 1,
		},
		{
			name: "panic exptected 2",
			args: args{
				fn: func(ctx context.Context) error {
					runCount = runCount + 1
					panic("panic 2 twice")
				},
				delay:             time.Microsecond * 10,
				maxRetryFromPanic: 2,
			},
			expectedRunCount: 2,
		},
		{
			name: "infinity panic",
			args: args{
				fn: func(ctx context.Context) error {
					if runCount < 10 {
						runCount = runCount + 1
					}
					panic("infinity panic")
				},
				delay:             time.Microsecond * 10,
				maxRetryFromPanic: -1,
			},
			expectedRunCount: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runCount = 0
			s := NewService()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			go func() {
				assert.NoError(t, s.Run(ctx))
			}()
			s.AddTask(tt.args.fn, tt.args.delay, tt.args.maxRetryFromPanic)
			<-ctx.Done()
			assert.Equal(t, tt.expectedRunCount, runCount)
		})
	}
}
