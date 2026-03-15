package templates

import "embed"

//go:embed python/* go/* rust/* bun/* shell/*
var FS embed.FS
