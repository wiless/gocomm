package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
	// "wiless/gocmm"
	// "reflect"
	"wiless/gocomm"
	// "wiless/gocomm/cdma"
	"wiless/gocomm/channel"
	"wiless/gocomm/chipset"
	"wiless/gocomm/modem"
	"wiless/gocomm/sink"
	"wiless/gocomm/sources"
)

var N int

var USERS int

func init() {
	N = 20
	USERS = 1
	rand.Seed(time.Now().UTC().UnixNano())

	flag.Parse()
	fmt.Printf("\n No of Args %d ", flag.NArg())
	if flag.NArg() == 0 {
		fmt.Printf("\n Defaults : Using a frame of %d bits with %d users \n", N, USERS)
	} else {

		N, _ = strconv.Atoi(flag.Args()[0])
		USERS, _ = strconv.Atoi(flag.Args()[1])
		fmt.Printf("\n Using a frame of %d bits with %d users \n", N, USERS)
	}

}

func main() {
	t := time.Now()

	bsrc := new(sources.BitSource)
	txmodem := modem.NewModem(2)
	rxmodem := modem.NewModem(2)
	var chem1, chem2 channel.ChannelEmulator

	/// Initialize each module
	bsrc.SetSize(N)
	bsrc.InitializeChip()

	chem1.InitializeChip()

	chem2.InitializeChip()
	chem2.SetNoise(0, 2)

	// var newwire chipset.Wire
	var chip []chipset.Chip
	var wire []chipset.Wire

	chip = make([]chipset.Chip, 5)
	wire = make([]chipset.Wire, len(chip)-1)

	// var bitch gocomm.BitChannel
	// fmt.Printf("\n TYPE OF CHannel is %v", reflect.TypeOf(bitch))
	chipcnt := 0

	chip[chipcnt] = bsrc
	chipcnt++
	chip[chipcnt] = txmodem
	chipcnt++
	chip[chipcnt] = chem1
	chipcnt++
	chip[chipcnt] = chem2
	chipcnt++
	chip[chipcnt] = rxmodem
	for i := 0; i < len(wire); i++ {
		wire[i].Join(chip[i], chip[i+1])
	}

	modules := [...]string{"genbit", "modulate", "fadingChannel", "awgn", "demodulate"}
	junctionwire := 2 // @output of chem1

	var success bool
	var outpin string

	bitch := (bsrc.PinByName("bitOut").Channel.(gocomm.BitChannel))
	outpin = chip[0].PinByID(chip[0].ModuleByName(modules[0]).OutPins[0]).Name
	go bsrc.GenBit(bitch)

	for i := 0; i < len(wire); i++ {
		fmt.Printf("\n Wire %d ", i)
		if i == junctionwire {
			wire[i].Split(2)
		}
		success, outpin = wire[i].ConnectPins(outpin, modules[i+1])
	}

	// success, outpin = wire1.ConnectPins(outpin, "modulate")
	// success, outpin = wire2.ConnectPins(outpin, "awgn")
	// success, outpin = wire2.ConnectPins(outpin, "awgn")
	// success, outpin = wire3.ConnectPins(outpin, "demodulate")
	go Sink(wire[junctionwire].GetProbe(0))
	if success {
		lastwire := wire[len(wire)-1]
		pin := lastwire.DestinationChip.PinByName(lastwire.RecentOutputPinName())
		Sink(pin)

	}

	return

	//< This waits till done is returned from all the USER's channel from the SinkData
	// for i := 0; i < USERS; i++ {
	// 	<-done[i]
	// }
	fmt.Println("\n", time.Now())
	fmt.Printf("\nTime Elaspsed %v \n", time.Since(t))

}

