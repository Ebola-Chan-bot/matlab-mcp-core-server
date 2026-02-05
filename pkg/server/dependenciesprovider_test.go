// Copyright 2026 The MathWorks, Inc.

package server_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDependenciesProvider_toInternal_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	expectedMessage := "test message"

	mockLogger.EXPECT().
		Info(expectedMessage).
		Once()

	type TestDependencies struct{}
	expectedDependencies := &TestDependencies{}

	provider := server.DependenciesProvider[*TestDependencies](func(resources server.DependenciesProviderResources) (*TestDependencies, error) {
		resources.Logger().Info(expectedMessage)
		return expectedDependencies, nil
	})

	// Act
	internalProvider := provider.ToInternal()
	dependencies, err := internalProvider(definition.DependenciesProviderResources{
		Logger: mockLogger,
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedDependencies, dependencies)
}

func TestDependenciesProvider_toInternal_NilProvider(t *testing.T) {
	// Arrange
	var provider server.DependenciesProvider[struct{}]

	// Act
	internalProvider := provider.ToInternal()
	result, err := internalProvider(definition.DependenciesProviderResources{})

	// Assert
	require.NoError(t, err)
	require.Nil(t, result)
}

func TestDependenciesProvider_toInternal_Error(t *testing.T) {
	// Arrange
	type TestDependencies struct{}
	expectedError := assert.AnError

	provider := server.DependenciesProvider[*TestDependencies](func(resources server.DependenciesProviderResources) (*TestDependencies, error) {
		return &TestDependencies{}, expectedError // Returning an actual pointer, to make sure we mask it, and return nil
	})

	// Act
	internalProvider := provider.ToInternal()
	dependencies, err := internalProvider(definition.DependenciesProviderResources{})

	// Assert
	require.ErrorIs(t, err, expectedError)
	require.Nil(t, dependencies)
}
