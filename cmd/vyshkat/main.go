package main

import (
	"flag"

	"vyshkat/di"
)

func main() {
	demoInspection := flag.Bool(
		"demo-rejection",
		false,
		"включить демонстрационные отказы техосмотра",
	)

	flag.Parse()

	di.NewContainer(*demoInspection).Run()
}
