package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ChrisGora/semaphore"
)

type buffer struct {
	b                 []int
	size, read, write int
}

func newBuffer(size int) buffer {
	return buffer{
		b:     make([]int, size),
		size:  size,
		read:  0,
		write: 0,
	}
}

func (buffer *buffer) get() int {
	x := buffer.b[buffer.read]
	fmt.Println("Get\t", x, "\t", buffer)
	buffer.read = (buffer.read + 1) % len(buffer.b)
	return x
}

func (buffer *buffer) put(x int) {

	buffer.b[buffer.write] = x
	fmt.Println("Put\t", x, "\t", buffer)
	buffer.write = (buffer.write + 1) % len(buffer.b)
}

func producer(buffer *buffer, start, delta int, space, work, mutex semaphore.Semaphore) {
	x := start
	for {
		space.Wait()
		mutex.Wait()
		buffer.put(x)
		x = x + delta
		work.Post()
		mutex.Post()
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}

func consumer(buffer *buffer, space, work, mutex semaphore.Semaphore) {
	for {
		work.Wait()
		mutex.Wait()
		_ = buffer.get()
		space.Post()
		mutex.Post()
		time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
	}
}

func main() {
	buffer := newBuffer(5)

	spaceAvailable := semaphore.Init(5, 5)
	workAvailable := semaphore.Init(5, 0)
	mutex := semaphore.Init(1, 1)

	go producer(&buffer, 1, 1, spaceAvailable, workAvailable, mutex)
	go producer(&buffer, 1000, -1, spaceAvailable, workAvailable, mutex)

	consumer(&buffer, spaceAvailable, workAvailable, mutex)
}
