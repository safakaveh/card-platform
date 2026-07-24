package health

import (
	"context"
	"time"
)

type Checker func(ctx context.Context) error

// internal checker definition
type HealthChecker struct {
	Name    string
	Check   Checker
	Timeout time.Duration
}

// response item
type CheckResult struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// API response
type Response struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Uptime    int64                  `json:"uptime"`
	Checks    map[string]CheckResult `json:"checks,omitempty"`
}
