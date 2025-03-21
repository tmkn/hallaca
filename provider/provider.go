package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PackageMetadata = map[string]interface{}

type Provider interface {
	GetPackageMetadata(name string, version string) (string, error)
	GetVersions(name string) ([]string, error)
}

type NPMProvider struct {
	RegistryUrl string
	Versions    map[string][]string
}

func (p *NPMProvider) GetPackageMetadata(name string, version string) (string, error) {
	var url string = p.RegistryUrl + "/" + name + "/" + version

	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (p *NPMProvider) FetchVersions(name string) ([]string, error) {
	var emptyList []string
	var metadata map[string]interface{}
	var url = p.RegistryUrl + "/" + name

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return emptyList, err
	}

	req.Header.Set("Accept", "application/vnd.npm.install-v1+json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return emptyList, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return emptyList, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return emptyList, err
	}

	err = json.Unmarshal(body, &metadata)
	if err != nil {
		return emptyList, fmt.Errorf("couldn't parse json %v", err)
	}

	// collect versions and put them into the string slice
	all_versions, ok := metadata["versions"].(map[string]interface{})
	var versions []string

	if !ok {
		return emptyList, fmt.Errorf("couldn't parse versions")
	}

	for key := range all_versions {
		versions = append(versions, key)
	}

	p.Versions[name] = versions

	return versions, nil
}

func (p *NPMProvider) GetVersions(name string) ([]string, error) {
	if availableVersion, ok := p.Versions[name]; ok {
		return availableVersion, nil
	}

	fetchedVersions, err := p.FetchVersions(name)
	if err != nil {
		return nil, err
	}

	return fetchedVersions, nil
}

func NewNPMProvider() *NPMProvider {
	return &NPMProvider{RegistryUrl: "https://registry.npmjs.org", Versions: make(map[string][]string)}
}
