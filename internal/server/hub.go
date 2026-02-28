package server

// client is a buffered channel that receives a signal when a reload is broadcast.
type client chan struct{}

// Hub tracks all connected WebSocket clients and fans out reload signals to them.
// It runs as a single goroutine (via hub.run) so the client map needs no mutex.
type Hub struct {
	register   chan client
	unregister chan client
	broadcast  chan struct{}
	clients    map[client]struct{}
}

func newHub() *Hub {
	return &Hub{
		register:   make(chan client),
		unregister: make(chan client),
		broadcast:  make(chan struct{}, 1),
		clients:    make(map[client]struct{}),
	}
}

// run is the hub's event loop. It must be called in its own goroutine.
func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = struct{}{}
		case c := <-h.unregister:
			delete(h.clients, c)
			close(c)
		case <-h.broadcast:
			for c := range h.clients {
				select {
				case c <- struct{}{}:
				default: // client is slow — skip this cycle
				}
			}
		}
	}
}

// Reload broadcasts a reload signal to every connected client.
// Non-blocking: if a previous signal hasn't been consumed it's coalesced.
func (h *Hub) Reload() {
	select {
	case h.broadcast <- struct{}{}:
	default:
	}
}
