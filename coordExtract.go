package main

import (
	"./systemCoordinate"
	"bytes"
	"encoding/binary"
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

	systems, err := systemCoordinate.LoadSystems(fileName)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: input error\n    %s", err)
		os.Exit(2)
	}

	writeTo(*outFileName, systems)
}

func writeTo(fileName string, systems []systemCoordinate.SystemCoord) {
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannnot open output file.\n    %s", err)
		os.Exit(3)
	}
	defer outFile.Close()

	for _, system := range systems {
		coord := system.Coord

		writeBytes(outFile, toBytes(coord.X))
		writeBytes(outFile, toBytes(coord.Y))
		writeBytes(outFile, toBytes(coord.Z))
	}
}

func toBytes(val float32) []byte {
	buf := new(bytes.Buffer)

	err_b := binary.Write(buf, binary.LittleEndian, val)
	if err_b != nil {
		fmt.Fprintf(os.Stderr, "Error: converting to binary\n    %s", err_b)
		os.Exit(4)
	}

	return buf.Bytes()
}

func writeBytes(outFile *os.File, bytes []byte) {
	_, err_w := outFile.Write(bytes)
	if err_w != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to write\n    %s", err_w)
		os.Exit(4)
	}
}
