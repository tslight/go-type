package books

import "embed"

//go:embed *.txt *.json
var EFS embed.FS
