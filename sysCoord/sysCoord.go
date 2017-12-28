package sysCoord

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
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

func LoadCoords(file_name string) []Coord {
	var coords []Coord

	file, err_f := os.Open(file_name)
	if err_f != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannnot open input file.\n    %s\n", err_f)
		os.Exit(1)
	}
	defer file.Close()

	buffer := make([]byte, 4*3)

	for {
		read_size, err_r := io.ReadFull(file, buffer)
		if err_r == io.EOF {
			break
		}
		if err_r != nil {
			fmt.Fprintf(os.Stderr, "Error: Cannnot read from file.\n    %s\n", err_r)
			os.Exit(1)
		}
		if read_size < 4*3 {
			fmt.Fprint(os.Stderr, "Error: read too few bytes.")
			os.Exit(1)
		}

		var coord Coord
		coord.X = decodeFloat32(buffer[0:4])
		coord.Y = decodeFloat32(buffer[4:8])
		coord.Z = decodeFloat32(buffer[8:12])

		coords = append(coords, coord)
	}

	return coords
}

func decodeFloat32(raw []byte) float32 {
	var val float32
	buf := bytes.NewReader(raw)
	err := binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		panic(err)
	}
	return val
}
