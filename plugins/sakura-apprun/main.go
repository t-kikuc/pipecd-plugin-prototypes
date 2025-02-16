package main

import (
	"log"

	"github.com/pipe-cd/pipecd/pkg/cli"
)

func main() {
	app := cli.NewApp(
		"pipecd-plugin-sakura-apprun-prototype",
		"Plugin component to deploy AppRun appliactions.",
	)
	app.AddCommands(
		newPluginCommand(),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
