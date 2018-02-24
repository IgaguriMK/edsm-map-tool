package main

import (
	"./sysCoord"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type Plane int

const (
	XZ = iota
	ZY
	XY
)

func main() {
	coords_file_name := flag.String("i", "coords.bin", "input file")
	image_file_name := flag.String("o", "", "output file")
	plane_name := flag.String("p", "xz", "dump plane (xz, zy, xy)")
	var chunk_size int
	flag.IntVar(&chunk_size, "s", 20, "pixcel size in LY")
	curve_name := flag.String("hc", "log", "heatmap curve (liner, log)")
	heatmap_name := flag.String("ht", "opaque", "heatmap type (colorful, noback, opaque, hard)")
	var heat_scale float64
	flag.Float64Var(&heat_scale, "hs", 0.1, "heatmap scale")
	var no_adjust bool
	flag.BoolVar(&no_adjust, "hna", false, "disable heatmap scale adjust")
	var scale_bar bool
	flag.BoolVar(&scale_bar, "bar", false, "enable scale bar")
	var sizeAdjust int
	flag.IntVar(&sizeAdjust, "multof", 0, "set image size to multiple of arg (0 is disable)")

	flag.Parse()

	if *image_file_name == "" {
		*image_file_name = fmt.Sprintf("%s_%d.png", *plane_name, chunk_size)
	}

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

	var curve func(float64) float64
	switch *curve_name {
	case "log":
		curve = func(v float64) float64 { return math.Log10(v + 1) }
	case "liner":
		curve = func(v float64) float64 { return v }
	default:
		fmt.Fprintln(os.Stderr, "Error: invalid curve name")
		os.Exit(1)
	}

	var heatmap func(*image.RGBA, int, int, int, int, int, int, int, float64)
	switch *heatmap_name {
	case "colorful":
		heatmap = coloful_heatmap
	case "noback":
		heatmap = coloful_noback_heatmap
	case "opaque":
		heatmap = coloful_opaque_heatmap
	case "hard":
		heatmap = hard_heatmap
	default:
		fmt.Fprintln(os.Stderr, "Error: invalid heatmap name")
		os.Exit(1)
	}

	//////////////

	coords := make([]sysCoord.Coord, 0, 1024)
	for c := range sysCoord.LoadCoords(*coords_file_name) {
		coords = append(coords, c)
	}
	max, min := maxMin(coords)

	s_max, t_max := getPosByPlane(plane, chunk_size, max)
	s_min, t_min := getPosByPlane(plane, chunk_size, min)

	s_size := s_max - s_min + 4
	t_size := t_max - t_min + 4

	if sizeAdjust > 0 {
		if s_size%sizeAdjust != 0 {
			s_size += sizeAdjust - s_size%sizeAdjust
		}
		if t_size%sizeAdjust != 0 {
			t_size += sizeAdjust - t_size%sizeAdjust
		}
	}

	var scale_bar_size int = 0
	if scale_bar {
		if t_size < 128 {
			scale_bar_size = 4
		} else if t_size > 1024 {
			scale_bar_size = 32
		} else {
			scale_bar_size = t_size / 32
		}
	}

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
	if !no_adjust && v_max > heat_scale {
		fmt.Println("Heat scale adjusted to heat max.")
		heat_scale = v_max
	}

	img := image.NewRGBA(image.Rect(0, 0, s_size, t_size+scale_bar_size))

	for t := 0; t < t_size; t++ {
		for s := 0; s < s_size; s++ {
			v := lines[t][s] / heat_scale
			if v > 1.0 {
				v = 1.0
			}
			if v < 0.0 {
				v = 0.0
			}
			heatmap(img, s, t, s_size, t_size+scale_bar_size, s_max, t_max, chunk_size, v)
		}
	}

	for t := t_size; t < (t_size + scale_bar_size); t++ {
		for s := 0; s < s_size; s++ {
			if t == t_size {
				img.Set(s, t_size-t, color.RGBA{255, 255, 255, 255})
				continue
			}
			v := float64(s) / float64(s_size)
			heatmap(img, s, t, s_size, t_size+scale_bar_size, s_max, t_max, chunk_size, v)
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

func coloful_heatmap(img *image.RGBA, s, t, s_size, t_size, s_min, t_min, chunk_size int, v float64) {
	if v > 0 {
		r := uint8(255 * v * v)
		g := uint8(255 * (1 - 4*(v-0.5)*(v-0.5)))
		b := uint8(255 * (1 - v*v))
		img.Set(s, t_size-t, color.RGBA{r, g, b, 255})
		return
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

func coloful_noback_heatmap(img *image.RGBA, s, t, s_size, t_size, s_min, t_min, chunk_size int, v float64) {
	if v > 0 {
		r := uint8(255 * v * v)
		g := uint8(255 * (1 - 4*(v-0.5)*(v-0.5)))
		b := uint8(255 * (1 - v*v))
		img.Set(s, t_size-t, color.RGBA{r, g, b, 255})
		return
	}

	img.Set(s, t_size-t, color.RGBA{0, 0, 0, 255})
}

func coloful_opaque_heatmap(img *image.RGBA, s, t, s_size, t_size, s_min, t_min, chunk_size int, v float64) {
	if v > 0 {
		r := uint8(255 * v * v)
		g := uint8(255 * (1 - 4*(v-0.5)*(v-0.5)))
		b := uint8(255 * (1 - v*v))
		img.Set(s, t_size-t, color.RGBA{r, g, b, 255})
		return
	}

	var a uint8
	so := s + s_min
	to := t + t_min
	switch 0 {
	case (so % (10000 / chunk_size)) * (to % (10000 / chunk_size)):
		a = 64
	case (so % (5000 / chunk_size)) * (to % (5000 / chunk_size)):
		a = 128
	case (so % (1000 / chunk_size)) * (to % (1000 / chunk_size)):
		a = 164
	case (so % (500 / chunk_size)) * (to % (500 / chunk_size)):
		a = 176
	case (so % (100 / chunk_size)) * (to % (100 / chunk_size)):
		a = 184
	default:
		a = 192
	}
	img.Set(s, t_size-t, color.RGBA{a, a, a, 255})
}

func hard_heatmap(img *image.RGBA, s, t, s_size, t_size, s_min, t_min, chunk_size int, v float64) {
	var r, g, b float64 = 0, 0, 0

	switch {
	case v == 0:
		r, g, b = 1, 1, 1
	case v <= 0.25:
		g = 4 * v
		b = 1
	case v <= 0.5:
		g = 1
		b = 1 - 4*(v-0.25)
	case v <= 0.75:
		r = 4 * (v - 0.5)
		g = 1
	default:
		r = 1
		g = 1 - 4*(v-0.75)
	}

	baseColor := color.RGBA{uint8(255 * r), uint8(255 * g), uint8(255 * b), 255}

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
	lineColor := color.RGBA{0, 0, 0, a}

	img.Set(s, t_size-t, blend(baseColor, lineColor))
}

func getPosByPlane(plane Plane, chunk_size int, coord sysCoord.Coord) (int, int) {
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

func maxMin(coords []sysCoord.Coord) (sysCoord.Coord, sysCoord.Coord) {
	var max, min sysCoord.Coord

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

func blend(back, front color.RGBA) color.RGBA {
	return color.RGBA{
		R: blend_single(back.R, front.R, front.A),
		G: blend_single(back.G, front.G, front.A),
		B: blend_single(back.B, front.B, front.A),
		A: blend_alpha(back.A, front.A),
	}
}

func blend_single(back, front, alpha uint8) uint8 {
	b := int16(back)
	f := int16(front)
	a := int16(alpha)
	return uint8((255*b + a*f - a*b) / 255)
}

func blend_alpha(back, front uint8) uint8 {
	b := int16(back)
	f := int16(front)
	return uint8((255*f + 255*b - b*f) / 255)
}
