// Nomes: Diogo Pessin Camargo e Giovanni Schardong Pereira
// por Fernando Dotti - PUCRS
// dado abaixo um exemplo de estrutura em arvore, uma arvore inicializada
// e uma operação de caminhamento, pede-se fazer:
//   1.a) a operação que soma todos elementos da arvore.
//        func soma(r *Nodo) int {...}
//   1.b) uma operação concorrente que soma todos elementos da arvore
//   2.a) a operação de busca de um elemento v, dizendo true se encontrou v na árvore, ou falso
//        func busca(r* Nodo, v int) bool {}...}
//   2.b) a operação de busca concorrente de um elemento, que informa imediatamente
//        por um canal se encontrou o elemento (sem acabar a busca), ou informa
//        que nao encontrou ao final da busca
//   3.a) a operação que escreve todos pares em um canal de saidaPares e
//        todos impares em um canal saidaImpares, e ao final avisa que acabou em um canal fin
//        func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}){...}
//   3.b) a versao concorrente da operação acima, ou seja, os varios nodos sao testados
//        concorrentemente se pares ou impares, escrevendo o valor no canal adequado
//
//  ABAIXO: RESPOSTAS A QUESTOES 1a e b
//  APRESENTE A SOLUÇÃO PARA AS DEMAIS QUESTÕES

package main

import (
	"fmt"
	"sync"
)

type Nodo struct {
	v int
	e *Nodo
	d *Nodo
}

func caminhaERD(r *Nodo) {
	if r != nil {
		caminhaERD(r.e)
		fmt.Print(r.v, ", ")
		caminhaERD(r.d)
	}
}

// -------- SOMA ----------
func soma(r *Nodo) int {
	if r != nil {
		return r.v + soma(r.e) + soma(r.d)
	}
	return 0
}

func somaConc(r *Nodo) int {
	s := make(chan int)
	go somaConcCh(r, s)
	return <-s
}

func somaConcCh(r *Nodo, s chan int) {
	if r != nil {
		s1 := make(chan int)
		go somaConcCh(r.e, s1)
		go somaConcCh(r.d, s1)
		s <- (r.v + <-s1 + <-s1)
	} else {
		s <- 0
	}
}

// -------- BUSCA ----------
func busca(r *Nodo, val int) bool {
	if r == nil {
		return false
	}
	if r.v == val {
		return true
	}
	return busca(r.e, val) || busca(r.d, val)
}

func buscaC(r *Nodo, val int) bool {
	ret := make(chan bool)
	go buscaConc(r, val, ret)
	return <-ret
}

func buscaConc(r *Nodo, val int, ret chan bool) {
	if r == nil {
		ret <- false
		return
	}
	if r.v == val {
		ret <- true
		return
	}
	ch := make(chan bool)
	go buscaConc(r.e, val, ch)
	go buscaConc(r.d, val, ch)
	ret <- <-ch || <-ch
}

// -------- SAIDAS PAR E IMPAR --------
func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) {
	if r == nil {
		return
	}

	if r.v%2 == 0 {
		saidaP <- r.v
	} else {
		saidaI <- r.v
	}

	retornaParImpar(r.e, saidaP, saidaI, fin)
	retornaParImpar(r.d, saidaP, saidaI, fin)

	// Só envia sinal de fim quando chegar na raiz novamente
	if r == root {
		fin <- struct{}{}
	}
}

func retornaParImparConc(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) {
	if r == nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if r.v%2 == 0 {
			saidaP <- r.v
		} else {
			saidaI <- r.v
		}
	}()

	go func() {
		defer wg.Done()
		retornaParImparConc(r.e, saidaP, saidaI, fin)
		retornaParImparConc(r.d, saidaP, saidaI, fin)
	}()

	wg.Wait()

	if r == root {
		fin <- struct{}{}
	}
}

var root *Nodo

func main() {
	root = &Nodo{v: 10,
		e: &Nodo{v: 5,
			e: &Nodo{v: 3,
				e: &Nodo{v: 1, e: nil, d: nil},
				d: &Nodo{v: 4, e: nil, d: nil}},
			d: &Nodo{v: 7,
				e: &Nodo{v: 6, e: nil, d: nil},
				d: &Nodo{v: 8, e: nil, d: nil}}},
		d: &Nodo{v: 15,
			e: &Nodo{v: 13,
				e: &Nodo{v: 12, e: nil, d: nil},
				d: &Nodo{v: 14, e: nil, d: nil}},
			d: &Nodo{v: 18,
				e: &Nodo{v: 17, e: nil, d: nil},
				d: &Nodo{v: 19, e: nil, d: nil}}}}

	// Teste a versão sequencial
	testarOperacoes(false)

	// Teste a versão concorrente
	testarOperacoes(true)
}

func testarOperacoes(concorrente bool) {
	saidaP := make(chan int, 15)
	saidaI := make(chan int, 15)
	fin := make(chan struct{})

	fmt.Println("\n=== TESTE", map[bool]string{true: "CONCORRENTE", false: "SEQUENCIAL"}[concorrente], "===")
	fmt.Println("Valores na árvore:")

	if concorrente {
		go retornaParImparConc(root, saidaP, saidaI, fin)
	} else {
		go retornaParImpar(root, saidaP, saidaI, fin)
	}

	fim := false
	for !fim {
		select {
		case par := <-saidaP:
			fmt.Println("Par:", par)
		case impar := <-saidaI:
			fmt.Println("Impar:", impar)
		case <-fin:
			fim = true
		}
	}

	fmt.Println("\nValores ordenados:")
	caminhaERD(root)
	fmt.Println("\n")

	fmt.Println("Soma total:", soma(root))
	fmt.Println("Soma concorrente:", somaConc(root))
	fmt.Println()

	fmt.Println("Busca por 17:", busca(root, 17))
	fmt.Println("Busca por 99:", busca(root, 99))
	fmt.Println()

	fmt.Println("Busca concorrente por 17:", buscaC(root, 17))
	fmt.Println("Busca concorrente por 99:", buscaC(root, 99))
	fmt.Println()
}
