// Copyright 2025 The MathWorks, Inc.

package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/wire"
)

func main() {
	modeSelector, err := wire.InitializeModeSelector()
	if err != nil {
		// As we failed to even initialize, we cannot use a LoggerFactory,
		// and we can't assume whatever failed had a logger factory to log the error either.
		// In this case, we use the default slog.
		slog.With("error", err).Error("Failed to initialize MATLAB MCP Core Server.")
		os.Exit(1)
	}

	ctx := context.Background()
	err = modeSelector.StartAndWaitForCompletion(ctx)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
