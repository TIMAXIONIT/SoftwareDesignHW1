package di

import (
	"vyshkat/app"
	"vyshkat/domain"
	"vyshkat/service"

	"go.uber.org/fx"
)

func NewContainer() *fx.App {
	return fx.New(
		fx.Provide(
			fx.Annotate(
				service.NewFleet,
				fx.As(new(domain.Inventory)),
			),
			fx.Annotate(
				service.NewServiceCenter,
				fx.As(new(domain.InspectionService)),
			),
			service.NewFleetService,
			app.NewConsole,
		),
		fx.Invoke(func(console *app.Console) {
			console.Run()
		}),
	)
}
