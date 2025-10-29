// Copyright 2025 The MathWorks, Inc.

package execsuite

import (
	"time"

	"github.com/stretchr/testify/suite"
)

const defaultConnectionTimeout = 5 * time.Second

type Suite struct {
	suite.Suite
	matlabMCPServerBinariesPath string
}

func NewSuite(matlabMCPServerBinariesPath string) *Suite {
	return &Suite{
		matlabMCPServerBinariesPath: matlabMCPServerBinariesPath,
	}
}
