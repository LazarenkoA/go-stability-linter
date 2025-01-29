package main

import (
	"fmt"
	"github.com/LazarenkoA/go-stability-linter/app"
	"github.com/pterm/pterm"
	"math"
	"strings"
)

// реши установленную задачу на Go

func print(tree []*app.PackageInfo, subStringExclude string) {
	for _, l := range tree {
		if len(l.Parents) > 0 && (subStringExclude == "" || !strings.Contains(l.ID, subStringExclude)) {
			helperPrint([]*app.PackageInfo{l}, 0)
			fmt.Print("\n")
		}
	}
}

func helperPrint(nodes []*app.PackageInfo, shift int) {
	if nodes == nil {
		return
	}

	for _, node := range nodes {
		r, g, b := HSVToRGB(mapRange(node.Stability, 1, 0, 0, 120), 100, 200)
		pterm.NewRGB(r, g, b).Printf("%s%s - %.2f\n", strings.Repeat("\t", shift), node.ID, node.Stability)
		helperPrint(node.Childs, shift+1)
	}

}

// mapRange масштабирует значение из одного диапазона в другой
// value - входное значение
// inMin, inMax - минимальное и максимальное значение входного диапазона
// outMin, outMax - минимальное и максимальное значение выходного диапазона
// Возвращает масштабированное значение
func mapRange(value, inMin, inMax, outMin, outMax float64) float64 {
	return (value-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func HSVToRGB(h, s, v float64) (uint8, uint8, uint8) {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	s = math.Max(0, math.Min(1, s))
	v = math.Max(0, math.Min(1, v))

	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var r, g, b float64

	switch {
	case 0 <= h && h < 60:
		r, g, b = c, x, 0
	case 60 <= h && h < 120:
		r, g, b = x, c, 0
	case 120 <= h && h < 180:
		r, g, b = 0, c, x
	case 180 <= h && h < 240:
		r, g, b = 0, x, c
	case 240 <= h && h < 300:
		r, g, b = x, 0, c
	case 300 <= h && h < 360:
		r, g, b = c, 0, x
	}

	r = (r + m) * 255
	g = (g + m) * 255
	b = (b + m) * 255

	return uint8(r), uint8(g), uint8(b)
}
