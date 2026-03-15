package templates

import "embed"

//go:embed python/* go/* rust/* bun/* shell/* java/* kotlin/*
var FS embed.FS
