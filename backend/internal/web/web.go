package web

import (
	"embed"
	"io/fs"
)

//go:embed static
var files embed.FS

// StaticFS returns the embedded frontend files, rooted so paths like
// "index.html" and "styles.css" resolve directly.
func StaticFS() fs.FS {
	sub, err := fs.Sub(files, "static")
	if err != nil {
		panic(err)
	}
	return sub
}
