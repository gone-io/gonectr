package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

func main() {
	// 定义选项列表和存储变量
	options := []string{"Go", "Python", "JavaScript", "Rust"}
	var selected []string

	// 配置多选提示
	prompt := &survey.MultiSelect{
		Message: "选择你擅长的编程语言:",
		Options: options,
		Default: []string{"Go"}, // 可设置默认选中项
	}

	// 执行交互
	if err := survey.AskOne(prompt, &selected); err != nil {
		fmt.Println("选择错误:", err)
		return
	}

	fmt.Printf("你选择了: %v\n", selected)
}
