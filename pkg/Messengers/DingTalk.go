package Messengers

import (
	"ForensicPro/utils"
	"fmt"
	"os"
	"path/filepath"
)

var messengerName = "DingTalk"

func DingTalkSave(path string) {

	// 构建存储数据库文件路径
	storageDBPath := utils.GetOperaPath("DingTalk\\globalStorage\\storage.db")

	// 检查存储数据库文件是否存在
	if _, err := os.Stat(storageDBPath); err == nil {
		// 构建其他相关文件路径
		storageDBShmPath := utils.GetOperaPath("DingTalk\\globalStorage\\storage.db-shm")
		storageDBWalPath := utils.GetOperaPath("DingTalk\\globalStorage\\storage.db-wal")

		// 构建目标文件夹路径
		targetPath := filepath.Join(path, messengerName) // 请替换 "path" 和 "MessengerName" 为实际路径和名称
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			fmt.Printf("创建目录失败: %v\n", err)
		}

		// 复制 storage.db 文件
		localdbdst := filepath.Join(targetPath, "storage.db")
		if err := utils.CopyFile(storageDBPath, localdbdst); err != nil {
			fmt.Printf("复制 storage.db 失败: %v\n", err)
		}

		// 检查并复制 storage.db-shm 文件
		localdbshmdst := filepath.Join(targetPath, "storage.db-shm")
		if _, err := os.Stat(storageDBShmPath); err == nil {
			if err := utils.CopyFile(storageDBShmPath, localdbshmdst); err != nil {
				fmt.Printf("复制 storage.db-shm 失败: %v\n", err)
			}
		}

		// 检查并复制 storage.db-wal 文件
		localdbwaldst := filepath.Join(targetPath, "storage.db-wal")
		if _, err := os.Stat(storageDBWalPath); err == nil {
			if err := utils.CopyFile(storageDBWalPath, localdbwaldst); err != nil {
				fmt.Printf("复制 storage.db-wal 失败: %v\n", err)
			}
		}
	} else {
		fmt.Printf("存储数据库文件不存在: %v\n", err)
	}
	fmt.Println("DingTalk 取证结束")

}
