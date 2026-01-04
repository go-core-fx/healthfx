package healthfx

import (
	"context"
	"runtime"
)

type health struct {
}

func newHealth() *health {
	return &health{}
}

// Name implements HealthProvider.
func (h *health) Name() string {
	return "system"
}

// LiveProbe implements HealthProvider.
func (h *health) LiveProbe(_ context.Context) (Checks, error) {
	const oneMiB uint64 = 1 << 20
	const memoryThreshold uint64 = 128 * oneMiB
	const goroutineThreshold = 100

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Basic runtime health checks
	goroutineCount := runtime.NumGoroutine()
	memoryUsage := m.Alloc / oneMiB

	goroutineCheck := CheckDetail{
		Description:   "Number of goroutines",
		ObservedValue: goroutineCount,
		ObservedUnit:  "goroutines",
		Status:        StatusPass,
	}

	memoryCheck := CheckDetail{
		Description:   "Memory usage",
		ObservedValue: memoryUsage,
		ObservedUnit:  "MiB",
		Status:        StatusPass,
	}

	// Check for potential memory issues
	if m.Alloc > memoryThreshold {
		memoryCheck.Status = StatusWarn
	}

	// Check for excessive goroutines
	if goroutineCount > goroutineThreshold {
		goroutineCheck.Status = StatusWarn
	}

	return Checks{"goroutines": goroutineCheck, "memory": memoryCheck}, nil
}

// ReadyProbe implements HealthProvider.
func (h *health) ReadyProbe(_ context.Context) (Checks, error) {
	return Checks{}, nil
}

// StartedProbe implements HealthProvider.
func (h *health) StartedProbe(_ context.Context) (Checks, error) {
	return Checks{}, nil
}

var _ Provider = (*health)(nil)
