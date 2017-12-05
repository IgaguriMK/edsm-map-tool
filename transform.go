package main

import (
	"./sysCoord"
	"flag"
	"fmt"
	"os"
)

func main() {
	input_file_name := flag.String("i", "coords.bin", "input file")
	output_file_name := flag.String("o", "trans.bin", "output file")

	flag.Parse()

	commands := flag.Args()

	fmt.Println(*input_file_name, ">>>", *output_file_name)

	coords := sysCoord.LoadCoords(*input_file_name)

	for len(commands) > 0 {
		switch command := pop(&commands); command {
		default:
			fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n", command)
		}
	}

	sysCoord.WriteCoords(*output_file_name, coords)
}

func pop(arr *[]string) string {
	if len(*arr) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no argument")
		os.Exit(1)
	}

	v := (*arr)[0]
	*arr = (*arr)[1:]
	return v
}
