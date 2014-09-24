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
	"wiless/gocomm"
	"wiless/gocomm/cdma"
	"wiless/gocomm/modem"
	"wiless/gocomm/sink"
	"wiless/gocomm/sources"
)

var N int

var USERS int

func init() {
	N = 32
	USERS = 10
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

	var done []chan bool = make([]chan bool, USERS)

	for user := 0; user < USERS; user++ {
		done[user] = make(chan bool, 1)

		/// Source objects
		bsrc := new(sources.BitSource)
		bsrc.SetSize(N)

		txqpsk := modem.NewModem(2) /// QPSK Modem
		rxqpsk := modem.NewModem(2) /// QPSK Modem
		cdma := CDMA.NewCDMA()

		/// Generate Bits
		bits2modemCH := make(chan uint8)
		go bsrc.GenBit(bits2modemCH)
		NextSize := N

		/// Modulate the block of bits into symbols
		modem2cdmaCH := make(gocomm.Complex128Channel)
		go txqpsk.ModulateBlock(NextSize, bits2modemCH, modem2cdmaCH)
		NextSize = txqpsk.GetOutputBlockSize(NextSize)

		// go SinkDataSample(user, NextSize, modem2cdmaCH, done[user])

		/// Spread the block of Symbols into Chips per frame
		cdma2Tx := make(gocomm.Complex128AChannel)
		go cdma.SpreadBlock(NextSize, modem2cdmaCH, cdma2Tx)
		NextSize = cdma.GetSpreadOutputBlockSize(NextSize)

		upCH := make(gocomm.Complex128Channel)
		go Vector2Sample(user, NextSize, cdma2Tx, upCH)

		tx2rxCH := make(gocomm.Complex128Channel)
		dualChannelA := make([]gocomm.Complex128Channel, 2)
		dualChannelA[0] = make(gocomm.Complex128Channel)
		dualChannelA[1] = make(gocomm.Complex128Channel)
		go ChannelDuplexer(NextSize, upCH, dualChannelA)

		go sink.CRO(0.125, NextSize, dualChannelA[0])
		go Channel(user, NextSize, dualChannelA[1], tx2rxCH)

		downCH := make(gocomm.Complex128AChannel)
		go Sample2Vector(user, NextSize, 1, tx2rxCH, downCH)

		/// Despread the input chips into symbol per frame
		despread2modem := make(gocomm.Complex128Channel)
		go cdma.DeSpreadBlock(NextSize, downCH, despread2modem)
		NextSize = cdma.GetDeSpreadOutputBlockSize(NextSize)

		SinkCH := make(gocomm.Complex128Channel)
		go rxqpsk.DeModulateBlock(NextSize, despread2modem, SinkCH)

		// /*
		// 	The following SinkData should be last function to be
		// 	called and waited for it to send the DONE flag for each user
		// */
		// go SinkDataVector(user, NextSize, SinkCH, done[user])

		// go SinkDataSample(user, NextSize, SinkCH, done[user])
		go SinkStreamDataSample(user, NextSize, SinkCH, done[user])

	}

	//< This waits till done is returned from all the USER's channel from the SinkData
	for i := 0; i < USERS; i++ {
		<-done[i]
	}
	fmt.Println("\n", time.Now())
	fmt.Printf("\nTime Elaspsed %v \n", time.Since(t))

}

/// Fading/AWGN Channel that operates on each sample
func Channel(uid int, NextSize int, InCH gocomm.Complex128Channel, OutCH gocomm.Complex128Channel) {
	if uid == 3 {
		dur := time.Millisecond * time.Duration(int64(rand.Intn(100)))
		time.Sleep(dur)
	}
	N0 := .01 /// 10dB SNR
	for i := 0; i < NextSize; i++ {
		sample := <-InCH
		/// Do the processing here
		// gain := 1 //sources.RandNC(1)
		noise := sources.RandNC(N0)
		///Fading
		// psample := sample * gain
		psample := sample + noise
		///
		OutCH <- psample
	}
	close(InCH)
}

/// Converts each Vector Sample to a Sample which can be processed at sample rate
/// This can be considered as Upsampler Each vector at rate Ts , is communicated to the next block at Ts/N samples
func Vector2Sample(uid int, NextSize int, InCH gocomm.Complex128AChannel, OutCH gocomm.Complex128Channel) {

	cnt := 0
	for i := 0; i < NextSize; i++ {
		indata := <-InCH
		veclen := len(indata)
		cnt += veclen
		for indx := 0; indx < veclen; indx++ {
			OutCH <- indata[indx]
		}
	}
	fmt.Printf("\n User%d : Closing", uid)
	close(InCH)
}

func ChannelDuplexer(NextSize int, InCH gocomm.Complex128Channel, OutCHA []gocomm.Complex128Channel) {
	Nchanels := len(OutCHA)

	for cnt := 0; cnt < NextSize; cnt++ {

		data := <-InCH
		// fmt.Printf("%d InputDuplexer : %v ", cnt, data)
		for i := 0; i < Nchanels; i++ {
			OutCHA[i] <- data
		}
	}
	close(InCH)
}

/// Converts each Vector Sample to a Sample which can be processed at sample rate
/// This can be considered as DownSample Each vector at rate Ts , is communicated to the next block at Ts/N samples
func Sample2Vector(uid int, NextSize int, factor int, InCH gocomm.Complex128Channel, OutCH gocomm.Complex128AChannel) {

	cnt := 0
	for i := 0; i < NextSize; i++ {

		vecdata := make([]complex128, factor)
		for i := 0; i < factor; i++ {
			vecdata[i] = <-InCH
		}
		cnt += factor
		OutCH <- vecdata
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
		data := <-InCH

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
		data := <-InCH
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
		data := <-InCH
		cnt += len(data)
		txsymbols = append(txsymbols, data...)
		fmt.Printf("\n %d : Received Sample (%d of %d)", randid, cnt, NextSize)
	}
	fmt.Printf("\nTxsymbols%d=%f", uid, txsymbols)
	close(InCH)
	done <- true

}
