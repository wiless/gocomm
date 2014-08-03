package sink

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	// "strconv"
	"time"
	"wiless/gocomm"
	// "wiless/gocmm"
	// "wiless/gocomm"
	// "wiless/gocomm/cdma"
	// "wiless/gocomm/modem"
)

type Metric struct {
	Name string
	Val  float64
	Time float64
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
	conn, _ := net.Dial("udp", "localhost:8080")
	var Ts float64 = 1.0
	for cnt := 0; cnt < NextSize; {
		time.Sleep(10 * time.Millisecond)

		data := (<-InCH).Ch
		// data = temp.Ch
		metric.Time = float64(cnt) * Ts * scale

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
