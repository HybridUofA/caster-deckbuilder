// Package localdata defines the shared on-disk card data used by every application.
package localdata

import "path/filepath"

// SharedApplicationID preserves the existing storage directory for installed users.
const SharedApplicationID = "io.github.hybriduofa.casterdeckbuilder"

// Paths contains every locally managed card database and artwork path.
type Paths struct {
	Root          string
	CardDatabase  string
	CardListHash  string
	Images        string
	Thumbnails    string
	SetupComplete string
}

// NewPaths derives every managed data path from the shared per-user root.
func NewPaths(root string) Paths {
	return Paths{
		Root:          root,
		CardDatabase:  filepath.Join(root, "cards.json"),
		CardListHash:  filepath.Join(root, "cardlist.sha256"),
		Images:        filepath.Join(root, "images"),
		Thumbnails:    filepath.Join(root, "thumbnails"),
		SetupComplete: filepath.Join(root, ".setup-complete-v1"),
	}
}
