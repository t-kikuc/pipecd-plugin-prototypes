package main

import (
	"log"

	"github.com/pipe-cd/pipecd/pkg/cli"
)

func main() {
	app := cli.NewApp(
		"pipecd-plugin-cdk-prototype",
		"Plugin component to deploy Lambda functions by cdk.",
	)
	app.AddCommands(
		newPluginCommand(),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
