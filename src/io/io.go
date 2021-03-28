package io

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func SaveConfig(filename string, data interface{}) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("../config/%s.json", filename), file, 0644)
}
