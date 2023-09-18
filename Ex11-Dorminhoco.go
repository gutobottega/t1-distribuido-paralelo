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
	"sync"
)

const NJ = 5 // numero de jogadores
const M = 4  // numero de cartas

var bateuChan = make(chan struct{}, NJ-1)
var ch [NJ]chan string // NJ canais de itens tipo carta

// Espera todos os jogadores baterem para finalizar o jogo
var wg sync.WaitGroup

func jogador(id int, in chan string, out chan string, cartasIniciais []string) {
	mao := cartasIniciais // estado local - as cartas na mao do jogador

	for {
		select {
		//case para receber carta
		case cartaRecebida := <-in:
			// adiciona a carta na mao e escolhe um aleatoria para passar para o proximo jogador
			mao = append(mao, cartaRecebida)
			remove := rand.Intn(len(mao))
			print("Jogador ", id, " recebeu a carta ", cartaRecebida, " e passou a carta ", mao[remove], "\n")
			carta := mao[remove]
			mao = append(mao[:remove], mao[remove+1:]...)
			//testa se todas as cartas na mao são iguais
			bate := true
			for i := 1; i < len(mao); i++ {
				if mao[i] != mao[0] {
					bate = false
					break
				}
			}
			quantidadeBateu := len(bateuChan)
			if quantidadeBateu > 0 {
				if quantidadeBateu == NJ-1 {
					// perdeu o jogo
					println("Jogador", id, "perdeu o jogo\n")
					wg.Done()
					return
				}
				bateuChan <- struct{}{}
				print("Jogador ", id, " bateu.\n")
				wg.Done()
				return
			}
			out <- carta
			if bate {
				print("Jogador ", id, " bateu com a mao: ")
				for _, carta := range mao {
					print(carta, " ")
				}
				print("\n")
				bateuChan <- struct{}{}
				wg.Done()
				return
			}
			print("Jogador ", id, " tem as cartas: ")
			for _, carta := range mao {
				print(carta, " ")
			}
			print("\n")
		default:
			quantidadeBateu := len(bateuChan)
			if quantidadeBateu > 0 {
				if quantidadeBateu == NJ-1 {
					// perdeu o jogo
					println("Jogador", id, "perdeu o jogo\n")
					wg.Done()
					return
				}
				bateuChan <- struct{}{}
				print("Jogador ", id, " bateu.\n")
				wg.Done()
				return
			}
		}
	}
}

func criaDeck() []string {
	deck := []string{}
	// cria um baralho com NJ tipo de cartas M vezes
	for i := 0; i < NJ; i++ {
		for j := 0; j < M; j++ {
			deck = append(deck, fmt.Sprintf("%d", i))
		}
	}
	// adiciona coringa no baralho
	deck = append(deck, "@")
	// embaralha o baralho
	for i := range deck {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
	return deck
}

func main() {
	// cria NJ canais de cartas
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan string)
	}
	// cria o deck
	deck := criaDeck()
	// cria NJ jogadores com M cartas cada
	for i := 0; i < NJ; i++ {
		deckIndex := i * M
		wg.Add(1)
		go jogador(i, ch[i], ch[(i+1)%NJ], deck[deckIndex:deckIndex+M]) // cria processos conectados circularmente
	}
	// Inicia o jogo mandando a ultima carta do deck para o primeiro jogador
	ch[0] <- deck[len(deck)-1]

	wg.Wait()
	print("Fim de jogo\n")
}
