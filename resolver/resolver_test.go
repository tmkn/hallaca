package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmkn/hallaca/provider"
)

type MockProvider struct {
	packages map[string]map[string]*provider.Package
}

func (m *MockProvider) GetPackageMetadata(name string, version string) (*provider.Package, error) {
	return m.packages[name][version], nil
}

func (m *MockProvider) GetVersions(name string) ([]string, error) {
	versions := make([]string, 0, len(m.packages[name]))
	for k := range m.packages[name] {
		versions = append(versions, k)
	}
	return versions, nil
}

func TestStandardResolver_Resolve(t *testing.T) {
	mockProvider := &MockProvider{
		packages: map[string]map[string]*provider.Package{
			"test-package": {
				"1.0.0": {
					Name:    "test-package",
					Version: "1.0.0",
					Dependencies: map[string]string{
						"dep1": "1.0.0",
					},
				},
			},
			"dep1": {
				"1.0.0": {
					Name:    "dep1",
					Version: "1.0.0",
				},
			},
		},
	}

	resolver := StandardResolver{}
	pkg, err := resolver.Resolve("test-package", "1.0.0", Options{
		Provider: mockProvider,
		Depth:    10,
	})

	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, "test-package", pkg.Name)
	assert.Equal(t, "1.0.0", pkg.Version)
	assert.Len(t, pkg.Dependencies, 1)
	assert.Equal(t, "dep1", pkg.Dependencies[0].Name)
	assert.Equal(t, "1.0.0", pkg.Dependencies[0].Version)
}
