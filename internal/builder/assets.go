package builder

import (
	"io/fs"
	"os"
	"path/filepath"
)

// copyAssets walks the embedded "assets" directory and copies every file into
// outputDir, preserving the directory structure. The result is
// outputDir/assets/theme.css, outputDir/assets/vendor/flexsearch.bundle.js, etc.
func copyAssets(embFS fs.FS, outputDir string) error {
	return fs.WalkDir(embFS, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		dest := filepath.Join(outputDir, filepath.FromSlash(path))
		if d.IsDir() {
			return os.MkdirAll(dest, 0o755)
		}
		data, err := fs.ReadFile(embFS, path)
		if err != nil {
			return err
		}
		return os.WriteFile(dest, data, 0o644)
	})
}