func Sink(pin chipset.PinInfo) {

	fmt.Printf("\n=======================\n  Will Sink DataOut from Pin %v", pin)
	count := 1
	switch pin.DataType.Name() {
	case "BitChannel":
		for i := 0; i < count; i++ {
			// fmt.Printf("\n Status of Channel %d = %#v ", i, pin.Channel)
			ddata := <-pin.Channel.(gocomm.BitChannel)
			// fmt.Printf(" SPECIAL MESSAGE %s", ddata.Message)
			if ddata.Message == "" {
				fmt.Printf("\nPin : %s - Read Bit %d = %v ", pin.Name, i, ddata.Ch)
			} else {
				fmt.Printf("\nPin : %s - Read Bit %d = %v : %s", pin.Name, i, ddata.Ch, ddata.Message)
			}

			count = ddata.MaxExpected
			// ddata := choutData.Ch
			// max = choutData.MaxExpected
			// fmt.Printf(" %d %d", uint8(real(ddata)), uint8(imag(ddata)))
			// fmt.Printf("\n %d @ max Symbols limit  = %d %s ", i, max, choutData.Message)

		}
	case "Complex128Channel":
		for i := 0; i < count; i++ {
			ddata := <-pin.Channel.(gocomm.Complex128Channel)
			// fmt.Printf(" SPECIAL MESSAGE %s", ddata.Message)
			if ddata.Message == "" {
				fmt.Printf("\nPin : %s - Read Complex %d = %v ", pin.Name, i, ddata.Ch)
			} else {
				fmt.Printf("\nPin : %s - Read Complex %d = %v : %s", pin.Name, i, ddata.Ch, ddata.Message)
			}
			count = ddata.MaxExpected
			// ddata := choutData.Ch
			// max = choutData.MaxExpected
			// fmt.Printf(" %d %d", uint8(real(ddata)), uint8(imag(ddata)))
			// fmt.Printf("\n %d @ max Symbols limit  = %d %s ", i, max, choutData.Message)

		}
	default:
		fmt.Printf("\n Unknown Data type")
	}

}

/// Converts each Vector Sample to a Sample which can be processed at sample rate
/// This can be considered as Upsampler Each vector at rate Ts , is communicated to the next block at Ts/N samples
func Vector2Sample(uid int, NextSize int, InCH gocomm.Complex128ChannelA, OutCH gocomm.Complex128Channel) {
	var chdataOut gocomm.SComplex128Channel
	var chdataIn gocomm.SComplex128ChannelA

	cnt := 0
	for i := 0; i < NextSize; i++ {
		chdataIn = <-InCH
		indata := chdataIn.Ch
		veclen := len(indata)
		cnt += veclen

		for indx := 0; indx < veclen; indx++ {
			chdataOut.Ch = indata[indx]
			OutCH <- chdataOut
		}
	}
	fmt.Printf("\n User%d : Closing", uid)

	close(InCH)
}

func ChannelDuplexer(InCH gocomm.Complex128Channel, OutCHA []gocomm.Complex128Channel) {
	Nchanels := len(OutCHA)
	var chdataIn gocomm.SComplex128Channel
	var chdataOut gocomm.SComplex128Channel
	NextSize := 1
	for cnt := 0; cnt < NextSize; cnt++ {
		chdataIn = <-InCH
		data := chdataIn.Ch
		NextSize = chdataIn.MaxExpected

		// fmt.Printf("%d InputDuplexer : %v ", cnt, data)
		for i := 0; i < Nchanels; i++ {
			chdataOut.Ch = data
			chdataOut.MaxExpected = NextSize
			chdataOut.Message = chdataIn.Message
			OutCHA[i] <- chdataOut
		}
	}
	close(InCH)
}

/// Converts each Vector Sample to a Sample which can be processed at sample rate
/// This can be considered as DownSample Each vector at rate Ts , is communicated to the next block at Ts/N samples
func Sample2Vector(uid int, NextSize int, factor int, InCH gocomm.Complex128Channel, OutCH gocomm.Complex128ChannelA) {
	var chdataOut gocomm.SComplex128ChannelA

	cnt := 0
	for i := 0; i < NextSize; i++ {

		vecdata := make([]complex128, factor)
		for i := 0; i < factor; i++ {
			vecdata[i] = (<-InCH).Ch
		}
		cnt += factor
		chdataOut.Ch = vecdata
		OutCH <- chdataOut
	}
	fmt.Printf("\n User%d : Closing", uid)
	close(InCH)
}

