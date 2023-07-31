package main

import (
	"encoding/json"
	parle "goParliament/src"
	"log"
	"os"
)

func main() {
	jsonFile := "./iam-tests/basic-duplicate-action.json"

	// read our opened jsonFile as a byte array.
	byteValue, _ := os.ReadFile(jsonFile)

	Policy := parle.PolicyDocument{}

	err := json.Unmarshal(byteValue, &Policy)

	if err != nil {
		log.Print(err)
	}

	log.Print(&Policy)
}
