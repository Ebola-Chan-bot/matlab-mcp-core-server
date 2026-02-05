// Copyright 2026 The MathWorks, Inc.

package definition_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	toolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools"
	basetoolmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	"github.com/stretchr/testify/require"
)

func TestDefinition_Name_HappyPath(t *testing.T) {
	// Arrange
	expectedName := "my-definition"
	def := definition.New(expectedName, "", "", nil, nil)

	// Act
	result := def.Name()

	// Assert
	require.Equal(t, expectedName, result)
}

func TestDefinition_Title_HappyPath(t *testing.T) {
	// Arrange
	expectedTitle := "My Definition Title"
	def := definition.New("", expectedTitle, "", nil, nil)

	// Act
	result := def.Title()

	// Assert
	require.Equal(t, expectedTitle, result)
}

func TestDefinition_Instructions_HappyPath(t *testing.T) {
	// Arrange
	expectedInstructions := "These are the instructions"
	def := definition.New("", "", expectedInstructions, nil, nil)

	// Act
	result := def.Instructions()

	// Assert
	require.Equal(t, expectedInstructions, result)
}

func TestDefinition_Dependencies_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	expectedDependencies := &struct{}{}
	expectedResources := definition.DependenciesProviderResources{
		Logger: mockLogger,
	}

	dependenciesProvider := func(resources definition.DependenciesProviderResources) (any, error) {
		require.Equal(t, expectedResources, resources)
		return expectedDependencies, nil
	}

	def := definition.New("", "", "", dependenciesProvider, nil)

	// Act
	result, err := def.Dependencies(expectedResources)

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedDependencies, result)
}

func TestDefinition_Dependencies_NilProvider(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	expectedResources := definition.DependenciesProviderResources{
		Logger: mockLogger,
	}
	def := definition.New("", "", "", nil, nil)

	// Act
	result, err := def.Dependencies(expectedResources)

	// Assert
	require.NoError(t, err)
	require.Nil(t, result)
}

func TestDefinition_Tools_HappyPath(t *testing.T) {
	// Arrange
	mockTool := &toolsmocks.MockTool{}
	defer mockTool.AssertExpectations(t)

	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	expectedResources := definition.ToolsProviderResources{
		LoggerFactory: mockLoggerFactory,
	}
	expectedTools := []tools.Tool{mockTool}

	toolsProvider := func(resources definition.ToolsProviderResources) []tools.Tool {
		require.Equal(t, expectedResources, resources)
		return expectedTools
	}

	def := definition.New("", "", "", nil, toolsProvider)

	// Act
	result := def.Tools(expectedResources)

	// Assert
	require.Equal(t, expectedTools, result)
}

func TestDefinition_Tools_NilProvider(t *testing.T) {
	// Arrange
	expectedResources := definition.ToolsProviderResources{}
	def := definition.New("", "", "", nil, nil)

	// Act
	result := def.Tools(expectedResources)

	// Assert
	require.Nil(t, result)
}
