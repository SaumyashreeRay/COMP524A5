package main

import (
  "fmt"
  "log"
  "net"
  "time"
  "bufio"
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
            for cli := range clients {
            	select {
								case cli <- msg:
								default:
            	}
            }
        case cli := <-entering:
            clients[cli] = true

        case cli := <-leaving:
            delete(clients, cli)
            close(cli)
        }
    }
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
    }
}

func handleConn(conn net.Conn) {
    ch := make(chan string) // outgoing client messages
    go clientWriter(conn, ch)

    who := conn.RemoteAddr().String()
    ch <- "You are " + who
    messages <- who + " has joined"
    entering <- ch
    timer := time.NewTimer(60 * time.Second)

    go func () {
  		<-timer.C
  		leaving <- ch
    	messages <- who + " has left"
  	}()

    input := bufio.NewScanner(conn)
    for input.Scan() {
        messages <- who + ": " + input.Text()
        timer.Reset(60 * time.Second)
    }

    leaving <- ch
    messages <- who + " has left"
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