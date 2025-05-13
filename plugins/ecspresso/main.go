package main

import (
	"log"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	"github.com/t-kikuc/pipecd-plugin-prototypes/ecspresso/deployment"
)

func main() {
	plugin, err := sdk.NewPlugin(
		"ecspresso", "0.0.1",
		sdk.WithDeploymentPlugin(&deployment.Plugin{}),
		// TODO: Add livestate plugin
	)
	if err != nil {
		log.Fatalln(err)
	}
	if err := plugin.Run(); err != nil {
		log.Fatalln(err)
	}
}
