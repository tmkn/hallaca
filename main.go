package main

import (
	"math"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/tmkn/hallaca/pkg"
	"github.com/tmkn/hallaca/provider"
	"github.com/tmkn/hallaca/resolver"
	"github.com/tmkn/hallaca/tar"
)

func main() {
	args := os.Args[1:]
	name := "react"
	version := "18.3.1"
	// name := "@tanstack/start"
	// version := "1.120.11"

	if len(args) == 2 {
		name = args[0]
		version = args[1]
	} else {
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
	}

	npmProvider := provider.NewNPMProvider()

	standardResolver := resolver.StandardResolver{}

	root, err := standardResolver.Resolve(name, version, resolver.Options{
		Provider: npmProvider,
		Depth:    math.MaxUint32,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Dependency count: %d\n", pkg.DependencyCount(root, false))

	tar.DownloadTar(root)

}
