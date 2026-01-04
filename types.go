package healthfx

import (
	"context"
)

type Status string
type statusLevel int

const (
	StatusPass Status = "pass"
	StatusWarn Status = "warn"
	StatusFail Status = "fail"

	levelPass statusLevel = 0
	levelWarn statusLevel = 1
	levelFail statusLevel = 2
)

// CheckResult represents the result of a set of health checks.
type CheckResult struct {
	// A map of check names to their respective details.
	Checks Checks
}

// Status returns the overall status of the application.
// It can be one of the following values: "pass", "warn", or "fail".
func (c CheckResult) Status() Status {
	// Determine overall status
	level := levelPass
	for _, detail := range c.Checks {
		switch detail.Status {
		case StatusPass:
		case StatusFail:
			level = max(level, levelFail)
		case StatusWarn:
			level = max(level, levelWarn)
		}
	}

	switch level {
	case levelPass:
		return StatusPass
	case levelWarn:
		return StatusWarn
	case levelFail:
		return StatusFail
	}

	return StatusFail
}

// CheckDetail of a health check.
type CheckDetail struct {
	// A human-readable description of the check.
	Description string
	// Unit of measurement for the observed value.
	ObservedUnit string
	// Observed value of the check.
	ObservedValue any
	// Status of the check.
	// It can be one of the following values: "pass", "warn", or "fail".
	Status Status
}

// Checks is a map of check names to their respective details.
type Checks map[string]CheckDetail

// Provider defines the interface for health check providers.
// Each provider can implement startup, readiness, and liveness probes.
// Probes should return an empty Checks map (not nil) when there are no checks to report.
// Returning an error indicates the probe failed to execute, not that checks failed.
type Provider interface {
	// Name returns the unique identifier for this provider.
	Name() string

	// StartedProbe checks if the provider has successfully started.
	StartedProbe(ctx context.Context) (Checks, error)
	// ReadyProbe checks if the provider is ready to handle requests.
	ReadyProbe(ctx context.Context) (Checks, error)
	// LiveProbe checks if the provider is still alive and functioning.
	LiveProbe(ctx context.Context) (Checks, error)
}
