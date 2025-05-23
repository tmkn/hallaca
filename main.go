package main

import (
	"math"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/tmkn/hallaca/provider"
	"github.com/tmkn/hallaca/resolver"
)

func main() {
	args := os.Args[1:]
	name := "react"
	version := "18.3.1"

	if len(args) == 2 {
		name = args[0]
		version = args[1]
	}

	form := huh.NewForm(huh.NewGroup(huh.NewInput().
		Title("Package?").
		Value(&name), huh.NewInput().
		Title("Version?").
		Value(&version)),
	).WithTheme(huh.ThemeCatppuccin())

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	npmProvider := provider.NewNPMProvider()

	standardResolver := resolver.StandardResolver{}

	root := standardResolver.Resolve(name, version, resolver.Options{
		Provider: npmProvider,
		Depth:    math.MaxUint32,
	})

	log.Info(root)

}
