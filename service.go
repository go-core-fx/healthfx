package healthfx

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

const providerTimeout = 5 * time.Second

type Service struct {
	providers []Provider
	mu        sync.RWMutex

	logger *zap.Logger
}

func NewService(providers []Provider, logger *zap.Logger) *Service {
	return &Service{
		providers: providers,
		mu:        sync.RWMutex{},

		logger: logger,
	}
}

func (s *Service) Register(p Provider) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.providers = append(s.providers, p)
}

func (s *Service) CheckReadiness(ctx context.Context) CheckResult {
	return s.checkProvider(ctx, func(ctx context.Context, p Provider) (Checks, error) {
		return p.ReadyProbe(ctx)
	})
}

func (s *Service) CheckLiveness(ctx context.Context) CheckResult {
	return s.checkProvider(ctx, func(ctx context.Context, p Provider) (Checks, error) {
		return p.LiveProbe(ctx)
	})
}

func (s *Service) CheckStartup(ctx context.Context) CheckResult {
	return s.checkProvider(ctx, func(ctx context.Context, p Provider) (Checks, error) {
		return p.StartedProbe(ctx)
	})
}

func (s *Service) checkProvider(
	ctx context.Context,
	probe func(context.Context, Provider) (Checks, error),
) CheckResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	check := CheckResult{
		Checks: map[string]CheckDetail{},
	}

	for _, p := range s.providers {
		select {
		case <-ctx.Done():
			return check
		default:
		}

		probeCtx, cancel := context.WithTimeout(ctx, providerTimeout)
		healthChecks, err := probe(probeCtx, p)
		cancel()

		if err != nil {
			s.logger.Error("failed check", zap.String("provider", p.Name()), zap.Error(err))
			check.Checks[p.Name()] = CheckDetail{
				Description:   "Failed check",
				ObservedUnit:  "",
				ObservedValue: 0,
				Status:        StatusFail,
			}
			continue
		}

		if len(healthChecks) == 0 {
			continue
		}

		for name, detail := range healthChecks {
			check.Checks[p.Name()+":"+name] = detail
		}
	}

	return check
}
