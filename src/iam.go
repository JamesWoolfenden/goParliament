package parlez

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type PolicyDocument struct {
	Version   string           `json:"Version"`
	Id        string           `json:"Id,omitempty"`
	Statement []StatementEntry `json:"Statement"`
}

type StatementEntry struct {
	Sid        *string     `json:"Sid,omitempty"`
	Effect     string      `json:"Effect"`
	Action     interface{} `json:"Action"`
	NotActions interface{} `json:"NotActions,omitempty"`
	Resource   interface{} `json:"Resource,omitempty"`
	Condition  *Condition  `json:"Condition,omitempty"`
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

	Id, ok := m["Id"].(string)

	if ok && Id != "" {
		policyDocument.Id = Id
	}

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

	result := reflect.DeepEqual(buffer.Bytes(), test)

	if !result {
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

	if basic["Sid"] != nil {
		Sid, ok := basic["Sid"].(string)
		if ok {
			entry.Sid = &Sid
		}
	}

	entry.Action = basic["Action"]

	value, ok := basic["Effect"].(string)
	if ok {
		entry.Effect = value
	}

	entry.Resource = basic["Resource"]

	return entry
}
