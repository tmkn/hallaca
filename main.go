package main

import (
	"encoding/json"
	"fmt"

	p "github.com/tmkn/hallaca/provider"
)

func main() {

	npm_provider := p.NPMProvider{RegistryUrl: "https://registry.npmjs.org"}
	name := "express"
	version := "4.17.1"

	var package_json = getPackageMetadata(name, version, &npm_provider)
	var dependencies = package_json["dependencies"]

	data, ok := dependencies.(map[string]interface{})

	if !ok {
		fmt.Println("Dependencies is not a string map!")
	}

	fmt.Printf("Name: %v@%v\n", package_json["name"], package_json["version"])
	fmt.Printf("Dependencies: %v\n", len(data))

	for key, value := range data {
		fmt.Printf("%v.| %v\n", key, value)
	}
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
