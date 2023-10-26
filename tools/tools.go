//go:build tools

package tools

import (
	_ "github.com/boumenot/gocover-cobertura"
	_ "github.com/magefile/mage/mage"
	_ "go.uber.org/mock/mockgen"
	_ "mvdan.cc/gofumpt"
)
