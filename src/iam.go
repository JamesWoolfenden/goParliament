package parlez

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

type AnonymousPolicyDocument struct {
	Version   string      `json:"Version"`
	Id        string      `json:"Id,omitempty"`
	Statement interface{} `json:"Statement"`
}

type PolicyDocument struct {
	Version   string           `json:"Version"`
	Id        string           `json:"Id,omitempty"`
	Statement []StatementEntry `json:"Statement"`
}

type StatementEntry struct {
	Sid          *string    `json:"Sid,omitempty"`
	Effect       string     `json:"Effect"`
	Action       []string   `json:"Action,omitempty"`
	NotActions   []string   `json:"NotActions,omitempty"`
	Resource     []string   `json:"Resource,omitempty"`
	NotResource  []string   `json:"NotResource,omitempty"`
	Principal    []string   `json:"Principal,omitempty"`
	NotPrincipal []string   `json:"NotPrincipal,omitempty"`
	Condition    *Condition `json:"Condition,omitempty"`
}

type Condition struct {
	StringLike struct {
		S3Prefix []string `json:"s3:prefix,omitempty"`
	} `json:"StringLike,omitempty"`
}

func VerifyJSON(b []byte) bool {

	var Policy AnonymousPolicyDocument
	var f interface{}
	err := json.Unmarshal(b, &f)

	if err != nil {
		log.Print(err)
	}

	m, ok := f.(map[string]interface{})

	if !ok {
		log.Print(ok)
	}

	Policy.Version = m["Version"].(string)

	Id, ok := m["Id"].(string)

	if ok && Id != "" {
		Policy.Id = Id
	}

	statements, ok := m["Statement"].([]interface{})
	if ok {
		Policy.Statement = statements
	} else {
		basic := m["Statement"].(map[string]interface{})

		Policy.Statement = basic
	}

	test, _ := json.Marshal(Policy)

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, b); err != nil {
		fmt.Println(err)
	}

	if len(buffer.String()) != len(string(test)) {
		log.Print("*** Json Mismatch error ***")
		log.Print("Original:")
		log.Print(buffer.String())
		log.Print("is not equal to Imported:")
		log.Print(string(test))
		return false
	}

	return true
}

func (policyDocument *PolicyDocument) UnmarshalJSON(b []byte) error {

	if !VerifyJSON(b) {
		return fmt.Errorf("failed to verify json import")
	}

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
			entry := newStatement(statement.(map[string]interface{}))
			entries = append(entries, entry)
		}

		policyDocument.Statement = entries
	} else {
		basic := m["Statement"].(map[string]interface{})

		entry := newStatement(basic)

		policyDocument.Statement = append(policyDocument.Statement, entry)
	}

	return nil
}

func newStatement(basic map[string]interface{}) StatementEntry {
	var entry StatementEntry

	if basic["Sid"] != nil {
		Sid, ok := basic["Sid"].(string)
		if ok {
			entry.Sid = &Sid
		}
	}

	if basic["Action"] != "" {
		action, ok := basic["Action"].(string)
		if ok {
			entry.Action = append(entry.Action, action)
		} else {
			actions, ok := basic["Action"].([]interface{})
			if ok {
				for _, action := range actions {
					entry.Action = append(entry.Action, action.(string))
				}
			}
		}
	}

	value, ok := basic["Effect"].(string)
	if ok {
		entry.Effect = value
	}

	if basic["Resource"] != "" {
		resource, ok := basic["Resource"].(string)
		if ok {
			entry.Resource = append(entry.Resource, resource)
		} else {
			resources, ok := basic["Resource"].([]interface{})
			if ok {
				for _, resource := range resources {
					entry.Resource = append(entry.Action, resource.(string))
				}
			}
		}
	}

	if basic["Condition"] != nil {
		myCondition := basic["Condition"].(map[string]interface{})
		myStringLike := myCondition["StringLike"].(map[string]interface{})

		var limit Condition

		prefixes := myStringLike["s3:prefix"].([]interface{})
		for _, prefix := range prefixes {
			limit.StringLike.S3Prefix = append(limit.StringLike.S3Prefix, prefix.(string))
		}

		entry.Condition = &limit
	}

	return entry
}
