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
// soma sequencial recursiva
func soma(r *Nodo) int {
    if r != nil {
        return r.v + soma(r.e) + soma(r.d)
    }
    return 0
}

// funcao "wraper" retorna valor
// internamente dispara recursao com somaConcCh
// usando canais
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

// Busca realiza a busca sequencial em uma árvore binária
func busca(r *Nodo, v int) bool {
    if r == nil {
        return false
    }

    if r.v == v {
        return true
    }

    return busca(r.e, v) || busca(r.d, v)
}

// BuscaConc realiza a busca concorrente em uma árvore binária
func buscaConc(r *Nodo, v int, encontrado chan bool) {
    if r == nil {
        encontrado <- false
        return
    }

    if r.v == v {
        encontrado <- true
        return
    }

    ch1 := make(chan bool)
    ch2 := make(chan bool)

    go func() {
        buscaConc(r.e, v, ch1)
    }()
    go func() {
        buscaConc(r.d, v, ch2)
    }()

    encontradoE := <-ch1
    encontradoD := <-ch2

    if encontradoE || encontradoD {
        encontrado <- true
    } else {
        encontrado <- false
    }
}

func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int) {
    defer close(saidaP)
    defer close(saidaI)

    var wg sync.WaitGroup

    var processa func(*Nodo)
    processa = func(r *Nodo) {
        defer wg.Done()
        if r != nil {
            if r.v%2 == 0 {
                saidaP <- r.v
            } else {
                saidaI <- r.v
            }
            wg.Add(1)
            go processa(r.e)
            wg.Add(1)
            go processa(r.d)
        }
    }

    wg.Add(1)
    go processa(r)

    wg.Wait()
}

func retornaParImparConc(r *Nodo, saidaP chan int, saidaI chan int) {
    var wg sync.WaitGroup

    var processaConc func(*Nodo)
    processaConc = func(r *Nodo) {
        defer wg.Done()
        if r != nil {
            if r.v%2 == 0 {
                saidaP <- r.v
            } else {
                saidaI <- r.v
            }
            wg.Add(1)
            go processaConc(r.e)
            wg.Add(1)
            go processaConc(r.d)
        }
    }

    wg.Add(1)
    go processaConc(r)

    wg.Wait()
}


func main() {
    root := &Nodo{v: 10,
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

    fmt.Println()
    fmt.Print("Valores na árvore: ")
    caminhaERD(root)
    fmt.Println()
    fmt.Println()

    fmt.Println("Soma: ", soma(root))
    fmt.Println()

    elemento := 14
    encontradoSeq := busca(root, elemento)
    fmt.Printf("Busca: %v\n", encontradoSeq)

    encontrado := make(chan bool)
    go buscaConc(root, elemento, encontrado)
    encontradoConc := <-encontrado
    fmt.Printf("BuscaConc: %v\n", encontradoConc)
    
    saidaP := make(chan int)
    saidaI := make(chan int)

    go retornaParImpar(root, saidaP, saidaI)
    go retornaParImparConc(root, saidaP, saidaI)

    // Leia os resultados dos canais de saída
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        for num := range saidaP {
            fmt.Println("Par:", num)
        }
    }()

    go func() {
        defer wg.Done()
        for num := range saidaI {
            fmt.Println("Ímpar:", num)
        }
    }()

    wg.Wait()
}