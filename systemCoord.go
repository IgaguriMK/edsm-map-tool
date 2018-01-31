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
	UITimeFormat   = "2006-01-02_15:04:05"
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

	coords := make([]sysCoord.Coord, len(systems))
	for i := 0; i < len(systems); i++ {
		utc, err := time.ParseInLocation(DumpTimeFormat, systems[i].Date, time.UTC)
		if err != nil {
			log.Fatal(err)
		}

		coords[i].X = systems[i].Coord.X
		coords[i].Y = systems[i].Coord.Y
		coords[i].Z = systems[i].Coord.Z
		coords[i].Date = utc.Unix()
	}

	sysCoord.WriteCoords(*outFileName, coords)
}
