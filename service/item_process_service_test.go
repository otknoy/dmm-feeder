package service

import (
	"encoding/json"
	"testing"

	"github.com/otknoy/dmm-feeder/model"
)

func createTestData() model.DmmItem {
	jsonStr := `
{
  "content_id": "hoge",
  "iteminfo": {
    "genre": [
      {
        "id": 4025,
        "name": "単体作品"
      },
      {
        "id": 6548,
        "name": "独占配信"
      }
    ],
    "maker": [
      {
        "id": 40072,
        "name": "ハヤブサ"
      }
    ],
    "actress": [
      {
        "id": 123,
        "name": "sample-name"
      },
      {
        "id": "123_ruby",
        "name": "さんぷるねーむ"
      },
      {
        "id": "123_classify",
        "name": "av"
      },
      {
        "id": 1011199,
        "name": "上原亜衣"
      },
      {
        "id": "1011199_ruby",
        "name": "うえはらあい"
      },
      {
        "id": "1011199_classify",
        "name": "av"
      }
    ]
  }
}`

	var dmmItem model.DmmItem
	json.Unmarshal(([]byte)(jsonStr), &dmmItem)

	return dmmItem
}

var dmmItem model.DmmItem = createTestData()

func TestNewItemProcessService(t *testing.T) {
	NewItemProcessService()
}

// func TestProcess(t *testing.T) {
// }

func TestParseActress(t *testing.T) {
	expected := parseActress(dmmItem)

	actual := []string{"sample-name", "上原亜衣"}

	if expected[0] != actual[0] || expected[1] != actual[1] {
		t.Errorf("fail: expected=%s, actual%s", expected, actual)
	}
}

func TestParseGenre(t *testing.T) {
	expected := parseGenre(dmmItem)
	actual := []string{"単体作品", "独占配信"}

	if expected[0] != actual[0] || expected[1] != actual[1] {
		t.Errorf("fail: expected=%s, actual%s", expected, actual)
	}
}
func TestParseMaker(t *testing.T) {
	expected := parseMaker(dmmItem)

	actual := []string{"ハヤブサ"}

	if expected[0] != actual[0] {
		t.Errorf("fail: expected=%s, actual%s", expected, actual)
	}
}
