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

func GetFinalShellInfo(finalShellPath string) [][]string {
	files, err := ioutil.ReadDir(finalShellPath)
	if err != nil {
		fmt.Println("读取目录失败:", err)
		return nil
	}
	var CSVData [][]string
	CSVData = append(CSVData, []string{"host", "port", "username", "password"})
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

		CSVData = append(CSVData, []string{host, port, userName, passDecode(password)})
	}
	return CSVData
}

func passDecode(input string) string {
	currentDir, _ := os.Getwd()

	exePath := filepath.Join(currentDir, "finalshellDC.exe")
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		return input
	}

	arg := input

	cmd := exec.Command(exePath, arg)

	output, err := cmd.Output()
	if err != nil {
		return input
	}

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
		if len(info) > 1 {
			targetPath := filepath.Join(path, FinalShellName)
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				log.Fatalf("创建目录失败: %v", err)
			}
			err := utils.WriteDataToCSV(targetPath+"\\FinalShell.csv", info)
			if err != nil {
				return
			}
		}

	}
	fmt.Println("FinalShell 取证结束")

}
