package FTPS

import (
	"ForensicPro/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var SnowflakeName = "Snowflake"

func SnowflakeSave(path string) {
	// 获取用户主目录
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		fmt.Println("无法获取 USERPROFILE 环境变量")
		return
	}

	// 构建源文件路径
	sourcePath := filepath.Join(userProfile, "snowflake-ssh", "session-store.json")

	// 检查文件是否存在
	if _, err := os.Stat(sourcePath); err == nil {

		// 构建目标文件路径
		targetPath := filepath.Join(path, SnowflakeName)
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}
		// 复制文件
		if err := utils.CopyFile(sourcePath, targetPath+"session-store.json"); err != nil {
			fmt.Printf("复制文件失败: %v\n", err)
			return
		}

	}

}
