package provider

import (
	"io"
	"net/http"
)

type Provider interface {
	GetPackageMetadata(name string, version string) (string, error)
}

type NPMProvider struct {
	RegistryUrl string
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
