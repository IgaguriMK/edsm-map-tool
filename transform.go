package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"

	"./sysCoord"
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
			command_cut(&coords, min, max, getX)
			fmt.Println(command, min, max)
		case ":cut-y":
			min := pop(&commands)
			max := pop(&commands)
			command_cut(&coords, min, max, getY)
			fmt.Println(command, min, max)
		case ":cut-z":
			min := pop(&commands)
			max := pop(&commands)
			command_cut(&coords, min, max, getZ)
			fmt.Println(command, min, max)
		case ":add":
			x := pop(&commands)
			y := pop(&commands)
			z := pop(&commands)
			command_add(&coords, x, y, z)
			fmt.Println(command, x, y, z)
		case ":sort":
			axis := pop(&commands)
			command_sort(&coords, axis)
			fmt.Println(command, axis)
		case ":print":
			command_print(&coords)
			fmt.Println(command)
		case ":exit":
			os.Exit(0)
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

func command_cut(coords *[]sysCoord.Coord, min_str, max_str string, getC func(sysCoord.Coord) float32) {
	min_d, err_mi := strconv.ParseFloat(min_str, 32)
	max_d, err_ma := strconv.ParseFloat(max_str, 32)
	if err_mi != nil || err_ma != nil {
		fmt.Fprintln(os.Stderr, "Error(:cut-*): invalid argument")
		os.Exit(1)
	}

	min, max := float32(min_d), float32(max_d)
	normalCut(coords, min, max, getC)
}

func normalCut(coords *[]sysCoord.Coord, min, max float32, getC func(sysCoord.Coord) float32) {
	filtered := make([]sysCoord.Coord, 0)

	for _, c := range *coords {
		v := getC(c)
		if min <= v && v <= max {
			filtered = append(filtered, c)
		}
	}

	*coords = filtered
}

func command_add(coords *[]sysCoord.Coord, xs, ys, zs string) {
	x, err_x := strconv.ParseFloat(xs, 32)
	y, err_y := strconv.ParseFloat(ys, 32)
	z, err_z := strconv.ParseFloat(zs, 32)
	if err_x != nil || err_y != nil || err_z != nil {
		fmt.Fprintln(os.Stderr, "Error(:add): invalid argument")
		os.Exit(1)
	}

	c := sysCoord.Coord{float32(x), float32(y), float32(z)}
	*coords = append(*coords, c)
}

func command_sort(coords *[]sysCoord.Coord, axis string) {
	var getC func(sysCoord.Coord) float32
	switch axis {
	case "x":
		getC = getX
	case "y":
		getC = getY
	case "z":
		getC = getZ
	default:
		fmt.Fprintln(os.Stderr, "Error(:sort): invalid argument")
		os.Exit(1)
	}

	ensureSorted(coords, getC)
}

func command_print(coords *[]sysCoord.Coord) {
	for _, c := range *coords {
		fmt.Printf("% 6.2f, % 6.2f, % 6.2f\n", c.X, c.Y, c.Z)
	}
}

func isSorted(coords *[]sysCoord.Coord, getC func(sysCoord.Coord) float32) bool {
	return sort.SliceIsSorted(*coords, func(i, j int) bool {
		return getC((*coords)[i]) < getC((*coords)[j])
	})
}

func ensureSorted(coords *[]sysCoord.Coord, getC func(sysCoord.Coord) float32) {
	if isSorted(coords, getC) {
		return
	}

	sort.Slice(*coords, func(i, j int) bool {
		return getC((*coords)[i]) < getC((*coords)[j])
	})
}

func getX(c sysCoord.Coord) float32 { return c.X }
func getY(c sysCoord.Coord) float32 { return c.Y }
func getZ(c sysCoord.Coord) float32 { return c.Z }
