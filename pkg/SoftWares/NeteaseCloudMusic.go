package SoftWares

import (
	"ForensicPro/utils"
	"log"
	"os"
	"path/filepath"
)

var CloudMusicName = "NeteaseCloudMusic"

func NeteaseCloudMusicSave(path string) {

	infoPath := utils.GetLocalAppDataPath("Netease\\CloudMusic\\info")

	if _, err := os.Stat(infoPath); err == nil {
		// 构建目标文件夹路径
		text, _ := utils.ReadFileContent(infoPath)
		targetPath := filepath.Join(path, CloudMusicName)
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}
		utils.WriteToFile(" [InternetShortcut]\r\nURL=https://music.163.com/#/user/home?id="+text, targetPath+"\\userinfo.url")
	}
}
