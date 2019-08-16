package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"stone/logger"
)

var loggerStore = logger.GetPrefixLogger("store")

// GetData read data from file
func GetData(filepath string) (map[string]interface{}, error) {
	valCache := GetCache(filepath)
	if valCache != nil {
		return valCache, nil
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		loggerStore.PrintError(err)
		return nil, err
	}
	var data map[string]interface{}
	e1 := json.Unmarshal(content, &data)
	if e1 != nil {
		loggerStore.PrintError(e1)
		return nil, e1
	}

	SetCache(filepath, data)
	return data, nil
}

func isFileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Save save data to file
func Save(filepath string, data map[string]interface{}) error {
	if data == nil {
		return errors.New("data is nil")
	}
	content, e := json.Marshal(data)
	if e != nil {
		loggerStore.PrintError(e)
		return e
	}

	if len(content) == 0 || string(content) == "{}" {
		return fmt.Errorf("data [%s] is empty", string(content))
	}
	dirname := path.Dir(filepath)
	if !isFileExists(dirname) {
		os.MkdirAll(dirname, 0777)
	}
	e1 := ioutil.WriteFile(filepath, content, 0777)
	if e1 != nil {
		loggerStore.PrintError(e1)
		return e
	}

	SetCache(filepath, data)
	return nil
}
