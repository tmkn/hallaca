package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PackageMetata = map[string]interface{}

type Provider interface {
	GetPackageMetadata(name string, version string) (string, error)
	GetVersions(name string) (map[string]interface{}, error)
}

type NPMProvider struct {
	RegistryUrl string
	MetadataMap map[string]PackageMetata
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

func (p *NPMProvider) GetVersions(name string) (map[string]interface{}, error) {
	var metadata map[string]interface{}
	var url = p.RegistryUrl + "/" + name

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return metadata, err
	}

	req.Header.Set("Accept", "application/vnd.npm.install-v1+json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return metadata, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return metadata, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return metadata, err
	}

	err = json.Unmarshal(body, &metadata)
	if err != nil {
		return metadata, fmt.Errorf("couldn't parse json %v", err)
	}

	return metadata, err
}
