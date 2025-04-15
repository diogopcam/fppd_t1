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
	"math/rand"
	"time"
)

const NJ = 5 // número de jogadores
const M = 4  // número de cartas por jogador

type carta string // carta é um string

var ch [NJ]chan carta    // canais de comunicação entre jogadores
var bateu [NJ]chan bool  // canais para sinalizar que bateu
var ordemBatida chan int // canal para registrar ordem de batida

func jogador(id int, in chan carta, out chan carta, bateuChan chan bool, cartasIniciais []carta) {
	mao := make([]carta, len(cartasIniciais))
	copy(mao, cartasIniciais)

	for {
		// Verifica se alguém já bateu
		select {
		case <-bateuChan:
			// Alguém bateu, então eu bato também
			ordemBatida <- id
			return
		default:
			// Ninguém bateu ainda, continua o jogo
		}

		if len(mao) > M {
			// Tem carta extra, deve jogar
			// Escolhe uma carta aleatória para passar
			idx := rand.Intn(len(mao))
			cartaParaSair := mao[idx]

			// Remove a carta da mão
			mao = append(mao[:idx], mao[idx+1:]...)

			// Verifica se pode bater (tem M cartas iguais)
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

			// Passa a carta para o próximo
			out <- cartaParaSair
			fmt.Printf("Jogador %d passou carta %s\n", id, cartaParaSair)
		} else {
			// Espera receber uma carta
			cartaRecebida := <-in
			fmt.Printf("Jogador %d recebeu carta %s\n", id, cartaRecebida)
			mao = append(mao, cartaRecebida)

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
		}
	}
}

func podeBater(mao []carta) bool {
	if len(mao) != M {
		return false
	}
	// Verifica se todas as cartas são iguais
	for i := 1; i < len(mao); i++ {
		if mao[i] != mao[0] {
			return false
		}
	}
	return true
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Inicializa canais
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan carta)
		bateu[i] = make(chan bool, 1) // Buffer para evitar deadlock
	}
	ordemBatida = make(chan int, NJ) // Buffer para todas as batidas

	// Cria baralho
	baralho := make([]carta, 0, NJ*M+1)
	for i := 0; i < NJ; i++ {
		c := carta(rune('A' + i))
		for j := 0; j < M; j++ {
			baralho = append(baralho, c)
		}
	}
	baralho = append(baralho, "@") // Joker

	// Embaralha
	rand.Shuffle(len(baralho), func(i, j int) {
		baralho[i], baralho[j] = baralho[j], baralho[i]
	})

	// Distribui cartas
	cartasPorJogador := make([][]carta, NJ)
	idx := 0
	for i := 0; i < NJ; i++ {
		cartasPorJogador[i] = make([]carta, M)
		for j := 0; j < M; j++ {
			cartasPorJogador[i][j] = baralho[idx]
			idx++
		}
	}

	// Um jogador (aleatório) recebe a carta extra (joker)
	jogadorInicial := rand.Intn(NJ)
	cartasPorJogador[jogadorInicial] = append(cartasPorJogador[jogadorInicial], baralho[idx])
	fmt.Printf("Jogador %d começa com carta extra\n", jogadorInicial)

	// Inicia jogadores
	for i := 0; i < NJ; i++ {
		go jogador(i, ch[i], ch[(i+1)%NJ], bateu[i], cartasPorJogador[i])
	}

	// Coleta ordem de batida
	ordem := make([]int, 0, NJ)
	for i := 0; i < NJ; i++ {
		id := <-ordemBatida
		ordem = append(ordem, id)
		fmt.Printf("Jogador %d bateu em %dº lugar\n", id, len(ordem))
	}

	// O último a bater é o dorminhoco
	dorminhoco := ordem[len(ordem)-1]
	fmt.Printf("Jogador %d é o dorminhoco!\n", dorminhoco)
}
