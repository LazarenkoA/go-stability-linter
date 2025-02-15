package app_module

import (
	"github.com/LazarenkoA/go-stability-linter/app"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type stabilityLint struct {
	mainName    string
	stdPackages []*packages.Package
	packages    []*app.RawPackageInfo
}

type CheckLogsInfo struct {
	Package string
	I       float64
}

func NewLint(rootDir string) (*stabilityLint, error) {
	//stdPackages, _ = packages.Load(nil, "std")

	pkg, err := packages.Load(&packages.Config{
		Mode:  packages.NeedImports | packages.NeedName,
		Tests: false, // Не включаем тестовые пакеты
		Dir:   rootDir,
	}, "./...")

	if err != nil {
		return nil, errors.Wrap(err, "packages load error")
	}

	lint := &stabilityLint{
		mainName: getProjectName(rootDir),
	}
	lint.packages = lint.cast(pkg)

	return lint, nil
}

func (l *stabilityLint) cast(pkg []*packages.Package) []*app.RawPackageInfo {
	result := make([]*app.RawPackageInfo, 0, len(pkg))

	for _, p := range pkg {
		if !strings.HasPrefix(p.PkgPath, l.mainName) {
			continue
		}

		result = append(result, &app.RawPackageInfo{
			ID:      p.ID,
			Path:    p.PkgPath,
			Name:    p.Name,
			Imports: l.cast(maps.Values(p.Imports)),
		})
	}

	return result
}

func (l *stabilityLint) Check() error {
	return app.Check(l.mainName, l.packages)
}

func (l *stabilityLint) GetDepList() []*CheckLogsInfo {
	pInfo := app.GetPackageInfo(l.mainName, l.packages)
	result := make([]*CheckLogsInfo, 0, len(pInfo))

	for _, v := range pInfo {
		result = append(result, &CheckLogsInfo{
			Package: v.ID,
			I:       v.Stability,
		})
	}

	return result
}

func (l *stabilityLint) GetPackageInfoTree() []*app.PackageInfo {
	return app.GetPackageInfo(l.mainName, l.packages)
}

func (l *stabilityLint) isStdPackage(pkgPath string) bool {
	_, ok := lo.Find(l.stdPackages, func(pkg *packages.Package) bool {
		return pkg.ID == pkgPath
	})

	return ok
}

func getProjectName(parent string) string {
	gomodPath := filepath.Join(parent, "go.mod")

	if f, err := os.OpenFile(gomodPath, os.O_RDONLY, 0); err == nil {
		data, _ := io.ReadAll(f)

		mod, err := modfile.Parse("go.mod", data, nil)
		if err != nil {
			return ""
		}

		return mod.Module.Mod.Path
	}

	return "."
}