func SinkStreamDataSample(uid int, NextSize int, InCH gocomm.Complex128Channel, done chan bool) {
	var txsymbols []complex128
	txsymbols = make([]complex128, 0, NextSize)
	// randid := uid

	// if err != nil {
	// 	fmt.Printf("Connection Found Error")
	// 	connection=false
	// } else {

	var metric sink.Metric
	metric.Name = fmt.Sprintf("EEEBCCCEE%d", uid)
	conn, _ := net.Dial("udp", "localhost:8080")
	var Ts float64 = 1.0
	fmt.Printf("\n New Sink")
	for cnt := 0; cnt < NextSize; {
		time.Sleep(10 * time.Millisecond)
		data := (<-InCH).Ch

		metric.Time = float64(cnt) * Ts

		// metric.Val = rand.NormFloat64() // float64(real(data))
		metric.Val = (2.*float64(real(data)) - 1) * 4
		// str := strconv.FormatFloat(real(data), 'f', 2, 64)

		// databyte := make([]byte, 1024)
		buf := new(bytes.Buffer)
		// fmt.Printf("Buffer is :%v ", buf)

		// binary.Write(buf, binary.BigEndian, real(data))
		binary.Write(buf, binary.BigEndian, []byte(metric.Name))
		// buf.WriteByte('0')
		// fmt.Printf("\n AFTER name Buffer is :%s ", buf.Bytes())

		binary.Write(buf, binary.LittleEndian, metric.Val)
		binary.Write(buf, binary.LittleEndian, metric.Time)
		// fmt.Printf("\n AFTER  Buffer is :%0 x  ", buf.Bytes())
		// fmt.Printf("\n METRIC  is :%v ", metric)

		// fmt.Fprintf(conn, "%c", str)
		// fmt.Printf("\n%f %v", real(data), buf.Bytes())
		// if 2 == 3 {
		if conn != nil {
			conn.Write(buf.Bytes())
			// conn.Close()
			fmt.Printf("\r %d Packet Set %d Total %v bytes", uid, cnt, cnt*buf.Len())
		}

		// }

		cnt++
		txsymbols = append(txsymbols, data)
		// fmt.Printf("\n %d : Total Received %d of %d : %f", randid, cnt, NextSize, data)

	}
	// }

	fmt.Printf("\nTxsymbols%d=%f", uid, txsymbols)
	close(InCH)
	done <- true
}

func SinkDataSample(uid int, NextSize int, InCH gocomm.Complex128Channel, done chan bool) {

	var txsymbols []complex128
	txsymbols = make([]complex128, 0, NextSize)
	randid := uid
	// fmt.Println(randid)
	for cnt := 0; cnt < NextSize; {
		data := (<-InCH).Ch
		cnt++
		txsymbols = append(txsymbols, data)
		fmt.Printf("\n %d : Total Received %d of %d : %f", randid, cnt, NextSize, data)
	}
	fmt.Printf("\nTxsymbols%d=%f", uid, txsymbols)
	close(InCH)
	done <- true
}

func SinkDataVector(uid int, NextSize int, InCH gocomm.Complex128ChannelA, done chan bool) {
	var txsymbols []complex128
	txsymbols = make([]complex128, 0, NextSize)
	randid := uid
	// fmt.Println(randid)
	for cnt := 0; cnt < NextSize; {
		data := (<-InCH).Ch
		cnt += len(data)
		txsymbols = append(txsymbols, data...)
		fmt.Printf("\n %d : Received Sample (%d of %d)", randid, cnt, NextSize)
	}
	fmt.Printf("\nTxsymbols%d=%f", uid, txsymbols)
	close(InCH)
	done <- true

}
