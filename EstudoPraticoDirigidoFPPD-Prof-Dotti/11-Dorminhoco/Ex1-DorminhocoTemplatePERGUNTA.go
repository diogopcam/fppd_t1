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
)

const NJ = 5           // numero de jogadores
const M = 4            // numero de cartas

type carta string      // carta é um strirng

var ch [NJ]chan carta  // NJ canais de itens tipo carta  

func jogador(id int, in chan carta, out chan carta, cartasIniciais []carta, ... ) {
	mao := cartasIniciais    // estado local - as cartas na mao do jogador
	nroDeCartas := M         // quantas cartas ele tem 
    cartaRecebida := " "     // carta recebida é vazia

	for {
		if // tenho que jogar
			{
			fmt.Println(id, " joga") // escreve seu identificador
			// escolhe alguma carta para passar adiante ...
			// guarda carta que entrou 
			// manda carta escolhida o proximo
			out <- cartaParaSair        
		} else {  
			// ...
			cartaRecebida := <-in   // recebe carta na entrada
			//  ...
			// e se alguem bate ?
		}
	}
}

func main() {
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan struct{})
	}
	// cria um baralho com NJ*M cartas
	for i := 0; i < NJ; i++ {
		// escolhe aleatoriamente (tira) cartas do baralho, passa cartas para jogador
		go jogador(i, ch[i], ch[(i+1)%N], cartasEscolhidas , ...) // cria processos conectados circularmente
	}
	
	<-make(chan struct{}) // bloqueia
}


