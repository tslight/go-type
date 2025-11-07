package godocs

import (
	"embed"
)

//go:embed *.txt
var EFS embed.FS
