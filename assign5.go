package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "time"
)

type client chan<- string // an outgoing message channel

var (
    entering = make(chan client)
    leaving  = make(chan client)
    messages = make(chan string) // all incoming client messages
)

func broadcaster() {
      clients := make(map[client]bool) // all connected clients
      for {
          select {
          
          case msg := <-messages:
            // 1. Send the message to all the clients
                for cli:=range clients{
                    cli <-msg
                }
          case cli := <-entering:
            // 2. Update the clients map
                clients[cli]=true
          case cli := <-leaving:
            // 3. Update the clients map and close the client channel
                delete(clients,cli)
                close(cli)
          }
      }
}

func handleConn(conn net.Conn) {
    ch := make(chan string) // outgoing client messages
    go clientWriter(conn, ch)

    timer1 := time.NewTimer(5 * time.Second)

    // Client's IP address
    who := conn.RemoteAddr().String()

    // 1. Send a message to the new client confirming their identity
    // e.g. "You are 1.2.3"
    ch <- "You are " + who

    // 2. Broadcast a message to all the clients that a new client has joined
    // e.g. "1.2.3. has joined"
    messages <- who + " has joined"

    // 3. Send the client to the entering channel
    entering <- ch

    // 4. Use a scanner (e.g. bufio.NewScanner) to read from the 
    //client and broadcast the incoming messages to all the clients
    input := bufio.NewScanner(conn)
    for input.Scan() {
        //broadcast message to all other clients
        messages <- who+":"+input.Text()
    }

    conn.Close()

    // 5. Send the client to the leaving channel
    leaving<-ch

    // 6. Broadcast a message to all clients that the client has left
    // e.g. "1.2.3. has left"
    messages<- who+" has left"
    conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
    for msg := range ch {
        fmt.Fprintln(conn, msg)
    }
}

func main() {
    listener, err := net.Listen("tcp", "localhost:8000")
    if err != nil {
        log.Fatal(err)
    }

    go broadcaster()
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Print(err)
            continue
        }
        go handleConn(conn)
    }
}