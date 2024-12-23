package FTPS

import (
	"ForensicPro/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var FileZillaName = "FileZilla"

func FileZillaSave(path string) {
	FileZillaXMLPath := utils.GetOperaPath("FileZilla\\recentservers.xml")
	// 检查存储数据库文件是否存在
	if _, err := os.Stat(FileZillaXMLPath); err == nil {
		// 构建目标文件夹路径
		targetPath := filepath.Join(path, FileZillaName)
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}
		ftpdst := filepath.Join(targetPath, FileZillaName+"_txt")
		utils.CopyFile(FileZillaXMLPath, ftpdst)

	}
	fmt.Println("FileZilla 取证结束.")
}
