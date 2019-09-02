package store

import (
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestSetCache(t *testing.T) {
	key := "key_test"
	val := map[string]interface{}{
		"name": "one",
	}
	SetCache(key, val)

	data := GetCache(key)
	assert.Assert(t, is.DeepEqual(data, val))

	val2 := map[string]interface{}{
		"name": "two",
	}

	SetCache(key, val2)

	data = GetCache(key)
	assert.Assert(t, is.DeepEqual(data, val2))

	DeleCache(key)
	data = GetCache(key)
	assert.Assert(t, is.Nil(data))
}
