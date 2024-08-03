package pkg

import "embed"

//go:embed tmpls/*.tmpl
var Templates embed.FS
