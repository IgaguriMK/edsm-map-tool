package main

import (
	"./sysCoord"
	"flag"
	"fmt"
	"os"
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
		coords[i] = systems[i].Coord
	}

	sysCoord.WriteCoords(*outFileName, coords)
}
