package create

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

	cp "github.com/otiai10/copy"
)

func createProjectFromTpl(tpl string, moduleName string, projectDir string) error {
	_, err := os.Stat(projectDir)
	if !errors.Is(err, os.ErrNotExist) {
		return errors.New("project dir already exists")
	}

	dir, isExample, err := getTemplateCodeDir(tpl, cacheDir)
	if err != nil {
		return err
	}
	originModule, err := checkDirIsGoModuleAndGetModuleName(dir)
	if err != nil {
		return err
	}

	err = cp.Copy(dir, projectDir)
	if err != nil {
		return err
	}
	if isExample {
		if err = fixReplace(projectDir); err != nil {
			return err
		}
	}

	if moduleName == "" {
		moduleName = path.Base(projectDir)
	}
	return replaceModuleName(projectDir, originModule, moduleName)
}

func fixReplace(dir string) error {
	modules, err := listDirModule(dir)
	if err != nil {
		return err
	}
	for _, module := range modules {
		err := fixModReplace(path.Join(dir, module, "go.mod"))
		if err != nil {
			return err
		}
	}
	return nil
}

func fixModReplace(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	content := processGoMod(string(data))

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}
	return nil
}

// processGoMod processes go.mod file content, removes local replace directives
func processGoMod(content string) string {
	// Process single-line replace directives
	content = processSingleLineReplace(content)

	// Process block-style replace directives
	content = processMultiLineReplace(content)

	// Clean up excessive empty lines
	content = cleanupEmptyLines(content)

	return content
}

// processSingleLineReplace processes single-line replace directives
func processSingleLineReplace(content string) string {
	// Match single-line local replace directives - relative paths
	relativePathRegex := regexp.MustCompile(`(?m)^replace\s+([^\s]+)\s*=>\s*\.{1,2}/[^\s]+\s*$`)
	content = relativePathRegex.ReplaceAllString(content, "")

	// Process absolute paths
	var absolutePathRegex *regexp.Regexp
	if runtime.GOOS == "windows" {
		// Windows路径: C:\path 或 D:\path 等
		absolutePathRegex = regexp.MustCompile(`(?m)^replace\s+([^\s]+)\s*=>\s*[A-Za-z]:\\[^\s]+\s*$`)
	} else {
		// Unix路径: /path
		absolutePathRegex = regexp.MustCompile(`(?m)^replace\s+([^\s]+)\s*=>\s*/[^\s]+\s*$`)
	}
	content = absolutePathRegex.ReplaceAllString(content, "")

	return content
}

// processMultiLineReplace processes block-style replace directives
func processMultiLineReplace(content string) string {
	// First find all replace blocks
	replaceBlockRegex := regexp.MustCompile(`(?ms)replace\s*\(\s*(.*?)\s*\)`)

	return replaceBlockRegex.ReplaceAllStringFunc(content, func(block string) string {
		// If this is an empty replace block, return empty string directly
		emptyBlockRegex := regexp.MustCompile(`replace\s*\(\s*\)`)
		if emptyBlockRegex.MatchString(block) {
			return ""
		}

		// Process block content
		lines := strings.Split(block, "\n")
		var newLines []string

		// Keep the first line "replace ("
		newLines = append(newLines, lines[0])

		// Process middle lines
		for i := 1; i < len(lines)-1; i++ {
			line := lines[i]
			// Check if this line contains local path
			if isLocalPathLine(line) {
				continue // Skip local path lines
			}
			newLines = append(newLines, line)
		}

		// If only "replace (" and ")" remain, return empty string
		if len(newLines) == 1 {
			if i := strings.Index(lines[len(lines)-1], ")"); i >= 0 {
				return ""
			}
		}

		// Add the last line ")"
		newLines = append(newLines, lines[len(lines)-1])

		return strings.Join(newLines, "\n")
	})
}

// isLocalPathLine checks if a line contains replace directive with local path
func isLocalPathLine(line string) bool {
	// Trim leading and trailing spaces
	trimmed := strings.TrimSpace(line)

	// Skip empty lines and comment lines
	if trimmed == "" || strings.HasPrefix(trimmed, "//") {
		return false
	}

	// Check if it contains relative paths (./ or ../)
	if strings.Contains(trimmed, " => ./") || strings.Contains(trimmed, " => ../") {
		return true
	}

	// Check if it contains absolute paths
	if runtime.GOOS == "windows" {
		// Windows path pattern: check for patterns like " => C:\"
		driveLetterPattern := regexp.MustCompile(` => [A-Za-z]:\\`)
		if driveLetterPattern.MatchString(trimmed) {
			return true
		}
	} else {
		// Unix path pattern: check for patterns like " => /"
		if strings.Contains(trimmed, " => /") {
			return true
		}
	}

	return false
}

// cleanupEmptyLines cleans up excessive empty lines
func cleanupEmptyLines(content string) string {
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))
	emptyLineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" {
			emptyLineCount++
			if emptyLineCount <= 1 {
				result.WriteString(line + "\n")
			}
		} else {
			emptyLineCount = 0
			result.WriteString(line + "\n")
		}
	}

	return result.String()
}
