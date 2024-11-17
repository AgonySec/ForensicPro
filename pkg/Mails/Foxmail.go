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

func getInstallPath() string {
	// 打开注册表项
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Classes\Foxmail.url.mailto\Shell\open\command`, registry.READ)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer key.Close()

	// 读取默认值
	values, _, err := key.GetStringValue("")
	if err != nil {
		return ""
	}

	// 处理字符串
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

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return true
	}
	return false
}
func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}
	return nil
}
func CopyDirectory(src, dst string, recursive bool) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if recursive {
				err = CopyDirectory(srcPath, dstPath, recursive)
				if err != nil {
					return err
				}
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyDirectory recursively copies a source directory to a destination
func CopyDirectory2(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyFile(path, targetPath)
	})
}

func FoxmailSave(path string) {
	installPath := getInstallPath()
	if dirExists(installPath) && dirExists(filepath.Join(installPath, "Storage")) {
		targetDir := filepath.Join(path, FoxmailName)
		os.MkdirAll(targetDir, 0755)

		storagePath := filepath.Join(installPath, "Storage")
		//directories, err := ioutil.ReadDir(storagePath)
		if _, err := os.Stat(storagePath); os.IsNotExist(err) {
			return
		}
		err := filepath.Walk(storagePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && filepath.Base(path) == "Accounts" {
				parentDir := filepath.Base(filepath.Dir(path))
				destination := filepath.Join(targetDir, parentDir, "Accounts")
				CopyDirectory2(path, destination)
			}

			return nil
		})
		if err != nil {
			fmt.Println("Failed to copy directories:", err)
			return
		}
		fmStorageListPath := filepath.Join(installPath, "FMStorage.list")
		if _, err := os.Stat(fmStorageListPath); err == nil {
			err = copyFile(fmStorageListPath, filepath.Join(targetDir, "FMStorage.list"))
			if err != nil {
				fmt.Println("Failed to copy FMStorage.list:", err)
			}
		}
	}
	fmt.Println("Foxmail邮箱取证结束")
}
