package deployment

import (
	"context"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"

	"github.com/t-kikuc/pipecd-plugin-prototypes/cdk/cli"
	config "github.com/t-kikuc/pipecd-plugin-prototypes/cdk/config"
)

func executeDeploy(
	ctx context.Context,
	dt *sdk.DeployTarget[config.DeployTargetConfig],
	input *sdk.ExecuteStageInput[config.ApplicationSpec],
) sdk.StageStatus {
	lp := input.Client.LogPersister()
	specInput := input.Request.TargetDeploymentSource.ApplicationConfig.Spec.Input

	cdkCmd, err := cli.NewCDK(ctx, input.Client.ToolRegistry(), input.Request.TargetDeploymentSource.ApplicationDirectory, dt.Config)
	if err != nil {
		lp.Errorf("failed to create cdk command: %v", err)
		return sdk.StageStatusFailure
	}

	// Get application config
	// TODO: uncomment after defining fields in CDKDeployStageOptions.
	// var stageCfg config.CDKDeployStageOptions
	// if err := json.Unmarshal(req.StageConfig, &stageCfg); err != nil {
	// 	lp.Errorf("failed to decode stage config: %v", err)
	// 	return sdk.StageStatusFailure
	// }

	if err := cdkCmd.Deploy(ctx, lp, specInput); err != nil {
		lp.Errorf("failed to execute 'cdk deploy': %v", err)
		return sdk.StageStatusFailure
	}

	lp.Successf("Successfully executed 'cdk deploy'")
	return sdk.StageStatusSuccess
}
