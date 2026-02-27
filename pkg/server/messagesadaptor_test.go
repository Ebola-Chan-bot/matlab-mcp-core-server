// Copyright 2026 The MathWorks, Inc.

package server_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/stretchr/testify/require"
)

func TestNewI18nErrorFactory_HappyPath(t *testing.T) {
	// Arrange
	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	// Act
	factory := server.NewI18nErrorFactory(mockMessageCatalog)

	// Assert
	require.NotNil(t, factory)
}

func TestI18nErrorFactory_FromInternalError_HappyPath(t *testing.T) {
	// Arrange
	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	expectedError := messages.AnError
	expectedMessage := "translated error message"

	mockMessageCatalog.EXPECT().
		GetFromError(expectedError).
		Return(expectedMessage).
		Once()

	factory := server.NewI18nErrorFactory(mockMessageCatalog)

	// Act
	result := factory.FromInternalError(expectedError)

	// Assert
	require.NotNil(t, result)
	require.Equal(t, expectedMessage, result.Error())
}
