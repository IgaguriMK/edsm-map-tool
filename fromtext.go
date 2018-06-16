package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/IgaguriMK/allStarMap/sysCoord"
)

func main() {
	var inputFileName string
	flag.StringVar(&inputFileName, "i", "", "input file")
	var outputFileName string
	flag.StringVar(&outputFileName, "o", "trans.bin", "output file")

	flag.Parse()

	if inputFileName == "" {
		log.Fatal("No input file.")
	}

	f, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal("Can't open input file: ", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	if !sc.Scan() {
		log.Fatal("No input line.")
	}

	var wg sync.WaitGroup
	writeCh := sysCoord.WriteCoords(outputFileName, &wg)

	lineNum := 1
	for sc.Scan() {
		line := sc.Text()
		lineNum++

		fields := strings.Split(line, "\t")

		utc, err := time.ParseInLocation(sysCoord.DumpTimeFormat, fields[3], time.UTC)
		if err != nil {
			log.Fatalf("line [%d] %s", lineNum, err)
		}

		writeCh <- sysCoord.Coord{
			X:    float32(parseFloat(fields[0])),
			Y:    float32(parseFloat(fields[1])),
			Z:    float32(parseFloat(fields[2])),
			Date: utc.Unix(),
		}
	}
	close(writeCh)

	wg.Wait()
}

func parseFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatal("Float parse error: ", str)
	}
	return f
}
