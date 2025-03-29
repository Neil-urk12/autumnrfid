package handlers

import (
	"fmt"
	"sync"
	"time"
)

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

type Message struct {
	Event string
	Data  string
}

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

func (b *Broadcaster) Close() {
	close(b.done)
}

func (b *Broadcaster) run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-b.done:
			for client := range b.clients {
				close(client.messages)
				delete(b.clients, client)
			}
			b.Close()
			return

		case client := <-b.register:
			b.mutex.Lock()
			b.clients[client] = true
			b.mutex.Unlock()
			fmt.Printf("Client registered, total clients: %d\n", len(b.clients))

		case client := <-b.unregister:
			b.mutex.Lock()
			if _, ok := b.clients[client]; ok {
				delete(b.clients, client)
				close(client.messages)
			}
			b.mutex.Unlock()
			fmt.Printf("Client unregistered, remaining clients: %d\n", len(b.clients))

		case message := <-b.broadcast:
			// Send to all clients concurrently
			b.mutex.RLock()
			for client := range b.clients {
				// Non-blocking send, skip clients with full buffers
				select {
				case client.messages <- message:
					fmt.Printf("Message sent successfully\n")
				default:
					// Client buffer full, consider unregistering
					go func(c *Client) {
						b.unregister <- c
					}(client)
				}
			}
			b.mutex.RUnlock()

		case <-ticker.C:
			// Periodic check for inactive clients and garbage collection
			// Clean up for leaked resources I think
			fmt.Printf("Active SSE clients: %d\n", len(b.clients))
		}
	}
}

func (b *Broadcaster) Broadcast(event string, data string) {
	if event == "" {
		event = "message"
	}

	message := Message{Event: event, Data: data}

	fmt.Printf("[SSE DEBUG] Raw message being sent:\n%s\n\n", message)

	select {
	case b.broadcast <- message:
		fmt.Printf("[SSE DEBUG] Message sent successfully\n\n")
	default:
		fmt.Println("Warning: SSE message buffer full, dropping message")
	}
}
