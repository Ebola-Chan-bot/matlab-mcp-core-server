// Copyright 2025-2026 The MathWorks, Inc.

package parser_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter/parser"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	parsermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/parameter/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse_BoolEnvVar(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "bool-param"
	paramEnvVar := "BOOL_ENV_VAR"

	mockParam := newMockParam(
		t,
		paramID,
		"bool-flag",
		paramEnvVar,
		false,
		"Test bool description",
		false,
		true,
	)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam}).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv(paramEnvVar).
		Return("true", true).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, true, result[paramID])
	assert.Equal(t, []entities.Parameter{mockParam}, parameters)
}

func TestParser_Parse_InactiveParameterEnvVarSkipped(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramID := "inactive-param"
	paramEnvVar := "INACTIVE_ENV_VAR"
	paramDefaultValue := "default-value"

	mockParam := newMockParam(
		t,
		paramID,
		"",
		paramEnvVar,
		paramDefaultValue,
		"",
		false,
		false,
	)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{mockParam}).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, err := p.Parse(args)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, paramDefaultValue, result[paramID])
	assert.Equal(t, []entities.Parameter{mockParam}, parameters)
}

func TestParser_Parse_BadEnvVarBoolValue(t *testing.T) {
	// Arrange
	mockOSLayer := &parsermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockDefaultParamFactory := &parsermocks.MockDefaultParameterFactory{}
	defer mockDefaultParamFactory.AssertExpectations(t)

	mockParamFactory := &parsermocks.MockParameterFactory{}
	defer mockParamFactory.AssertExpectations(t)

	paramEnvVar := "BOOL_ENV_VAR"
	badEnvValue := "notabool"

	mockParam := newMockParam(
		t,
		"bool-param",
		"bool-flag",
		paramEnvVar,
		false,
		"Test bool description",
		false,
		true,
	)

	mockDefaultParamFactory.EXPECT().
		DefaultParameters().
		Return([]entities.Parameter{mockParam}).
		Once()

	mockParamFactory.EXPECT().
		Parameters().
		Return([]entities.Parameter{}).
		Once()

	mockOSLayer.EXPECT().
		LookupEnv(paramEnvVar).
		Return(badEnvValue, true).
		Once()

	args := []string{}

	// Act
	p := parser.New(mockOSLayer, mockDefaultParamFactory, mockParamFactory)
	parameters, result, err := p.Parse(args)

	// Assert
	expectedError := messages.New_StartupErrors_BadValueForEnvVar_Error(badEnvValue, paramEnvVar)
	require.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.Nil(t, parameters)
}
