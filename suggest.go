package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// suggestDirs возвращает список папок, подходящих для автодополнения value.
// value разбивается на базовый путь и префикс; возвращаются подпапки базового
// пути, начинающиеся с префикса (case-insensitive).
func suggestDirs(root, value string) []string {
	value = filepath.ToSlash(value)
	baseDir, prefix := splitPath(value)

	absBase := filepath.Join(root, baseDir)
	entries, err := os.ReadDir(absBase)
	if err != nil {
		return nil
	}

	lower := strings.ToLower(prefix)
	var matches []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasPrefix(strings.ToLower(name), lower) {
			continue
		}
		match := baseDir + name
		matches = append(matches, filepath.ToSlash(match))
	}
	sort.Strings(matches)
	return matches
}

func splitPath(value string) (baseDir, prefix string) {
	value = filepath.ToSlash(value)
	idx := strings.LastIndex(value, "/")
	if idx == -1 {
		return "", value
	}
	return value[:idx+1], value[idx+1:]
}
