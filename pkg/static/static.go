package static

import (
	"embed"
	_ "embed"
)

var (
	//go:embed shared
	Resources embed.FS
)
