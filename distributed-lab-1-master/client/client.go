package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func read(conn net.Conn) {
	//TODO In a continuous loop, read a message from the server and display it.
	reader := bufio.NewReader(conn)
	for {
		msg, _ := reader.ReadString('\n')
		fmt.Print("\nmessage from another client: ", msg)
		fmt.Print("Enter text->")
	}
}

func write(conn net.Conn) {
	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Enter text->")
		msg, _ := stdin.ReadString('\n')
		fmt.Fprintln(conn, msg)
	}
	//TODO Continually get input from the user and send messages to the server.
}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()
	//TODO Try to connect to the server
	conn, err := net.Dial("tcp", *addrPtr)
	if err != nil {
		fmt.Println(err)
		return
	}
	//TODO Start asynchronously reading and displaying messages
	//TODO Start getting and sending user messages.
	go write(conn)
	go read(conn)
	select {}
}
