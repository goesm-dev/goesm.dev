// Package content is the Markdown of the sites, one directory per locale:
// en and ja for goesm, gosfc/en and gosfc/ja for gosfc. Directories starting
// with "_" hold files that pages include.
package content

import "embed"

//go:embed all:en all:ja all:gosfc
var FS embed.FS
