package config

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadFile(t *testing.T) {
	tests := []struct {
		desc string
		path string
		expected Config
		expectedErr error
	}{
		{"basic", "basic.yml", Config{
			AppProviders: []AppProviderConfig{
				{Name: "provider-a", Host: "https://provider-a.com"},
				{Name: "provider-b", Host: "https://provider-b.com"},
			},
		}, nil},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			_, testFilePath, _, _ := runtime.Caller(0)
			testDir := filepath.Dir(testFilePath)
			fileName := filepath.Join(testDir, "test_configs", test.path)

			res, err := LoadFrom(fileName)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nexpected err: %v\n   got %v", test.expectedErr, err)
			}

			if !cmp.Equal(test.expected, res) {
				tt.Errorf("\ndiff between configs: \n%s\n", cmp.Diff(test.expected, res))
			}
		})
	}

}