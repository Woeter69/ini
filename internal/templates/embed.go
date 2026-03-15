package templates

import "embed"

//go:embed python/* go/* rust/* bun/* shell/* java/* kotlin/* c/* cpp/*
var FS embed.FS
