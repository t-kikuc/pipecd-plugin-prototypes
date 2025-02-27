package deployment

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/pipe-cd/pipecd/pkg/plugin/signalhandler"
	cdkconfig "github.com/t-kikuc/pipecd-plugin-prototypes/cdk/config"
	"github.com/t-kikuc/pipecd-plugin-prototypes/cdk/toolregistry"
)

type toolClient interface {
	InstallTool(ctx context.Context, name, version, script string) (path string, err error)
}

type toolRegistry interface {
	CDK(ctx context.Context, nodeVersion, version string) (path string, err error)
}

type logPersister interface {
	StageLogPersister(deploymentID, stageID string) logpersister.StageLogPersister
}

type DeploymentServiceServer struct {
	deployment.UnimplementedDeploymentServiceServer

	// this field is set with the plugin configuration
	// the plugin configuration is sent from piped while initializing the plugin
	pluginConfig *config.PipedPlugin
	// deployTargetConfig might be empty. e.g. when it's not specified in the piped config.
	// For now, this plugin supports up to one deploy target.
	deployTargetConfig cdkconfig.CDKDeployTargetConfig

	logger       *zap.Logger
	toolRegistry toolRegistry
	logPersister logPersister
}

// NewDeploymentServiceServer creates a new DeploymentServiceServer of cdk plugin.
func NewDeploymentServiceServer(
	config *config.PipedPlugin,
	logger *zap.Logger,
	toolClient toolClient,
	logPersister logPersister,
) (*DeploymentServiceServer, error) {
	toolRegistry := toolregistry.NewRegistry(toolClient)

	deployTargetConfig := cdkconfig.CDKDeployTargetConfig{}
	if len(config.DeployTargets) > 0 {
		var err error
		if deployTargetConfig, err = cdkconfig.ParseDeployTargetConfig(config.DeployTargets[0]); err != nil {
			return nil, err
		}
	}

	return &DeploymentServiceServer{
		pluginConfig:       config,
		deployTargetConfig: deployTargetConfig,
		logger:             logger.Named("cdk-plugin"),
		toolRegistry:       toolRegistry,
		logPersister:       logPersister,
	}, nil
}

// Register registers all handling of this service into the specified gRPC server.
func (s *DeploymentServiceServer) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, s)
}

// DetermineStrategy implements deployment.DeploymentServiceServer.
func (s *DeploymentServiceServer) DetermineStrategy(ctx context.Context, request *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
	cfg, err := config.DecodeYAML[*cdkconfig.CDKApplicationSpec](request.GetInput().GetTargetDeploymentSource().GetApplicationConfig())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	strategy, summary, err := determineStrategy(*cfg.Spec)
	if err != nil {
		return nil, err
	}
	return &deployment.DetermineStrategyResponse{
		SyncStrategy: strategy,
		Summary:      summary,
	}, nil
}

// DetermineVersions implements deployment.DeploymentServiceServer.
func (s *DeploymentServiceServer) DetermineVersions(ctx context.Context, request *deployment.DetermineVersionsRequest) (*deployment.DetermineVersionsResponse, error) {
	return &deployment.DetermineVersionsResponse{
		Versions: nil,
	}, nil
}

// BuildPipelineSyncStages implements deployment.DeploymentServiceServer.
func (s *DeploymentServiceServer) BuildPipelineSyncStages(ctx context.Context, request *deployment.BuildPipelineSyncStagesRequest) (*deployment.BuildPipelineSyncStagesResponse, error) {
	now := time.Now()
	stages := buildPipelineStages(request.GetStages(), request.GetRollback(), now)
	return &deployment.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}

// BuildQuickSyncStages implements deployment.DeploymentServiceServer.
func (s *DeploymentServiceServer) BuildQuickSyncStages(ctx context.Context, request *deployment.BuildQuickSyncStagesRequest) (*deployment.BuildQuickSyncStagesResponse, error) {
	now := time.Now()
	stages := buildQuickSyncStages(request.GetRollback(), now)
	return &deployment.BuildQuickSyncStagesResponse{
		Stages: stages,
	}, nil
}

// FetchDefinedStages implements deployment.DeploymentServiceServer.
func (s *DeploymentServiceServer) FetchDefinedStages(context.Context, *deployment.FetchDefinedStagesRequest) (*deployment.FetchDefinedStagesResponse, error) {
	return &deployment.FetchDefinedStagesResponse{
		Stages: allStages,
	}, nil
}

// ExecuteStage performs stage-defined tasks.
// It returns stage status after execution without error.
// An error will be returned only if the given stage is not supported.
func (s *DeploymentServiceServer) ExecuteStage(ctx context.Context, request *deployment.ExecuteStageRequest) (response *deployment.ExecuteStageResponse, _ error) {
	lp := s.logPersister.StageLogPersister(request.GetInput().GetDeployment().GetId(), request.GetInput().GetStage().GetId())
	defer func() {
		// When termination signal received and the stage is not completed yet, we should not mark the log persister as completed.
		// This can occur when the piped is shutting down while the stage is still running.
		if !response.GetStatus().IsCompleted() && signalhandler.Terminated() {
			return
		}
		lp.Complete(time.Minute)
	}()

	status, err := s.executeStage(ctx, lp, request.GetInput())
	if err != nil {
		return nil, err
	}
	return &deployment.ExecuteStageResponse{
		Status: status,
	}, nil
}
