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

type packageInfo struct {
	ID        string
	name      string
	Stability float64 // О - соответствует максимальной устойчивости компонента, 1 - максимальной неустойчивости.
	inCount   int
	outCount  int
	parent    *packageInfo
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

func GetPackageInfo(currentApp string, packages []*RawPackageInfo) []*packageInfo {
	return maps.Values(getPackageInfo(currentApp, packages))
}

func getPackageInfo(currentApp string, packages []*RawPackageInfo) map[string]*packageInfo {
	result := make(map[string]*packageInfo, len(packages))
	for _, m := range packages {
		if p := createItem(currentApp, result, m); p != nil {
			result[p.ID] = p
		}
	}

	calcStability(result)

	return result
}

func calcStability(existing map[string]*packageInfo) {
	for _, pkg := range existing {
		helperCalcStability(pkg)
	}
}

func helperCalcStability(pkg *packageInfo) {
	if pkg == nil {
		return
	}
	pkg.Stability = float64(pkg.outCount) / float64(pkg.outCount+pkg.inCount)

	if math.IsNaN(pkg.Stability) {
		pkg.Stability = 0
	}
}

func createItem(currentApp string, existing map[string]*packageInfo, pkg *RawPackageInfo) *packageInfo {
	if pkg == nil || !strings.HasPrefix(pkg.Path, currentApp) {
		return nil
	}

	if v, ok := existing[pkg.ID]; ok {
		return v
	}

	pinfo := &packageInfo{
		ID:   pkg.ID,
		name: pkg.Name,
	}

	existing[pkg.ID] = pinfo

	for _, imp := range pkg.Imports {
		if p := createItem(currentApp, existing, imp); p != nil {
			pinfo.outCount++
			p.inCount++
			p.parent = pinfo
		}
	}

	return pinfo
}

func check(pkgs map[string]*packageInfo) (err error) {
	for _, pkg := range pkgs {
		if pkg.parent != nil && pkg.Stability > pkg.parent.Stability {
			err = errors.Join(err, fmt.Errorf("%s (%s:%.2f) -> %s (%s:%.2f)\n", pkg.parent.ID, pkg.parent.name, pkg.parent.Stability, pkg.ID, pkg.name, pkg.Stability))
		}
	}

	return
}
