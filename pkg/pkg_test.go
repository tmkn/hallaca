package pkg

import (
	"testing"
)

func TestValidPackageSpecs(t *testing.T) {
	tests := []struct {
		input   string
		name    string
		version string
	}{
		{"foo@1.2.3", "foo", "1.2.3"},
		{"@foo/bar@4.5.6", "@foo/bar", "4.5.6"},
		{"foo@^3.0.0", "foo", "^3.0.0"},
	}

	for _, testCase := range tests {
		name, version, error := ParsePackageSpec(testCase.input)

		if error != nil {
			t.Errorf("Didn't expect error for %q", testCase.input)
		}

		if name != testCase.name {
			t.Errorf("Expected name to be %q but got %q", testCase.name, name)
		}

		if version != testCase.version {
			t.Errorf("Expected version to be %q but got %q", testCase.version, version)
		}
	}
}

func TestInValidPackageSpecs(t *testing.T) {
	tests := []struct {
		input   string
		name    string
		version string
	}{
		{"invalidspec", "", ""},
		{"foo@", "", ""},
		{"", "", ""},
	}

	for _, testCase := range tests {
		name, version, error := ParsePackageSpec(testCase.input)

		if error == nil {
			t.Errorf("Expected error for %q", testCase.input)
		}

		if name != testCase.name {
			t.Errorf("Expected name to be %q but got %q", testCase.name, name)
		}

		if version != testCase.version {
			t.Errorf("Expected version to be %q but got %q", testCase.version, version)
		}
	}
}
