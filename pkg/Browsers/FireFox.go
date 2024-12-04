package Browsers

import (
	"ForensicPro/utils"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var masterPassword = ""

var FirefoxPath = GetOperaPath("Mozilla\\Firefox\\Profiles")
var FirefoxBrowserName = "FireFox"

func FirefoxPasswords() (string, error) {
	// todo 本机firefox解密不了。。。
	return "", nil
}
func FirefoxHistory(targetDir string) (string, error) {
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

		utils.CopyFile(filePath, targetDir+"\\"+"places.sqlite")

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
func FirefoxCookies(targetDir string) ([][]string, error) {

	var CSVData [][]string

	// 获取 BrowserPath 下的所有目录
	directories, err := os.ReadDir(FirefoxPath)
	if err != nil {
		fmt.Printf("Failed to read directory %s: %v\n", BrowserPath, err)
		return nil, nil
	}
	CSVData = append(CSVData, []string{"host", "name", "value"})
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
		utils.CopyFile(filePath, targetDir+"\\"+"cookies.sqlite")

		CSVData, _ = utils.ReadSQLiteDB(filePath, "SELECT host,name,value FROM moz_cookies ", CSVData)

		//builder.WriteString(result)
	}
	return CSVData, nil
}
func FirefoxSave(path string) {
	if _, err := os.Stat(FirefoxPath); os.IsNotExist(err) {
		return
	}
	targetDir := filepath.Join(path, FirefoxBrowserName)
	os.MkdirAll(targetDir, os.ModePerm)
	history, _ := FirefoxHistory(targetDir)
	books, _ := FirefoxBooks()
	cookies, _ := FirefoxCookies(targetDir)
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
	if len(cookies) > 1 {
		outputFile := FirefoxBrowserName + "_cookies.csv"
		if err := utils.WriteDataToCSV(targetDir+"\\"+outputFile, cookies); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

	}
	if passwords != "" {
		outputFile := FirefoxBrowserName + "_passwords.txt"
		if err := utils.WriteToFile(passwords, targetDir+"\\"+outputFile); err != nil {
			fmt.Println()
		}

	}
	// 获取目录列表
	directories, err := ioutil.ReadDir(FirefoxPath)
	if err != nil {
		fmt.Printf("读取目录时出错: %v\n", err)
		return
	}

	// 遍历目录
	for _, dir := range directories {
		path2 := filepath.Join(FirefoxPath, dir.Name())

		// 检查 storage-sync-v2.sqlite 文件是否存在
		sqlitePath := filepath.Join(path2, "storage-sync-v2.sqlite")
		if _, err := os.Stat(sqlitePath); err == nil {
			// 复制 storage-sync-v2.sqlite 文件
			destSqlitePath := filepath.Join(targetDir, "storage-sync-v2.sqlite")
			if err := utils.CopyFile(sqlitePath, destSqlitePath); err != nil {
				fmt.Printf("复制文件时出错: %v\n", err)
				continue
			}

			// 检查并复制 storage-sync-v2.sqlite-shm 文件
			shmPath := filepath.Join(path2, "storage-sync-v2.sqlite-shm")
			if _, err := os.Stat(shmPath); err == nil {
				destShmPath := filepath.Join(targetDir, "storage-sync-v2.sqlite-shm")
				if err := utils.CopyFile(shmPath, destShmPath); err != nil {
					fmt.Printf("复制文件时出错: %v\n", err)
				}
			}

			// 检查并复制 storage-sync-v2.sqlite-wal 文件
			walPath := filepath.Join(path2, "storage-sync-v2.sqlite-wal")
			if _, err := os.Stat(walPath); err == nil {
				destWalPath := filepath.Join(targetDir, "storage-sync-v2.sqlite-wal")
				if err := utils.CopyFile(walPath, destWalPath); err != nil {
					fmt.Printf("复制文件时出错: %v\n", err)
				}
			}

			// 找到并复制了所需的文件，跳出循环
			break
		}
	}
	fmt.Println("Firefox浏览器取证结束")

}
