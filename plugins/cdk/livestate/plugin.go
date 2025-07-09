package livestate

import (
	"bytes"
	"context"
	"fmt"
	"time"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/t-kikuc/pipecd-plugin-prototypes/cdk/cli"
	"github.com/t-kikuc/pipecd-plugin-prototypes/cdk/config"
	"go.uber.org/zap"
)

type Plugin struct{}

// GetLivestate implements sdk.LivestatePlugin.
func (p Plugin) GetLivestate(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[config.DeployTargetConfig], input *sdk.GetLivestateInput[config.ApplicationSpec]) (*sdk.GetLivestateResponse, error) {
	if len(dts) != 1 {
		return nil, fmt.Errorf("only 1 deploy target is allowed but got %d", len(dts))
	}

	dt := dts[0]
	cfg, err := input.Request.DeploymentSource.AppConfig()
	if err != nil {
		input.Logger.Error("Failed while loading application config", zap.Error(err))
		return nil, err
	}

	// Create CDK client
	cdk, err := cli.NewCDK(
		ctx,
		input.Client.ToolRegistry(),
		input.Request.DeploymentSource.ApplicationDirectory,
		dt.Config,
	)
	if err != nil {
		input.Logger.Error("Failed to create CDK client", zap.Error(err))
		return nil, err
	}

	// Run cdk diff to detect changes
	var buf bytes.Buffer
	if err := cdk.Diff(ctx, &buf, cfg.Spec.Input, "--quiet"); err != nil {
		input.Logger.Error("Failed to run cdk diff", zap.Error(err))
		return nil, err
	}
	// TODO: Check the drift between CFn stack and the actual resource??

	return toResult(buf.Bytes(), dt.Name), nil
}

func toResult(diffOutput []byte, dtName string) *sdk.GetLivestateResponse {
	hasDiff := hasDiff(string(diffOutput))

	// Consider the Stack as one resource since it's troublesome to fetch all resources for now.
	return &sdk.GetLivestateResponse{
		LiveState: sdk.ApplicationLiveState{
			// TODO: Create Resources for each stack.
			Resources: []sdk.ResourceState{
				{
					ID:           "cdk-stack",
					Name:         "CDK Stack",
					DeployTarget: dtName,
					CreatedAt:    time.Now(),
					HealthStatus: sdk.ResourceHealthStateHealthy, // TODO: Improve?
				},
			},
		},
		SyncState: sdk.ApplicationSyncState{
			Status:      getSyncStatus(hasDiff),
			ShortReason: getSyncShortReason(hasDiff),
			Reason:      getDriftDetail(hasDiff, diffOutput),
		},
	}
}

func hasDiff(diffOutput string) bool {
	return !bytes.Contains([]byte(diffOutput), []byte("Number of stacks with differences: 0"))
}

func getSyncStatus(hasDiff bool) sdk.ApplicationSyncStatus {
	if hasDiff {
		return sdk.ApplicationSyncStateOutOfSync
	}
	return sdk.ApplicationSyncStateSynced
}

func getSyncShortReason(hasDiff bool) string {
	if hasDiff {
		return "Differences detected between the live state and the desired state"
	}
	return "Number of stacks with differences: 0"
}

func getDriftDetail(hasDiff bool, diffOutput []byte) string {
	if !hasDiff {
		return "No drift detected"
	}
	return string(diffOutput)
}
