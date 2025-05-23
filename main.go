package main

import (
	"fmt"
	"math"

	"github.com/tmkn/hallaca/provider"
	"github.com/tmkn/hallaca/resolver"
)

func main() {
	npm_provider := provider.NewNPMProvider()
	name := "react"
	version := "18.3.1"

	standardResolver := resolver.StandardResolver{}

	root := standardResolver.Resolve(name, version, resolver.Options{
		Provider: npm_provider,
		Depth:    math.MaxUint32,
	})

	fmt.Println(root)

}
