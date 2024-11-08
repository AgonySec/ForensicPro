package SoftWares

import (
	"ForensicPro/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var VSCodeName = "VSCode"

func VSCodeSave(path string) {

	VSCodeHistory := utils.GetOperaPath("Code\\User\\History")

	if _, err := os.Stat(VSCodeHistory); err == nil {
		// 构建目标文件夹路径
		targetPath := filepath.Join(path, VSCodeName)
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}
		utils.CopyDirectory(VSCodeHistory, targetPath+"\\History")
	}
	fmt.Println("VSCode取证结束")
}
