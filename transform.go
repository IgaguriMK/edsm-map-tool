package main

import (
	"./sysCoord"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	input_file_name := flag.String("i", "coords.bin", "input file")
	output_file_name := flag.String("o", "trans.bin", "output file")

	flag.Parse()

	commands := flag.Args()

	fmt.Println(*input_file_name, ">>>", *output_file_name)

	coords := sysCoord.LoadCoords(*input_file_name)
	fmt.Println("Loaded.")

	for len(commands) > 0 {
		switch command := pop(&commands); command {
		case ":cut-x":
			min := pop(&commands)
			max := pop(&commands)
			command_cut(&coords, min, max, func(c sysCoord.Coord) float32 { return c.X })
			fmt.Println(command, min, max)
		case ":cut-y":
			min := pop(&commands)
			max := pop(&commands)
			command_cut(&coords, min, max, func(c sysCoord.Coord) float32 { return c.Y })
			fmt.Println(command, min, max)
		case ":cut-z":
			min := pop(&commands)
			max := pop(&commands)
			command_cut(&coords, min, max, func(c sysCoord.Coord) float32 { return c.Z })
			fmt.Println(command, min, max)
		default:
			fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n", command)
		}
	}

	fmt.Println("Writing...")
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

func command_cut(coords *[]sysCoord.Coord, min_str, max_str string, axis func(sysCoord.Coord) float32) {
	min_d, err_mi := strconv.ParseFloat(min_str, 32)
	max_d, err_ma := strconv.ParseFloat(max_str, 32)
	if err_mi != nil || err_ma != nil {
		fmt.Fprintln(os.Stderr, "Error(:cut-*): invalid argument")
		os.Exit(1)
	}

	min, max := float32(min_d), float32(max_d)

	filtered := make([]sysCoord.Coord, 0)

	for _, c := range *coords {
		a := axis(c)
		if min <= a && a <= max {
			filtered = append(filtered, c)
		}
	}

	*coords = filtered
}
