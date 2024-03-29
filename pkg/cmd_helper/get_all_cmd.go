package cmdhelper

import (
	"encoding/json"
	"io/ioutil"
)

func GetAllCmd(filename string) ([]Commands, error) {
	var cmd []Commands

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}