// Copyright 2025 The MathWorks, Inc.

package embeddedconnector_test

import (
	"context"
	"net/http"
	"testing"

	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	httpclientfactorymocks "github.com/matlab/matlab-mcp-core-server/mocks/utils/httpclientfactory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_Ping_DoErrors(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientfactorymocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	mockHttpClient.EXPECT().
		Do(mock.AnythingOfType("*http.Request")).
		Return(nil, assert.AnError).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)
	client.SetPingRetry(20 * time.Millisecond)
	client.SetPingTimeout(30 * time.Millisecond)

	ctx := t.Context()

	// Act
	response := client.Ping(ctx, mockLogger)

	// Assert
	assert.False(t, response.IsAlive)
}

func TestClient_Ping_ContextPropagation(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientfactorymocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	type contextKeyType string
	const contextKey contextKeyType = "uniqueKey"
	const contextKeyValue = "uniqueValue"

	expectedContext := context.WithValue(t.Context(), contextKey, contextKeyValue)

	mockHttpClient.EXPECT().
		Do(mock.MatchedBy(func(request *http.Request) bool {
			return request.Context().Value(contextKey) == contextKeyValue
		})).
		Return(nil, assert.AnError).
		Once()

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)
	client.SetPingRetry(20 * time.Millisecond)
	client.SetPingTimeout(30 * time.Millisecond)

	// Act
	response := client.Ping(expectedContext, mockLogger)

	// Assert
	assert.False(t, response.IsAlive)
}

func TestClient_Ping_Timeout(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockHttpClient := &httpclientfactorymocks.MockHttpClient{}
	defer mockHttpClient.AssertExpectations(t)

	var callCount int
	mockHttpClient.EXPECT().
		Do(mock.AnythingOfType("*http.Request")).
		Run(func(_ *http.Request) { callCount++ }).
		Return(nil, assert.AnError)

	client := embeddedconnector.Client{}
	client.SetHttpClient(mockHttpClient)
	client.SetPingRetry(10 * time.Millisecond)
	client.SetPingTimeout(45 * time.Millisecond)

	ctx := t.Context()

	// Act
	start := time.Now()
	response := client.Ping(ctx, mockLogger)
	duration := time.Since(start)

	// Assert
	assert.False(t, response.IsAlive)
	assert.GreaterOrEqual(t, callCount, 3, "Should have retried at least 3 times")
	assert.LessOrEqual(t, callCount, 5, "Should not retry more than 5 times")
	assert.GreaterOrEqual(t, duration, 45*time.Millisecond, "Should have waited for at least the timeout duration")
}
