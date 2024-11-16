package SystemInfos

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var WindowsLogsName = "WindowsLogs"

// WindowsLogsSave 导出系统日志并保存到指定路径下的文件中
func WindowsLogsSave(path string) {
	// 指定日志文件夹路径
	logsFolderPath := "C:\\Windows\\System32\\winevt\\Logs"
	// 写入文件
	targetPath := filepath.Join(path, WindowsLogsName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}

	// 遍历日志文件夹中的所有文件
	err := filepath.Walk(logsFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 如果是文件，将文件内容复制到目标路径
		if !info.IsDir() {
			relativePath, err := filepath.Rel(logsFolderPath, path)
			if err != nil {
				return err
			}
			targetFilePath := filepath.Join(targetPath, relativePath)

			// 确保目标文件夹存在
			if err := os.MkdirAll(filepath.Dir(targetFilePath), os.ModePerm); err != nil {
				return err
			}

			// 打开源文件
			sourceFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer sourceFile.Close()

			// 创建目标文件
			targetFile, err := os.Create(targetFilePath)
			if err != nil {
				return err
			}
			defer targetFile.Close()

			// 将文件内容复制到目标文件
			_, err = io.Copy(targetFile, sourceFile)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("遍历日志文件夹时出错:", err)
		return
	}
	fmt.Println("日志文件已成功全部导出")

}
