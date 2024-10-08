package main

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"time"
)

func foo(channel chan string) {
	for {
		channel <- "ping"
		fmt.Println("foo is sending: ping")
		message := <-channel
		fmt.Println("foo has recieved: " + message)
		fmt.Println()
	}
}

func bar(channel chan string) {
	for {
		message := <-channel
		fmt.Println("bar has recieved: " + message)
		channel <- "pong"
		fmt.Println("bar is sending: pong")
	}

}

func pingPong() {
	channel := make(chan string)
	go foo(channel)
	go bar(channel)
	// TODO: make channel of type string and pass it to foo and bar
	go foo(nil) // Nil is similar to null. Sending or receiving from a nil chan blocks forever.
	go bar(nil)
	time.Sleep(500 * time.Millisecond)
}

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	pingPong()
}
