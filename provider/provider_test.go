package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNPMProvider_GetPackageMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"versions": {
				"1.0.0": {
					"name": "test-package",
					"version": "1.0.0",
					"dependencies": {
						"dep1": "1.0.0"
					}
				}
			}
		}`))
	}))
	defer server.Close()

	provider := NewNPMProvider()
	provider.registryUrl = server.URL

	pkg, err := provider.GetPackageMetadata("test-package", "1.0.0")
	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, "test-package", pkg.Name)
	assert.Equal(t, "1.0.0", pkg.Version)
	assert.Equal(t, "1.0.0", pkg.Dependencies["dep1"])
}
