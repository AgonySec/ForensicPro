package Browsers

import (
	"ForensicPro/utils"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func IEHistory() (string, error) {
	var builder strings.Builder

	// 打开注册表项
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Internet Explorer\TypedURLs`, registry.READ)
	if err != nil {
		log.Fatal(err)
	}
	defer key.Close()

	// 读取所有值
	values, err := key.ReadValueNames(-1) // 传入 -1 读取所有值
	if err != nil {
		log.Fatal(err)
	}

	// 输出每个值
	for _, valueName := range values {
		if strings.Contains(valueName, "url") {
			value, _, err := key.GetStringValue(valueName)
			if err != nil {
				log.Printf("无法读取值 %s: %v\n", valueName, err)
				continue
			}
			builder.WriteString(value)
		}
	}

	return builder.String(), nil
}

func IEBooks() (string, error) {
	var builder strings.Builder

	// 获取收藏夹目录
	username := os.Getenv("USERNAME")
	// 获取收藏夹目录
	favoritesPath := filepath.Join("C:\\Users", username, "Favorites")

	// 递归获取所有 .url 文件
	files, err := filepath.Glob(filepath.Join(favoritesPath, "*.url"))
	if err != nil {
		return "", fmt.Errorf("获取文件列表失败: %w", err)
	}

	// 遍历每个文件
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			// 读取文件内容
			content, err := ioutil.ReadFile(file)
			if err != nil {
				log.Printf("读取文件 %s 失败: %v", file, err)
				continue
			}

			// 使用正则表达式匹配 URL
			re := regexp.MustCompile(`URL=(.*?)\n`)
			match := re.FindStringSubmatch(string(content))

			// 添加到 StringBuilder
			//builder.WriteString(file)
			//builder.WriteString("\n")
			if len(match) > 1 {
				builder.WriteString(match[1])
			} else {
				builder.WriteString("\t" + "No URL found")
			}
			builder.WriteString("\n")
		}
	}

	return builder.String(), nil
}

func IEPasswords() string {
	return ""
}

func IESave(path string) {
	BrowserName := "IE"
	targetDir := filepath.Join(path, BrowserName)
	os.MkdirAll(targetDir, os.ModePerm)
	history, err := IEHistory()
	books, err := IEBooks()
	if err != nil {
		log.Printf("获取IE历史记录失败: %v", err)
		fmt.Println("获取IE历史记录失败，请检查日志文件以获取更多信息。")
		return
	}
	if history != "" {
		outputFile := BrowserName + "_history.txt"
		err := utils.WriteToFile(history, filepath.Join(targetDir, outputFile))
		if err != nil {
			log.Printf("写入文件失败: %v", err)
			fmt.Println("写入文件失败，请检查日志文件以获取更多信息。")
			return
		}
	}
	if books != "" {
		outputFile := BrowserName + "_books.txt"
		err := utils.WriteToFile(books, filepath.Join(targetDir, outputFile))
		if err != nil {
			log.Printf("写入文件失败: %v", err)
			fmt.Println("写入文件失败，请检查日志文件以获取更多信息。")
			return
		}
	}
	fmt.Println("IE浏览器取证结束")
}
