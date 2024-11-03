package utils

import (
	"fmt"
	"io/ioutil"
	"os"
)

// WriteToFile 函数
func WriteToFile(content string, filePath string) error {
	// 将内容写入指定的文件
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}
func ReadFileContent(filePath string) (string, error) {
	// 读取文件的全部内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("无法读取文件: %w", err)
	}

	// 将字节切片转换为字符串
	content := string(data)

	return content, nil
}
