package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
)

type PackageMetadata = map[string]interface{}

type Provider interface {
	GetPackageMetadata(name string, version string) (PackageMetadata, error)
	GetVersions(name string) ([]string, error)
}

type NPMProvider struct {
	registryUrl string
	cache       map[string]map[string]PackageMetadata
}

func NewNPMProvider() *NPMProvider {
	return &NPMProvider{
		registryUrl: "https://registry.npmjs.org",
		cache:       make(map[string]map[string]PackageMetadata),
	}
}

func (p *NPMProvider) GetPackageMetadata(name string, version string) (PackageMetadata, error) {
	if _, exists := p.cache[name]; !exists {
		if err := p.populateCache(name); err != nil {
			return nil, err
		}
	}

	if cacheMetadata, exists := p.cache[name][version]; exists {
		return cacheMetadata, nil
	}

	return nil, fmt.Errorf("couldn't get metadata for %s %s", name, version)
}

func (p *NPMProvider) GetVersions(name string) ([]string, error) {
	cache, exists := p.cache[name]

	if !exists {
		if err := p.populateCache(name); err != nil {
			return nil, err
		}

		cache = p.cache[name]
	}

	var versions []string
	for version := range maps.Keys(cache) {
		versions = append(versions, version)
	}

	return versions, nil
}

func (p *NPMProvider) populateCache(name string) error {
	if _, exists := p.cache[name]; exists {
		return nil
	}

	var url = p.registryUrl + "/" + name

	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var metadata map[string]interface{}

	err = json.Unmarshal(body, &metadata)
	if err != nil {
		return fmt.Errorf("couldn't parse json %v", err)
	}

	allVersions, ok := metadata["versions"].(map[string]interface{})

	if !ok {
		return fmt.Errorf("couldn't parse versions json %v", err)
	}

	p.cache[name] = make(map[string]PackageMetadata)

	for version, _metadata := range allVersions {
		metadata, ok := _metadata.(PackageMetadata)

		if ok {
			p.cache[name][version] = metadata
		}
	}

	return nil
}
