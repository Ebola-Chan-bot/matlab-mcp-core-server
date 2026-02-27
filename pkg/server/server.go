// Copyright 2026 The MathWorks, Inc.

package server

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/wire/adaptor"
)

type Parameter interface {
	GetID() string
	GetFlagName() string
	GetHiddenFlag() bool
	GetEnvVarName() string
	GetDescription() string
	GetDefaultValue() any

	GetActive() bool
	GetRecordToLog() bool
}

type Definition[Dependencies any] struct {
	Name         string
	Title        string
	Instructions string

	Features Features

	Parameters Parameters

	DependenciesProvider DependenciesProvider[Dependencies]

	ToolsProvider ToolsProvider[Dependencies]
}

type Server[Dependencies any] struct {
	serverDefinition Definition[Dependencies]

	applicationFactory adaptor.ApplicationFactory
	errorWriter        entities.Writer
}

func New[Dependencies any](thisDefinition Definition[Dependencies]) *Server[Dependencies] {
	// Cloning parameters to avoid unexpected mutations
	thisDefinition.Parameters = slices.Clone(thisDefinition.Parameters)

	return &Server[Dependencies]{
		applicationFactory: adaptor.NewFactory(),

		serverDefinition: thisDefinition,
		errorWriter:      os.Stderr,
	}
}

func (s *Server[Dependencies]) StartAndWaitForCompletion(ctx context.Context) int {
	serverDefinition := definition.New(
		s.serverDefinition.Name,
		s.serverDefinition.Title,
		s.serverDefinition.Instructions,
		s.serverDefinition.Features.toInternal(),
		s.serverDefinition.Parameters.ToInternal(),
		s.serverDefinition.DependenciesProvider.toInternal(),
		s.serverDefinition.ToolsProvider.toInternal(),
	)
	application := s.applicationFactory.New(serverDefinition)

	if err := application.ModeSelector().StartAndWaitForCompletion(ctx); err != nil {
		errorMessage, ok := application.MessageCatalog().GetFromGeneralError(err)
		if ok {
			fmt.Fprintf(s.errorWriter, "%s\n", errorMessage) //nolint:errcheck // Nothing we can do then
			return 1
		}

		fallbackMessage := application.MessageCatalog().Get(messages.StartupErrors_GenericInitializeFailure)
		fmt.Fprintf(s.errorWriter, "%s\n", fallbackMessage) //nolint:errcheck // Nothing we can do then
		return 1
	}

	return 0
}

type Parameters []Parameter

func (p Parameters) ToInternal() []entities.Parameter {
	if len(p) == 0 {
		return nil
	}

	internalParameters := make([]entities.Parameter, len(p))

	for i, parameter := range p {
		internalParameters[i] = parameter
	}

	return internalParameters
}
