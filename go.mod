module github.com/becoded/go-harvest

go 1.25.1

require (
	github.com/google/go-querystring v1.1.0
	github.com/magefile/mage v1.15.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.11.1
	golang.org/x/tools v0.38.0
)

require (
	github.com/boumenot/gocover-cobertura v1.4.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/mock v0.6.0 // indirect
	golang.org/x/mod v0.29.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mvdan.cc/gofumpt v0.9.1 // indirect
)

tool (
	github.com/boumenot/gocover-cobertura
	github.com/magefile/mage/mage
	go.uber.org/mock/mockgen
	mvdan.cc/gofumpt
)
