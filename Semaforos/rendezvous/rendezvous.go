package main

import (
	"fmt"
	"sync"
)

var (
	aDone = make(chan struct{}, 1)
	bDone = make(chan struct{}, 1)
	wg    sync.WaitGroup
)

func threadA() {
	defer wg.Done()
	fmt.Println("A: Parte 1")

	// Sinaliza que A terminou a parte 1
	aDone <- struct{}{}

	// Espera B terminar a parte 1
	<-bDone

	fmt.Println("A: Parte 2")
}

func threadB() {
	defer wg.Done()
	fmt.Println("B: Parte 1")

	// Sinaliza que B terminou a parte 1
	bDone <- struct{}{}

	// Espera A terminar a parte 1
	<-aDone

	fmt.Println("B: Parte 2")
}

func main() {
	wg.Add(2)
	go threadA()
	go threadB()
	wg.Wait()
}