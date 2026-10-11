package di

import (
	"go.uber.org/fx"

	"vyshkat/app"
	"vyshkat/domain"
	"vyshkat/service"
)

func NewContainer(demoInspection bool) *fx.App {
	return fx.New(
		fx.Provide(
			fx.Annotate(
				service.NewFleet,
				fx.As(new(domain.Inventory)),
			),
			func() domain.InspectionService {
				if demoInspection {
					return service.NewDemoServiceCenter()
				}

				return service.NewServiceCenter()
			},
			fx.Annotate(
				service.NewFleetService,
				fx.As(new(app.FleetOperations)),
			),
			app.NewConsole,
		),
		fx.Invoke(func(console *app.Console) {
			console.Run()
		}),
	)
}
