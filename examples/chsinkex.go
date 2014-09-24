package main

import (
	"fmt"
	"time"
)

func bitGenerator(size int, CH chan int) {
	for i := 0; i < size; i++ {
		fmt.Printf("\nGEnerated %d", i)
		CH <- i
	}

}

func dumpDoubleValue(maxsize int, CH chan int) {
	for i := 0; i < maxsize; i++ {
		data := <-CH
		fmt.Printf("\n\t \t %d  *** DoubleValue :fn(%d)=%d", i, data, 2*data)
	}
}

func dumpSquareValue(maxsize int, CH chan int) {
	for i := 0; i < maxsize; i++ {
		data := <-CH
		fmt.Printf("\n\t \t%d === SquareValue: fn(%d)=%d", i, data, data*data)
	}
	// close(CH)
}

func main() {
	CH := make(chan int, 1)
	go bitGenerator(10, CH)
	go dumpDoubleValue(10, CH)
	go dumpSquareValue(10, CH)
	time.Sleep(1000 * time.Millisecond)
	close(CH)
}
