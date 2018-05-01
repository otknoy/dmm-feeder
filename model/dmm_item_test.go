package model

import (
	"encoding/json"
	"testing"
)

type TestData struct {
	Number NumberOrString `json:"number"`
	String NumberOrString `json:"string"`
}

func TestNumberOrString_UnmarshalJSON(t *testing.T) {
	jsonStr := "{\"number\":1234, \"string\":\"1234_ruby\"}"

	var testData TestData
	err := json.Unmarshal(([]byte)(jsonStr), &testData)
	if err != nil {
		t.Error(err)
	}
	if testData.Number != "1234" {
		t.Errorf("fail: testData.Number (%s) is not equal to \"1234\"", testData.Number)
	}

	if testData.String != "1234_ruby" {
		t.Errorf("fail: testData.String (%s) is not equal to \"1234_ruby\"", testData.String)
	}
}
