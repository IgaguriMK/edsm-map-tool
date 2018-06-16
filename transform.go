package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/IgaguriMK/allStarMap/sysCoord"
)

const (
	uiTimeFormat      = "2006-01-02_15:04:05"
	uiTimeFormatShort = "2006-01-02"
)

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "i", "coords.bin", "input file")
	var outputFile string
	flag.StringVar(&outputFile, "o", "trans.bin", "output file")

	flag.Parse()

	fmt.Println(inputFile, ">>>", outputFile)

	var wg sync.WaitGroup
	ch := sysCoord.LoadCoords(inputFile)
	outCh := sysCoord.WriteCoords(outputFile, &wg)

	parser := commandParser{
		":cut-x":  cut{getX},
		":cut-y":  cut{getY},
		":cut-z":  cut{getZ},
		":add":    add{},
		":after":  after{},
		":before": before{},
	}

	args := newArgList()
	ch = parser.Exec(ch, args)

	go func() {
		for c := range ch {
			outCh <- c
		}
		close(outCh)
	}()

	wg.Wait()
}

type commandParser map[string]command

func (c *commandParser) Exec(ch <-chan sysCoord.Coord, args *argList) <-chan sysCoord.Coord {
	for !args.Empty() {
		commandName := args.PopString()

		command, ok := (*c)[commandName]
		if !ok {
			log.Fatalf("Unknown command %q", commandName)
		}

		ch = command.Filter(ch, args)
	}

	return ch
}

type argList struct {
	args []string
}

func newArgList() *argList {
	return &argList{
		args: flag.Args(),
	}
}

func (a *argList) Empty() bool {
	return len(a.args) == 0
}

func (a *argList) PopString() string {
	if a.Empty() {
		log.Fatal("Too few args")
	}
	r := a.args[0]
	a.args = a.args[1:]
	return r
}

func (a *argList) PopFloat32() float32 {
	str := a.PopString()
	f, err := strconv.ParseFloat(str, 32)
	if err != nil {
		log.Fatalf("Can't Parse %q to float32", str)
	}

	return float32(f)
}

func (a *argList) PopUnix() int64 {
	str := a.PopString()

	if len(str) == len(uiTimeFormatShort) {
		str = str + "_00:00:00"
	}

	t, err := time.ParseInLocation(uiTimeFormat, str, time.UTC)
	if err != nil {
		log.Fatalf("Can't parse date %q to time: %s", str, err)
	}

	return t.Unix()
}

type command interface {
	Filter(ch <-chan sysCoord.Coord, argList *argList) <-chan sysCoord.Coord
}

type cut struct {
	Axis func(sysCoord.Coord) float32
}

func (cc cut) Filter(ch <-chan sysCoord.Coord, args *argList) <-chan sysCoord.Coord {
	min := args.PopFloat32()
	max := args.PopFloat32()

	filtered := make(chan sysCoord.Coord, sysCoord.StreamBufferSize)

	go func() {
		for c := range ch {
			v := cc.Axis(c)
			if min <= v && v <= max {
				filtered <- c
			}
		}
		close(filtered)
	}()

	return filtered
}

type add struct{}

func (a add) Filter(ch <-chan sysCoord.Coord, args *argList) <-chan sysCoord.Coord {
	x := args.PopFloat32()
	y := args.PopFloat32()
	z := args.PopFloat32()

	added := make(chan sysCoord.Coord, sysCoord.StreamBufferSize)

	go func() {
		added <- sysCoord.Coord{x, y, z, 0}

		for c := range ch {
			added <- c
		}
		close(added)
	}()

	return added
}

type after struct{}

func (a after) Filter(ch <-chan sysCoord.Coord, args *argList) <-chan sysCoord.Coord {
	thres := args.PopUnix()

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

type before struct{}

func (b before) Filter(ch <-chan sysCoord.Coord, args *argList) <-chan sysCoord.Coord {
	thres := args.PopUnix()

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

func getX(c sysCoord.Coord) float32 { return c.X }
func getY(c sysCoord.Coord) float32 { return c.Y }
func getZ(c sysCoord.Coord) float32 { return c.Z }
