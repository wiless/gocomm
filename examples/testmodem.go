package main

import (
	"fmt"
	"github.com/wiless/gocomm/modem"
	"math/rand"

	// "strconv"

	"github.com/wiless/gocomm/sources"
	"github.com/wiless/vlib"

	"time"
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

	fmt.Print("\nRandNoise  = ", sources.RandNCVec(N, 1))
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

	bpskModem := modem.NewModem(1, "BPSK")
	// bpskModem.Init(1)
	fmt.Printf("\n%v", bpskModem)
	fmt.Printf("\n%f", bpskModem.Constellation)

	print("\n") // Legacy
	// qpskModem := modem.NewModem(2, "QPSK")

	qpskModem := new(modem.Modem)
	qpskModem.Init(2, "QPSK")
	fmt.Printf("\n%v", qpskModem)
	fmt.Printf("\n%f", qpskModem.Constellation)
	qpskModem.ModulateBits(bits2)
}
