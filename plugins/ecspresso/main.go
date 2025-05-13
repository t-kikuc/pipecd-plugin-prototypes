package main

import (
	"log"

	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func main() {
	plugin, err := sdk.NewPlugin("ecspresso", "0.0.1", sdk.WithStagePlugin(&plugin{}))
	if err != nil {
		log.Fatalln(err)
	}

	if err := plugin.Run(); err != nil {
		log.Fatalln(err)
	}
}
