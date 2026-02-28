package server

import (
	"crypto/sha1"
	"encoding/base64"
	"net"
	"net/http"
)

// wsGUID is the fixed magic string from the WebSocket RFC (RFC 6455 §1.3).
const wsGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// wsHandler returns an http.HandlerFunc that upgrades the connection to
// WebSocket, registers the client with the hub, and writes a reload frame
// each time the hub broadcasts.
func wsHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") != "websocket" {
			http.Error(w, "websocket upgrade required", http.StatusBadRequest)
			return
		}
		key := r.Header.Get("Sec-WebSocket-Key")
		if key == "" {
			http.Error(w, "missing Sec-WebSocket-Key", http.StatusBadRequest)
			return
		}

		// Compute the accept token: base64(SHA1(key + wsGUID)).
		h := sha1.New()
		h.Write([]byte(key + wsGUID))
		accept := base64.StdEncoding.EncodeToString(h.Sum(nil))

		// Hijack the connection so we can write raw WebSocket frames.
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "hijack unsupported", http.StatusInternalServerError)
			return
		}
		conn, rw, err := hj.Hijack()
		if err != nil {
			return
		}

		// Send the 101 Switching Protocols handshake.
		resp := "HTTP/1.1 101 Switching Protocols\r\n" +
			"Upgrade: websocket\r\n" +
			"Connection: Upgrade\r\n" +
			"Sec-WebSocket-Accept: " + accept + "\r\n\r\n"
		if _, err := rw.WriteString(resp); err != nil {
			conn.Close()
			return
		}
		if err := rw.Flush(); err != nil {
			conn.Close()
			return
		}

		// Register this client and wait for reload signals.
		ch := make(client, 1)
		hub.register <- ch
		defer func() { hub.unregister <- ch }()

		for range ch {
			if err := sendReloadFrame(conn); err != nil {
				return
			}
		}
	}
}

// sendReloadFrame writes a minimal, unmasked WebSocket text frame carrying the
// literal string "reload". Server frames are never masked (RFC 6455 §5.1).
// Frame layout: 0x81 (FIN+text), 0x06 (6-byte payload), then "reload".
func sendReloadFrame(conn net.Conn) error {
	frame := [8]byte{0x81, 0x06, 'r', 'e', 'l', 'o', 'a', 'd'}
	_, err := conn.Write(frame[:])
	return err
}
