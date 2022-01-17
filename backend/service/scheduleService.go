package service

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

var quit = make(chan bool)

type ScheduleService struct {
	task     func()
	duration time.Duration
}

type ScheduleServices []ScheduleService

func NewScheduleService(task func(), duration time.Duration) *ScheduleService {
	return &ScheduleService{task, duration}
}

func NewScheduleServices(scheduleServices ...ScheduleService) *ScheduleServices {
	res := ScheduleServices{}

	for _, s := range scheduleServices {
		res = append(res, s)
	}
	return &res
}

func (s *ScheduleService) Start() error {
	if s.task == nil {
		return errors.New("nil task")
	}

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				methodName := runtime.FuncForPC(reflect.ValueOf(s.task).Pointer()).Name()
				log.Infof("Task %s start.", methodName)
				s.task()
				log.Infof("Task %s executed.", methodName)
				time.Sleep(s.duration)
			}
		}
	}()

	return nil
}

func (s *ScheduleServices) Start() error {
	for _, scheduleService := range *s {
		if err := scheduleService.Start(); err != nil {
			stopAll()
			return fmt.Errorf("%v, stop all tasks", err)
		}
	}
	return nil
}

func stopAll() {
	quit <- true
}
