package localdata

import (
	"path/filepath"
	"testing"
)

// TestNewPathsKeepsManagedDataUnderSharedRoot verifies applications cannot diverge their caches.
func TestNewPathsKeepsManagedDataUnderSharedRoot(t *testing.T) {
	root := t.TempDir()
	paths := NewPaths(root)
	want := map[string]string{
		"database":   filepath.Join(root, "cards.json"),
		"hash":       filepath.Join(root, "cardlist.sha256"),
		"images":     filepath.Join(root, "images"),
		"thumbnails": filepath.Join(root, "thumbnails"),
		"setup":      filepath.Join(root, ".setup-complete-v1"),
	}
	got := map[string]string{
		"database":   paths.CardDatabase,
		"hash":       paths.CardListHash,
		"images":     paths.Images,
		"thumbnails": paths.Thumbnails,
		"setup":      paths.SetupComplete,
	}
	for name, expected := range want {
		if got[name] != expected {
			t.Fatalf("%s path = %q, want %q", name, got[name], expected)
		}
	}
}
