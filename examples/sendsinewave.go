package main

import (
	"fmt"
	"net"
	"strconv"
)

func main() {

	SendSinWaveForm(1000)
}

func SendSinWaveForm(nsamples int) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Found Error")
	}
	// fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	for i := 0; i < nsamples; i++ {
		str := strconv.FormatInt(int64(i), 10)
		fmt.Fprintf(conn, "%s\n", str)
	}

}
