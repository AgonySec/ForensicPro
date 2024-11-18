package Mails

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var FoxmailName = "Foxmail"

// 获取 Foxmail 安装路径
func getInstallPath() string {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Classes\Foxmail.url.mailto\Shell\open\command`, registry.READ)
	if err != nil {
		log.Printf("Failed to open registry key: %v", err)
		return ""
	}
	defer key.Close()

	values, _, err := key.GetStringValue("")
	if err != nil {
		log.Printf("Failed to get registry value: %v", err)
		return ""
	}

	text := strings.Replace(values, "\"", "", -1)
	index := strings.LastIndex(text, "Foxmail.exe")
	if index != -1 {
		text = text[:index]
	}

	return text
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return true
	}
	return false
}

// 复制文件
func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", src, err)
	}
	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", dst, err)
	}
	return nil
}

// 递归复制目录
func copyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking the path %s: %w", path, err)
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("error getting relative path for %s: %w", path, err)
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyFile(path, targetPath)
	})
}

// 保存 Foxmail 数据
func FoxmailSave(path string) {
	fmt.Println("开始取证Foxmail邮箱")
	installPath := getInstallPath()
	if installPath == "" {
		fmt.Println("未找到 Foxmail 安装路径")
		return
	}

	storagePath := filepath.Join(installPath, "Storage")
	if !dirExists(storagePath) {
		fmt.Println("未找到 Foxmail 存储目录")
		return
	}

	targetDir := filepath.Join(path, FoxmailName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Fatalf("Failed to create target directory %s: %v", targetDir, err)
	}

	err := filepath.Walk(storagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking the path %s: %w", path, err)
		}

		if info.IsDir() && filepath.Base(path) == "Accounts" {
			parentDir := filepath.Base(filepath.Dir(path))
			destination := filepath.Join(targetDir, parentDir, "Accounts")
			if err := copyDirectory(path, destination); err != nil {
				return fmt.Errorf("failed to copy directory %s to %s: %w", path, destination, err)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to copy directories: %v\n", err)
		return
	}

	fmStorageListPath := filepath.Join(installPath, "FMStorage.list")
	if _, err := os.Stat(fmStorageListPath); err == nil {
		if err := copyFile(fmStorageListPath, filepath.Join(targetDir, "FMStorage.list")); err != nil {
			fmt.Printf("Failed to copy FMStorage.list: %v\n", err)
		}
	}

	fmt.Println("Foxmail邮箱取证结束")
}
