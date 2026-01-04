package healthfx

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"healthfx",
		logger.WithNamedLogger("healthfx"),
		fx.Provide(
			fx.Annotate(NewService, fx.ParamTags(`group:"health-providers"`)),
		),
		fx.Invoke(func(svc *Service) {
			svc.Register(newHealth())
		}),
	)
}

func AsProvider(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Provider)),
		fx.ResultTags(`group:"health-providers"`),
	)
}
