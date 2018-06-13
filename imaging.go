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

type plane int

const (
	xz plane = iota
	zy
	xy
)

func main() {
	// flags
	var coordsFileName string
	flag.StringVar(&coordsFileName, "i", "coords.bin", "input file")

	var imageFileName string
	flag.StringVar(&imageFileName, "o", "", "output file")

	var planeName string
	flag.StringVar(&planeName, "p", "xz", "dump plane (xz, zy, xy)")

	var chunkSize int
	flag.IntVar(&chunkSize, "s", 20, "pixcel size in LY")

	var curveName string
	flag.StringVar(&curveName, "hc", "log", "heatmap curve (liner, log)")

	var heatmapName string
	flag.StringVar(&heatmapName, "ht", "opaque", "heatmap type (colorful, noback, opaque, hard)")

	var heatStcale float64
	flag.Float64Var(&heatStcale, "hs", 0.1, "heatmap scale")

	var noAdjust bool
	flag.BoolVar(&noAdjust, "hna", false, "disable heatmap scale adjust")

	var scaleVar bool
	flag.BoolVar(&scaleVar, "bar", false, "enable scale bar")

	var sizeAdjust int
	flag.IntVar(&sizeAdjust, "multof", 0, "set image size to multiple of arg (0 is disable)")

	// flag parse
	flag.Parse()

	// main
	if imageFileName == "" {
		imageFileName = fmt.Sprintf("%s_%d.png", planeName, chunkSize)
	}

	var plane plane
	switch planeName {
	case "xz":
		plane = xz
	case "zy":
		plane = zy
	case "xy":
		plane = xy
	default:
		fmt.Fprintf(os.Stderr, "Error: Invalid plane name'%s'", planeName)
		os.Exit(1)
	}

	var curve func(float64) float64
	switch curveName {
	case "log":
		curve = func(v float64) float64 { return math.Log10(v + 1) }
	case "liner":
		curve = func(v float64) float64 { return v }
	default:
		fmt.Fprintln(os.Stderr, "Error: invalid curve name")
		os.Exit(1)
	}

	var heatmap func(*image.RGBA, int, int, int, int, int, int, int, float64)
	switch heatmapName {
	case "colorful":
		heatmap = colofulHeatmap
	case "noback":
		heatmap = colofulNobackHeatmap
	case "opaque":
		heatmap = colofulOpaqueHeatmap
	case "hard":
		heatmap = hardHeatmap
	default:
		fmt.Fprintln(os.Stderr, "Error: invalid heatmap name")
		os.Exit(1)
	}

	//////////////

	coords := make([]sysCoord.Coord, 0, 1024)
	for c := range sysCoord.LoadCoords(coordsFileName) {
		coords = append(coords, c)
	}
	max, min := maxMin(coords)

	sMax, tMax := getPosByPlane(plane, chunkSize, max)
	sMin, tMin := getPosByPlane(plane, chunkSize, min)

	sSize := sMax - sMin + 4
	tSize := tMax - tMin + 4

	if sizeAdjust > 0 {
		if sSize%sizeAdjust != 0 {
			sSize += sizeAdjust - sSize%sizeAdjust
		}
		if tSize%sizeAdjust != 0 {
			tSize += sizeAdjust - tSize%sizeAdjust
		}
	}

	scaleVarSize := 0
	if scaleVar {
		if tSize < 128 {
			scaleVarSize = 4
		} else if tSize > 1024 {
			scaleVarSize = 32
		} else {
			scaleVarSize = tSize / 32
		}
	}

	lines := make([][]float64, tSize)
	for t := 0; t < tSize; t++ {
		lines[t] = make([]float64, sSize)
	}

	for _, coord := range coords {
		s, t := getPosByPlane(plane, chunkSize, coord)
		s -= sMin
		t -= tMin
		lines[t][s]++
	}

	var vMax float64
	for t := 0; t < tSize; t++ {
		for s := 0; s < sSize; s++ {
			v := curve(lines[t][s])
			lines[t][s] = v
			if v > vMax {
				vMax = v
			}
		}
	}
	fmt.Println("Heat max:", vMax)
	if !noAdjust && vMax > heatStcale {
		fmt.Println("Heat scale adjusted to heat max.")
		heatStcale = vMax
	}

	img := image.NewRGBA(image.Rect(0, 0, sSize, tSize+scaleVarSize))

	for t := 0; t < tSize; t++ {
		for s := 0; s < sSize; s++ {
			v := lines[t][s] / heatStcale
			if v > 1.0 {
				v = 1.0
			}
			if v < 0.0 {
				v = 0.0
			}
			heatmap(img, s, t, sSize, tSize+scaleVarSize, sMin, tMin, chunkSize, v)
		}
	}

	for t := tSize; t < (tSize + scaleVarSize); t++ {
		for s := 0; s < sSize; s++ {
			if t == tSize {
				img.Set(s, tSize-t, color.RGBA{255, 255, 255, 255})
				continue
			}
			v := float64(s) / float64(sSize)
			heatmap(img, s, t, sSize, tSize+scaleVarSize, sMin, tMin, chunkSize, v)
		}
	}

	imgFile, err := os.Create(imageFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannnot open file.\n    %s", err)
		os.Exit(1)
	}
	defer imgFile.Close()

	if err := png.Encode(imgFile, img); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannnot write image.\n    %s", err)
		os.Exit(1)
	}
}

