package make

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

// addRegistryItem 添加注册项到指定列表
func addRegistryItem(file, markerText, itemText string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	item := []byte(itemText)
	if bytes.Contains(content, item) {
		return nil
	}

	marker := []byte(markerText)
	start := bytes.Index(content, marker)
	if start < 0 {
		return fmt.Errorf("未找到注册列表")
	}

	newline := "\n"
	if bytes.Contains(content, []byte("\r\n")) {
		newline = "\r\n"
	}

	bodyStart := start + len(marker)
	closeMarker := []byte(newline + "\t}")
	end := bytes.Index(content[bodyStart:], closeMarker)
	if end < 0 {
		return fmt.Errorf("注册列表格式不正确")
	}

	insertAt := bodyStart + end
	itemLine := []byte(newline + "\t\t" + string(item) + ",")
	result := make([]byte, 0, len(content)+len(itemLine))
	result = append(result, content[:insertAt]...)
	result = append(result, itemLine...)
	result = append(result, content[insertAt:]...)

	return os.WriteFile(file, result, 0644)
}

// addRegistryImport 添加注册文件导入
func addRegistryImport(file, importPath string) (string, error) {
	alias := registryImportAlias(importPath)
	if err := addImport(file, importPath, alias, true); err != nil {
		return "", err
	}

	return alias, nil
}

// addBlankImport 添加空白导入
func addBlankImport(file, importPath string) error {
	return addImport(file, importPath, "_", false)
}

// addImport 添加导入
func addImport(file, importPath, alias string, ensureAlias bool) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	quotedPath := fmt.Sprintf("%q", importPath)
	if bytes.Contains(content, []byte(quotedPath)) {
		if !ensureAlias || alias == "_" {
			return nil
		}

		_, err = ensureRegistryImportAlias(file, content, quotedPath, alias)
		return err
	}

	newline := "\n"
	if bytes.Contains(content, []byte("\r\n")) {
		newline = "\r\n"
	}

	marker := []byte("import (")
	start := bytes.Index(content, marker)
	if start < 0 {
		return fmt.Errorf("未找到导入列表")
	}

	bodyStart := start + len(marker)
	closeMarker := []byte(newline + ")")
	end := bytes.Index(content[bodyStart:], closeMarker)
	if end < 0 {
		return fmt.Errorf("导入列表格式不正确")
	}

	insertAt := bodyStart + end
	importLine := []byte(newline + "\t" + alias + " " + quotedPath)
	result := make([]byte, 0, len(content)+len(importLine))
	result = append(result, content[:insertAt]...)
	result = append(result, importLine...)
	result = append(result, content[insertAt:]...)

	if err = os.WriteFile(file, result, 0644); err != nil {
		return err
	}

	return nil
}

// ensureRegistryImportAlias 确保导入使用别名
func ensureRegistryImportAlias(file string, content []byte, quotedPath, alias string) (string, error) {
	newline := "\n"
	if bytes.Contains(content, []byte("\r\n")) {
		newline = "\r\n"
	}

	lines := strings.Split(string(content), newline)
	for index, line := range lines {
		if !strings.Contains(line, quotedPath) {
			continue
		}

		current := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(line), quotedPath))
		if current != "" {
			return current, nil
		}

		indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
		lines[index] = indent + alias + " " + quotedPath
		if err := os.WriteFile(file, []byte(strings.Join(lines, newline)), 0644); err != nil {
			return "", err
		}

		return alias, nil
	}

	return alias, nil
}

// registryImportAlias 生成注册导入别名
func registryImportAlias(importPath string) string {
	alias := lo.CamelCase(path.Base(importPath))
	if alias == "" {
		return "dynamic"
	}

	return alias
}

// registryImportPath 生成注册导入路径
func registryImportPath(file string) string {
	dir := filepath.ToSlash(filepath.Dir(file))
	dir = strings.TrimPrefix(dir, "./")
	return "gin/" + strings.TrimPrefix(dir, "/")
}
