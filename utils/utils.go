package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintResInJson(input interface{}) {
	out, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling: %v", err)
	}
	fmt.Println(string(out))
}
