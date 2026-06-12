package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	root := t.TempDir()
	rootDir = root
	composer = Composer{}
	composer.Autoload.PSR4 = map[string]string{
		"App\\": "app/",
	}

	composerPath := filepath.Join(root, "composer.json")
	if err := os.WriteFile(composerPath, []byte(`{"autoload":{"psr-4":{"App\\":"app/"}}}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Class с final и strict_types
	if err := generate("class", "app/Services", "UserService", true); err != nil {
		t.Fatal(err)
	}
	content := readFile(t, filepath.Join(root, "app/Services/UserService.php"))
	if !strings.Contains(content, "declare(strict_types=1);") {
		t.Errorf("expected strict_types declaration, got:\n%s", content)
	}
	if !strings.Contains(content, "final class UserService") {
		t.Errorf("expected final class, got:\n%s", content)
	}

	// Interface: final должен игнорироваться
	if err := generate("interface", "app/Contracts", "UserRepository", true); err != nil {
		t.Fatal(err)
	}
	content = readFile(t, filepath.Join(root, "app/Contracts/UserRepository.php"))
	if !strings.Contains(content, "declare(strict_types=1);") {
		t.Errorf("expected strict_types declaration, got:\n%s", content)
	}
	if strings.Contains(content, "final interface") {
		t.Errorf("interface should not be final, got:\n%s", content)
	}
	if !strings.Contains(content, "interface UserRepository") {
		t.Errorf("expected interface, got:\n%s", content)
	}

	// Enum: final игнорируется
	if err := generate("enum", "app/Enums", "Status", true); err != nil {
		t.Fatal(err)
	}
	content = readFile(t, filepath.Join(root, "app/Enums/Status.php"))
	if strings.Contains(content, "final enum") {
		t.Errorf("enum should not be final, got:\n%s", content)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
