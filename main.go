package main

import (

	// "fmt"

	"math"

	"github.com/tmkn/hallaca/provider"
	"github.com/tmkn/hallaca/resolver"
)

func main() {
	npm_provider := provider.NewNPMProvider()
	name := "express"
	version := "4.17.1"

	standardResolver := resolver.StandardResolver{}

	standardResolver.Resolve(name, version, resolver.Options{
		Provider: npm_provider,
		Depth:    math.MaxUint32,
	})

	// return

	// var package_json = getPackageMetadata(name, version, npm_provider)
	// var dependencies = package_json["dependencies"]

	// data, ok := dependencies.(map[string]interface{})

	// if !ok {
	// 	log.Println("Dependencies is not a string map!")
	// }

	// log.Printf("Name: %v@%v\n", package_json["name"], package_json["version"])
	// log.Printf("Dependencies: %v\n", len(data))

	// for dep_name, dep_version := range data {
	// 	version_array, err := npm_provider.GetVersions(dep_name)

	// 	if err != nil {
	// 		log.Println("couldn't get versions")
	// 	}

	// 	log.Printf("Available versions for %v: %v", dep_name, version_array)

	// 	depVersionStr, ok := dep_version.(string)

	// 	if !ok {
	// 		log.Fatalf("Couldn't convert string")
	// 	}

	// 	constraint, err := semver.NewConstraint(depVersionStr)

	// 	var matchingVersions []*semver.Version

	// 	for _, versionString := range version_array {
	// 		version, err := semver.NewVersion(versionString)
	// 		if err != nil {
	// 			log.Fatalf("invalid available version: %w", err)
	// 		}

	// 		if constraint.Check(version) {
	// 			matchingVersions = append(matchingVersions, version)
	// 		}

	// 	}

	// 	sort.Sort(sort.Reverse(semver.Collection(matchingVersions)))
	// 	// result := make([]string, len(matchingVersions))

	// 	if len(matchingVersions) > 0 {
	// 		log.Println(dep_name, " @ ", dep_version, " -> ", matchingVersions[0], "\n")
	// 	} else {
	// 		log.Println("No version found for", dep_name, "@", dep_version, len(version_array), "\n")
	// 	}
	// }

}

// func getPackageMetadata(name, version string, provider p.Provider) map[string]interface{} {
// 	var package_json map[string]interface{}
// 	result, e := provider.GetPackageMetadata(name, version)

// 	if e != nil {
// 		return package_json
// 	}

// 	err := json.Unmarshal([]byte(result), &package_json)

// 	if err != nil {
// 		return package_json
// 	}

// 	return package_json
// }
