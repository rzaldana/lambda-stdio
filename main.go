package main

import (
	"fmt"
	"encoding/json"
)

func main() {
	fmt.Println("Hello world!")
	b, err := json.Marshal(json.RawMessage("not valid json\""))
	if err != nil {
		fmt.Printf("ERROR: couldn't unmarshall string. Error: '%v'\n", err)
	}

	fmt.Printf("Marshalled JSON: '%v'\n", string(b))
}
