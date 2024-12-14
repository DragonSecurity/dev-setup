package efs

import "embed"

//go:embed scripts dotfiles
var EmbeddedFiles embed.FS
