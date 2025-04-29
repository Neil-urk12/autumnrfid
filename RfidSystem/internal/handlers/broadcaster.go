package handlers

import (
	"log"
	"sync"
	"time"
)

// Broadcaster manages Server-Sent Events (SSE) clients and broadcasts messages to them.
type Broadcaster struct {
	// clients tracks all connected SSE clients with a bool indicating active status
	clients map[*Client]bool
	// register receives new clients that want to subscribe to SSE events
	register chan *Client
	// unregister receives clients that should be removed from broadcasting
	unregister chan *Client
	// broadcast channel receives messages that should be sent to all connected clients
	broadcast chan Message
	// mutex protects concurrent access to the clients map
	mutex sync.RWMutex
	// done channel signals when the broadcaster should shut down
	done chan struct{}
}

// Message represents a Server-Sent Events (SSE) message with an event type and data.
type Message struct {
	Event string
	Data  string
}

// Client represents a single connected Server-Sent Events (SSE) client.
type Client struct {
	// messages is a buffered channel that receives SSE messages for this specific client
	// Buffer size of 20 allows for message queuing before potential backpressure
	messages chan Message
	// done signals when this client's connection should be terminated
	// closed when the client disconnects or encounters an error
	done chan struct{}
}

var (
	broadcaster *Broadcaster
	once        sync.Once
)

// GetBroadcaster returns the singleton instance of the Broadcaster.
// It initializes the Broadcaster the first time it is called.
func GetBroadcaster() *Broadcaster {
	once.Do(func() {
		broadcaster = &Broadcaster{
			clients:    make(map[*Client]bool),
			register:   make(chan *Client, 10),
			unregister: make(chan *Client, 10),
			broadcast:  make(chan Message, 100),
			done:       make(chan struct{}),
		}
		go broadcaster.run()
	})
	return broadcaster
}

// Close closes the Broadcaster, shutting down all connected client connections.
func (b *Broadcaster) Close() {
	close(b.done)
}

func (b *Broadcaster) run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-b.done:
			// cleanup on shutdown with lock
			b.mutex.Lock()
			for client := range b.clients {
				close(client.messages)
				delete(b.clients, client)
			}
			b.mutex.Unlock()
			return

		case client := <-b.register:
			b.mutex.Lock()
			b.clients[client] = true
			b.mutex.Unlock()
			log.Printf("Client registered, total clients: %d\n", len(b.clients))

		case client := <-b.unregister:
			b.mutex.Lock()
			if _, ok := b.clients[client]; ok {
				delete(b.clients, client)
				close(client.messages)
			}
			b.mutex.Unlock()
			log.Printf("Client unregistered, remaining clients: %d\n", len(b.clients))

		case message := <-b.broadcast:
			// Send to all clients concurrently
			b.mutex.RLock()
			for client := range b.clients {
				// Non-blocking send, skip clients with full buffers
				select {
				case client.messages <- message:
					log.Printf("Message sent successfully\n")
				default:
					// client buffer full: unregister inline
					b.unregister <- client
				}
			}
			b.mutex.RUnlock()
			log.Printf("[SSE] broadcast event=%s to %d clients", message.Event, len(b.clients))

		case <-ticker.C:
			// Periodic check for inactive clients and garbage collection
			// Clean up for leaked resources I think
			log.Printf("Active SSE clients: %d\n", len(b.clients))
		}
	}
}

// Broadcast sends a message with the given event type and data to all connected SSE clients.
// If the event type is empty, it defaults to "message".
func (b *Broadcaster) Broadcast(event string, data string) {
	if event == "" {
		event = "message"
	}

	message := Message{Event: event, Data: data}

	log.Printf("[SSE DEBUG] Raw message being sent:\n%s\n\n", message)

	select {
	case b.broadcast <- message:
		log.Printf("[SSE DEBUG] Message sent successfully\n\n")
	default:
		log.Println("Warning: SSE message buffer full, dropping message")
	}
}
