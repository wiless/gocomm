package main

import (
	"fmt"
	"time"
)

type Info struct {
	val int
	id  int
}

func Speaker(spid int, Nwords int, CH chan Info) {
	var info Info
	info.id = spid
	for i := 0; i < Nwords; i++ {
		info.val = i + spid*100
		fmt.Printf("\nI am %d , Shouting %d ", spid, info.val)

		CH <- info
	}
}

func Listener(Nwords int, CH chan Info) {
	for i := 0; i < Nwords; i++ {
		rxinfo := <-CH
		fmt.Printf("\n\t\t\t==========I heard %v ``", rxinfo)
	}
}

func main() {
	CH := make(chan Info, 5)
	Nspeakers := 2
	Nwords := 10
	for i := 0; i < Nspeakers; i++ {
		go Speaker(i, Nwords, CH)
	}

	go Listener(Nspeakers*Nwords, CH)

	time.Sleep(1000 * time.Millisecond)
}
