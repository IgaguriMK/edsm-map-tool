package sysCoord

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
	"unsafe"
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

func LoadSystems(fileName string) <-chan SystemCoord {
	ch := make(chan SystemCoord, StreamBufferSize)

	go func() {

		f, err := os.Open(fileName)
		if err != nil {
			log.Fatal("IO error:", err)
		}
		defer f.Close()

		dec := json.NewDecoder(f)

		_, err = dec.Token()
		if err != nil {
			log.Fatal(err)
		}

		for dec.More() {
			var system SystemCoord
			err := dec.Decode(&system)
			if err != nil {
				log.Fatal("JSON error:", err)
			}

			ch <- system
		}

		_, err = dec.Token()
		if err != nil {
			log.Fatal(err)
		}

		close(ch)
	}()

	return ch
}

func WriteCoords(fileName string, wg *sync.WaitGroup) chan<- Coord {
	ch := make(chan Coord, StreamBufferSize)
	wg.Add(1)

	go func() {
		defer wg.Done()

		outFile, err := os.Create(fileName)
		if err != nil {
			log.Fatal("Error: Cannnot open output file.\n    %s\n", err)
		}
		defer outFile.Close()

		for c := range ch {
			binary.Write(outFile, binary.LittleEndian, c)
		}
	}()

	return ch
}

func LoadCoords(fileName string) <-chan Coord {
	ch := make(chan Coord, StreamBufferSize)

	go func() {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal("Error: Cannnot open input file.\n    %s\n", err)
			os.Exit(1)
		}
		defer file.Close()

		for {
			var coord Coord
			err := binary.Read(file, binary.LittleEndian, &coord)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			ch <- coord
		}
		close(ch)
	}()

	return ch
}

func LoadCoordsAll(fileName string) []Coord {
	fInfo, err := os.Stat(fileName)
	if err != nil {
		log.Fatal("Error: Cannnot get file size:", err)
	}
	coordCount := fInfo.Size() / int64(unsafe.Sizeof(Coord{}))
	coords := make([]Coord, coordCount)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error: Cannnot open input file:", err)
	}
	defer file.Close()

	binary.Read(file, binary.LittleEndian, coords)

	return coords
}
