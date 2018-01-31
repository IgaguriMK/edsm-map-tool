package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"./sysCoord"
)

const (
	DumpTimeFormat = "2006-01-02 15:04:05"
)

func main() {
	outFileName := flag.String("o", "coords.bin", "output file")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no argument.")
		os.Exit(1)
	}

	fileName := args[0]

	systems, err := sysCoord.LoadSystems(fileName)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: input error\n    %s", err)
		os.Exit(2)
	}

	coords := make([]sysCoord.Coord, 0, len(systems))
	for _, sys := range systems {
		utc, err := time.ParseInLocation(DumpTimeFormat, sys.Date, time.UTC)
		if err != nil {
			log.Fatal(err)
		}

		coord := sysCoord.Coord{
			X:    sys.Coord.X,
			Y:    sys.Coord.Y,
			Z:    sys.Coord.Z,
			Date: utc.Unix(),
		}
		coords = append(coords, coord)
	}

	sysCoord.WriteCoords(*outFileName, coords)
}
