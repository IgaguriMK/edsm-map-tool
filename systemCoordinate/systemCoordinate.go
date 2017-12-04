package systemCoordinate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Coordinate struct {
	X float32 "json:x"
	Y float32 "json:y"
	Z float32 "json:z"
}

type SystemCoordinate struct {
	//Id         uint32     `json:"id"`
	//Id64       uint64     `json:"id64"`
	//Name       string     `json:"name"`
	Coordinate Coordinate `json:"coords"`
	//Date       string     `json:"date"`
}

func LoadSystems(fileName string) ([]SystemCoordinate, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("IO error: %s", err)
	}

	var systemCoordinates []SystemCoordinate
	if err := json.Unmarshal(bytes, &systemCoordinates); err != nil {
		return nil, fmt.Errorf("JSON error: %s", err)
	}

	return systemCoordinates, nil
}
