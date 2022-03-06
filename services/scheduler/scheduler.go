package scheduler

import (
	"context"
	"log"
	"time"

	"xrate/services"
)

// Task is to hold task function
type TaskFn func(ctx context.Context) error

// ISchedulerService base methods scheduler service interface
type ISchedulerService interface {
	services.Service

	// Add adding task to scheduler by given task and delay time
	AddTask(t TaskFn, delay time.Duration, maxRetryFromPanic int)
}

type task struct {
	fn                TaskFn
	delay             time.Duration
	maxRetryFromPanic int
	countPanic        int
}

type schedulerService struct {
	// buffer of task channel
	ch chan *task
}

// NewService returns new scheduler service
func NewService() ISchedulerService {
	return &schedulerService{
		// create task buffer max 10 channel
		ch: make(chan *task, 10),
	}
}

// AddTask implementing ISchedulerService.AddTask
func (s *schedulerService) AddTask(fn TaskFn, delay time.Duration, maxRetryFromPanic int) {
	t := &task{
		fn:                fn,
		delay:             delay,
		maxRetryFromPanic: maxRetryFromPanic,
	}

	s.ch <- t
}

// Run implementing services.Service
func (s *schedulerService) Run(ctx context.Context) error {
	log.Println("Scheduler service started...")
	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-s.ch:
			go s.createWorker(ctx, t)
		default:
			continue
		}
	}
}

// createWorker create new worker by given task and delay
func (s *schedulerService) createWorker(ctx context.Context, t *task) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(t.delay)
			s.runFn(ctx, t)
			if t.countPanic > t.maxRetryFromPanic {
				return
			}
		}
	}
}

func (s *schedulerService) runFn(ctx context.Context, t *task) {
	defer func(t *task) {
		p := recover()
		if p == nil {
			return
		}
		log.Println("job raised panic")
		t.countPanic = t.countPanic + 1
	}(t)

	if err := t.fn(ctx); err != nil {
		log.Println("failed when run task with error", err)
	}
}
