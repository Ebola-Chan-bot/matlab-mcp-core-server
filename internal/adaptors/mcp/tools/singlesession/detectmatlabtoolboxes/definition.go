// Copyright 2025 The MathWorks, Inc.

package detectmatlabtoolboxes

const (
	name        = "detect_matlab_toolboxes"
	title       = "Detect MATLAB Toolboxes"
	description = "List installed MATLAB toolboxes with their versions and installation status."
)

type Args struct {
}

type ReturnArgs struct {
	InstallationInfo string `json:"installation_info" jsonschema:"Output of the 'ver' command showing installed MATLAB toolboxes."`
}
