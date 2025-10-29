// Copyright 2025 The MathWorks, Inc.

package config

import "testing"

func SetVersionLikeLDFLAGSWould(t *testing.T, newVersion string) {
	t.Cleanup(func() {
		version = unsetVersion
	})
	version = newVersion
}
