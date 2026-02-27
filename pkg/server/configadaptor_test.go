// Copyright 2026 The MathWorks, Inc.

package server_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/stretchr/testify/require"
)

func TestNewConfigAdaptor_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	// Act
	adaptor := server.NewConfigAdaptor(mockConfig, mockMessageCatalog)

	// Assert
	require.NotNil(t, adaptor)
}

func TestConfigAdaptor_Get_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	expectedKey := "test-key"
	expectedValue := "test-value"

	mockConfig.EXPECT().
		Get(expectedKey).
		Return(expectedValue, nil).
		Once()

	adaptor := server.NewConfigAdaptor(mockConfig, mockMessageCatalog)

	// Act
	result, err := adaptor.Get(expectedKey, "")

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedValue, result)
}

func TestConfigAdaptor_Get_ConfigError(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	expectedKey := "missing-key"
	expectedErrorMessage := "translated error message"

	mockConfig.EXPECT().
		Get(expectedKey).
		Return(nil, messages.AnError).
		Once()

	mockMessageCatalog.EXPECT().
		GetFromError(messages.AnError).
		Return(expectedErrorMessage).
		Once()

	adaptor := server.NewConfigAdaptor(mockConfig, mockMessageCatalog)

	// Act
	result, err := adaptor.Get(expectedKey, "")

	// Assert
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, expectedErrorMessage, err.Error())
}

func TestConfigAdaptor_Get_TypeMismatch(t *testing.T) {
	// Arrange
	mockConfig := &configmocks.MockGenericConfig{}
	defer mockConfig.AssertExpectations(t)

	mockMessageCatalog := &definitionmocks.MockMessageCatalog{}
	defer mockMessageCatalog.AssertExpectations(t)

	expectedKey := "test-key"
	returnedValue := 123
	expectedInternalError := messages.New_StartupErrors_InvalidParameterType_Error(expectedKey, "string")
	expectedErrorMessage := "type mismatch error"

	mockConfig.EXPECT().
		Get(expectedKey).
		Return(returnedValue, nil).
		Once()

	mockMessageCatalog.EXPECT().
		GetFromError(expectedInternalError).
		Return(expectedErrorMessage).
		Once()

	adaptor := server.NewConfigAdaptor(mockConfig, mockMessageCatalog)

	// Act
	result, err := adaptor.Get(expectedKey, "")

	// Assert
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, expectedErrorMessage, err.Error())
}
