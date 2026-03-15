package templates

import "embed"

//go:embed python/* go/* rust/* bun/* shell/* java/*
var FS embed.FS
