package main

import (
	"log"

	"github.com/t-kikuc/pipecd-plugin-prototypes/cdk/deployment"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
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
