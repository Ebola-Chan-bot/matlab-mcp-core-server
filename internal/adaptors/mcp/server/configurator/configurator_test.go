// Copyright 2025 The MathWorks, Inc.

package configurator_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/codingguidelines"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/configurator"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	evalmatlabmultisession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/listavailablematlabs"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/startmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/stopmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/checkmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/detectmatlabtoolboxes"
	evalmatlabsinglesession "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabfile"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabtestfile"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/server/configurator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}

	// Act
	result := configurator.New(
		mockConfig,
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInSingleSessionTool,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		codingGuidelinesResource,
	)

	// Assert
	require.NotNil(t, result, "Configurator should not be nil")
}

func TestConfigurator_GetToolsToAdd_MultipleMATLABSession_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(false).
		Once()

	c := configurator.New(
		mockConfig,
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInSingleSessionTool,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		codingGuidelinesResource,
	)

	// Act
	toolsToAdd := c.GetToolsToAdd()

	// Assert
	assert.ElementsMatch(t, toolsToAdd, []tools.Tool{
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
	}, "GetToolsToAdd should return all the injected tools for multi session")
}

func TestConfigurator_GetToolsToAdd_SingleMATLABSession_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	c := configurator.New(
		mockConfig,
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInSingleSessionTool,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		codingGuidelinesResource,
	)

	// Act
	toolsToAdd := c.GetToolsToAdd()

	// Assert
	assert.ElementsMatch(t, toolsToAdd, []tools.Tool{
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		detectMATLABToolboxesInSingleSessionTool,
	}, "GetToolsToAdd should all injected tools for single session")
}

func TestConfigurator_GetResourcesToAdd_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &mocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	listAvailableMATLABsTool := &listavailablematlabs.Tool{}
	startMATLABSessionTool := &startmatlabsession.Tool{}
	stopMATLABSessionTool := &stopmatlabsession.Tool{}
	evalInMATLABSessionTool := &evalmatlabmultisession.Tool{}
	evalInGlobalMATLABSessionTool := &evalmatlabsinglesession.Tool{}
	checkMATLABCodeInGlobalMATLABSession := &checkmatlabcode.Tool{}
	detectMATLABToolboxesInSingleSessionTool := &detectmatlabtoolboxes.Tool{}
	runMATLABFileInGlobalMATLABSessionTool := &runmatlabfile.Tool{}
	runMATLABTestFileInGlobalMATLABSessionTool := &runmatlabtestfile.Tool{}
	codingGuidelinesResource := &codingguidelines.Resource{}

	c := configurator.New(
		mockConfig,
		listAvailableMATLABsTool,
		startMATLABSessionTool,
		stopMATLABSessionTool,
		evalInMATLABSessionTool,
		evalInGlobalMATLABSessionTool,
		checkMATLABCodeInGlobalMATLABSession,
		detectMATLABToolboxesInSingleSessionTool,
		runMATLABFileInGlobalMATLABSessionTool,
		runMATLABTestFileInGlobalMATLABSessionTool,
		codingGuidelinesResource,
	)

	// Act
	result := c.GetResourcesToAdd()

	// Assert
	assert.ElementsMatch(t, []resources.Resource{codingGuidelinesResource}, result)
}
