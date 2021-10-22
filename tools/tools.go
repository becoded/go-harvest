//go:build tools

package tools

import (
	_ "github.com/boumenot/gocover-cobertura"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/magefile/mage/mage"
	_ "mvdan.cc/gofumpt"
)
