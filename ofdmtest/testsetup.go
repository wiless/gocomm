package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"time"
	"wiless/gocomm/chipset"
	"wiless/gocomm/core"
)

func main() {

	fmt.Print("GOPROCS:=", runtime.GOMAXPROCS(8))
	runtime.SetCPUProfileRate(-1)
	start := time.Now()
	user1 := core.NewSetup()
	user2 := core.NewSetup()
	// user3 := core.NewSetup()
	// user4 := core.NewSetup()
	// user5 := core.NewSetup()
	// fmt.Printf("\nLink %v", user1)
	// fmt.Printf("\nLink %v", user2)
	data, err := ioutil.ReadFile("settings.json")
	if err != nil {
		log.Print("Unable to Read File : ", err)
	}
	result := chipset.GetMetaInfo(data, "Modem1")
	fmt.Print("Found Setting : ", result, "len = ", len(result))
	var jsons string = `{"NBlocks":100,"snr":"0:2:16","SF":1}`
	var mymodem core.Modem
	mymodem.SetName("Modem2")
	mymodem.SetJson(data)

	fmt.Print("SOMETHING", string(mymodem.GetJson()))
	user1.Set(jsons)
	user2.Set(jsons)
	// user3.Set(jsons)
	// user4.Set(jsons)
	// user5.Set(jsons)
	fmt.Printf("Starting simulation ...")

	go user1.Run()
	// go user2.Run()
	// go user3.Run()
	// go user4.Run()

	// 	user1.Run()
	// 	user2.Run()
	// 	user3.Run()
	// 	user4.Run()
	fmt.Printf("\n started user 2")
	user2.Run()

	// time.Sleep(10 * time.Second)
	fmt.Print("\n Elapsed : ", time.Since(start))

}
