package templates

import "embed"

//go:embed python/* go/* rust/* bun/*
var FS embed.FS
