package Mails

import (
	"ForensicPro/utils"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var mailMasterName = "MailMaster"

func getMailMasterDataPath() []string {
	var list []string
	dataPath := utils.GetLocalAppDataPath("Netease\\MailMaster\\data\\app.db")
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		fmt.Printf("文件 %s 不存在\n", dataPath)
		return list
	}
	// 创建一个临时文件
	tempFile, err := os.CreateTemp("", "temp_sqlite-*.db")
	if err != nil {
		fmt.Printf("创建临时文件失败: %v\n", err)
		return nil
	}
	defer os.Remove(tempFile.Name()) // 确保临时文件在函数结束时被删除

	// 将数据库文件复制到临时文件
	if err := utils.CopyFile(dataPath, tempFile.Name()); err != nil {
		fmt.Printf("复制文件失败: %v\n", err)
		return nil
	}

	// 检查临时文件是否正确创建
	if _, err := os.Stat(tempFile.Name()); os.IsNotExist(err) {
		fmt.Printf("临时文件 %s 不存在\n", tempFile.Name())
		return nil
	}
	// 打开临时数据库文件
	db, err := sql.Open("sqlite3", tempFile.Name())
	//db, err := sql.Open("sqlite3", dataPath)
	if err != nil {
		fmt.Printf("打开数据库文件失败: %v\n", err)
		return nil
	}
	defer db.Close()

	// 查询 Account 表中的 DataPath 字段
	rows, err := db.Query("SELECT DataPath FROM Account")
	if err != nil {
		fmt.Printf("查询数据库失败: %v\n", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var dataPath string
		if err := rows.Scan(&dataPath); err != nil {
			fmt.Printf("扫描数据失败: %v\n", err)
			continue
		}

		list = append(list, dataPath)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("遍历数据失败: %v\n", err)
		return nil
	}

	return list
}
func MailMasterSave(path string) {
	mailMasterPath := utils.GetLocalAppDataPath("Netease\\MailMaster\\data")
	if _, err := os.Stat(mailMasterPath); err == nil {
		targetPath := filepath.Join(path, mailMasterName)
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}
		dataPath := getMailMasterDataPath()
		for _, data := range dataPath {
			// 获取路径中的最后一个目录名
			dirName := filepath.Base(data)
			// 构建目标路径
			// 删除 _4332 部分
			parts := strings.Split(dirName, "_")
			if len(parts) > 1 {
				dirName = strings.Join(parts[:len(parts)-1], "_")
			}
			destPath := filepath.Join(targetPath, dirName)
			err := utils.CopyDirectory(data, destPath)
			if err != nil {
				fmt.Printf("Error copying directory: %v\n", err)
				return
			}

		}

	}
	fmt.Println(mailMasterName + " 邮箱信息已保存")
}
