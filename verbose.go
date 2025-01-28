package main

import (
	"github.com/LazarenkoA/go-stability-linter/app_module"
	"github.com/pterm/pterm"
	"math"
)

func print(log []*app_module.CheckLogsInfo) {
	for _, l := range log {
		if l.I == 0 || math.IsNaN(l.I) {
			pterm.NewRGB(0, 128, 0).Printf("%s - %.2f\n", l.Package, l.I)
		} else if l.I == 1 {
			pterm.NewRGB(255, 215, 0).Printf("%s - %.2f\n", l.Package, l.I)
		} else {
			pterm.NewRGB(128, 128, 128).Printf("%s - %.2f\n", l.Package, l.I)
		}
	}
}
