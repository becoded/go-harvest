//go:build mage

// Run "mage help" for prerequisites
package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"golang.org/x/tools/go/packages"
)

const (
	ldflags      = "-ldflags=-extldflags -static"
	debugGCflags = "-gcflags='all=-N -l'"
)

type (
	Test    mg.Namespace
	Build   mg.Namespace
	Lint    mg.Namespace
	IntTest mg.Namespace
)

var Aliases = map[string]interface{}{
	"test": Test.Unit,
}

// Show dependencies, useful tips.
func Help() {
	h := `Dependencies:
	- golang: brew install golang
	- golangci-lint: brew install golangci-lint
	- go run mage.go tools

	Add 'export MAGEFILE_ENABLE_COLOR=1' to env for colors
`
	fmt.Println(h)
}

// Run unit tests.
func (Test) Unit() error {
	return test()
}

// Run tests and create coverage profile.
func (Test) Coverage() error {
	pkgs, err := getPackages()
	if err != nil {
		return err
	}

	err = test("-coverprofile=coverage", "-covermode=atomic", "-coverpkg", strings.Join(pkgs, ","))
	if err != nil {
		return err
	}

	return sh.RunV("bash", "-c", "go tool cover -func=coverage | tail -n1")
}

// Show test coverage in html (opens browser).
func (Test) CoverHTML() error {
	mg.Deps(Test.Coverage)

	return sh.RunV("go", "tool", "cover", "-html=coverage")
}

// Save test coverage as cobertura (coverage.xml).
func (Test) CoverXML() error {
	mg.Deps(Test.Coverage)

	return sh.RunV("bash", "-c", "gocover-cobertura < coverage > coverage.xml")
}

// Run golangci-lint with --fix.
func (Lint) GoFix() error {
	return sh.RunV("golangci-lint", "run", "--timeout", "5m", "--fix")
}

// Run golangci-lint and output as code-climate: gl-code-quality-report.json.
func (Lint) Go() error {
	return sh.RunV("bash", "-c", "golangci-lint run --timeout 5m --out-format code-climate > gl-code-quality-report.json")
}

// Remove all build, coverage and linting artifacts.
func Clean() error {
	return sh.RunV(
		"bash", "-c", "rm -f "+
			"coverage "+
			"coverage.xml "+
			"gl-code-quality-report.json ",
	)
}

// Update go dependencies.
func Update() error {
	err := sh.RunV("go", "get", "-u", "-t", "./...")
	if err != nil {
		return err
	}

	return sh.RunV("go", "mod", "tidy")
}

// Generate
func Generate() error {
	return sh.RunV("go", "generate", "./...")
}

// Check if auto-generated code is up-2-date. Ignores comments starting with //..
func DiffGen() error {
	mg.Deps(Generate)

	return sh.RunV("git", "diff", "-G", "^[^/]", "--exit-code")
}

// Install all tools needed to work with this project.
func Tools() error {
	tools := []string{
		"github.com/boumenot/gocover-cobertura",
		"github.com/golang/mock/mockgen",
		"github.com/magefile/mage/mage",
		"mvdan.cc/gofumpt",
	}

	for _, tool := range tools {
		err := sh.RunV("go", "install", tool)
		if err != nil {
			return err
		}
	}

	return nil
}

func getPackages() ([]string, error) {
	pkgs, err := packages.Load(nil, "./...")
	if err != nil {
		return []string{}, err
	}

	var pkgPaths []string

	prx := regexp.MustCompile("mage|mocks")

	for _, pkg := range pkgs {
		matched := prx.MatchString(pkg.PkgPath)
		if !matched {
			pkgPaths = append(pkgPaths, pkg.PkgPath)
		}
	}

	return pkgPaths, nil
}

func test(args ...string) error {
	a := []string{"test", "./...", "-test.short", "-race"}
	a = append(a, args...)

	return sh.RunV("go", a...)
}
