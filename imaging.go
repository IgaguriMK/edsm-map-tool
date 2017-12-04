package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
)

type Coord struct {
	X float32
	Y float32
	Z float32
}

type Plane int

const (
	XZ = iota
	ZY
	XY
)

func main() {
	coords_file_name := flag.String("i", "coords.bin", "input file")
	image_file_name := flag.String("o", "xz.png", "output file")
	plane_name := flag.String("p", "xz", "dump plane (xz, zy, xy)")
	var chunk_size int
	flag.IntVar(&chunk_size, "s", 20, "pixcel size")
	var heat_scale float64
	flag.Float64Var(&heat_scale, "h", 4, "heatmap scale")

	flag.Parse()

	var plane Plane
	switch *plane_name {
	case "xz":
		plane = XZ
	case "zy":
		plane = ZY
	case "xy":
		plane = XY
	default:
		fmt.Fprintf(os.Stderr, "Error: Invalid plane name'%s'", *plane_name)
		os.Exit(1)
	}

	coords := loadCoords(*coords_file_name)
	max, min := maxMin(coords)

	s_max, t_max := getPosByPlane(plane, chunk_size, max)
	s_min, t_min := getPosByPlane(plane, chunk_size, min)

	s_size := s_max - s_min + 1
	t_size := t_max - t_min + 1

	lines := make([][]float64, t_size)
	for t := 0; t < t_size; t++ {
		lines[t] = make([]float64, s_size)
	}

	for _, coord := range coords {
		s, t := getPosByPlane(plane, chunk_size, coord)
		s -= s_min
		t -= t_min
		lines[t][s] += 1
	}

	var v_max float64 = 0.0
	for t := 0; t < t_size; t++ {
		for s := 0; s < s_size; s++ {
			v := curve(lines[t][s])
			lines[t][s] = v
			if v > v_max {
				v_max = v
			}
		}
	}
	fmt.Println("Heat max:", v_max)
	if v_max > heat_scale {
		fmt.Println("Warning: Heat scale adjusted to heat max.")
		heat_scale = v_max
	}

	img := image.NewRGBA(image.Rect(0, 0, s_size, t_size))

	for t := 0; t < t_size; t++ {
		for s := 0; s < s_size; s++ {
			v := lines[t][s] / heat_scale

			if v > 0 {
				r := uint8(255 * v * v)
				g := uint8(255 * (1 - 4*(v-0.5)*(v-0.5)))
				b := uint8(255 * (1 - v*v))
				img.Set(s, t_size-t, color.RGBA{r, g, b, 255})
				continue
			}

			var a uint8
			so := s + s_min
			to := t + t_min
			switch 0 {
			case (so % (10000 / chunk_size)) * (to % (10000 / chunk_size)):
				a = 192
			case (so % (5000 / chunk_size)) * (to % (5000 / chunk_size)):
				a = 128
			case (so % (1000 / chunk_size)) * (to % (1000 / chunk_size)):
				a = 92
			case (so % (500 / chunk_size)) * (to % (500 / chunk_size)):
				a = 80
			case (so % (100 / chunk_size)) * (to % (100 / chunk_size)):
				a = 72
			default:
				a = 64
			}
			img.Set(s, t_size-t, color.RGBA{0, 0, 0, a})
		}
	}

	img_file, err_f := os.Create(*image_file_name)
	if err_f != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannnot open file.\n    %s", err_f)
		os.Exit(1)
	}
	defer img_file.Close()

	if err_io := png.Encode(img_file, img); err_io != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannnot write image.\n    %s", err_io)
		os.Exit(1)
	}
}

func curve(val float64) float64 {
	return math.Log10(val + 1)
}

func getPosByPlane(plane Plane, chunk_size int, coord Coord) (int, int) {
	if plane == XZ {
		return chunk(chunk_size, coord.X), chunk(chunk_size, coord.Z)
	} else if plane == ZY {
		return chunk(chunk_size, coord.Z), chunk(chunk_size, coord.Y)
	}
	return chunk(chunk_size, coord.X), chunk(chunk_size, coord.Y)
}

func chunk(chunk_size int, val float32) int {
	return int(val / float32(chunk_size))
}

func maxMin(coords []Coord) (Coord, Coord) {
	var max, min Coord

	for _, c := range coords {
		if max.X < c.X {
			max.X = c.X
		}
		if max.Y < c.Y {
			max.Y = c.Y
		}
		if max.Z < c.Z {
			max.Z = c.Z
		}
		if min.X > c.X {
			min.X = c.X
		}
		if min.Y > c.Y {
			min.Y = c.Y
		}
		if min.Z > c.Z {
			min.Z = c.Z
		}
	}

	return max, min
}

func loadCoords(file_name string) []Coord {
	var coords []Coord

	file, err_f := os.Open(file_name)
	if err_f != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannnot open input file.\n    %s", err_f)
		os.Exit(1)
	}
	defer file.Close()

	buffer := make([]byte, 4*3)

	for {
		read_size, err_r := io.ReadFull(file, buffer)
		if err_r == io.EOF {
			break
		}
		if err_r != nil {
			fmt.Fprintf(os.Stderr, "Error: Cannnot read from file.\n    %s", err_r)
			os.Exit(1)
		}
		if read_size < 4*3 {
			fmt.Fprint(os.Stderr, "Error: read too few bytes.")
			os.Exit(1)
		}

		var coord Coord
		coord.X = decodeFloat32(buffer[0:4])
		coord.Y = decodeFloat32(buffer[4:8])
		coord.Z = decodeFloat32(buffer[8:12])

		coords = append(coords, coord)
	}

	return coords
}

func decodeFloat32(raw []byte) float32 {
	var val float32
	buf := bytes.NewReader(raw)
	err := binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		panic(err)
	}
	return val
}
