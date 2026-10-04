// Package content is the site's Markdown, one directory per locale.
// Directories starting with "_" hold files that pages include.
package content

import "embed"

//go:embed all:en all:ja
var FS embed.FS
