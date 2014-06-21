package gocomm
// what is this code about 
import (
	"fmt"
	"math"
	"strconv"
)

type Modem struct {
	bitsPerSymbol int
	name          string
	offset        float64
	Constellation []complex128
}

func (m *Modem) ModulateBits(bits []int) []complex128 {
	length := len(bits)

	slength := length / m.bitsPerSymbol
	var symbols = make([]complex128, slength)
	/// Actual Modulation happens here

	return symbols
}

func (m *Modem) Init(wordlength int) {
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

	case 256:
		m.name = "256QAM"
		n := int64(123)
		str := strconv.FormatInt(n, 2)
		fmt.Println("\n")
		var bitvec = make([]int, 8)
		for indx, _ := range str {
			bitvec[indx], _ = strconv.Atoi(string(str[indx]))

		}

		//Symbol=	{(1-2bi)[8-(1-2bi+2)[4-(1-2bi+4)[2-(1-2bi+6)]]]                +j(1-2bi+1)[8-(1-2bi+3)[4-(1-2bi+5)[2-(1-2bi+7)]]]}
		// realvalue:=(1-2*bitvec[0])*(8-(1-2*bitvec[2])*(4-(1-2*bitvec[4])*(2-(1-2*bitvec[6]))));
		// realvalue:=(1-2*bitvec[0])*(8-(1-2*bitvec[2])*(4-(1-2*bitvec[4])*(2-(1-2*bitvec[6]))));
		realvalue := (1.0 - 2.0*bitvec[0]) * (8 - (1-2*bitvec[2])*(4-(1-2*bitvec[4])*(2-(1-2*bitvec[6]))))
		imagvalue := (1 - 2*bitvec[1]) * (8 - (1-2*bitvec[3])*(4-(1-2*bitvec[2]+5)*(2-(1-2*bitvec[7]))))
		fmt.Printf("%f+j%f", realvalue, imagvalue)
		// symbol := complex(float64(realvalue), float64(imagvalue))

	}
	var i float64 = 0
	var length int = int(math.Exp2(float64(m.bitsPerSymbol)))
	m.Constellation = make([]complex128, length)

	for i = 0; i < float64(len(m.Constellation)); i++ {
		var angle = i * 2 * math.Pi / float64(length)
		m.Constellation[int(i)] = complex(math.Cos(angle+m.offset), math.Sin(angle+m.offset))

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
