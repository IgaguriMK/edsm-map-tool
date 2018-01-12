package sysCoord

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	sw "github.com/IgaguriMK/allStarMap/stopwatch"
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
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("IO error: %s\n", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	systemCoords := make([]SystemCoord, 0)

	for sc.Scan() {
		line := sc.Text()

		if line == "[" {
			continue
		}
		if line == "]" {
			break
		}

		line = strings.TrimPrefix(line, "    ")
		line = strings.TrimSuffix(line, ",")
		bytes := []byte(line)

		var system SystemCoord
		if err := json.Unmarshal(bytes, &system); err != nil {
			return nil, fmt.Errorf("JSON error: %s\n", err)
		}
		systemCoords = append(systemCoords, system)
	}

	return systemCoords, nil
}

func WriteCoords(fileName string, coords []Coord) {
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannnot open output file.\n    %s\n", err)
		os.Exit(3)
	}
	defer outFile.Close()

	binary.Write(outFile, binary.LittleEndian, coords)
}

func LoadCoords(fileName string) []Coord {
	sw.StartTier(`START LoadCoords(` + fileName + `)`)
	defer sw.StartTier(`END LoadCoords(` + fileName + `)`)

	fInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannnot get file size.\n    %s\n", err)
		os.Exit(1)
	}
	coordCount := fInfo.Size() / (4 * 3)
	coords := make([]Coord, coordCount)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannnot open input file.\n    %s\n", err)
		os.Exit(1)
	}
	defer file.Close()
	sw.Mark("Open file")

	binary.Read(file, binary.LittleEndian, coords)

	return coords
}
