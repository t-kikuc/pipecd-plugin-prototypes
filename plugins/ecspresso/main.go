package main

import (
	"log"

	"github.com/pipe-cd/pipecd/pkg/cli"
)

func main() {
	app := cli.NewApp(
		"pipecd-plugin-ecspresso-prototype",
		"Plugin component to deploy ECS services by ecspresso.",
	)
	app.AddCommands(
		newPluginCommand(),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
