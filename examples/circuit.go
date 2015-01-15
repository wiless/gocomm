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
	// "github.com/wiless/gocmm"
	"reflect"
	"github.com/wiless/gocomm"

	// "github.com/wiless/gocomm/cdma"
	"github.com/wiless/gocomm/channel"
	"github.com/wiless/gocomm/chipset"
	"github.com/wiless/gocomm/modem"
	"github.com/wiless/gocomm/sink"
	"github.com/wiless/gocomm/sources"
)

var N int

var USERS int

func init() {
	N = 5120
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

	fmt.Printf("\n 1. Type of Modem %v", reflect.TypeOf(txmodem).Name())
	fmt.Printf("\n 2. Type of channel %v", reflect.TypeOf(chem1).Name())
	fmt.Printf("\n 2. Type of Bitsrc %v", reflect.TypeOf(*bsrc).Name())
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
	// junctionwire := 0 // @output of genbit

	var success bool
	var outpin string

	bitch := (bsrc.PinByName("bitOut").Channel.(gocomm.BitChannel))
	outpin = chip[0].PinByID(chip[0].ModuleByName(modules[0]).OutPins[0]).Name
	go bsrc.GenBit(bitch)

	for i := 0; i < len(wire); i++ {
		fmt.Printf("\n Wire %d ", i)
		// if i == junctionwire {
		// 	wire[i].Split(2)
		// }
		success, outpin = wire[i].ConnectPins(outpin, modules[i+1])
	}

	// success, outpin = wire1.ConnectPins(outpin, "modulate")
	// success, outpin = wire2.ConnectPins(outpin, "awgn")
	// success, outpin = wire2.ConnectPins(outpin, "awgn")
	// success, outpin = wire3.ConnectPins(outpin, "demodulate")
	// go Sink(wire[junctionwire].GetProbe(0))
	// fmt.Printf("%v ", wire[junctionwire].GetProbe(0).DataType)

	// go sink.CROremote(wire[junctionwire].GetProbe(0))

	// go Sink(wire[junctionwire].GetProbe(0))

	if success {
		lastwire := wire[len(wire)-1]
		pin := lastwire.DestinationChip.PinByName(lastwire.RecentOutputPinName())
		//Sink(pin)
		sink.CROremote(pin)

	}

	return

	//< This waits till done is returned from all the USER's channel from the SinkData
	// for i := 0; i < USERS; i++ {
	// 	<-done[i]
	// }
	fmt.Println("\n", time.Now())
	fmt.Printf("\nTime Elaspsed %v \n", time.Since(t))

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

func SinkDataVector(uid int, NextSize int, InCH gocomm.Complex128AChannel, done chan bool) {
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
