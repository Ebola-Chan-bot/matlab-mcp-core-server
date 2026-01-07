// Copyright 2025-2026 The MathWorks, Inc.

package sdk_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/sdk"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/server/sdk"
	"github.com/stretchr/testify/assert"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	// Act
	factory := sdk.NewFactory(mockConfig)

	// Assert
	assert.NotNil(t, factory, "Factory should not be nil")
}

func TestFactory_New_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	expectedVersion := "1.0.0"
	expectedName := "test-server"
	expectedInstructions := "test instructions"

	mockConfig.EXPECT().
		Version().
		Return(expectedVersion).
		Once()

	factory := sdk.NewFactory(mockConfig)

	// Act
	server := factory.NewServer(expectedName, expectedInstructions)

	// Assert
	assert.NotNil(t, server, "Server should not be nil") // We can't easily qualify that the correct version is used. Further qualifications take place in system tests.
}
