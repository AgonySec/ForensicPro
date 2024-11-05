package SystemInfos

import (
	"ForensicPro/utils"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var InstalledAppName = "InstalledApp"

func GetInfo() string {
	var result strings.Builder

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `Software\Microsoft\Windows\CurrentVersion\Uninstall`, registry.READ)
	if err != nil {
		return ""
	}
	defer key.Close()

	subKeyNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return ""
	}

	for _, name := range subKeyNames {
		subKey, err := registry.OpenKey(key, name, registry.READ)
		if err != nil {
			continue
		}
		defer subKey.Close()

		displayName, _, err := subKey.GetStringValue("DisplayName")
		if err == nil && displayName != "" && !strings.Contains(displayName, "Windows") {
			result.WriteString(displayName + "\n")
		}
	}

	return result.String()
}
func InstalledAppSave(path string) {

	result := GetInfo()
	targetPath := filepath.Join(path, InstalledAppName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(result, targetPath+"\\"+InstalledAppName+".txt")
}
