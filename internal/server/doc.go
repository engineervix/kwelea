// Package server implements the live-reloading development server for kwelea.
//
// The entry point is [Start], which:
//
//  1. Runs an initial [builder.Build] to populate the output directory.
//  2. Binds to the configured port (auto-incrementing up to +10 if occupied).
//  3. Starts a [Hub] goroutine that broadcasts WebSocket reload signals to all
//     connected browser tabs.
//  4. Starts an fsnotify watcher over the docs directory and kwelea.toml; on
//     any change it rebuilds and calls [Hub.Reload].
//  5. Serves the output directory over HTTP and handles WebSocket upgrades at
//     /_kwelea/ws.
//  6. Optionally opens the browser at the resolved URL.
//
// The server shuts down cleanly on SIGTERM or SIGINT.
package server
