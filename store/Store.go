package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"stone/logger"
	"strconv"
)

var loggerStore = logger.GetPrefixLogger("store")

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// GetData read data from file
func GetData(filepath string) (map[string]interface{}, error) {
	valCache := GetCache(filepath)
	if valCache != nil {
		loggerStore.PrintInfof("%d %s from cache", getGID(), filepath)
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

	loggerStore.PrintInfof("%d %s from file", getGID(), filepath)
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
	dataOld, _ := GetData(filepath)
	if dataOld != nil {
		for key, val := range data {
			dataOld[key] = val
		}
	} else {
		dataOld = data
	}

	content, e := json.Marshal(dataOld)
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

	SetCache(filepath, dataOld)
	return nil
}

// Del delete date
func Del(filepath string) error {
	if isFileExists(filepath) {
		DeleCache(filepath)
		return os.Remove(filepath)
	}
	return nil
}
