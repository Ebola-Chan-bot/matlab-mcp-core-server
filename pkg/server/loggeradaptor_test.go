// Copyright 2026 The MathWorks, Inc.

package server_test

import (
	"errors"
	"testing"

	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestLoggerAdaptor_Debug(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	expectedMessage := "debug message"

	mockLogger.EXPECT().
		Debug(expectedMessage).
		Once()

	adaptor := server.NewLoggerAdaptor(mockLogger)

	// Act
	adaptor.Debug(expectedMessage)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestLoggerAdaptor_Info(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	expectedMessage := "info message"

	mockLogger.EXPECT().
		Info(expectedMessage).
		Once()

	adaptor := server.NewLoggerAdaptor(mockLogger)

	// Act
	adaptor.Info(expectedMessage)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestLoggerAdaptor_Warn(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	expectedMessage := "warn message"

	mockLogger.EXPECT().
		Warn(expectedMessage).
		Once()

	adaptor := server.NewLoggerAdaptor(mockLogger)

	// Act
	adaptor.Warn(expectedMessage)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestLoggerAdaptor_Error(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	expectedMessage := "error message"

	mockLogger.EXPECT().
		Error(expectedMessage).
		Once()

	adaptor := server.NewLoggerAdaptor(mockLogger)

	// Act
	adaptor.Error(expectedMessage)

	// Assert
	// Assertions are verified via deferred mock expectations.
}

func TestLoggerAdaptor_With(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockNewLogger := &entitiesmocks.MockLogger{}
	defer mockNewLogger.AssertExpectations(t)

	expectedKey := "request-id"
	expectedValue := "abc123"
	expectedMessage := "test message"

	mockLogger.EXPECT().
		With(expectedKey, expectedValue).
		Return(mockNewLogger).
		Once()

	mockNewLogger.EXPECT().
		Info(expectedMessage).
		Once()

	adaptor := server.NewLoggerAdaptor(mockLogger)

	// Act
	newAdaptor := adaptor.With(expectedKey, expectedValue)

	// Assert
	assert.NotNil(t, newAdaptor)
	newAdaptor.Info(expectedMessage)
}

func TestLoggerAdaptor_WithError(t *testing.T) {
	// Arrange
	mockLogger := &entitiesmocks.MockLogger{}
	defer mockLogger.AssertExpectations(t)

	mockNewLogger := &entitiesmocks.MockLogger{}
	defer mockNewLogger.AssertExpectations(t)

	expectedError := errors.New("test error")
	expectedMessage := "test message"

	mockLogger.EXPECT().
		WithError(expectedError).
		Return(mockNewLogger).
		Once()

	mockNewLogger.EXPECT().
		Info(expectedMessage).
		Once()

	adaptor := server.NewLoggerAdaptor(mockLogger)

	// Act
	newAdaptor := adaptor.WithError(expectedError)

	// Assert
	assert.NotNil(t, newAdaptor)
	newAdaptor.Info(expectedMessage)
}
