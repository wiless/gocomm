package main

import (
	"fmt"
	"wiless/gocomm"
)

func main() {

	bpskModem := new(gocomm.Modem)
	bpskModem.Init(1)
	fmt.Printf("\n%v", bpskModem)
	fmt.Printf("\n%f", bpskModem.Constellation)

	print("\n") // Legacy

	qpskModem := new(gocomm.Modem)
	qpskModem.Init(2)
	fmt.Printf("\n%v", qpskModem)
	fmt.Printf("\n%f", qpskModem.Constellation)

}
