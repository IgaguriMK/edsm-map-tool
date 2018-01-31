package sysCoord

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unsafe"

	sw "github.com/IgaguriMK/allStarMap/stopwatch"
)

const (
	StreamBufferSize = 1024
)

type Coord struct {
	X    float32
	Y    float32
	Z    float32
	Date int64
}

type SystemCoord struct {
	//Id         uint32     `json:"id"`
	//Id64       uint64     `json:"id64"`
	//Name       string     `json:"name"`
	Coord Coord  `json:"coords"`
	Date  string `json:"date"`
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
	sw.StartTier(`START WriteCoords(` + fileName + `)`)
	defer sw.EndTier(`END WriteCoords(` + fileName + `)`)

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
	defer sw.EndTier(`END LoadCoords(` + fileName + `)`)

	fInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannnot get file size.\n    %s\n", err)
		os.Exit(1)
	}
	coordCount := fInfo.Size() / int64(unsafe.Sizeof(Coord{}))
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

func NewCoordStream() (*CoordWriter, *CoordReader) {
	return NewCoordStreamCap(StreamBufferSize)
}

func NewCoordStreamCap(cap int) (*CoordWriter, *CoordReader) {
	ch := make(chan Coord, cap)

	cw := &CoordWriter{
		ch: ch,
	}
	cr := &CoordReader{
		ch:  ch,
		buf: nil,
	}

	return cw, cr
}

type CoordWriter struct {
	ch chan<- Coord
}

func (cw *CoordWriter) Write(coord Coord) {
	cw.ch <- coord
}

func (cw *CoordWriter) Close() {
	close(cw.ch)
}

type CoordReader struct {
	ch  <-chan Coord
	buf *Coord
}

func (cr *CoordReader) Next() bool {
	c, ok := <-cr.ch
	if !ok {
		return false
	}

	cr.buf = &c
	return true
}

func (cr *CoordReader) Read() Coord {
	if cr.buf == nil {
		panic("Read() must called after Next()")
	}

	c := *cr.buf
	cr.buf = nil
	return c
}
