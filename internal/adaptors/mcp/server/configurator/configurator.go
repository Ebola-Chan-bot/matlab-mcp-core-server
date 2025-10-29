// Copyright 2025 The MathWorks, Inc.

package configurator

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	evalmatlabcodemultisession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/listavailablematlabs"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/startmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/stopmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/checkmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/detectmatlabtoolboxes"
	evalmatlabcodesinglesession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabfile"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabtestfile"
)

type Config interface {
	UseSingleMATLABSession() bool
}

type Configurator struct {
	config Config

	// Multi Session
	listAvailableMATLABsTool tools.Tool
	startMATLABSessionTool   tools.Tool
	stopMATLABSessionTool    tools.Tool
	evalInMATLABSessionTool  tools.Tool

	// Single Session
	evalInGlobalMATLABSessionTool                  tools.Tool
	checkMATLABCodeInGlobalMATLABSessionTool       tools.Tool
	detectMATLABToolboxesInGlobalMATLABSessionTool tools.Tool
	runMATLABFileInGlobalMATLABSessionTool         tools.Tool
	runMATLABTestFileInGlobalMATLABSessionTool     tools.Tool
}

func New(
	config Config,

	listAvailableMATLABsTool *listavailablematlabs.Tool,
	startMATLABSessionTool *startmatlabsession.Tool,
	stopMATLABSessionTool *stopmatlabsession.Tool,
	evalInMATLABSessionTool *evalmatlabcodemultisession.Tool,

	evalInGlobalMATLABSessionTool *evalmatlabcodesinglesession.Tool,
	checkMATLABCodeInGlobalMATLABSession *checkmatlabcode.Tool,
	detectMATLABToolboxesInGlobalMATLABSessionTool *detectmatlabtoolboxes.Tool,
	runMATLABFileInGlobalMATLABSessionTool *runmatlabfile.Tool,
	runMATLABTestFileInGlobalMATLABSessionTool *runmatlabtestfile.Tool,
) *Configurator {
	return &Configurator{
		config: config,

		listAvailableMATLABsTool: listAvailableMATLABsTool,
		startMATLABSessionTool:   startMATLABSessionTool,
		stopMATLABSessionTool:    stopMATLABSessionTool,
		evalInMATLABSessionTool:  evalInMATLABSessionTool,

		evalInGlobalMATLABSessionTool:                  evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSessionTool:       checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInGlobalMATLABSessionTool: detectMATLABToolboxesInGlobalMATLABSessionTool,
		runMATLABFileInGlobalMATLABSessionTool:         runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool:     runMATLABTestFileInGlobalMATLABSessionTool,
	}
}

func (c *Configurator) GetToolsToAdd() []tools.Tool {
	// Choose which tool to expose

	if c.config.UseSingleMATLABSession() {
		return []tools.Tool{
			c.evalInGlobalMATLABSessionTool,
			c.checkMATLABCodeInGlobalMATLABSessionTool,
			c.detectMATLABToolboxesInGlobalMATLABSessionTool,
			c.runMATLABFileInGlobalMATLABSessionTool,
			c.runMATLABTestFileInGlobalMATLABSessionTool,
		}
	}

	return []tools.Tool{
		c.listAvailableMATLABsTool,
		c.startMATLABSessionTool,
		c.stopMATLABSessionTool,
		c.evalInMATLABSessionTool,
	}
}
