package parlez_test

import (
	parlez "goParliament/src"
	"os"
	"testing"
)

func TestPolicyDocument_UnmarshalJSON(t *testing.T) {

	type fields struct {
		Version   string
		Statement []parlez.StatementEntry
	}

	type args struct {
		b []byte
	}

	basic, _ := os.ReadFile("../iam-tests/basic.json")
	dupe, _ := os.ReadFile("../iam-tests/basic-duplicate-action.json")
	notPretty, _ := os.ReadFile("../iam-tests/basic-notpretty.json")
	withID, _ := os.ReadFile("../iam-tests/basic-withid.json")
	withIDSquare, _ := os.ReadFile("../iam-tests/basic-withidsquare.json")
	withCondition, _ := os.ReadFile("../iam-tests/basic-withcondition.json")
	withSid, _ := os.ReadFile("../iam-tests/basic-withsid.json")

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"basic", fields{}, args{basic}, false},
		{"dupe", fields{}, args{dupe}, true},
		{"not-pretty", fields{}, args{notPretty}, false},
		{"withID", fields{}, args{withID}, false},
		{"withSID", fields{}, args{withSid}, false},
		{"withCondition", fields{}, args{withCondition}, false},
		{"todo", fields{}, args{withIDSquare}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policyDocument := &parlez.PolicyDocument{
				Version:   tt.fields.Version,
				Statement: tt.fields.Statement,
			}
			if err := policyDocument.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
