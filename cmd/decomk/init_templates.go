package main

import _ "embed"

var (
	// initDevcontainerJSONTemplate is the stage-0 template for .devcontainer/devcontainer.json.
	//
	//go:embed templates/devcontainer.json.tmpl
	initDevcontainerJSONTemplate string

	// initStage0ScriptTemplate is the stage-0 template for .devcontainer/decomk-stage0.sh.
	//
	//go:embed templates/decomk-stage0.sh.tmpl
	initStage0ScriptTemplate string
)
