package main

import (
	"log"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/t-kikuc/pipecd-plugin-prototypes/ecschedule/deployment"
)

func main() {
	plugin, err := sdk.NewPlugin(
		"0.0.1",
		sdk.WithDeploymentPlugin(&deployment.Plugin{}),
	)
	if err != nil {
		log.Fatalln(err)
	}
	if err := plugin.Run(); err != nil {
		log.Fatalln(err)
	}
}
