package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	NCL           = 100
	Pool          = 10
	TotalRequests = 50 // Número total de requisições a serem processadas
)

type Request struct {
	v      int
	ch_ret chan int
}

func cliente(i int, req chan Request, wg *sync.WaitGroup, clientWG *sync.WaitGroup) {
	defer wg.Done()
	my_ch := make(chan int)
	requestsToMake := TotalRequests / NCL
	if i < TotalRequests%NCL {
		requestsToMake++
	}

	fmt.Printf("[CLIENTE %02d] Iniciando (%d requisições)\n", i, requestsToMake)

	for j := 0; j < requestsToMake; j++ {
		v := rand.Intn(1000)
		fmt.Printf("[CLIENTE %02d] Enviando requisição %d: valor=%d\n", i, j+1, v)
		req <- Request{v, my_ch}
		r := <-my_ch
		fmt.Printf("[CLIENTE %02d] Recebida resposta %d: %d*2=%d\n", i, j+1, v, r)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	fmt.Printf("[CLIENTE %02d] Terminou\n", i)
	clientWG.Done()
}

func trataReq(id int, req Request, sem chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		sem <- struct{}{}
		fmt.Printf("[WORKER %02d] Liberando vaga (vagas disponíveis: %d)\n", id, Pool-len(sem)+1)
	}()

	fmt.Printf("[WORKER %02d] Processando valor %d\n", id, req.v)
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond) // Simula trabalho
	req.ch_ret <- req.v * 2
	fmt.Printf("[WORKER %02d] Resposta enviada: %d\n", id, req.v*2)
}

func servidorLimitado(in chan Request, wg *sync.WaitGroup) {
	sem := make(chan struct{}, Pool)

	// Inicializa o semáforo
	fmt.Printf("[SERVER] Inicializando pool com %d workers\n", Pool)
	for i := 0; i < Pool; i++ {
		sem <- struct{}{}
	}

	var workerID int
	for req := range in {
		workerID++
		fmt.Printf("[SERVER] Requisição recebida (worker %02d)\n", workerID)

		// Espera por vaga no pool
		startWait := time.Now()
		<-sem
		waitTime := time.Since(startWait)
		if waitTime > 10*time.Millisecond {
			fmt.Printf("[SERVER] Worker %02d esperou %v por vaga\n", workerID, waitTime)
		}

		wg.Add(1)
		fmt.Printf("[SERVER] Iniciando worker %02d (vagas disponíveis: %d)\n",
			workerID, Pool-len(sem))
		go trataReq(workerID, req, sem, wg)
	}

	fmt.Println("[SERVER] Canal de requisições fechado. Encerrando...")
}

func main() {
	fmt.Println("[MAIN] Iniciando servidor com pool limitado")
	fmt.Printf("[MAIN] Configuração: Clientes=%d, Workers=%d, Requisições=%d\n",
		NCL, Pool, TotalRequests)

	serv_chan := make(chan Request)
	var wg sync.WaitGroup
	var clientWG sync.WaitGroup

	// Inicia o servidor
	go func() {
		servidorLimitado(serv_chan, &wg)
	}()

	// Inicia os clientes
	clientWG.Add(NCL)
	wg.Add(NCL)
	for i := 0; i < NCL; i++ {
		go cliente(i, serv_chan, &wg, &clientWG)
	}

	// Quando todos os clientes terminarem, fecha o canal
	go func() {
		clientWG.Wait()
		fmt.Println("[MAIN] Todos os clientes terminaram. Fechando canal...")
		close(serv_chan)
	}()

	startTime := time.Now()
	wg.Wait()
	executionTime := time.Since(startTime)

	fmt.Println("\n[MAIN] Todos os workers terminaram")
	fmt.Printf("[MAIN] Tempo total de execução: %v\n", executionTime)
	fmt.Println("[MAIN] Programa encerrado com sucesso")
}
