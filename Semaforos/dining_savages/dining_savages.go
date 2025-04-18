package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numSavages = 5
	capacity   = 3
)

var (
	mutex        sync.Mutex
	emptyPot     = make(chan struct{}, 1) // apenas 1 chamado para cozinheiro por vez
	fullPot      = make(chan struct{})    // canal de sinalização para reabastecimento
	servingsLeft = 0
)

func savage(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 2; i++ {
		for {
			mutex.Lock()
			if servingsLeft == 0 {
				select {
				case emptyPot <- struct{}{}: // somente o primeiro a ver a panela vazia chama o cozinheiro
					fmt.Printf("Savage %d: Panela vazia! Chamando o cozinheiro...\n", id)
				default:
					// outro savage já chamou o cozinheiro
				}
				mutex.Unlock()
				<-fullPot // aguarda o cozinheiro reabastecer
				continue  // tenta comer novamente
			}

			servingsLeft--
			fmt.Printf("Savage %d: Comeu! Porções restantes: %d\n", id, servingsLeft)
			mutex.Unlock()
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			break
		}
	}
}

func cook() {
	for {
		_, ok := <-emptyPot
		if !ok {
			return // canal fechado, encerrar goroutine do cozinheiro
		}
		fmt.Println("Cozinheiro: Reabastecendo a panela...")
		time.Sleep(1 * time.Second)

		mutex.Lock()
		servingsLeft = capacity
		mutex.Unlock()

		// acorda todos os selvagens esperando
		for i := 0; i < capacity; i++ {
			fullPot <- struct{}{}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	wg.Add(numSavages)

	go cook()

	for i := 0; i < numSavages; i++ {
		go savage(i, &wg)
	}

	wg.Wait()

	// Finalizar goroutine do cozinheiro
	close(emptyPot)
	fmt.Println("Todos os selvagens comeram.")
}