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
// When cfgPath itself changes, the config is reloaded from disk and overrides
// are re-applied before the next rebuild.
// watch runs until the watcher is closed; it is intended to run in a goroutine.
func watch(cfg *config.Config, embFS fs.FS, cfgPath string, hub *Hub, applyOverrides func(*config.Config) error) {
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
	currentCfg := cfg

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
			// When kwelea.toml changes, reload it and re-apply CLI overrides.
			if isConfigFile(event.Name, cfgPath) {
				if fresh, err := reloadConfig(cfgPath, applyOverrides); err != nil {
					log.Printf("watcher: config reload: %v", err)
				} else {
					currentCfg = fresh
					log.Println("→ config reloaded")
				}
			}
			// When a new directory is created inside docs/, start watching it too.
			if event.Has(fsnotify.Create) {
				if fi, err := os.Stat(event.Name); err == nil && fi.IsDir() {
					_ = w.Add(event.Name)
				}
			}
			// Debounce: reset the timer on every event.
			// Snapshot currentCfg so the closure captures the value at this moment,
			// avoiding a data race if currentCfg is updated before the timer fires.
			if timer != nil {
				timer.Stop()
			}
			snap := currentCfg
			timer = time.AfterFunc(debounce, func() {
				rebuildAndReload(snap, embFS, hub)
			})

		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Printf("watcher error: %v", err)
		}
	}
}

// isConfigFile reports whether eventName refers to the same file as cfgPath
// after resolving both to absolute paths.
func isConfigFile(eventName, cfgPath string) bool {
	a, err1 := filepath.Abs(eventName)
	b, err2 := filepath.Abs(cfgPath)
	if err1 != nil || err2 != nil {
		return false
	}
	return a == b
}

// reloadConfig loads cfgPath from disk and applies overrides on top.
// Returns the fresh config, or an error if loading or overrides fail.
func reloadConfig(cfgPath string, applyOverrides func(*config.Config) error) (*config.Config, error) {
	fresh, err := config.Load(cfgPath)
	if err != nil {
		return nil, err
	}
	if err := applyOverrides(fresh); err != nil {
		return nil, err
	}
	return fresh, nil
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
