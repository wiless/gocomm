package main

import (
	"fmt"
	//	"time"
	"math/rand"
)

var gcounter int

func genBit(blockSize int, cs chan int) {
	for i := 0; i < blockSize; i++ {
		//temp:=rand.Intn(100)
		gcounter++
		fmt.Printf("\nTx(%d) := %d \t ", gcounter, gcounter)
		cs <- gcounter //send a strawberry cake
		fmt.Println(".. exiting  ...", gcounter)
	}
	fmt.Printf("Transmitter closing..")

}

func genSymbol(N int, cs chan int, csout chan []int) {
	var symbol [3]int
	//	N:=10
	for j := 0; j < N; j++ {
		for i := 0; i < 3; i++ {

			fmt.Printf("\n %d Reading ..%d", j, i)
			symbol[i] = <-cs //get whatever cake is on the channel
			//gcounter++

		}

		fmt.Printf("\n Sending Symbol %d %v", j, symbol)
		csout <- symbol[:]
		fmt.Printf("\n Ready to Read Next Pair %d", j)

	}

	fmt.Printf("Symbol Generator  closing..")

}

func genFrame(cs chan []int) {
	var symbol [4]int
	cnt := 0
	indx := 0
	for cnt <= len(symbol) {
		tmp := <-cs
		copy(symbol[0:len(tmp)], tmp)
		fmt.Printf("Reading Pair : %d %v ", indx, tmp)
		cnt = cnt + len(tmp)
	}

	result := 0
	for i := 0; i < 4; i++ {
		result += symbol[i]
	}
	fmt.Println("\nCum SUM : ", result)
}

func main() {
	fmt.Print(rand.Int())
	cs := make(chan int)
	chout := make(chan []int)
	N := 30

	go genBit(N, cs)
	go genSymbol(10, cs, chout)

	for i := 0; i < 10; i++ {
		fmt.Printf("\n %d Reading Final Pair %v ", i, <-chout)
	}
	close(cs)
	close(chout)
	//time.Sleep(100 * 1e9)

}
