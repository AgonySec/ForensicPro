package utils

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"io/ioutil"
	"os"
	"path/filepath"
)

// WriteToFile 函数
func WriteToFile(content string, filePath string) error {
	// 将内容写入指定的文件
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}
func ReadFileContent(filePath string) (string, error) {
	// 读取文件的全部内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("无法读取文件: %w", err)
	}

	// 将字节切片转换为字符串
	content := string(data)

	return content, nil
}
func GetShortcutTargetPath(shortcutPath string) string {

	err := ole.CoInitialize(0)
	if err != nil {
		return ""

	}
	defer ole.CoUninitialize()

	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		fmt.Println("创建WScript.Shell对象失败:", err)
		return ""

	}
	defer oleShellObject.Release()

	shellObject, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		fmt.Println("获取IDispatch接口失败:", err)
		return ""

	}
	defer shellObject.Release()

	shortcut, err := oleutil.CallMethod(shellObject, "CreateShortcut", shortcutPath)
	if err != nil {
		return ""

	}

	targetPath, err := oleutil.GetProperty(shortcut.ToIDispatch(), "TargetPath")
	if err != nil {
		fmt.Println("获取目标路径失败:", err)
		return ""
	}
	fileName := filepath.Base(shortcutPath)
	result := fileName + ":			" + targetPath.ToString()
	return result

}
