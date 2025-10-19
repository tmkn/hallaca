package pkg

import (
	"fmt"
	"strings"

	"github.com/tmkn/hallaca/provider"
)

type Pkg struct {
	Parent       *Pkg
	Name         string
	AliasedFrom  string
	Version      string
	Dependencies []*Pkg
	IsLoop       bool
	Metadata     *provider.Package
}

func (pkg *Pkg) String() string {
	return pkg.Name + "@" + pkg.Version
}

func DependencyCount(root *Pkg, includeRoot bool) (count int) {
	var queue = []*Pkg{root}

	for len(queue) > 0 {
		count++

		queue = append(queue[1:], queue[0].Dependencies...)
	}

	if !includeRoot {
		return count - 1
	}

	return
}

func GetDepth(start *Pkg) (depth uint) {
	item := start

	for item != nil {
		depth++

		item = item.Parent
	}

	return
}

func ParsePackageSpec(spec string) (string, string, error) {
	if spec == "" {
		return "", "", fmt.Errorf("empty package spec")
	}

	at := strings.LastIndex(spec, "@")
	if at <= 0 || at == len(spec)-1 {
		return "", "", fmt.Errorf("invalid package spec: %s", spec)
	}

	name := spec[:at]
	version := spec[at+1:]
	return name, version, nil
}
