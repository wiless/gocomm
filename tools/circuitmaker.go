package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"
	// "wiless/gocomm/chipset"
	"runtime"
	// "wiless/gocomm/core"
)

type ChipObj struct {
	Type string
	Name string
}

type ConnectionObj struct {
	Source              string   `json:"src"`
	SourcePins          []string `json:"srcOutputPin"`
	Destination         string   `json:"destination"`
	DestinationFunction string   `json:"Modulate"`
}

type Circuit struct {
	Name        string
	Chips       []ChipObj
	Connections []ConnectionObj
}

var toolspath string
var inputfile string

func init() {
	flag.StringVar(&toolspath, "t", "", "The template dir where connection.txt is available")
	flag.StringVar(&inputfile, "i", "", "The template dir where connection.txt is available")
}

func main() {
	flag.Parse()
	if inputfile == "" {
		flag.PrintDefaults()
		return
	}
	fmt.Print("GOPROCS:=", runtime.GOMAXPROCS(8))
	start := time.Now()

	filepath := toolspath + "./" + inputfile

	bytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Printf("\n err = %v", err)
		return
	}

	// fmt.Printf("\n Data  = %s", bytes)
	var obj Circuit

	json.Unmarshal(bytes, &obj)
	fmt.Printf("\n Structure = %+v", obj)

	/// execute template

	t, terr := template.ParseFiles(toolspath + "./connection.tmpl")
	if terr != nil {
		fmt.Printf("\t Error  %v ", terr)
		return
	}
	outfilename := strings.ToLower(obj.Name + ".go")
	fd, ferr := os.Create(outfilename)

	// fmt.Printf("\n===== AUTO GENERATED %s \n", outfilename)
	// packagestr := "package " + packagename
	// fd.WriteString(packagestr)

	if ferr == nil {
		t.Execute(fd, obj)
	}

	fmt.Print("\n Elapsed : ", time.Since(start))

}
