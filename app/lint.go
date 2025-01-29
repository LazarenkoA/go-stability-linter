package app

import (
	"errors"
	"fmt"
	"golang.org/x/exp/maps"
	"math"
	"strings"
)

// Принцип устойчивых зависимостей (SDP) говорит, что метрика компонента должна быть больше метрик компонентов которые от него зависят.
// То есть метрики должны уменьшаться в направлении зависимости.

type PackageInfo struct {
	ID        string
	name      string
	Stability float64 // О - соответствует максимальной устойчивости компонента, 1 - максимальной неустойчивости.
	inCount   int
	outCount  int
	Parent    *PackageInfo
	Childs    []*PackageInfo
}

type RawPackageInfo struct {
	ID      string
	Path    string
	Name    string
	Imports []*RawPackageInfo
}

func Check(currentApp string, packages []*RawPackageInfo) error {
	return check(getPackageInfo(currentApp, packages))
}

func GetPackageInfo(currentApp string, packages []*RawPackageInfo) []*PackageInfo {
	return maps.Values(getPackageInfo(currentApp, packages))
}

func getPackageInfo(currentApp string, packages []*RawPackageInfo) map[string]*PackageInfo {
	result := make(map[string]*PackageInfo, len(packages))
	for _, m := range packages {
		if p := createItem(currentApp, result, m); p != nil {
			result[p.ID] = p
		}
	}

	calcStability(result)

	return result
}

func calcStability(existing map[string]*PackageInfo) {
	for _, pkg := range existing {
		helperCalcStability(pkg)
	}
}

func helperCalcStability(pkg *PackageInfo) {
	if pkg == nil {
		return
	}
	pkg.Stability = float64(pkg.outCount) / float64(pkg.outCount+pkg.inCount)

	if math.IsNaN(pkg.Stability) {
		pkg.Stability = 0
	}
}

func createItem(currentApp string, existing map[string]*PackageInfo, pkg *RawPackageInfo) *PackageInfo {
	if pkg == nil || !strings.HasPrefix(pkg.Path, currentApp) {
		return nil
	}

	if v, ok := existing[pkg.ID]; ok {
		return v
	}

	pinfo := &PackageInfo{
		ID:   pkg.ID,
		name: pkg.Name,
	}

	existing[pkg.ID] = pinfo

	for _, imp := range pkg.Imports {
		if p := createItem(currentApp, existing, imp); p != nil {
			pinfo.outCount++
			p.inCount++

			p.Parent = pinfo
			pinfo.Childs = append(pinfo.Childs, p)
		}
	}

	return pinfo
}

func check(pkgs map[string]*PackageInfo) (err error) {
	for _, pkg := range pkgs {
		if pkg.Parent == nil {
			fmt.Println(pkg.ID)
		}

		if pkg.Parent != nil && pkg.Stability > pkg.Parent.Stability {
			err = errors.Join(err, fmt.Errorf("%s (%s:%.2f) -> %s (%s:%.2f)\n", pkg.Parent.ID, pkg.Parent.name, pkg.Parent.Stability, pkg.ID, pkg.name, pkg.Stability))
		}
	}

	return
}
