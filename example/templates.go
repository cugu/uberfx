package example

import _ "embed"

//go:embed minimal/uberfx.hcl
var ConfigTemplate string

//go:embed minimal/main.go
var MainTemplate string
