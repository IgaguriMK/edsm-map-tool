package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/IgaguriMK/allStarMap/sysCoord"
)

const (
	UITimeFormat      = "2006-01-02_15:04:05"
	UITimeFormatShort = "2006-01-02"
)

func main() {
	input_file_name := flag.String("i", "coords.bin", "input file")
	output_file_name := flag.String("o", "trans.bin", "output file")

	flag.Parse()

	commands := flag.Args()

	fmt.Println(*input_file_name, ">>>", *output_file_name)

	var wg sync.WaitGroup

	ch := sysCoord.LoadCoords(*input_file_name)
	outCh := sysCoord.WriteCoords(*output_file_name, &wg)

	for len(commands) > 0 {
		switch command := pop(&commands); command {
		case ":cut-x":
			min := pop(&commands)
			max := pop(&commands)
			ch = command_cut(ch, min, max, getX)
			fmt.Println(command, min, max)
		case ":cut-y":
			min := pop(&commands)
			max := pop(&commands)
			ch = command_cut(ch, min, max, getY)
			fmt.Println(command, min, max)
		case ":cut-z":
			min := pop(&commands)
			max := pop(&commands)
			ch = command_cut(ch, min, max, getZ)
			fmt.Println(command, min, max)
		case ":add":
			x := pop(&commands)
			y := pop(&commands)
			z := pop(&commands)
			ch = command_add(ch, x, y, z)
			fmt.Println(command, x, y, z)
		case ":after":
			date := pop(&commands)
			ch = command_after(ch, date)
			fmt.Println(command, date)
		case ":before":
			date := pop(&commands)
			ch = command_before(ch, date)
			fmt.Println(command, date)
		default:
			fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n", command)
		}
	}

	go func() {
		for c := range ch {
			outCh <- c
		}
		close(outCh)
	}()

	wg.Wait()
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

func command_cut(ch <-chan sysCoord.Coord, min_str, max_str string, getC func(sysCoord.Coord) float32) <-chan sysCoord.Coord {
	min_d, err_mi := strconv.ParseFloat(min_str, 32)
	max_d, err_ma := strconv.ParseFloat(max_str, 32)
	if err_mi != nil || err_ma != nil {
		fmt.Fprintln(os.Stderr, "Error(:cut-*): invalid argument")
		os.Exit(1)
	}

	min, max := float32(min_d), float32(max_d)

	filtered := make(chan sysCoord.Coord, sysCoord.StreamBufferSize)

	go func() {
		for c := range ch {
			v := getC(c)
			if min <= v && v <= max {
				filtered <- c
			}
		}
		close(filtered)
	}()

	return filtered
}

func command_add(ch <-chan sysCoord.Coord, xs, ys, zs string) <-chan sysCoord.Coord {
	x, err_x := strconv.ParseFloat(xs, 32)
	y, err_y := strconv.ParseFloat(ys, 32)
	z, err_z := strconv.ParseFloat(zs, 32)
	if err_x != nil || err_y != nil || err_z != nil {
		fmt.Fprintln(os.Stderr, "Error(:add): invalid argument")
		os.Exit(1)
	}

	added := make(chan sysCoord.Coord, sysCoord.StreamBufferSize)

	go func() {
		added <- sysCoord.Coord{float32(x), float32(y), float32(z), 0}

		for c := range ch {
			added <- c
		}
		close(added)
	}()

	return added
}

func command_after(ch <-chan sysCoord.Coord, date string) <-chan sysCoord.Coord {
	thres := getThres(date)

	filtered := make(chan sysCoord.Coord, sysCoord.StreamBufferSize)

	go func() {
		cnt := 0
		for c := range ch {
			if c.Date >= thres {
				filtered <- c
				cnt++
			}
		}
		log.Println("Hit:", cnt)
		close(filtered)
	}()

	return filtered
}

func command_before(ch <-chan sysCoord.Coord, date string) <-chan sysCoord.Coord {
	thres := getThres(date)

	filtered := make(chan sysCoord.Coord, sysCoord.StreamBufferSize)

	go func() {
		cnt := 0
		for c := range ch {
			if c.Date < thres {
				filtered <- c
				cnt++
			}
		}
		log.Println("Hit:", cnt)
		close(filtered)
	}()

	return filtered
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
