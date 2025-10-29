// Copyright 2025 The MathWorks, Inc.

// This package implements a watchdog to run alongside the MATLAB MCP Core Server.
// As the MATLAB MCP Core Server may start multiple external dependencies (e.g. MATLAB processes), in the event of forceful shutdown,
// this watchdog ensures proper termination and cleanup of the external dependencies.
package watchdog
