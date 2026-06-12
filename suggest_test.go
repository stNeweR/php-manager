package main

import (
	"os"
	"path/filepath"
	"testing"
)


func TestSuggestDirs(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, "app", "Services"))
	mustMkdir(t, filepath.Join(root, "app", "Contracts"))
	mustMkdir(t, filepath.Join(root, "app", "Http", "Controllers"))
	mustMkdir(t, filepath.Join(root, "tests", "Feature"))

	cases := []struct {
		value string
		want  []string
	}{
		{"app/", []string{"app/Contracts", "app/Http", "app/Services"}},
		{"app/S", []string{"app/Services"}},
		{"app/Http/", []string{"app/Http/Controllers"}},
		{"tests/", []string{"tests/Feature"}},
		{"unknown/", nil},
		{"zzz", nil},
	}

	for _, c := range cases {
		got := suggestDirs(root, c.value)
		if !sliceEqual(got, c.want) {
			t.Errorf("suggestDirs(%q) = %v, want %v", c.value, got, c.want)
		}
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
