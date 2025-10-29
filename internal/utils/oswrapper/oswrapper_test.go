// Copyright 2025 The MathWorks, Inc.

package oswrapper_test

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/utils/oswrapper"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	oswrappermocks "github.com/matlab/matlab-mcp-core-server/mocks/utils/oswrapper"
	"github.com/stretchr/testify/assert"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &oswrappermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockOSLayer.EXPECT().
		GOOS().
		Return("linux").
		Once()

	// Act
	wrapper := oswrapper.New(mockOSLayer)

	// Assert
	assert.NotNil(t, wrapper, "OSWrapper instance should not be nil")
}

func TestOSWrapper_FindProcess_Unix_HappyPath(t *testing.T) {
	for _, goos := range []string{"linux", "darwin"} {
		t.Run(goos, func(t *testing.T) {
			// Arrange
			mockOSLayer := &oswrappermocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockProcess := &osfacademocks.MockProcess{}
			defer mockProcess.AssertExpectations(t)

			processPid := 1234

			mockOSLayer.EXPECT().
				GOOS().
				Return(goos).
				Once()

			mockOSLayer.EXPECT().
				FindProcess(processPid).
				Return(mockProcess, nil).
				Once()

			mockProcess.EXPECT().
				Signal(syscall.Signal(0)).
				Return(nil).
				Once()

			wrapper := oswrapper.New(mockOSLayer)

			// Act
			result := wrapper.FindProcess(processPid)

			// Assert
			assert.Equal(t, mockProcess, result, "Should return the found process")
		})
	}
}

func TestOSWrapper_FindProcess_Windows_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &oswrappermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcess := &osfacademocks.MockProcess{}
	defer mockProcess.AssertExpectations(t)

	processPid := 1234

	mockOSLayer.EXPECT().
		GOOS().
		Return("windows").
		Once()

	mockOSLayer.EXPECT().
		FindProcess(processPid).
		Return(mockProcess, nil).
		Once()

	wrapper := oswrapper.New(mockOSLayer)

	// Act
	result := wrapper.FindProcess(processPid)

	// Assert
	assert.Equal(t, mockProcess, result, "Should return the found process")
}

func TestOSWrapper_FindProcess_OSLayerFindProcessError(t *testing.T) {
	// Arrange
	mockOSLayer := &oswrappermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	processPid := 1234
	expectedError := errors.New("process not found")

	mockOSLayer.EXPECT().
		GOOS().
		Return("linux").
		Once()

	mockOSLayer.EXPECT().
		FindProcess(processPid).
		Return(nil, expectedError).
		Once()

	wrapper := oswrapper.New(mockOSLayer)

	// Act
	result := wrapper.FindProcess(processPid)

	// Assert
	assert.Nil(t, result, "Should return nil when OSLayer.FindProcess returns error")
}

func TestOSWrapper_FindProcess_Unix_ProcessSignalError(t *testing.T) {
	for _, goos := range []string{"linux", "darwin"} {
		t.Run(goos, func(t *testing.T) {
			// Arrange
			mockOSLayer := &oswrappermocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockProcess := &osfacademocks.MockProcess{}
			defer mockProcess.AssertExpectations(t)

			processPid := 1234
			signalError := errors.New("process not accessible")

			mockOSLayer.EXPECT().
				GOOS().
				Return(goos).
				Once()

			mockOSLayer.EXPECT().
				FindProcess(processPid).
				Return(mockProcess, nil).
				Once()

			mockProcess.EXPECT().
				Signal(syscall.Signal(0)).
				Return(signalError).
				Once()

			wrapper := oswrapper.New(mockOSLayer)

			// Act
			result := wrapper.FindProcess(processPid)

			// Assert
			assert.Nil(t, result, "Should return nil when process signal check fails on Unix")
		})
	}
}

func TestOSWrapper_WaitForProcessToComplete_Unix_HappyPath(t *testing.T) {
	for _, goos := range []string{"linux", "darwin"} {
		t.Run(goos, func(t *testing.T) {
			// Arrange
			mockOSLayer := &oswrappermocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockProcess := &osfacademocks.MockProcess{}
			defer mockProcess.AssertExpectations(t)

			processPid := 1234

			mockOSLayer.EXPECT().
				GOOS().
				Return(goos).
				Once()

			mockOSLayer.EXPECT().
				FindProcess(processPid).
				Return(mockProcess, nil).
				Times(3) // Simulate 3 ticks, always finding the process

			mockProcess.EXPECT().
				Signal(syscall.Signal(0)).
				Return(nil).
				Twice() // Twice the process responds

			mockProcess.EXPECT().
				Signal(syscall.Signal(0)).
				Return(fmt.Errorf("process not found")).
				Once() // But the third times, it does not

			wrapper := oswrapper.New(mockOSLayer)

			tickInterval := 1 * time.Millisecond
			wrapper.SetCheckParentAliveInterval(tickInterval)

			// Act
			startTime := time.Now()
			wrapper.WaitForProcessToComplete(processPid)
			duration := time.Since(startTime)

			// Assert
			assert.GreaterOrEqual(t, duration, 3*tickInterval, "Should wait three intervals before process ends")
		})
	}
}

func TestOSWrapper_WaitForProcessToComplete_Windows_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &oswrappermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcess := &osfacademocks.MockProcess{}
	defer mockProcess.AssertExpectations(t)

	processPid := 1234

	mockOSLayer.EXPECT().
		GOOS().
		Return("windows").
		Once()

	mockOSLayer.EXPECT().
		FindProcess(processPid).
		Return(mockProcess, nil).
		Once()

	processRunTime := 10 * time.Millisecond
	mockProcess.EXPECT().
		Wait().
		RunAndReturn(func() (*os.ProcessState, error) {
			<-time.After(processRunTime)
			return nil, nil
		}).
		Once()

	wrapper := oswrapper.New(mockOSLayer)

	// Act
	startTime := time.Now()
	wrapper.WaitForProcessToComplete(processPid)
	duration := time.Since(startTime)

	// Assert
	assert.GreaterOrEqual(t, duration, processRunTime, "Should wait for process to complete on Windows")
}
