// Copyright 2025 The MathWorks, Inc.

package evalmatlabcode

const (
	name        = "evaluate_matlab_code"
	title       = "Evaluate MATLAB Code"
	description = "Evaluate arbitrary MATLAB code (`code`) within a specified project directory (`project_path`) context in an existing MATLAB session. Returns the command window output from code execution."
)

type Args struct {
	ProjectPath string `json:"project_path" jsonschema:"The full path to the project directory - Becomes MATLAB's working directory during execution - Folder must exist - Example: C:\\Users\\username\\matlab-project or /home/user/research."`
	Code        string `json:"code"         jsonschema:"The MATLAB code to evaluate."`
}
