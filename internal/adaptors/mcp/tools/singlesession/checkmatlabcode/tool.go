// Copyright 2025 The MathWorks, Inc.

package checkmatlabcode

import (
	"context"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/annotations"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/checkmatlabcode"
)

type Usecase interface {
	Execute(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, request checkmatlabcode.Args) (checkmatlabcode.ReturnArgs, error)
}

type Tool struct {
	basetool.ToolWithStructuredContentOutput[Args, ReturnArgs]
}

func New(
	loggerFactory basetool.LoggerFactory,
	usecase Usecase,
	globalMATLAB entities.GlobalMATLAB,
) *Tool {
	return &Tool{
		ToolWithStructuredContentOutput: basetool.NewToolWithStructuredContent(name, title, description, annotations.NewReadOnlyAnnotations(), loggerFactory, Handler(usecase, globalMATLAB)),
	}
}

func Handler(usecase Usecase, globalMATLAB entities.GlobalMATLAB) basetool.HandlerWithStructuredContentOutput[Args, ReturnArgs] {
	return func(ctx context.Context, sessionLogger entities.Logger, inputs Args) (ReturnArgs, error) {
		sessionLogger.Info("Executing Check MATLAB code tool")
		defer sessionLogger.Info("Done - Executing Check MATLAB code tool")

		// Not returning nil for empty slices, to comply with MCP spec.
		mcpCompliantZeroValue := ReturnArgs{
			CheckCodeOutput: []string{},
		}

		client, err := globalMATLAB.Client(ctx, sessionLogger)
		if err != nil {
			return mcpCompliantZeroValue, err
		}

		checkcodeResponse, err := usecase.Execute(ctx, sessionLogger, client, checkmatlabcode.Args{
			ScriptPath: inputs.ScriptPath,
		})
		if err != nil {
			return mcpCompliantZeroValue, err
		}

		if checkcodeResponse.CheckCodeOutput == nil {
			return mcpCompliantZeroValue, nil
		}

		return ReturnArgs{
			CheckCodeOutput: checkcodeResponse.CheckCodeOutput,
		}, nil
	}
}
