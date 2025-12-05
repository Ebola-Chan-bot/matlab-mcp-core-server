// Copyright 2025 The MathWorks, Inc.

package codingguidelines

import (
	"context"
	_ "embed"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/baseresource"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

//go:embed codingguidelines.md
var codingGuidelines string

type Resource struct {
	*baseresource.Resource
}

func New(loggerFactory baseresource.LoggerFactory) (*Resource, error) {
	baseRes, err := baseresource.New(
		name,
		title,
		description,
		mimeType,
		estimatedSize,
		uri,
		loggerFactory,
		Handler(),
	)
	if err != nil {
		return nil, err
	}

	return &Resource{
		Resource: baseRes,
	}, nil
}

func Handler() baseresource.ResourceHandler {
	return func(_ context.Context, logger entities.Logger) (*baseresource.ReadResourceResult, error) {
		logger.Info("Returning MATLAB coding guidelines resource")

		return &baseresource.ReadResourceResult{
			Contents: []baseresource.ResourceContents{
				{
					MIMEType: mimeType,
					Text:     codingGuidelines,
				},
			},
		}, nil
	}
}
