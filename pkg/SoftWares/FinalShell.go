package SoftWares

import (
	"ForensicPro/utils"
	_ "crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var FinalShellName = "FinalShell"

func GetFinalShellInfo(finalShellPath string) string {
	var builder strings.Builder
	files, err := ioutil.ReadDir(finalShellPath)
	if err != nil {
		fmt.Println("读取目录失败:", err)
		return ""
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), "_connect_config.json") {
			continue
		}

		filePath := filepath.Join(finalShellPath, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("读取文件失败:", err)
			continue
		}

		input := string(content)
		userName := extractValue(input, `"user_name":"(.*?)"`)
		password := extractValue(input, `"password":"(.*?)"`)
		host := extractValue(input, `"host":"(.*?)"`)
		port := extractValue(input, `"port":(.*?),`)

		builder.WriteString(fmt.Sprintf("host: %s\n", host))
		builder.WriteString(fmt.Sprintf("port: %s\n", port))
		builder.WriteString(fmt.Sprintf("user_name: %s\n", userName))
		//builder.WriteString(fmt.Sprintf("password: %s\n", password))
		builder.WriteString(fmt.Sprintf("password: %s\n", passDecode(password)))
		builder.WriteString("\n")
	}
	return builder.String()
}

func passDecode(input string) string {
	currentDir, _ := os.Getwd()

	exePath := filepath.Join(currentDir, "finalshellDC.exe")
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		log.Fatalf("文件 %s 不存在", exePath)
	}

	// 获取当前目录
	//fmt.Printf("当前目录: %s\n", currentDir)
	arg := input

	cmd := exec.Command(exePath, arg)

	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("执行命令时出错: %v", err)
	}
	//打印输出结果
	//fmt.Println(string(output))
	return string(output)
}
func extractValue(input, pattern string) string {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func FinalShellSave(path string) {

	finalShellPath := utils.GetLocalAppDataPath("finalshell\\conn")

	if _, err := os.Stat(finalShellPath); err == nil {
		info := GetFinalShellInfo(finalShellPath)
		if info != "" {
			targetPath := filepath.Join(path, FinalShellName)
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				log.Fatalf("创建目录失败: %v", err)
			}
			err := utils.WriteToFile(info, targetPath+"\\FinalShell.txt")
			if err != nil {
				return
			}
		}

	}
	fmt.Println("FinalShell 数据保存成功")

}
