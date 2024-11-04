package Browsers

import (
	"ForensicPro/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var masterPassword = ""

var FirefoxPath = getOperaPath("Mozilla\\Firefox\\Profiles")
var FirefoxBrowserName = "FireFox"

func FirefoxPasswords() (string, error) {
	// todo
	return "", nil
}
func FirefoxHistory() (string, error) {
	var builder strings.Builder
	// 获取 BrowserPath 下的所有目录
	directories, err := os.ReadDir(FirefoxPath)
	if err != nil {
		fmt.Printf("Failed to read directory %s: %v\n", BrowserPath, err)
		return "", nil
	}
	for _, dirEntry := range directories {
		if !dirEntry.IsDir() {
			continue
		}
		dirPath := filepath.Join(FirefoxPath, dirEntry.Name())
		filePath := filepath.Join(dirPath, "places.sqlite")

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}
		result, err := utils.ReadSQLiteDB_url(filePath, "SELECT url FROM moz_places")
		if err != nil {
			fmt.Println("Error reading database:", err)
			continue
		}
		builder.WriteString(result)
	}
	return builder.String(), nil
}

func FirefoxBooks() (string, error) {
	var builder strings.Builder
	// 获取 BrowserPath 下的所有目录
	directories, err := os.ReadDir(FirefoxPath)
	if err != nil {
		fmt.Printf("Failed to read directory %s: %v\n", BrowserPath, err)
		return "", nil
	}
	for _, dirEntry := range directories {
		if !dirEntry.IsDir() {
			continue
		}
		dirPath := filepath.Join(FirefoxPath, dirEntry.Name())
		filePath := filepath.Join(dirPath, "places.sqlite")

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}
		result, _ := utils.ReadSQLiteDB_url2(filePath, "SELECT fk FROM moz_bookmarks ")

		builder.WriteString(result)
	}
	return builder.String(), nil
}
func FirefoxCookies() (string, error) {
	var builder strings.Builder
	// 获取 BrowserPath 下的所有目录
	directories, err := os.ReadDir(FirefoxPath)
	if err != nil {
		fmt.Printf("Failed to read directory %s: %v\n", BrowserPath, err)
		return "", nil
	}
	for _, dirEntry := range directories {
		if !dirEntry.IsDir() {
			continue
		}
		dirPath := filepath.Join(FirefoxPath, dirEntry.Name())
		filePath := filepath.Join(dirPath, "cookies.sqlite")

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}
		result, _ := utils.ReadSQLiteDB(filePath, "SELECT host,name,value FROM moz_cookies ")

		builder.WriteString(result)
	}
	return builder.String(), nil
}
func FirefoxSave(path string) {
	if _, err := os.Stat(FirefoxPath); os.IsNotExist(err) {
		return
	}
	targetDir := filepath.Join(path, FirefoxBrowserName)
	os.MkdirAll(targetDir, os.ModePerm)
	history, _ := FirefoxHistory()
	books, _ := FirefoxBooks()
	cookies, _ := FirefoxCookies()
	passwords, _ := FirefoxPasswords()
	if history != "" {
		// 将历史记录写入到文件
		outputFile := FirefoxBrowserName + "_history.txt"
		if err := utils.WriteToFile(history, targetDir+"\\"+outputFile); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	if books != "" {
		// 将书签写入到文件
		outputFile := FirefoxBrowserName + "_books.txt"
		if err := utils.WriteToFile(books, targetDir+"\\"+outputFile); err != nil {
			fmt.Println()
		}
	}
	if cookies != "" {
		outputFile := FirefoxBrowserName + "_cookies.txt"
		if err := utils.WriteToFile(cookies, targetDir+"\\"+outputFile); err != nil {
			fmt.Println()
		}

	}
	if passwords != "" {
		outputFile := FirefoxBrowserName + "_passwords.txt"
		if err := utils.WriteToFile(passwords, targetDir+"\\"+outputFile); err != nil {
			fmt.Println()
		}

	}

	fmt.Println("Firefox浏览器取证结束")

}
