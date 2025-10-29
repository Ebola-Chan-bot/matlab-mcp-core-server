// Copyright 2025 The MathWorks, Inc.

package evalmatlabcode

const (
	name        = "eval_in_matlab_session"
	title       = "Evaluate MATLAB Code in a MATLAB Session"
	description = "Evaluate arbitrary MATLAB code (`code`) within a specified project directory (`project_path`) context in an existing MATLAB session, given its session ID (`session_id`)."
)

type Args struct {
	SessionID   int    `json:"session_id"   jsonschema:"The ID of the MATLAB session in which to evaluate the code."`
	ProjectPath string `json:"project_path" jsonschema:"The full path to the project directory - Becomes MATLAB's working directory during execution - Folder must exist - Example: C:\\Users\\username\\matlab-project or /home/user/research."`
	Code        string `json:"code"         jsonschema:"The MATLAB code to evaluate."`
}
