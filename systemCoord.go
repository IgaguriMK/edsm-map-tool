package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/IgaguriMK/allStarMap/sysCoord"
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

	var wg sync.WaitGroup

	sysCh := sysCoord.LoadSystems(fileName)
	writeCh := sysCoord.WriteCoords(*outFileName, &wg)

	for sys := range sysCh {
		utc, err := time.ParseInLocation(sysCoord.DumpTimeFormat, sys.Date, time.UTC)
		if err != nil {
			log.Fatal(err)
		}

		coord := sysCoord.Coord{
			X:    sys.Coord.X,
			Y:    sys.Coord.Y,
			Z:    sys.Coord.Z,
			Date: utc.Unix(),
		}
		writeCh <- coord
	}
	close(writeCh)

	wg.Wait()
}
