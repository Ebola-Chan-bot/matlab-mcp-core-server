// Copyright 2025-2026 The MathWorks, Inc.

package configurator

import (
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/codingguidelines"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/plaintextlivecodegeneration"
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
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type Configurator struct {
	configFactory ConfigFactory

	// Multi Session tools
	listAvailableMATLABsTool tools.Tool
	startMATLABSessionTool   tools.Tool
	stopMATLABSessionTool    tools.Tool
	evalInMATLABSessionTool  tools.Tool

	// Single Session tools
	evalInGlobalMATLABSessionTool                  tools.Tool
	checkMATLABCodeInGlobalMATLABSessionTool       tools.Tool
	detectMATLABToolboxesInGlobalMATLABSessionTool tools.Tool
	runMATLABFileInGlobalMATLABSessionTool         tools.Tool
	runMATLABTestFileInGlobalMATLABSessionTool     tools.Tool

	// Resources
	codingGuidelinesResource            resources.Resource
	plaintextlivecodegenerationResource resources.Resource
}

func New(
	configFactory ConfigFactory,

	listAvailableMATLABsTool *listavailablematlabs.Tool,
	startMATLABSessionTool *startmatlabsession.Tool,
	stopMATLABSessionTool *stopmatlabsession.Tool,
	evalInMATLABSessionTool *evalmatlabcodemultisession.Tool,

	evalInGlobalMATLABSessionTool *evalmatlabcodesinglesession.Tool,
	checkMATLABCodeInGlobalMATLABSession *checkmatlabcode.Tool,
	detectMATLABToolboxesInGlobalMATLABSessionTool *detectmatlabtoolboxes.Tool,
	runMATLABFileInGlobalMATLABSessionTool *runmatlabfile.Tool,
	runMATLABTestFileInGlobalMATLABSessionTool *runmatlabtestfile.Tool,

	codingGuidelinesResource *codingguidelines.Resource,
	plaintextlivecodegenerationResource *plaintextlivecodegeneration.Resource,
) *Configurator {
	return &Configurator{
		configFactory: configFactory,

		listAvailableMATLABsTool: listAvailableMATLABsTool,
		startMATLABSessionTool:   startMATLABSessionTool,
		stopMATLABSessionTool:    stopMATLABSessionTool,
		evalInMATLABSessionTool:  evalInMATLABSessionTool,

		evalInGlobalMATLABSessionTool:                  evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSessionTool:       checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInGlobalMATLABSessionTool: detectMATLABToolboxesInGlobalMATLABSessionTool,
		runMATLABFileInGlobalMATLABSessionTool:         runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool:     runMATLABTestFileInGlobalMATLABSessionTool,

		codingGuidelinesResource:            codingGuidelinesResource,
		plaintextlivecodegenerationResource: plaintextlivecodegenerationResource,
	}
}

func (c *Configurator) GetToolsToAdd() ([]tools.Tool, error) {
	config, err := c.configFactory.Config()
	if err != nil {
		return nil, err
	}

	if config.UseSingleMATLABSession() {
		return []tools.Tool{
			c.evalInGlobalMATLABSessionTool,
			c.checkMATLABCodeInGlobalMATLABSessionTool,
			c.detectMATLABToolboxesInGlobalMATLABSessionTool,
			c.runMATLABFileInGlobalMATLABSessionTool,
			c.runMATLABTestFileInGlobalMATLABSessionTool,
		}, nil
	}

	return []tools.Tool{
		c.listAvailableMATLABsTool,
		c.startMATLABSessionTool,
		c.stopMATLABSessionTool,
		c.evalInMATLABSessionTool,
	}, nil
}

func (c *Configurator) GetResourcesToAdd() []resources.Resource {
	return []resources.Resource{
		c.codingGuidelinesResource,
		c.plaintextlivecodegenerationResource,
	}
}
