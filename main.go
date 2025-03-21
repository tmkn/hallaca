package main

import (
	"encoding/json"
	"sort"

	// "fmt"
	"log"

	"github.com/Masterminds/semver/v3"
	p "github.com/tmkn/hallaca/provider"
)

func main() {

	npm_provider := p.NewNPMProvider()
	name := "express"
	version := "4.17.1"

	var package_json = getPackageMetadata(name, version, npm_provider)
	var dependencies = package_json["dependencies"]

	data, ok := dependencies.(map[string]interface{})

	if !ok {
		log.Println("Dependencies is not a string map!")
	}

	log.Printf("Name: %v@%v\n", package_json["name"], package_json["version"])
	log.Printf("Dependencies: %v\n", len(data))

	for dep_name, dep_version := range data {
		version_array, err := npm_provider.GetVersions(dep_name)

		if err != nil {
			log.Println("couldn't get versions")
		}

		log.Printf("Available versions for %v: %v", dep_name, version_array)

		depVersionStr, ok := dep_version.(string)

		if !ok {
			log.Fatalf("Couldn't convert string")
		}

		constraint, err := semver.NewConstraint(depVersionStr)

		var matchingVersions []*semver.Version

		for _, versionString := range version_array {
			version, err := semver.NewVersion(versionString)
			if err != nil {
				log.Fatalf("invalid available version: %w", err)
			}

			if constraint.Check(version) {
				matchingVersions = append(matchingVersions, version)
			}

		}

		sort.Sort(sort.Reverse(semver.Collection(matchingVersions)))
		// result := make([]string, len(matchingVersions))

		if len(matchingVersions) > 0 {
			log.Println(dep_name, " @ ", dep_version, " -> ", matchingVersions[0], "\n")
		} else {
			log.Println("No version found for", dep_name, "@", dep_version, len(version_array), "\n")
		}
	}

	// meta, err := npm_provider.GetVersions(name)

	// if err != nil {
	// 	log.Println("Couldn't get versions")
	// }

	// versions := meta["versions"]
	// data2, ok2 := versions.(map[string]interface{})
	// if !ok2 {
	// 	log.Println("Couldn't cast dependencies!")
	// }
	// foo := getVersions(data2)

	// log.Printf("Available versions for %v: %v", name, foo)
}

func getPackageMetadata(name, version string, provider p.Provider) map[string]interface{} {
	var package_json map[string]interface{}
	result, e := provider.GetPackageMetadata(name, version)

	if e != nil {
		return package_json
	}

	err := json.Unmarshal([]byte(result), &package_json)

	if err != nil {
		return package_json
	}

	return package_json
}

// func getVersions(versions map[string]any) []string {
// 	// var versions = metadata["versions"]

// 	// version_data, ok := versions.(map[string]interface{})

// 	// if !ok {
// 	// 	log.Println("Dependencies is not a string map!")
// 	// }

// 	keys := make([]string, 0, len(versions))

// 	for key := range versions {
// 		keys = append(keys, key)
// 	}

// 	return keys
// }
