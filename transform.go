package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	sw "github.com/IgaguriMK/allStarMap/stopwatch"
	"github.com/IgaguriMK/allStarMap/sysCoord"
)

const (
	UITimeFormat      = "2006-01-02_15:04:05"
	UITimeFormatShort = "2006-01-02"
)

func main() {
	sw.StartTier("START main()")
	defer sw.Close("END main()")

	input_file_name := flag.String("i", "coords.bin", "input file")
	output_file_name := flag.String("o", "trans.bin", "output file")

	flag.Parse()

	commands := flag.Args()

	sw.Mark("Flag parse")

	fmt.Println(*input_file_name, ">>>", *output_file_name)

	coords := sysCoord.LoadCoords(*input_file_name)

	sw.Mark("Coords load")

	for len(commands) > 0 {
		switch command := pop(&commands); command {
		case ":cut-x":
			min := pop(&commands)
			max := pop(&commands)
			coords = command_cut(coords, min, max, getX)
			fmt.Println(command, min, max)
			sw.Mark("cut-x")
		case ":cut-y":
			min := pop(&commands)
			max := pop(&commands)
			coords = command_cut(coords, min, max, getY)
			fmt.Println(command, min, max)
			sw.Mark("cut-y")
		case ":cut-z":
			min := pop(&commands)
			max := pop(&commands)
			coords = command_cut(coords, min, max, getZ)
			fmt.Println(command, min, max)
			sw.Mark("cut-z")
		case ":add":
			x := pop(&commands)
			y := pop(&commands)
			z := pop(&commands)
			command_add(&coords, x, y, z)
			fmt.Println(command, x, y, z)
			sw.Mark("add")
		case ":after":
			date := pop(&commands)
			command_after(&coords, date)
			fmt.Println(command, date)
			sw.Mark("after")
		case ":before":
			date := pop(&commands)
			command_before(&coords, date)
			fmt.Println(command, date)
			sw.Mark("before")
		case ":print":
			command_print(&coords)
			fmt.Println(command)
			sw.Mark("print")
		case ":exit":
			sw.Close("Exit")
			os.Exit(0)
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

func command_cut(coords []sysCoord.Coord, min_str, max_str string, getC func(sysCoord.Coord) float32) []sysCoord.Coord {
	min_d, err_mi := strconv.ParseFloat(min_str, 32)
	max_d, err_ma := strconv.ParseFloat(max_str, 32)
	if err_mi != nil || err_ma != nil {
		fmt.Fprintln(os.Stderr, "Error(:cut-*): invalid argument")
		os.Exit(1)
	}

	min, max := float32(min_d), float32(max_d)
	return normalCut(coords, min, max, getC)
}

func normalCut(coords []sysCoord.Coord, min, max float32, getC func(sysCoord.Coord) float32) []sysCoord.Coord {
	filtered := make([]sysCoord.Coord, 0, len(coords))

	for _, c := range coords {
		v := getC(c)
		if min <= v && v <= max {
			filtered = append(filtered, c)
		}
	}

	return filtered
}

func command_add(coords *[]sysCoord.Coord, xs, ys, zs string) {
	x, err_x := strconv.ParseFloat(xs, 32)
	y, err_y := strconv.ParseFloat(ys, 32)
	z, err_z := strconv.ParseFloat(zs, 32)
	if err_x != nil || err_y != nil || err_z != nil {
		fmt.Fprintln(os.Stderr, "Error(:add): invalid argument")
		os.Exit(1)
	}

	c := sysCoord.Coord{float32(x), float32(y), float32(z), 0}
	*coords = append(*coords, c)
}

func command_print(coords *[]sysCoord.Coord) {
	for _, c := range *coords {
		fmt.Printf("% 6.2f, % 6.2f, % 6.2f\n", c.X, c.Y, c.Z)
	}
}

func command_after(coords *[]sysCoord.Coord, date string) {
	thres := getThres(date)

	filtered := make([]sysCoord.Coord, 0, len(*coords))
	for _, c := range *coords {
		if c.Date >= thres {
			filtered = append(filtered, c)
		}
	}

	log.Println("Hit:", len(filtered))
	*coords = filtered
}

func command_before(coords *[]sysCoord.Coord, date string) {
	thres := getThres(date)

	filtered := make([]sysCoord.Coord, 0, len(*coords))

	for _, c := range *coords {
		if c.Date < thres {
			filtered = append(filtered, c)
		}
	}

	log.Println("Hit:", len(filtered))
	*coords = filtered
}

func getThres(date string) int64 {
	if len(date) == len(UITimeFormatShort) {
		date = date + "_00:00:00"
	}

	thresDate, err := time.ParseInLocation(UITimeFormat, date, time.UTC)
	if err != nil {
		log.Fatalf("Invalid date format[%s]: %e\n", date, err)
	}

	return thresDate.Unix()
}

func getX(c sysCoord.Coord) float32 { return c.X }
func getY(c sysCoord.Coord) float32 { return c.Y }
func getZ(c sysCoord.Coord) float32 { return c.Z }
