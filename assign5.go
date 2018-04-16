package main

import(
"net"
"fmt"
"os"
)

func main() {
	
	//server is listening for new potential connections
	l, err := net.Listen("localhost", ":8080")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// listener is closed when program closes.
	defer l.Close()

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error: ", err.Error())
			os.Exit(1)
		}
		// Handle connections 
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn){
	// buffer hold incoming data from client
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
  	if err != nil {
    	fmt.Println("Error:", err.Error())
	  }
	conn.Write([]byte("Message received."))
	conn.Close()
}