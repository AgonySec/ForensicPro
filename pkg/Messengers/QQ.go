package Messengers

import (
	"ForensicPro/utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

var messengerQQName = "QQ"

func isQQAccount(s string) bool {
	// 定义正则表达式，匹配5到12位数字
	re := regexp.MustCompile(`^\d{5,12}$`)
	return re.MatchString(s)
}
func QQSave(path string) {
	var builder strings.Builder

	userDataInfoPath := "C:\\Users\\Public\\Documents\\Tencent\\QQ\\UserDataInfo.ini"
	// 获取当前用户
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("获取当前用户失败:", err)
	}

	// 构建文档目录路径
	TencentFilesPath := filepath.Join(currentUser.HomeDir, "Documents", "Tencent Files")

	if _, err := os.Stat(userDataInfoPath); err == nil {
		targetPath := filepath.Join(path, messengerQQName) // 请替换 "path" 和 "MessengerName" 为实际路径和名称
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}

		// 读取扩展目录下的所有子目录
		directories, err := ioutil.ReadDir(TencentFilesPath)
		if err != nil {
			return // 如果读取失败，跳过
		}
		builder.WriteString("All QQ number:\n")

		for _, dir := range directories {
			if dir.IsDir() { // 确保是目录
				dirName := dir.Name()
				if isQQAccount(dirName) {
					builder.WriteString(dirName + "\n")
				}
			}
		}
		qqNumbers := builder.String()
		if qqNumbers != "" {
			outputFile := messengerQQName + "_numbers.txt"
			if err := utils.WriteToFile(qqNumbers, targetPath+"\\"+outputFile); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}

		}
	}
	fmt.Println("QQ取证结束.")

}
