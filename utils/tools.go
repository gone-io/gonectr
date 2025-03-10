package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ExtractPackageArg 从命令行参数中提取 package 参数
func ExtractPackageArg(args []string) string {
	var packageArg string
	skipNext := false

	for _, arg := range args {
		// 如果需要跳过下一个参数（例如 -exec 的值），则跳过
		if skipNext {
			skipNext = false
			continue
		}

		// 判断是否为 -exec 标志
		if arg == "-exec" {
			skipNext = true // 跳过 -exec 的下一个参数
			continue
		}

		// 判断是否为 build flag
		if len(arg) > 1 && arg[0] == '-' {
			continue // 跳过 build flag
		}

		// 第一个非 flag 参数即为 package
		packageArg = arg
		break
	}

	return packageArg
}

type ModuleInfo struct {
	ModuleName string
	ModulePath string
}

// parseModuleName 读取 go.mod 文件并解析出 module 名称
func parseModuleName(goModPath string) (string, error) {
	file, err := os.Open(goModPath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("关闭文件出错:", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// 检查行是否以 "module" 开头
		if strings.HasPrefix(line, "module ") {
			// 提取模块名称
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			return moduleName, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取文件出错: %w", err)
	}

	return "", fmt.Errorf("未找到 module 声明")
}

// FindGoModFile 从指定目录向上逐层搜索 "go.mod" 文件
func FindGoModFile(dir string) (string, error) {
	for {
		goModPath := filepath.Join(dir, "go.mod")

		// 检查当前目录是否有 "go.mod" 文件
		if _, err := os.Stat(goModPath); err == nil {
			return filepath.Dir(goModPath), nil
		}

		// 获取上级目录
		parentDir := filepath.Dir(dir)

		// 如果已经到达根目录，就退出
		if parentDir == dir {
			return "", fmt.Errorf("未找到 go.mod 文件")
		}

		// 更新目录为上级目录，继续搜索
		dir = parentDir
	}
}

// FindModuleInfo 获取当前目录所在的模块信息
func FindModuleInfo(dir string) (*ModuleInfo, error) {
	modulePath, err := FindGoModFile(dir)
	if err != nil {
		return nil, err
	}

	moduleName, err := parseModuleName(path.Join(modulePath, "go.mod"))
	if err != nil {
		return nil, err
	}

	modulePath, err = filepath.Abs(modulePath)
	if err != nil {
		return nil, err
	}
	return &ModuleInfo{
		ModuleName: moduleName,
		ModulePath: modulePath,
	}, nil
}

// FindFirstGoGenerateLine 扫描目录，找到第一个匹配 //go:generate gonectr generate 的注释行
func FindFirstGoGenerateLine(dir string) (targetPath string, targetNumber int, targetContent string, err error) {
	// 正则表达式匹配 //go:generate gonectr generate 开头的注释行，忽略前后空格
	re := regexp.MustCompile(`^\s*//\s*go:generate\s+gonectr\s+generate`)

	var End = errors.New("end")

	// 遍历目录中的所有文件
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 .go 文件
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			// 查找文件中的匹配行
			line, lineNum, err := findGenerateLineInFile(path, re)
			if err != nil {
				return err
			}
			if line != "" {
				targetPath = path
				targetNumber = lineNum
				targetContent = line

				// 找到匹配的行，返回文件路径、行号和内容
				return End
			}
		}
		return nil
	})

	// 如果找到了匹配的行，err会被设置为"found"，提取该信息
	if err != nil {
		if errors.Is(End, err) {
			return targetPath, targetNumber, targetContent, nil
		}
		return
	}
	return
}

// findGenerateLineInFile 查找文件中的 //go:generate gonectr generate 注释行
func findGenerateLineInFile(filePath string, re *regexp.Regexp) (string, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			println("关闭文件出错:", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// 去除前后空格并检查是否匹配
		if re.MatchString(strings.TrimSpace(line)) {
			return line, lineNum, nil // 找到匹配行
		}
	}

	if err := scanner.Err(); err != nil {
		return "", 0, err
	}

	return "", 0, nil // 没有找到匹配行
}

func Command(command string, args []string) error {
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			fmt.Println("Error copying stdout:", err)
		}
	}()

	go func() {
		if _, err := io.Copy(os.Stderr, stderr); err != nil {
			fmt.Println("Error copying stderr:", err)
		}
	}()

	return cmd.Wait()
}

// GetExternalIP 获取用户的外网 IP 地址
func GetExternalIP() (string, error) {
	// 创建带有超时的 HTTP 客户端
	var httpClient = &http.Client{
		Timeout: 3 * time.Second, // 设置超时时间为 5 秒
	}

	resp, err := httpClient.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", fmt.Errorf("failed to get external IP: %w", err)
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(ip), nil
}

// IsIPInChina 通过 ip-api.com 检查 IP 是否在国内
func IsIPInChina(ip string) (bool, error) {
	// 创建带有超时的 HTTP 客户端
	var httpClient = &http.Client{
		Timeout: 3 * time.Second, // 设置超时时间为 5 秒
	}

	resp, err := httpClient.Get("https://ip-api.com/json/" + ip)
	if err != nil {
		return false, fmt.Errorf("failed to get IP info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	// IPInfo 是 ip-api 的响应结构体
	type IPInfo struct {
		Country string `json:"country"`
	}

	var ipInfo IPInfo
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return false, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// 检查国家是否为中国
	return ipInfo.Country == "China", nil
}

func IsInChina() bool {
	// 获取用户外网 IP
	ip, err := GetExternalIP()
	if err != nil {
		return true
	}

	// 判断 IP 是否在国内
	inChina, err := IsIPInChina(ip)
	if err != nil {
		return true
	}
	return inChina
}