func colofulHeatmap(img *image.RGBA, s, t, sSize, tSize, sMin, tMin, chunkSize int, v float64) {
	if v > 0 {
		r := uint8(255 * v * v)
		g := uint8(255 * (1 - 4*(v-0.5)*(v-0.5)))
		b := uint8(255 * (1 - v*v))
		img.Set(s, tSize-t, color.RGBA{r, g, b, 255})
		return
	}

	var a uint8
	so := s + sMin
	to := t + tMin
	switch 0 {
	case (so % (10000 / chunkSize)) * (to % (10000 / chunkSize)):
		a = 192
	case (so % (5000 / chunkSize)) * (to % (5000 / chunkSize)):
		a = 128
	case (so % (1000 / chunkSize)) * (to % (1000 / chunkSize)):
		a = 92
	case (so % (500 / chunkSize)) * (to % (500 / chunkSize)):
		a = 80
	case (so % (100 / chunkSize)) * (to % (100 / chunkSize)):
		a = 72
	default:
		a = 64
	}
	img.Set(s, tSize-t, color.RGBA{0, 0, 0, a})
}

func colofulNobackHeatmap(img *image.RGBA, s, t, sSize, tSize, sMin, tMin, chunkSize int, v float64) {
	if v > 0 {
		r := uint8(255 * v * v)
		g := uint8(255 * (1 - 4*(v-0.5)*(v-0.5)))
		b := uint8(255 * (1 - v*v))
		img.Set(s, tSize-t, color.RGBA{r, g, b, 255})
		return
	}

	img.Set(s, tSize-t, color.RGBA{0, 0, 0, 255})
}

func colofulOpaqueHeatmap(img *image.RGBA, s, t, sSize, tSize, sMin, tMin, chunkSize int, v float64) {
	if v > 0 {
		r := uint8(255 * v * v)
		g := uint8(255 * (1 - 4*(v-0.5)*(v-0.5)))
		b := uint8(255 * (1 - v*v))
		img.Set(s, tSize-t, color.RGBA{r, g, b, 255})
		return
	}

	var a uint8
	so := s + sMin
	to := t + tMin
	switch 0 {
	case (so % (10000 / chunkSize)) * (to % (10000 / chunkSize)):
		a = 64
	case (so % (5000 / chunkSize)) * (to % (5000 / chunkSize)):
		a = 128
	case (so % (1000 / chunkSize)) * (to % (1000 / chunkSize)):
		a = 164
	case (so % (500 / chunkSize)) * (to % (500 / chunkSize)):
		a = 176
	case (so % (100 / chunkSize)) * (to % (100 / chunkSize)):
		a = 184
	default:
		a = 192
	}
	img.Set(s, tSize-t, color.RGBA{a, a, a, 255})
}

func hardHeatmap(img *image.RGBA, s, t, sSize, tSize, sMin, tMin, chunkSize int, v float64) {
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
	so := s + sMin
	to := t + tMin
	switch 0 {
	case (so % (10000 / chunkSize)) * (to % (10000 / chunkSize)):
		a = 192
	case (so % (5000 / chunkSize)) * (to % (5000 / chunkSize)):
		a = 128
	case (so % (1000 / chunkSize)) * (to % (1000 / chunkSize)):
		a = 92
	case (so % (500 / chunkSize)) * (to % (500 / chunkSize)):
		a = 80
	case (so % (100 / chunkSize)) * (to % (100 / chunkSize)):
		a = 72
	default:
		a = 0
	}
	lineColor := color.RGBA{0, 0, 0, a}

	img.Set(s, tSize-t, blend(baseColor, lineColor))
}

func getPosByPlane(plane plane, chunkSize int, coord sysCoord.Coord) (int, int) {
	switch plane {
	case xz:
		return chunk(chunkSize, coord.X), chunk(chunkSize, coord.Z)
	case zy:
		return chunk(chunkSize, coord.Z), chunk(chunkSize, coord.Y)
	case xy:
		return chunk(chunkSize, coord.X), chunk(chunkSize, coord.Y)
	default:
		panic("Inlvalid plane value")
	}
}

func chunk(chunkSize int, val float32) int {
	return int(val / float32(chunkSize))
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
		R: brendSingle(back.R, front.R, front.A),
		G: brendSingle(back.G, front.G, front.A),
		B: brendSingle(back.B, front.B, front.A),
		A: blendAlpha(back.A, front.A),
	}
}

func brendSingle(back, front, alpha uint8) uint8 {
	b := int16(back)
	f := int16(front)
	a := int16(alpha)
	return uint8((255*b + a*f - a*b) / 255)
}

func blendAlpha(back, front uint8) uint8 {
	b := int16(back)
	f := int16(front)
	return uint8((255*f + 255*b - b*f) / 255)
}
