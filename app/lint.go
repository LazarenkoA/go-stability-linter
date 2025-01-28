package app

import (
	"errors"
	"fmt"
	"strings"
)

// Принцип устойчивых зависимостей (SDP) говорит, что метрика компонента должна быть больше метрик компонентов которые от него зависят.
// То есть метрики должны уменьшаться в направлении зависимости.

type packageInfo struct {
	id        string
	name      string
	stability float64 // О - соответствует максимальной устойчивости компонента, 1 - максимальной неустойчивости.
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
	result := make(map[string]*packageInfo, len(packages))
	for _, m := range packages {
		if p := createItem(currentApp, result, m); p != nil {
			result[p.id] = p
		}
	}

	calcStability(result)
	return check(result)
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
	pkg.stability = float64(pkg.outCount) / float64(pkg.outCount+pkg.inCount)
}

func createItem(currentApp string, existing map[string]*packageInfo, pkg *RawPackageInfo) *packageInfo {
	if pkg == nil || !strings.HasPrefix(pkg.Path, currentApp) {
		return nil
	}

	if v, ok := existing[pkg.ID]; ok {
		return v
	}

	pinfo := &packageInfo{
		id:   pkg.ID,
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
		if pkg.parent != nil && pkg.stability > pkg.parent.stability {
			err = errors.Join(err, fmt.Errorf("%s (%.2f) -> %s (%.2f)\n", pkg.parent.id, pkg.parent.stability, pkg.id, pkg.stability))
		}
	}

	return
}
