package parlez

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

type PolicyDocument struct {
	Version   string           `json:"Version"`
	Id        string           `json:"Id"`
	Statement []StatementEntry `json:"Statement"`
}

type StatementEntry struct {
	Effect    string     `json:"Effect"`
	Action    []string   `json:"Action"`
	Resource  string     `json:"Resource"`
	Condition *Condition `json:"Condition,omitempty"`
}

type Condition struct {
	StringLike struct {
		S3Prefix []string `json:"s3:prefix,omitempty"`
	} `json:"StringLike,omitempty"`
}

func (policyDocument *PolicyDocument) UnmarshalJSON(b []byte) error {

	var f interface{}
	err := json.Unmarshal(b, &f)

	if err != nil {
		log.Print(err)
	}

	m, ok := f.(map[string]interface{})

	if !ok {
		log.Print(ok)
	}

	policyDocument.Version = m["Version"].(string)
	policyDocument.Id = m["Id"].(string)

	var entries []StatementEntry

	statements, ok := m["Statement"].([]interface{})
	if ok {
		for _, statement := range statements {
			entry := NewStatement(statement.(map[string]interface{}))
			entries = append(entries, entry)
		}

		policyDocument.Statement = entries
	} else {
		basic := m["Statement"].(map[string]interface{})

		entry := NewStatement(basic)

		policyDocument.Statement = append(policyDocument.Statement, entry)
	}

	test, _ := json.Marshal(policyDocument)

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		fmt.Println(err)
	}
	result := bytes.Compare(buffer.Bytes(), test)

	if result != 0 {
		log.Print("*** Json Mismatch error ***")
		log.Print("Original:")
		log.Print(string(buffer.Bytes()))
		log.Print("is not equal to Imported:")
		log.Printf(string(test))
		return fmt.Errorf("json import parsing failure detected")
	}

	return nil
}

func NewStatement(basic map[string]interface{}) StatementEntry {
	var entry StatementEntry

	if basic["Action"] != nil {
		Actions, ok := basic["Action"].([]interface{})
		if ok {
			for _, action := range Actions {
				entry.Action = append(entry.Action, action.(string))
			}
		} else {
			entry.Action = append(entry.Action, basic["Action"].(string))
		}

	}

	value, ok := basic["Effect"]
	if ok {
		entry.Effect = value.(string)
	}

	value, ok = basic["Resource"]
	if ok {
		entry.Resource = value.(string)
	}
	return entry
}
