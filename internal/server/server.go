package server

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/engineervix/kwelea/internal/builder"
	"github.com/engineervix/kwelea/internal/config"
	"github.com/engineervix/kwelea/internal/nav"
)

const (
	wsPath  = "/_kwelea/ws"
	maxPort = 10 // try up to base+10 before giving up
)

// Start is the entry point for `kwelea serve`. It:
//  1. Runs an initial build in devMode
//  2. Resolves an available port (base..base+10)
//  3. Starts the WebSocket hub and file watcher
//  4. Serves the output directory as a static site
//  5. Opens the browser (unless disabled)
//  6. Blocks until SIGINT/SIGTERM, then shuts down cleanly
func Start(cfg *config.Config, embFS fs.FS, cfgPath string) error {
	// Initial build so there's something to serve immediately.
	if err := initialBuild(cfg, embFS); err != nil {
		return err
	}

	listener, port, err := resolvePort(cfg.Serve.Port)
	if err != nil {
		return err
	}

	hub := newHub()
	go hub.run()

	mux := http.NewServeMux()
	mux.Handle(wsPath, wsHandler(hub))
	mux.Handle("/", http.FileServer(http.Dir(cfg.Build.OutputDir)))

	srv := &http.Server{Handler: mux}

	go watch(cfg, embFS, cfgPath, hub)

	addr := fmt.Sprintf("http://localhost:%d", port)
	fmt.Printf("→ starting dev server on %s\n", addr)

	if cfg.Serve.OpenBrowser {
		openBrowser(addr)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("\n→ shutting down")
	return srv.Shutdown(context.Background())
}

// resolvePort tries to bind to base and each subsequent port up to base+maxPort.
// It returns the listener (already bound) and the resolved port number.
func resolvePort(base int) (net.Listener, int, error) {
	for offset := 0; offset <= maxPort; offset++ {
		port := base + offset
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			if offset > 0 {
				fmt.Printf("  port %d busy — using %d instead\n", base, port)
			}
			return ln, port, nil
		}
	}
	return nil, 0, fmt.Errorf(
		"ports %d–%d are all in use; free one or set [serve] port in kwelea.toml",
		base, base+maxPort,
	)
}

// initialBuild runs a full devMode build before the server starts accepting
// requests, so there's a complete site ready to serve immediately.
func initialBuild(cfg *config.Config, embFS fs.FS) error {
	site, err := nav.NewSite(cfg)
	if err != nil {
		return fmt.Errorf("initial build: nav: %w", err)
	}
	site.BasePath = "" // dev server always serves from /, never a subpath
	if err := builder.Build(site, embFS, true); err != nil {
		return fmt.Errorf("initial build: %w", err)
	}
	return nil
}

// openBrowser opens url in the system default browser without blocking.
func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		cmd, args = "open", []string{url}
	case "windows":
		cmd, args = "cmd", []string{"/c", "start", url}
	default: // linux, bsd, etc.
		cmd, args = "xdg-open", []string{url}
	}
	if err := exec.Command(cmd, args...).Start(); err != nil {
		log.Printf("could not open browser: %v", err)
	}
}
