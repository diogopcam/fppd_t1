// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// ------- Voce deve ter feito o Ex3 desta serie para entender e continuar aqui ------
// Aqui está a resposta do Ex3.
// PROBLEMA: Suponha agora que cada nodo que recebe uma mensagem deva mandar uma resposta ao origem
//           usando o mesmo protocolo de inundacao.
//           Este sistema suporta ?
// SOLUCAO:  a solucao foi: assim que a resposta chegou ao destino, fazer um broadcast para o origem.
//           o mesmo sistema funciona.
// ATENCAO: por simplicidade, adotei que a de resposta tem o mesmo identificador da recebida, só que negativo
//
// EXERCÍCIO:
//          como proximo passo, implemente que durante a inundação de ida a rota vai sendo
//          gravada.  A rota é a sequencia de nodos por onde a mensagem passa.
//          Pode ser uma pilha de inteiros.  Cada nodo antes de repassar, empilha seu id.
//          Desta forma, a resposta pode ser enviada somente pela rota de retorno.
//          Ou seja, a mensagem trafega pela rota reversa.  Basta que cada nodo intermediario
//          desempilhe o identificador do proximo e repasse a mensagem para este.


// Enfrentamos problemas com deadlock! (todos canais bloqueados)
package main

import (
	"fmt"
	"time"
)

const N = 10
const channelBufferSize = 5

type Topology [N][N]int

type Message struct {
	id       int
	source   int
	receiver int
	data     string
	route    []int  // Pilha para armazenar a rota
}

type inputChan [N]chan Message

type nodeStruct struct {
	id               int
	topo             Topology
	inCh             inputChan
	received         map[int]bool
	receivedMessages []Message
}

func (n *nodeStruct) broadcast(m Message) {
	for j := 0; j < N; j++ {
		if n.topo[n.id][j] == 1 {
			n.inCh[j] <- m
		}
	}
}

func (n *nodeStruct) sendToNext(m Message, next int) {
	n.inCh[next] <- m
}

func (n *nodeStruct) nodo() {
	fmt.Printf("Nó %d iniciado\n", n.id)
	for {
		select {
		case m := <-n.inCh[n.id]:
			if m.receiver == n.id {
				n.receivedMessages = append(n.receivedMessages, m)
				if m.id > 0 { // Mensagem de ida
					fmt.Printf("Nó %d (destino) recebeu mensagem %d de %d. Rota: %v\n", 
						n.id, m.id, m.source, m.route)
					
					// Prepara resposta com rota reversa
					resp := Message{
						id:       -m.id,
						source:   n.id,
						receiver: m.source,
						data:     "resp",
						route:    m.route, // Usa a mesma rota (será desempilhada)
					}
					
					if len(resp.route) > 0 {
						next := resp.route[len(resp.route)-1]
						resp.route = resp.route[:len(resp.route)-1] // Desempilha
						fmt.Printf("Nó %d enviando resposta para nó %d. Rota restante: %v\n",
							n.id, next, resp.route)
						go n.sendToNext(resp, next)
					}
				} else { // Mensagem de resposta
					fmt.Printf("Nó %d recebeu resposta %d de %d\n", 
						n.id, m.id, m.source)
				}
			} else if !n.received[m.id] {
				n.received[m.id] = true
				
				// Adiciona seu ID à rota antes de repassar
				newRoute := make([]int, len(m.route))
				copy(newRoute, m.route)
				newRoute = append(newRoute, n.id)
				m.route = newRoute
				
				fmt.Printf("Nó %d repassando mensagem %d. Rota atual: %v\n",
					n.id, m.id, m.route)
				go n.broadcast(m)
			}

		case <-time.After(2 * time.Second):
			fmt.Printf("Nó %d: tempo limite atingido\n", n.id)
			return
		}
	}
}

func carga(nodoInicial int, inCh chan Message, numMensagens int) {
	for i := 1; i <= numMensagens; i++ {
		msg := Message{
			id:       nodoInicial*1000 + i,
			source:   nodoInicial,
			receiver: (nodoInicial + i) % N,
			data:     fmt.Sprintf("req%d", i),
			route:    []int{}, // Rota inicial vazia
		}
		inCh <- msg
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	topo := Topology{
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 0 ↔ 1
		{1, 0, 1, 0, 0, 0, 0, 0, 0, 0}, // 1 ↔ 2
		{0, 1, 0, 1, 0, 0, 0, 0, 0, 0}, // 2 ↔ 3
		{0, 0, 1, 0, 1, 0, 0, 0, 0, 0}, // 3 ↔ 4
		{0, 0, 0, 1, 0, 1, 0, 0, 0, 1}, // 4 ↔ 5, 4 ↔ 9
		{0, 0, 0, 0, 1, 0, 1, 0, 0, 0}, // 5 ↔ 6
		{0, 0, 0, 0, 0, 1, 0, 1, 0, 0}, // 6 ↔ 7
		{0, 0, 0, 0, 0, 0, 1, 0, 1, 0}, // 7 ↔ 8
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 1}, // 8 ↔ 9
		{0, 0, 0, 0, 1, 0, 0, 0, 1, 0}, // 9 ↔ 4
	}

	var inCh inputChan
	for i := 0; i < N; i++ {
		inCh[i] = make(chan Message, channelBufferSize)
	}

	for id := 0; id < N; id++ {
		n := nodeStruct{
			id:       id,
			topo:     topo,
			inCh:     inCh,
			received: make(map[int]bool),
		}
		go n.nodo()
	}

	go carga(0, inCh[0], 2) // Nó 0 envia 2 mensagens
	go carga(3, inCh[3], 1) // Nó 3 envia 1 mensagem

	time.Sleep(5 * time.Second)
	fmt.Println("Simulação concluída")
}