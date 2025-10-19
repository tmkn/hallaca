package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"sync"
	"time"
)

type Package struct {
	Name         string
	Version      string
	Dependencies map[string]string
}

type Provider interface {
	GetPackageMetadata(name string, version string) (*Package, error)
	GetVersions(name string) ([]string, error)
}

type NPMProvider struct {
	registryUrl string
	cache       map[string]map[string]*Package
	mutex       sync.RWMutex
	client      *http.Client
}

func NewNPMProvider() *NPMProvider {
	return &NPMProvider{
		registryUrl: "https://registry.npmjs.org",
		cache:       make(map[string]map[string]*Package),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *NPMProvider) GetPackageMetadata(name string, version string) (*Package, error) {
	p.mutex.RLock()
	cachedVersions, exists := p.cache[name]
	p.mutex.RUnlock()

	if !exists {
		if err := p.populateCache(name); err != nil {
			return nil, err
		}
		p.mutex.RLock()
		cachedVersions = p.cache[name]
		p.mutex.RUnlock()
	}

	if cacheMetadata, exists := cachedVersions[version]; exists {
		return cacheMetadata, nil
	}

	return nil, fmt.Errorf("couldn't get metadata for %s %s", name, version)
}

func (p *NPMProvider) GetVersions(name string) ([]string, error) {
	p.mutex.RLock()
	cache, exists := p.cache[name]
	p.mutex.RUnlock()

	if !exists {
		if err := p.populateCache(name); err != nil {
			return nil, err
		}
		p.mutex.RLock()
		cache = p.cache[name]
		p.mutex.RUnlock()
	}

	var versions []string
	for version := range maps.Keys(cache) {
		versions = append(versions, version)
	}

	return versions, nil
}

func (p *NPMProvider) populateCache(name string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, exists := p.cache[name]; exists {
		return nil
	}

	var url = p.registryUrl + "/" + name

	resp, err := p.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch package metadata: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var registryResponse struct {
		Versions map[string]Package `json:"versions"`
	}

	err = json.Unmarshal(body, &registryResponse)
	if err != nil {
		return fmt.Errorf("couldn't parse json %v", err)
	}

	p.cache[name] = make(map[string]*Package)
	for version, pkg := range registryResponse.Versions {
		p.cache[name][version] = &pkg
	}

	return nil
}
