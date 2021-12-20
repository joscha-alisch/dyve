package config

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadFile(t *testing.T) {
	tests := []struct {
		desc        string
		path        string
		expected    Config
		expectedErr error
	}{
		{"basic", "basic.yml", Config{
			Providers: []ProviderConfig{
				{Id: "provider-a", Host: "https://provider-a.com", Name: "Provider A", Features: []provider.Type{
					provider.TypeApps, provider.TypePipelines,
				}},
				{Id: "provider-b", Host: "https://provider-b.com", Name: "Provider B", Features: []provider.Type{
					provider.TypeGroups,
				}},
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
