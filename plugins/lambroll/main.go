package main

import (
	"log"

	"github.com/pipe-cd/pipecd/pkg/cli"
)

func main() {
	app := cli.NewApp(
		"pipecd-plugin-lambroll-prototype",
		"Plugin component to deploy Lambda functions by lambroll.",
	)
	app.AddCommands(
		newPluginCommand(),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
