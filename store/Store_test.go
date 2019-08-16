package store

import (
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestGetDataNotExists(t *testing.T) {
	_, e := GetData("./not_exist.json")
	if e == nil {
		t.Error("file not exists can not read data")
	}
}
func TestGetDataExists(t *testing.T) {
	d, e := GetData("../testdata/1.json")
	assert.Assert(t, is.Nil(e))
	if e != nil {
		t.Error("read json file should get json")
	} else {
		data := make(map[string]interface{})
		data["name"] = "test"
		data["age"] = float64(10)
		data["arr"] = []interface{}{float64(1), float64(2), float64(3)}
		data["child"] = []interface{}{
			map[string]interface{}{
				"name": "one",
			},
		}
		assert.Assert(t, is.DeepEqual(data, d))
	}
}

func TestSave(t *testing.T) {
	data := make(map[string]interface{})
	data["name"] = "test"
	data["age"] = float64(10)
	data["arr"] = []interface{}{float64(1), float64(2), float64(3)}
	data["child"] = []interface{}{
		map[string]interface{}{
			"name": "one",
		},
	}
	filepath := "../cache/1.json"
	e := Save(filepath, data)
	if e != nil {
		t.Errorf("save get error %v", e)
	}
	d, _ := GetData(filepath)
	assert.Assert(t, is.DeepEqual(data, d))

	// test cache
	valCache := GetCache(filepath)
	assert.Assert(t, is.DeepEqual(data, valCache))
}

func TestSaveNil(t *testing.T) {
	filepath := "../cache/1.json"
	e := Save(filepath, nil)
	assert.Assert(t, is.Error(e, "data is nil"))

	data := make(map[string]interface{})
	e1 := Save(filepath, data)
	assert.Assert(t, is.Error(e1, "data [{}] is empty"))
}
