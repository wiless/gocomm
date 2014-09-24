package sink

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	// "strconv"
	"time"
	"wiless/gocomm"
	"wiless/gocomm/chipset"
	// "wiless/gocmm"
	// "wiless/gocomm"
	// "wiless/gocomm/cdma"
	// "wiless/gocomm/modem"
)

type Metric struct {
	Name string
	Val  float64
	Time float64
	Ts   float64
}

func CRO(scale float64, NextSize int, InCH gocomm.Complex128Channel) {
	var txsymbols []complex128
	txsymbols = make([]complex128, 0, NextSize)
	// randid := scale

	// if err != nil {
	// 	fmt.Printf("Connection Found Error")
	// 	connection=false
	// } else {

	var metric Metric
	metric.Name = fmt.Sprintf("EEEBCCCEEN")
	conn, _ := net.Dial("udp", "192.168.0.24:8080")
	var Ts float64 = 1.0
	for cnt := 0; cnt < NextSize; {
		time.Sleep(10 * time.Millisecond)

		chdata := (<-InCH)
		data := chdata.Ch
		// data = temp.Ch
		metric.Time = float64(cnt) * Ts * scale
		metric.Ts = Ts
		// metric.Time = chdata.TimeStamp
		// metric.Time=chdata.

		// if uid > 0 {
		// metric.Time = metric.Time * scale
		// }
		// metric.Val = rand.NormFloat64() // float64(real(data))
		metric.Val = float64(real(data))
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
		binary.Write(buf, binary.LittleEndian, metric.Ts)
		// fmt.Printf("\n AFTER  Buffer is :%0 x  ", buf.Bytes())
		// fmt.Printf("\n METRIC  is :%v ", metric)

		// fmt.Fprintf(conn, "%c", str)
		// fmt.Printf("\n%f %v", real(data), buf.Bytes())
		// if 2 == 3 {
		if conn != nil {
			conn.Write(buf.Bytes())
			// conn.Close()
		}

		// }

		cnt++
		txsymbols = append(txsymbols, data)
		// fmt.Printf("\n %f : Total Received %d of %d : %f", randid, cnt, NextSize, data)

	}
	// }

	// fmt.Printf("\nTxsymbols%f=%f", scale, txsymbols)
}

func CROcomplex(InCH gocomm.Complex128Channel, labels ...string) {

	var metric Metric
	// metric.Name = fmt.Sprintf("EEEBCCCEEN")
	if len(labels) == 0 {
		metric.Name = fmt.Sprintf("EEEBCCCEEN")
	} else {
		metric.Name = labels[0]
		if labels[0] == "" {
			metric.Name = "PlotCurve1"
		}
		if len(metric.Name) > 10 {
			metric.Name = metric.Name[0:10]
		} else {
			metric.Name = metric.Name + strings.Repeat("*", 10-len(metric.Name))
		}
	}
	conn, _ := net.Dial("udp", "192.168.0.24:8080")
	var Ts float64 = 0.252
	NextSize := 1
	// packetbuf := new(bytes.Buffer)
	buf := new(bytes.Buffer)

	for cnt := 0; cnt < NextSize; cnt++ {

		data := <-InCH
		NextSize = data.MaxExpected

		// data = temp.Ch
		metric.Time = float64(cnt) * Ts
		metric.Val = float64(real(data.Ch))
		metric.Ts = Ts
		// str := strconv.FormatFloat(real(data), 'f', 2, 64)

		// databyte := make([]byte, 1024)
		// fmt.Printf("Buffer is :%v ", buf)

		// binary.Write(buf, binary.BigEndian, real(data))
		binary.Write(buf, binary.BigEndian, []byte(metric.Name))
		// buf.WriteByte('0')
		// fmt.Printf("\n AFTER name Buffer is :%s ", buf.Bytes())

		binary.Write(buf, binary.LittleEndian, metric.Val)
		binary.Write(buf, binary.LittleEndian, metric.Time)
		binary.Write(buf, binary.LittleEndian, metric.Ts)
		// fmt.Printf("\n AFTER  Buffer is :%0 x  ", buf.Bytes())
		// fmt.Printf("\n METRIC  is :%v ", metric)

		// fmt.Fprintf(conn, "%c", str)
		// fmt.Printf("\n%f %v", real(data), buf.Bytes())
		// if 2 == 3 {
		// packetbuf.ReadFrom(buf)

		// if math.Mod(float64(cnt), 20.0) == 0 {

		if buf.Len() >= 2040 {

			if conn != nil {
				conn.Write(buf.Bytes())
				// conn.Write(packetbuf.Bytes())
				fmt.Printf("\n Sent %f %v bytes", metric.Time, buf.Len())
			}
			buf.Reset()
		}
		// packetbuf.Reset()

		// The sleep is only to allow the slow replot in Qt applicaiton
		time.Sleep(2 * time.Millisecond)

	}
	// }

}
func CROBitCh(InCH gocomm.BitChannel, labels ...string) {

	var metric Metric
	// metric.Name = fmt.Sprintf("EEEBCCCEEN")
	if len(labels) == 0 {
		metric.Name = fmt.Sprintf("EEEBCCCEEN")
	} else {
		metric.Name = labels[0]
		if labels[0] == "" {
			metric.Name = "PlotCurve1"
		}
		if len(metric.Name) > 10 {
			metric.Name = metric.Name[0:10]
		} else {
			metric.Name = metric.Name + strings.Repeat(" ", 10-len(metric.Name))
		}
	}
	conn, _ := net.Dial("udp", "localhost:8080")
	var Ts float64 = 1.0
	NextSize := 1

	for cnt := 0; cnt < NextSize; cnt++ {

		data := <-InCH
		NextSize = data.MaxExpected

		// data = temp.Ch
		metric.Time = float64(cnt) * Ts
		metric.Val = float64(data.Ch)
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
		}

		// The sleep is only to allow the slow replot in Qt applicaiton
		time.Sleep(10 * time.Millisecond)

	}
}
func CROremote(InCHPin chipset.PinInfo) {
	switch InCHPin.DataType.Name() {

	case "Complex128Channel":
		CROcomplex(InCHPin.Channel.(gocomm.Complex128Channel), InCHPin.Name)

	case "BitChannel":
		CROBitCh(InCHPin.Channel.(gocomm.BitChannel))
	default:
		fmt.Printf("\n unknown Channel Type in the Pin %v", InCHPin)
	}
}
