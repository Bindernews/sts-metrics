package main

import (
	"encoding/json"
	"log"
)

type Foo struct {
	A []*string `json:"a"`
}

func main() {
	foo := Foo{}
	const test1 = `{ "a": ["q", "w", "e", null, "r", "t", "y"]}`
	if err := json.Unmarshal([]byte(test1), &foo); err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", foo)
}
