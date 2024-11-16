package SystemInfos

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var RegEditName = "RegEdit"

func RegEditSave(path string) {
	// 创建目标目录
	targetPath := filepath.Join(path, RegEditName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		return
	}

	// 定义注册表根键
	regHives := []string{"HKEY_LOCAL_MACHINE", "HKEY_CURRENT_USER", "HKEY_CLASSES_ROOT", "HKEY_USERS", "HKEY_CURRENT_CONFIG"}

	for _, hive := range regHives {
		// 构建导出命令
		cmd := exec.Command("reg", "export", hive, filepath.Join(targetPath, hive+".reg"), "/y")

		// 执行命令
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("导出注册表 %s 失败: %v, 错误信息: %s", hive, err, stderr.String())
			return
		}
	}

	fmt.Println("所有注册表信息已成功导出")
}
