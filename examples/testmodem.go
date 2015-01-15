package main

import (
	"fmt"
	"math/rand"

	// "strconv"
	"time"
	"github.com/wiless/gocomm"
	"github.com/wiless/gocomm/sources"
	"github.com/wiless/vlib"
)

func main() {

	rand.Seed(time.Now().Unix())
	rand.Seed(1)

	// for i := 0; i < 256; i++ {
	// 	fmt.Printf("\n %d %c ", i, i)
	// }

	var N = 10
	bits2 := vlib.NewOnesB(10)

	fmt.Print("\nRandBits  = ", sources.Randsrc(N, 2))
	fmt.Print("\nRandNoise  = ", sources.Noise(N, 1))
	// msg := sources.RandChars(N)
	msg := sources.RandReadableChars(N)
	strmsg := string(msg)
	// strmsg = "Hello world, this is a message sent"

	// var msg = "HELLO SENDIL"

	fmt.Printf("\nRandChars = %s ", strmsg)
	fmt.Print("\nBit messages  = ", sources.BitsFromMessage(strmsg))
	bits3 := sources.BitsFromMessage(strmsg)
	fmt.Printf("\nBits  : %v", bits2)
	fmt.Printf("\nBits  : %v ", bits3)

	bpskModem := new(gocomm.Modem)
	bpskModem.Init(1)
	fmt.Printf("\n%v", bpskModem)
	fmt.Printf("\n%f", bpskModem.Constellation)

	print("\n") // Legacy

	qpskModem := new(gocomm.Modem)
	qpskModem.Init(2)
	fmt.Printf("\n%v", qpskModem)
	fmt.Printf("\n%f", qpskModem.Constellation)
	qpskModem.ModulateBits(bits2)
}
