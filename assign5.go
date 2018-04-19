package main

import (
  "bufio"
  "fmt"
  "log"
  "net"
)

type client chan<- string // an outgoing message channel

var (
  entering = make(chan client)
  leaving  = make(chan client)
  messages = make(chan string) // all incoming client messages
)

// func echo(c net.Conn, shout string) {
//   fmt.Fprintln(c, shout)
// }

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
  	for {
		select {
		// 1. Send the message to all the clients
			case msg := <-messages:
		// 2. Update the clients map
			case cli := <-entering:
		// 3. Update the clients map and close the client channel
			case cli := <-leaving:
		}
  	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

  // Client's IP address
  	who := conn.RemoteAddr().String()

  // 1. Send a message to the new client confirming their identity
  // e.g. "You are 1.2.3"

  // 2. Broadcast a message to all the clients that a new client has joined
  // e.g. "1.2.3. has joined"

  // 3. Send the client to the entering channel

  // 4. Use a scanner (e.g. bufio.NewScanner) to read from the client and broadcast the incoming messages to all the clients

  // 5. Send the client to the leaving channel

  // 6. Broadcast a message to all clients that the client has left
  // e.g. "1.2.3. has left"

	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
    		fmt.Fprintln(conn, msg)
  	}
}

// func handleConn(c net.Conn) {
//   input := bufio.NewScanner(c)
//   for input.Scan() {
//     echo(c, "Message received: " + input.Text())
//   }
//   c.Close()
// }

func main() {
  l, err := net.Listen("tcp", "localhost:8000")
  if err != nil {
    log.Fatal(err)
  }
  for {
    conn, err := l.Accept()
    if err != nil {
      log.Print(err)
      continue
    }
    handleConn(conn)
  }
}