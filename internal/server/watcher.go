package server

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/engineervix/kwelea/internal/builder"
	"github.com/engineervix/kwelea/internal/config"
	"github.com/engineervix/kwelea/internal/nav"
)

// watch starts an fsnotify watcher over docsDir (recursively) and cfgPath.
// File-change events are debounced: after a 300 ms quiet period the site is
// rebuilt in devMode and hub.Reload() broadcasts to connected browsers.
// watch runs until the watcher is closed; it is intended to run in a goroutine.
func watch(cfg *config.Config, embFS fs.FS, cfgPath string, hub *Hub) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("watcher: %v", err)
		return
	}
	defer w.Close()

	addDirRecursive(w, cfg.Build.DocsDir)

	// Watch the directory that contains the config file so we catch writes to it.
	if dir := filepath.Dir(cfgPath); dir != "" {
		_ = w.Add(dir)
	}

	const debounce = 300 * time.Millisecond
	var timer *time.Timer

	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}
			// Ignore chmod-only events; they don't change content.
			if event.Has(fsnotify.Chmod) {
				continue
			}
			// When a new directory is created inside docs/, start watching it too.
			if event.Has(fsnotify.Create) {
				if fi, err := os.Stat(event.Name); err == nil && fi.IsDir() {
					_ = w.Add(event.Name)
				}
			}
			// Debounce: reset the timer on every event.
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(debounce, func() {
				rebuildAndReload(cfg, embFS, hub)
			})

		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Printf("watcher error: %v", err)
		}
	}
}

// addDirRecursive adds root and every subdirectory under it to the watcher.
func addDirRecursive(w *fsnotify.Watcher, root string) {
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if d.IsDir() {
			_ = w.Add(path)
		}
		return nil
	})
}

// rebuildAndReload re-runs the full build pipeline in dev mode and, on success,
// signals all connected browsers to reload.
func rebuildAndReload(cfg *config.Config, embFS fs.FS, hub *Hub) {
	site, err := nav.NewSite(cfg)
	if err != nil {
		log.Printf("rebuild: nav: %v", err)
		return
	}
	site.BasePath = "" // dev server always serves from /, never a subpath
	if err := builder.Build(site, embFS, true); err != nil {
		log.Printf("rebuild: %v", err)
		return
	}
	log.Println("→ rebuilt")
	hub.Reload()
}
