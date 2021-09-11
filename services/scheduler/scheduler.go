package scheduler

import (
	"context"
	"log"
	"time"
	"xrate/services"
)

//Task is to hold task function
type TaskFn func(ctx context.Context) error

//ISchedulerService base methods scheduler service interface
type ISchedulerService interface {
	services.Service

	//Add adding task to scheduler by given task and delay time
	AddTask(t TaskFn, delay time.Duration)
}

type task struct {
	fn    TaskFn
	delay time.Duration
}

type schedulerService struct {
	//buffer of task channel
	ch chan *task
}

//NewService returns new scheduler service
func NewService() ISchedulerService {
	return &schedulerService{
		//create task buffer max 10 channel
		ch: make(chan *task, 10),
	}
}

//AddTask implementing ISchedulerService.AddTask
func (s *schedulerService) AddTask(fn TaskFn, delay time.Duration) {
	t := &task{
		fn:    fn,
		delay: delay,
	}

	s.ch <- t
}

//Run implementing services.Service
func (s *schedulerService) Run(ctx context.Context) error {
	log.Println("Scheduler service started...")
	chErr := make(chan error, 1)
	go func() {
		//wait until done
		<-ctx.Done()
		log.Println("scheduler:gracefully shutdown...")
		//close all task channel
		close(s.ch)
		close(chErr)
	}()

	for t := range s.ch {
		task := t
		go s.createWorker(ctx, task)
	}

	//block until get close signal from chErr
	return <-chErr
}

//createWorker create new worker by given task and delay
func (s *schedulerService) createWorker(ctx context.Context, t *task) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := t.fn(ctx); err != nil {
				log.Println("failed when run task with error", err)
			}
			time.Sleep(t.delay)
		}
	}
}
