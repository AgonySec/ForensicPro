package SoftWares

import (
	"ForensicPro/utils"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var NavicatName = "Navicat"
var databaseTypes = map[string]string{
	"Navicat":        "MySql",
	"NavicatMSSQL":   "SQL Server",
	"NavicatOra":     "Oracle",
	"NavicatPG":      "pgsql",
	"NavicatMARIADB": "MariaDB",
	"NavicatMONGODB": "MongoDB",
	"NavicatSQLite":  "SQLite",
}

func DecryptPwd(ciphertext string) string {
	cipher, err := utils.NewNavicat11Cipher()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	plaintext, err := cipher.DecryptString(ciphertext)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return plaintext
}

func GetNavicatInfo() string {
	var builder strings.Builder

	for databaseType, _ := range databaseTypes {
		// 打开注册表项
		key, err := registry.OpenKey(registry.CURRENT_USER, `Software\PremiumSoft\`+databaseType+`\Servers`, registry.READ)
		if err != nil {
			//log.Fatal(err)
			continue
		}
		defer key.Close()

		subKeyNames, err := key.ReadSubKeyNames(-1)
		if err != nil {
			return ""
		}
		builder.WriteString("DatabaseName: " + databaseTypes[databaseType])
		builder.WriteString("\n")

		for _, name := range subKeyNames {
			subKey, err := registry.OpenKey(key, name, registry.READ)
			if err != nil {
				continue
			}
			defer subKey.Close()
			Host, _, err := subKey.GetStringValue("Host")
			UserName, _, err := subKey.GetStringValue("UserName")
			Pwd, _, err := subKey.GetStringValue("Pwd")
			Pwd = DecryptPwd(Pwd)
			//fmt.Println(Host, UserName, Pwd)
			builder.WriteString("ConnectName: " + name + "\n")

			if Host != "" {
				builder.WriteString("Host: " + Host + "\n")
			}
			if UserName != "" {
				builder.WriteString("UserName: " + UserName + "\n")
			}
			if Pwd != "" {
				builder.WriteString("password: " + Pwd + "\n")
			}
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func NavicatSave(path string) {
	_, err := registry.OpenKey(registry.CURRENT_USER, `Software\PremiumSoft`, registry.READ)
	if err != nil {
		return
	}
	targetPath := filepath.Join(path, NavicatName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	text := GetNavicatInfo()
	err = utils.WriteToFile(text, targetPath+"\\"+NavicatName+".txt")
	if err != nil {
		return
	}
	fmt.Println("Navicat 取证结束")
}
