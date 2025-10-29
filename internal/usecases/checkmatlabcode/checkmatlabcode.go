// Copyright 2025 The MathWorks, Inc.

package checkmatlabcode

import (
	"context"
	"fmt"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

type Args struct {
	ScriptPath string
}

type ReturnArgs struct {
	CheckCodeOutput []string
}

type PathValidator interface {
	ValidateMATLABScript(filePath string) (string, error)
}

type Usecase struct {
	pathValidator PathValidator
}

func New(
	pathValidator PathValidator,
) *Usecase {
	return &Usecase{
		pathValidator: pathValidator,
	}
}

func (u *Usecase) Execute(ctx context.Context, sessionLogger entities.Logger, client entities.MATLABSessionClient, checkcodeRequest Args) (ReturnArgs, error) {
	sessionLogger.Debug("Entering CheckMATLABCode Usecase")
	defer sessionLogger.Debug("Exiting CheckMATLABCode Usecase")

	validatedPath, err := u.pathValidator.ValidateMATLABScript(checkcodeRequest.ScriptPath)
	if err != nil {
		return ReturnArgs{}, fmt.Errorf("path validation failed: %w", err)
	}

	request := entities.EvalRequest{
		Code: fmt.Sprintf("checkcode('%s')", strings.ReplaceAll(validatedPath, "'", "''")), // Escape single quotes
	}
	response, err := client.EvalWithCapture(ctx, sessionLogger, request)
	if err != nil {
		return ReturnArgs{}, err
	}

	output := response.ConsoleOutput
	if output == "" {
		output = "No issues found by checkcode"
	}

	return ReturnArgs{
		CheckCodeOutput: splitAndCleanLines(output),
	}, nil
}

func splitAndCleanLines(text string) []string {
	lines := strings.Split(text, "\n")
	result := []string{}

	for _, line := range lines {
		// Trim whitespace and skip empty lines
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
