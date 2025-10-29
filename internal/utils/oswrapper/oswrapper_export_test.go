// Copyright 2025 The MathWorks, Inc.

package oswrapper

import "time"

func (w *OSWrapper) SetCheckParentAliveInterval(interval time.Duration) {
	w.checkParentAliveInterval = interval
}
