// Copyright 2025 The MathWorks, Inc.

package annotations

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Annotations represents tool safety classification metadata.
// All fields are required and use plain bool types to ensure complete specification.
// This design insulates the codebase from MCP SDK's optional field semantics.
type Annotations struct {
	readOnly    bool
	destructive bool
	idempotent  bool
	openWorld   bool
}

// ToToolAnnotations converts to the MCP SDK protocol type.
// Handles the SDK's use of *bool for certain fields.
func (a Annotations) ToToolAnnotations() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		ReadOnlyHint:    a.readOnly,
		DestructiveHint: &a.destructive,
		IdempotentHint:  a.idempotent,
		OpenWorldHint:   &a.openWorld,
	}
}

// NewReadOnlyAnnotations creates annotations for tools that perform inspection
// or query operations without modifying state or executing user code.
func NewReadOnlyAnnotations() Annotations {
	return Annotations{
		readOnly:    true,
		destructive: false,
		idempotent:  false,
		openWorld:   false,
	}
}

// NewDestructiveAnnotations creates annotations for tools that execute code,
// modify state, or interact with external services.
func NewDestructiveAnnotations() Annotations {
	return Annotations{
		readOnly:    false,
		destructive: true,
		idempotent:  false,
		openWorld:   true,
	}
}
