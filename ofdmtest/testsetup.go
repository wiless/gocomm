package main

import (
	"fmt"
	"time"
	// "wiless/gocomm/chipset"
	"runtime"
	"wiless/gocomm/core"
)

func main() {
	fmt.Print("GOPROCS:=", runtime.GOMAXPROCS(8))

	start := time.Now()
	user1 := core.NewSetup()
	user2 := core.NewSetup()
	// fmt.Printf("\nLink %v", user1)
	// fmt.Printf("\nLink %v", user2)
	var jsons string = `{"NBlocks":100,"snr":"0:0","SF":1}`
	user1.Set(jsons)
	fmt.Printf("Starting simulation ...")

	jsons = `{"NBlocks":100,"snr":"1:2","SF":1}`
	user2.Set(jsons)
	fmt.Printf("Starting simulation ...")

	go user1.Run()
	fmt.Printf("\n started user 2")
	user2.Run()

	// time.Sleep(10 * time.Second)
	fmt.Print("\n Elapsed : ", time.Since(start))

}
