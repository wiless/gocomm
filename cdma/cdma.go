package SS

import (

	// "fmt"
	// "log"
	"math/cmplx"

	"github.com/wiless/gocomm"
	"github.com/wiless/vlib"
)

type CDMA struct {
	SpreadSequence []complex128
}

func NewCDMA() (cdma *CDMA) {
	cdma = new(CDMA)
	cdma.SpreadSequence = []complex128{1, 1, 1, 1, -1, -1, -1, -1}
	return cdma
	// cdma.SpreadSequence = []complex128{1, 1, 1, 1, -1, -1, -1, -1}
}

func (c *CDMA) GetSpreadOutputBlockSize(N int) int {
	return len(c.SpreadSequence) * N
}
func (c *CDMA) GetDeSpreadOutputBlockSize(N int) int {
	return N / len(c.SpreadSequence)
}

func (cdma *CDMA) DeSpreadBlock(expectedInputSize int, chInway gocomm.Complex128AChannel, OutCH gocomm.Complex128Channel) {

	despcode := vlib.Conj(cdma.SpreadSequence)

	SF := len(despcode)

	despcode = despcode.Scale(1. / (float64(SF)))

	if SF == 0 {
		panic("Spreading Code not Set")
	}
	// maxSymbols := expectedInputSize / SF
	// rxsymbols := vlib.NewVectorC(maxSymbols)
	var recentBuffer vlib.VectorC
	for cnt := 0; cnt < expectedInputSize; {

		data := <-chInway
		rxlen := len(data.Ch)
		// log.Printf("\n Received %d samples out of %d/%d ", rxlen, cnt, expectedInputSize)
		cnt += rxlen
		recentBuffer = append(recentBuffer, data.Ch...)
		for {
			if recentBuffer.Size() < SF {

				break

			} else {
				// log.Printf("\n Symbol %d Ready to Despread with %d", sym, cnt)
				rxchips := recentBuffer[0:SF]
				recentBuffer = recentBuffer[SF:]
				rxsymbols := vlib.DotC(despcode, rxchips)
				var chdataOut gocomm.SComplex128Obj
				chdataOut.Ch = rxsymbols
				OutCH <- chdataOut

			}
		}
	}

	close(chInway)
}

func (cdma *CDMA) DeSpread(InCH gocomm.Complex128Channel, OutCH gocomm.Complex128Channel) {

	despcode := vlib.Conj(cdma.SpreadSequence)
	var result complex128
	for i := 0; i < len(despcode); i++ {
		rxchips := (<-InCH).Ch
		result += rxchips * cmplx.Conj(despcode[i])
	}
	var chdataOut gocomm.SComplex128Obj
	chdataOut.Ch = result
	OutCH <- chdataOut
}

func (cdma *CDMA) SpreadBlock(expectedInputSymbols int, chInway gocomm.Complex128Channel, OutCH gocomm.Complex128AChannel) {

	spcode := cdma.SpreadSequence
	if len(spcode) == 0 {
		panic("Spreading Code not Set")
	}
	for i := 0; i < expectedInputSymbols; i++ {

		indata := <-chInway
		insymbol := indata.Ch
		var result = make([]complex128, len(spcode))
		for i := 0; i < len(spcode); i++ {
			result[i] = insymbol * spcode[i]
		}
		var chdataOut gocomm.SComplex128AObj
		chdataOut.Ch = result
		OutCH <- chdataOut
	}
	close(chInway)
}

func (cdma *CDMA) Spread(chInway gocomm.Complex128Channel, OutCH gocomm.Complex128AChannel) {

	indata := <-chInway
	insymbol := indata.Ch
	spcode := cdma.SpreadSequence
	var result = make([]complex128, len(spcode))
	for i := 0; i < len(spcode); i++ {
		result[i] = insymbol * spcode[i]
	}
	var chdataOut gocomm.SComplex128AObj
	chdataOut.Ch = result
	chdataOut.MaxExpected = indata.MaxExpected / len(spcode)
	// chdataOut.MaxExpected = int(math.Floor(float64(indata.MaxExpected) / float64(len(spcode))))
	OutCH <- chdataOut
}
