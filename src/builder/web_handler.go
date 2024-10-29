package builder

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
)

type WebHandler struct {
	port        int
	connections sync.Map
}

func (wh *WebHandler) StartServer() {
	http.HandleFunc("/reload", wh.serveClient)
	http.Handle("/", http.FileServer(http.Dir("public")))

	wh.port = 8080
	for {
		addr := fmt.Sprintf(":%d", wh.port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			wh.port++
			continue // The port is busy, try the next one
		}

		fmt.Printf("Listening on http://localhost%s\n", addr)
		go func() {
			if err := http.Serve(listener, nil); err != nil {
				log.Fatal(err)
			}
		}()

		break
	}
}

func (wh *WebHandler) SendContent(content FileContent) {
	wh.connections.Range(func(_, value any) bool {
		conn := value.(http.ResponseWriter)

		if conn != nil {
			data, err := json.Marshal(map[string]string{
				"block":   content.Name,
				"content": content.Content,
			})
			if err != nil {
				fmt.Println("Error marshaling data:", err)

				return false
			}

			_, err = fmt.Fprintf(conn, "data: %s\n\n", data)
			conn.(http.Flusher).Flush()

			if err != nil {
				fmt.Println("Error sending event data:", err)

				return false
			}
		}

		fmt.Println("Sent event data successfully")

		return true
	})
}

func (wh *WebHandler) serveClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, "data: connected\n\n"); err != nil {
		log.Printf("Error sending initial message: %v", err)
		return
	}
	flusher.Flush()

	wh.connections.Store(r.RemoteAddr, w)
	defer wh.connections.Delete(r.RemoteAddr)

	<-ctx.Done()
	log.Printf("Client %s disconnected: %v", r.RemoteAddr, ctx.Err())
}
