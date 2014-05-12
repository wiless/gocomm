package main

import (
	"fmt"
	"math"
)

type Modem struct {
	bitsPerSymbol int
	name          string
	offset        float64
	constellation []complex128
}

func (m *Modem) modulateBits(bits []int) []complex128 {
	length := len(bits)

	slength := length / m.bitsPerSymbol
	var symbols = make([]complex128, slength)
	return symbols
}

func (m *Modem) init(wordlength int) {
	m.bitsPerSymbol = wordlength
	switch wordlength {
	case 1:
		m.name = "BPSK"
	case 2:
		m.name = "QPSK"
		m.offset = math.Pi / 4.0
	case 3:
		m.name = "8PSK"
	case 4:
		m.name = "16PSK"
	}
	var i float64 = 0
	var length int = int(math.Exp2(float64(m.bitsPerSymbol)))
	m.constellation = make([]complex128, length)

	for i = 0; i < float64(len(m.constellation)); i++ {
		var angle = i * 2 * math.Pi / float64(length)
		m.constellation[int(i)] = complex(math.Cos(angle+m.offset), math.Sin(angle+m.offset))

	}
}

func (m Modem) String() string {
	return fmt.Sprintf("A %s Modem with %d bits/Symbol", m.name, m.bitsPerSymbol)
}

func shift(data []complex128, offset complex128) []complex128 {
	var result = make([]complex128, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = (0.707/2.0)*data[i] + offset
	}
	return result
}

func main() {

	bpskModem := new(Modem)
	bpskModem.init(1)
	fmt.Printf("\n%v", bpskModem)
	fmt.Printf("\n%f", bpskModem.constellation)

	qpskModem := new(Modem)
	qpskModem.init(2)
	fmt.Printf("\n%v", qpskModem)
	fmt.Printf("\n%f", qpskModem.constellation)

	qpsk16Modem := new(Modem)
	qpsk16Modem.init(4)
	fmt.Printf("\n%v", qpsk16Modem)
	fmt.Printf("\n%f", qpsk16Modem.constellation)

}
