package main

import (
	"fmt"
	"time"
	// "wiless/gocomm/chipset"
	"runtime"
	"wiless/gocomm/core"
)

func main() {

	fmt.Print("GOPROCS:=", runtime.GOMAXPROCS(16))
	runtime.SetCPUProfileRate(-1)
	start := time.Now()
	user1 := core.NewSetup()
	user2 := core.NewSetup()
	user3 := core.NewSetup()
	user4 := core.NewSetup()
	user5 := core.NewSetup()
	// fmt.Printf("\nLink %v", user1)
	// fmt.Printf("\nLink %v", user2)
	var jsons string = `{"NBlocks":100,"snr":"0:2:16","SF":1}`

	user1.Set(jsons)
	user2.Set(jsons)
	user3.Set(jsons)
	user4.Set(jsons)
	user5.Set(jsons)
	fmt.Printf("Starting simulation ...")

	go user1.Run()
	go user2.Run()
	go user3.Run()
	go user4.Run()

// 	user1.Run()
// 	user2.Run()
// 	user3.Run()
// 	user4.Run()
	fmt.Printf("\n started user 5")
	user5.Run()

	// time.Sleep(10 * time.Second)
	fmt.Print("\n Elapsed : ", time.Since(start))

}
