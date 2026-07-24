package health

import (
	"time"
)

type Service struct {
	startTime time.Time
	checkers  []HealthChecker
}

func NewService() *Service {
	return &Service{
		startTime: time.Now(),
		checkers:  make([]HealthChecker, 0),
	}
}

func (s *Service) Register(name string, timeout time.Duration, check Checker) {
	s.checkers = append(s.checkers, HealthChecker{
		Name:    name,
		Timeout: timeout,
		Check:   check,
	})
}
