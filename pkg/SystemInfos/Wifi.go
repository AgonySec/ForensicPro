package SystemInfos

import (
	"ForensicPro/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var WifiInfoName = "Wifi"

func WifiInfoSave(path string) {

	// 执行netsh命令获取WiFi信息
	cmd := exec.Command("netsh", "wlan", "show", "profiles")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}
	if len(output) == 0 {
		return
	}

	// 将GBK编码的输出转换为UTF-8编码
	outputStr, err := utils.ConvertGBKToUTF8(output)

	// 将输出转换为字符串并分割为行
	profiles := strings.Split(outputStr, "\n")
	// 准备存储WiFi数据的切片
	var wifiData [][]string
	// 添加CSV头
	wifiData = append(wifiData, []string{"SSID", "Password"})

	// 遍历每个WiFi配置文件
	for _, profile := range profiles {
		// 查找WiFi名称
		if strings.Contains(profile, "所有用户配置文件") {
			// 提取WiFi名称
			name := strings.Split(profile, ":")[1]
			name = strings.TrimLeft(name, " ")
			name = strings.TrimRight(name, "\r") // 去除名称右侧的空格和回车符
			// 获取WiFi密码
			cmd = exec.Command("netsh", "wlan", "show", "profile", name, "key=clear")
			output, err = cmd.Output()
			outputStr, err := utils.ConvertGBKToUTF8(output)
			if err != nil {
				fmt.Println("获取WiFi密码失败:", err)
				continue
			}

			// 查找密码
			password := ""
			for _, line := range strings.Split(outputStr, "\n") {
				if strings.Contains(line, "关键内容") { // 关键内容是指密码所在的行
					password = strings.Split(line, ":")[1]
					password = strings.TrimLeft(password, " ")
					password = strings.TrimRight(password, "\r")
					break
				}
			}

			// 添加WiFi数据到切片
			wifiData = append(wifiData, []string{name, password})

		}
	}
	// 写入文件
	targetPath := filepath.Join(path, WifiInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	csvFilePath := filepath.Join(targetPath, "wifi_passwords.csv")

	if len(wifiData) > 1 {
		// 调用函数将数据写入CSV文件
		if err := utils.WriteDataToCSV(csvFilePath, wifiData); err != nil {
			fmt.Println("Error writing WiFi data to CSV:", err)
			return
		}
	}

	fmt.Println("WiFi 取证结束")

}
