// Copyright 2025 The MathWorks, Inc.

package checkmatlabcode

const (
	name        = "check_matlab_code"
	title       = "Check MATLAB Code"
	description = "Perform static code analysis on a MATLAB script (`script_path`) using MATLAB's built-in checkcode function in an existing MATLAB session. Returns warnings about coding style, potential errors, deprecated functions, performance issues, and best practice violations. This is a non-destructive, read-only operation that helps identify code quality issues without executing the script."
)

type Args struct {
	ScriptPath string `json:"script_path" jsonschema:"The full absolute path to the MATLAB script file to analyze - Must be a .m file that exists - File is not modified during analysis - Example: C:\\Users\\username\\matlab\\myFunction.m or /home/user/scripts/analysis.m."`
}

type ReturnArgs struct {
	CheckCodeOutput []string `json:"checkcode_messages" jsonschema:"List of code style and correctness warnings."`
}
