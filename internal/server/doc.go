// Package server implements the live-reloading development server for kwelea.
//
// Phase 6 work: HTTP file server over output_dir, port auto-increment logic,
// fsnotify watcher, WebSocket broadcast for browser live reload.
//
// Dependency (added in Phase 6):
//   - github.com/fsnotify/fsnotify
package server
