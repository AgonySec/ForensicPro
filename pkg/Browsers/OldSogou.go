package Browsers

import (
	"ForensicPro/utils"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var SogouName = "OldSogouExplorer"
var SogouPath = GetOperaPath("SogouExplorer")
var sougoDefaultPath = filepath.Join(SogouPath, "\\Webkit\\Default")

func SogouCookies() string {
	var builder strings.Builder
	sougoMasterKeyPath := filepath.Join(SogouPath, "\\Webkit")
	sougoMasterkey, _ := utils.GetMasterKey(sougoMasterKeyPath)
	if sougoMasterkey == nil {
		return ""
	}
	text := filepath.Join(sougoDefaultPath, "Cookies")
	if _, err := os.Stat(text); os.IsNotExist(err) {
		return ""
	}
	// 创建临时文件
	tempFileName, err := os.CreateTemp("", "temp-*.db")
	if err != nil {
		return ""
	}
	defer os.Remove(tempFileName.Name()) // 确保临时文件在函数结束时被删除

	// 复制文件
	if err := utils.CopyFile(text, tempFileName.Name()); err != nil {
		return ""
	}
	// 连接 SQLite 数据库
	db, err := sql.Open("sqlite3", tempFileName.Name())
	if err != nil {
		return ""
	}
	defer db.Close()

	// 查询 UserRankUrl 表
	rows, err := db.Query("SELECT host_key,name,encrypted_value FROM cookies")
	if err != nil {
		return ""
	}
	defer rows.Close()
	// 遍历查询结果
	for rows.Next() {
		//var hostKey, name, encryptedValue string
		var hostKey, name string
		var encryptedValue []byte
		if err := rows.Scan(&hostKey, &name, &encryptedValue); err != nil {
			return ""
		}

		// 解密 encrypted_value
		decryptedValue, err := utils.DecryptAESGCM(encryptedValue, sougoMasterkey)
		if err != nil {
			continue
		}

		// 构建输出字符串
		builder.WriteString(fmt.Sprintf("[%s] \t {%s}={%s}\n", hostKey, name, decryptedValue))
	}

	if err := rows.Err(); err != nil {
		return ""
	}
	return builder.String()
}
func SogouHistory() string {
	var builder strings.Builder
	// 构建 HistoryUrl3.db 的路径
	text := filepath.Join(SogouPath, "HistoryUrl3.db")

	// 检查文件是否存在
	if _, err := os.Stat(text); os.IsNotExist(err) {
		return ""
	}

	// 创建临时文件
	tempFileName, err := os.CreateTemp("", "temp-*.db")
	if err != nil {
		return ""
	}
	defer os.Remove(tempFileName.Name()) // 确保临时文件在函数结束时被删除

	// 复制文件
	if err := utils.CopyFile(text, tempFileName.Name()); err != nil {
		return ""
	}
	// 连接 SQLite 数据库
	db, err := sql.Open("sqlite3", tempFileName.Name())
	if err != nil {
		return ""
	}
	defer db.Close()

	// 查询 UserRankUrl 表
	rows, err := db.Query("SELECT id FROM UserRankUrl")
	if err != nil {
		return ""
	}
	defer rows.Close()
	// 遍历查询结果
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return ""
		}
		builder.WriteString(id)
		builder.WriteString("\n")
	}

	if err := rows.Err(); err != nil {
		return ""
	}

	return builder.String()

}

func SogouSave(path string) {

	if _, err := os.Stat(SogouPath); err == nil {
		sogouHistory := SogouHistory()
		sogouCookies := SogouCookies()
		targetPath := filepath.Join(path, SogouName)
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}
		if sogouHistory != "" {
			outputFile := SogouName + "_history.txt"
			utils.WriteToFile(sogouHistory, targetPath+"\\"+outputFile)
		}
		if sogouCookies != "" {
			outputFile := SogouName + "_cookies.txt"
			utils.WriteToFile(sogouCookies, targetPath+"\\"+outputFile)
		}
		// 检查并复制Local Storage文件
		sougoLSPath := filepath.Join(sougoDefaultPath, "Local Storage")
		if _, err := os.Stat(sougoLSPath); err == nil {
			destPath := filepath.Join(targetPath, "Local Storage")
			if err := utils.CopyDirectory(sougoLSPath, destPath); err != nil {
				fmt.Printf("复制文件时出错: %v\n", err)
			}
		}
	}
	fmt.Println("搜狗浏览器取证结束")
}
