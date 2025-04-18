// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// PROBLEMA:
//   o dorminhoco especificado no arquivo Ex1-ExplanacaoDoDorminhoco.pdf nesta pasta
// ESTE ARQUIVO
//   Um template para criar um anel generico.
//   Adapte para o problema do dorminhoco.
//   Nada está dito sobre como funciona a ordem de processos que batem.
//   O ultimo leva a rolhada ...
//   ESTE  PROGRAMA NAO FUNCIONA.    É UM RASCUNHO COM DICAS.

package main

import (
	"fmt"
	"time"
)

const NJ = 5
const M = 4
const channelBufferSize = 2 // Definido para corrigir o erro

type carta string

var ch [NJ]chan carta
var bateu [NJ]chan bool
var ordemBatida chan int

func jogador(id int, in chan carta, out chan carta, bateuChan chan bool, cartasIniciais []carta) {
	mao := make([]carta, len(cartasIniciais))
	copy(mao, cartasIniciais)

	for {
		select {
		case <-bateuChan:
			ordemBatida <- id
			return
		default:
			if len(mao) > M {
				// Escolhe uma carta aleatória para passar
				idx := 0 // Simplificação: sempre passa a primeira carta
				cartaParaSair := mao[idx]
				mao = append(mao[:idx], mao[idx+1:]...)

				// Verifica se pode bater
				if podeBater(mao) {
					fmt.Printf("Jogador %d bateu!\n", id)
					for i := 0; i < NJ; i++ {
						if i != id {
							bateu[i] <- true
						}
					}
					ordemBatida <- id
					return
				}

				out <- cartaParaSair
				fmt.Printf("Jogador %d passou carta %s\n", id, cartaParaSair)
			} else {
				cartaRecebida := <-in
				mao = append(mao, cartaRecebida)
				fmt.Printf("Jogador %d recebeu carta %s\n", id, cartaRecebida)

				if podeBater(mao) {
					fmt.Printf("Jogador %d bateu!\n", id)
					for i := 0; i < NJ; i++ {
						if i != id {
							bateu[i] <- true
						}
					}
					ordemBatida <- id
					return
				}
			}
		}
	}
}

func podeBater(mao []carta) bool {
	if len(mao) != M {
		return false
	}
	for i := 1; i < len(mao); i++ {
		if mao[i] != mao[0] {
			return false
		}
	}
	return true
}

func main() {
	// Inicializa canais
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan carta, channelBufferSize)
		bateu[i] = make(chan bool, 1)
	}
	ordemBatida = make(chan int, NJ)

	// Cria e distribui cartas
	baralho := []carta{"A", "A", "A", "A", "B", "B", "B", "B", "C", "C", "C", "C", "D", "D", "D", "D", "E", "E", "E", "E", "@"}

	// Distribui cartas (simplificado)
	cartasPorJogador := make([][]carta, NJ)
	for i := 0; i < NJ; i++ {
		cartasPorJogador[i] = make([]carta, M)
		for j := 0; j < M; j++ {
			cartasPorJogador[i][j] = baralho[i*M+j]
		}
	}
	// Jogador 0 recebe carta extra
	cartasPorJogador[0] = append(cartasPorJogador[0], "@")

	// Inicia jogadores
	for i := 0; i < NJ; i++ {
		go jogador(i, ch[i], ch[(i+1)%NJ], bateu[i], cartasPorJogador[i])
	}

	// Coleta ordem de batida
	go func() {
		for i := 0; i < NJ; i++ {
			id := <-ordemBatida
			fmt.Printf("Jogador %d bateu em %dº lugar\n", id, i+1)
			if i == NJ-1 {
				fmt.Printf("Jogador %d é o DORMINHOCO!\n", id)
				return
			}
		}
	}()

	// Tempo para execução
	time.Sleep(5 * time.Second)
	fmt.Println("Fim do jogo")
}