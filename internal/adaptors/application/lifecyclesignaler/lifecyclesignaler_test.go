// Copyright 2025 The MathWorks, Inc.

package lifecyclesignaler_test

import (
	"context"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/lifecyclesignaler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange

	// Act
	signaler := lifecyclesignaler.New()

	// Assert
	assert.NotNil(t, signaler, "Signaler should not be nil")
}

func TestLifecycleSignaler_RequestShutdown_HappyPath(t *testing.T) {
	// Arrange
	signaler := lifecyclesignaler.New()

	// Act
	signaler.RequestShutdown()
	err := signaler.WaitForShutdownToComplete()

	// Assert
	require.NoError(t, err)
}

func TestLifecycleSignaler_AddShutdownFunction_SingleFunction(t *testing.T) {
	// Arrange
	signaler := lifecyclesignaler.New()
	functionCalled := false

	// Act & Assert
	signaler.AddShutdownFunction(func() error {
		functionCalled = true
		return nil
	})

	assert.False(t, functionCalled, "Function called before shutdown requested")

	signaler.RequestShutdown()
	err := signaler.WaitForShutdownToComplete()

	require.NoError(t, err)
	assert.True(t, functionCalled, "Function not called after shutdown requested")
}

func TestLifecycleSignaler_AddShutdownFunction_MultipleFunctions(t *testing.T) {
	// Arrange
	signaler := lifecyclesignaler.New()
	function1Called := false
	function2Called := false

	// Act
	signaler.AddShutdownFunction(func() error {
		function1Called = true
		return nil
	})
	signaler.AddShutdownFunction(func() error {
		function2Called = true
		return nil
	})

	signaler.RequestShutdown()
	err := signaler.WaitForShutdownToComplete()

	// Assert
	require.NoError(t, err)
	assert.True(t, function1Called)
	assert.True(t, function2Called)
}

func TestLifecycleSignaler_AddShutdownFunction_WithError(t *testing.T) {
	// Arrange
	signaler := lifecyclesignaler.New()
	expectedError := assert.AnError

	// Act
	signaler.AddShutdownFunction(func() error {
		return expectedError
	})

	signaler.RequestShutdown()
	err := signaler.WaitForShutdownToComplete()

	// Assert
	assert.ErrorIs(t, err, expectedError)
}

func TestLifecycleSignaler_AddShutdownFunction_MultipleWithOneError(t *testing.T) {
	// Arrange
	signaler := lifecyclesignaler.New()
	expectedError := assert.AnError
	function1Called := false
	function2Called := false
	function3Called := false

	// Act
	signaler.AddShutdownFunction(func() error {
		function1Called = true
		return nil
	})
	signaler.AddShutdownFunction(func() error {
		function2Called = true
		return expectedError
	})
	signaler.AddShutdownFunction(func() error {
		function3Called = true
		return nil
	})

	signaler.RequestShutdown()
	err := signaler.WaitForShutdownToComplete()

	// Assert
	assert.True(t, function1Called)
	assert.True(t, function2Called)
	assert.True(t, function3Called)
	require.ErrorIs(t, err, expectedError)
}

func TestLifecycleSignaler_WaitForShutdownToComplete_WaitAndDoesNotTimeOutWithoutShutdownRequested(t *testing.T) {
	// Arrange
	signaler := lifecyclesignaler.New()
	const shutdownTimeout = 100 * time.Millisecond
	signaler.SetShutdownTimeout(shutdownTimeout)

	const additionalDelay = 100 * time.Millisecond

	// Act & Assert
	errC := make(chan error)
	go func() {
		errC <- signaler.WaitForShutdownToComplete()
	}()

	time.Sleep(shutdownTimeout + additionalDelay)

	select {
	case <-errC:
		t.Fatal("WaitForShutdownToComplete should be blocking until RequestShutdown is called")
	default:
		// No error received, continue
	}

	signaler.RequestShutdown()

	err := <-errC

	assert.NoError(t, err, "There should be no error")
}

func TestLifecycleSignaler_WaitForShutdownToComplete_TimesOut(t *testing.T) {
	// Arrange
	signaler := lifecyclesignaler.New()
	signaler.SetShutdownTimeout(50 * time.Millisecond)

	signaler.AddShutdownFunction(func() error {
		time.Sleep(signaler.ShutdownTimeout() + (100 * time.Millisecond))
		return nil
	})

	start := time.Now()

	// Act
	signaler.RequestShutdown()
	err := signaler.WaitForShutdownToComplete()

	// Assert
	duration := time.Since(start)

	assert.Equal(t, context.DeadlineExceeded, err)
	assert.GreaterOrEqual(t, duration, signaler.ShutdownTimeout(), "Shutdown should timeout after the specified timeout duration")
}

func TestLifecycleSignaler_RequestShutdown_MultipleRequestIsOK(t *testing.T) {
	// Arrange
	signaler := lifecyclesignaler.New()
	callCount := 0

	signaler.AddShutdownFunction(func() error {
		callCount++
		return nil
	})

	// Act
	signaler.RequestShutdown()
	signaler.RequestShutdown()
	err := signaler.WaitForShutdownToComplete()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, callCount, "Function should only be called once")
}
