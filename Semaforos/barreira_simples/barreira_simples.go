package main

import (
	"fmt"
	"sync"
	"time"
)

const N = 4

var arrived int
var mutex sync.Mutex
var barrier = make(chan struct{})

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Goroutine %d: Parte 1\n", id)
	time.Sleep(time.Duration(id) * 200 * time.Millisecond)

	mutex.Lock()
	arrived++
	if arrived == N {
		close(barrier)
	}
	mutex.Unlock()

	<-barrier
	fmt.Printf("Goroutine %d: Parte 2\n", id)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go worker(i, &wg)
	}

	wg.Wait()
}