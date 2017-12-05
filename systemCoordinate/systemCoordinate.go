package systemCoordinate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Coord struct {
	X float32 "json:x"
	Y float32 "json:y"
	Z float32 "json:z"
}

type SystemCoord struct {
	//Id         uint32     `json:"id"`
	//Id64       uint64     `json:"id64"`
	//Name       string     `json:"name"`
	Coord Coord `json:"coords"`
	//Date       string     `json:"date"`
}

func LoadSystems(fileName string) ([]SystemCoord, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("IO error: %s", err)
	}

	var systemCoords []SystemCoord
	if err := json.Unmarshal(bytes, &systemCoords); err != nil {
		return nil, fmt.Errorf("JSON error: %s", err)
	}

	return systemCoords, nil
}
