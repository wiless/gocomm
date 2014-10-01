package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name       string
	Age        int
	MaleFemale bool
	Phone      string ",omitempty"
}

func main() {
	var obj Person
	obj.Name = "Bob"
	obj.Age = 34
	obj.MaleFemale = true
	data, err := bson.Marshal(&obj)
	if err != nil {
		panic(err)
	}
	jdata, _ := json.Marshal(&obj)
	fmt.Printf("\n%s : size %d", data, len(data))
	fmt.Printf("\n%s : size %d ", jdata, len(jdata))
}
