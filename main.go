package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Composer struct {
	Autoload struct {
		PSR4 map[string]string `json:"psr-4"`
	} `json:"autoload"`
}

var (
	rootDir  string
	composer Composer
	allDirs  []string
)

func main() {
	cwd, _ := os.Getwd()
	rootDir = findComposerRoot(cwd)
	if rootDir == "" {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render("❌ composer.json не найден. Запусти из папки проекта."))
		os.Exit(1)
	}

	// Парсим composer
	data, _ := os.ReadFile(filepath.Join(rootDir, "composer.json"))
	json.Unmarshal(data, &composer)

	// Собираем все папки для автокомплита
	allDirs = collectDirs(rootDir, composer.Autoload.PSR4)

	// Форма
	var kind, dir, name string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Что создаём?").
				Description("Выбери тип файла").
				Options(
					huh.NewOption("PHP Class", "class"),
					huh.NewOption("PHP Interface", "interface"),
					huh.NewOption("PHP Trait", "trait"),
					huh.NewOption("PHP Enum", "enum"),
					huh.NewOption("Laravel Controller", "controller"),
					huh.NewOption("Laravel Model", "model"),
					huh.NewOption("Laravel Middleware", "middleware"),
					huh.NewOption("Laravel FormRequest", "request"),
					huh.NewOption("Laravel Resource", "resource"),
				).
				Value(&kind),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Папка").
				Description("Tab — автокомплит, стрелки — выбор, Enter — подтвердить").
				Placeholder("app/Services").
				Suggestions(allDirs).
				Value(&dir),
			huh.NewInput().
				Title("Имя файла / класса").
				Placeholder("UserService").
				Value(&name),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println("Отменено.")
			return
		}
		fmt.Println("Ошибка:", err)
		os.Exit(1)
	}

	if err := generate(kind, dir, name); err != nil {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render("❌ " + err.Error()))
		os.Exit(1)
	}
}

func findComposerRoot(start string) string {
	for {
		if _, err := os.Stat(filepath.Join(start, "composer.json")); err == nil {
			return start
		}
		parent := filepath.Dir(start)
		if parent == start {
			return ""
		}
		start = parent
	}
}

func collectDirs(root string, psr4 map[string]string) []string {
	seen := map[string]bool{}
	var out []string

	for _, prefix := range psr4 {
		base := filepath.Join(root, prefix)
		filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			rel, _ := filepath.Rel(root, path)
			rel = filepath.ToSlash(rel)
			if !seen[rel] {
				seen[rel] = true
				out = append(out, rel)
			}
			return nil
		})
	}
	return out
}

func resolveNamespace(root, filePath string) string {
	rel, _ := filepath.Rel(root, filePath)
	relDir := filepath.Dir(rel)

	for nsPrefix, pathPrefix := range composer.Autoload.PSR4 {
		pathPrefix = filepath.Clean(pathPrefix) + string(filepath.Separator)
		if strings.HasPrefix(relDir+string(filepath.Separator), pathPrefix) {
			suffix := strings.TrimPrefix(relDir+string(filepath.Separator), pathPrefix)
			suffix = strings.Trim(suffix, string(filepath.Separator))
			ns := nsPrefix + strings.ReplaceAll(suffix, string(filepath.Separator), "\\")
			return strings.Trim(ns, "\\")
		}
	}
	return ""
}

func generate(kind, dir, name string) error {
	name = strings.TrimSuffix(name, ".php")
	if name == "" {
		return fmt.Errorf("имя не может быть пустым")
	}

	filePath := filepath.Join(rootDir, dir, name+".php")
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	ns := resolveNamespace(rootDir, filePath)

	var (
		kindStr string
		extends string
		imports string
		body    string
	)

	switch kind {
	case "controller":
		kindStr, extends, imports = "class", "Controller", "use App\\Http\\Controllers\\Controller;\n"
		body = `    public function index()
    {
        //
    }`
	case "model":
		kindStr, extends, imports = "class", "Model", "use Illuminate\\Database\\Eloquent\\Model;\n"
	case "middleware":
		kindStr = "class"
		imports = "use Closure;\nuse Illuminate\\Http\\Request;\nuse Symfony\\Component\\HttpFoundation\\Response;\n"
		body = `    public function handle(Request $request, Closure $next): Response
    {
        return $next($request);
    }`
	case "request":
		kindStr, extends, imports = "class", "FormRequest", "use Illuminate\\Foundation\\Http\\FormRequest;\n"
		body = `    public function authorize(): bool
    {
        return true;
    }

    public function rules(): array
    {
        return [
            //
        ];
    }`
	case "resource":
		kindStr, extends, imports = "class", "JsonResource", "use Illuminate\\Http\\Request;\nuse Illuminate\\Http\\Resources\\Json\\JsonResource;\n"
		body = `    public function toArray(Request $request): array
    {
        return [
            //
        ];
    }`
	case "interface":
		kindStr, body = "interface", "    //"
	case "trait":
		kindStr, body = "trait", "    //"
	case "enum":
		kindStr, body = "enum", "    case Example = 'example';"
	default:
		kindStr, body = "class", "    //"
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "<?php\n\n")
	if ns != "" {
		fmt.Fprintf(f, "namespace %s;\n\n", ns)
	}
	if imports != "" {
		fmt.Fprintf(f, "%s", imports)
	}
	fmt.Fprintf(f, "%s %s", kindStr, name)
	if extends != "" {
		fmt.Fprintf(f, " extends %s", extends)
	}
	fmt.Fprintf(f, "\n{\n%s\n}\n", body)

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Render
	fmt.Printf("%s %s\n", green("✅ Создан:"), filePath)
	if ns != "" {
		fmt.Printf("   %s %s\n", green("Namespace:"), ns)
	}
	return nil
}
