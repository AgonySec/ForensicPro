package SoftWares

import (
	"ForensicPro/utils"
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	_ "syscall"
	_ "unsafe"
)

var XmanagerName = "Xmanager"

var sessionFiles []string

func GetAllAccessibleFiles(rootPath string) {
	directories, err := ioutil.ReadDir(rootPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, dir := range directories {
		if dir.IsDir() {
			dirPath := filepath.Join(rootPath, dir.Name())
			GetAllAccessibleFiles(dirPath)
		}
	}

	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(rootPath, file.Name())
			if filepath.Ext(filePath) == ".xsh" || filepath.Ext(filePath) == ".xfp" {
				sessionFiles = append(sessionFiles, filePath)
			}
		}
	}
}

// DecryptSessions 解密会话文件
func DecryptSessions() [][]string {

	// 获取当前用户信息
	var currentUser, err = user.Current()
	if err != nil {
		fmt.Println("获取当前用户信息失败:", err)
		return nil
	}
	username := currentUser.Username
	// 对用户名进行切割
	// 对用户名进行切割
	parts := strings.Split(username, "\\")

	// 检查切割结果
	if len(parts) > 1 {
		// 取切割后的第二个部分
		username = parts[1]
	}
	sid, _ := utils.GetCurrentUserSID()
	//text := fmt.Sprintf("%d", currentUser)
	// 准备存储WiFi数据的切片
	var CSVData [][]string
	// 添加CSV头
	CSVData = append(CSVData, []string{"Session File", "Version", "Host", "UserName", "rawPass", "UserName", "Sid", "Decrypt rawPass"})
	for _, sessionFile := range sessionFiles {
		list := ReadConfigFile(sessionFile)
		if len(list) >= 4 {
			decryptPass := Decrypt(username, sid, list[3], strings.ReplaceAll(list[0], "\r", ""))
			CSVData = append(CSVData, []string{sessionFile, list[0], list[1], list[2], list[3], username, sid, decryptPass})
		}
	}
	return CSVData
}

func Decrypt(username, sid, rawPass, ver string) string {
	if strings.HasPrefix(ver, "5.0") || strings.HasPrefix(ver, "4") || strings.HasPrefix(ver, "3") || strings.HasPrefix(ver, "2") {
		array, _ := base64.StdEncoding.DecodeString(rawPass)
		key := sha256.New()
		key.Write([]byte("!X@s#h$e%l^l&"))
		keyHash := key.Sum(nil)

		array2 := array[:len(array)-32]
		bytes := utils.Decrypt(keyHash, array2)

		return string(bytes)
	}

	if strings.HasPrefix(ver, "5.1") || strings.HasPrefix(ver, "5.2") {
		array3, _ := base64.StdEncoding.DecodeString(rawPass)
		key2 := sha256.New()
		key2.Write([]byte(sid))
		keyHash2 := key2.Sum(nil)

		array4 := array3[:len(array3)-32]
		bytes2 := utils.Decrypt(keyHash2, array4)

		return string(bytes2)
	}

	if strings.HasPrefix(ver, "5") || strings.HasPrefix(ver, "6") || strings.HasPrefix(ver, "7.0") {
		array5, _ := base64.StdEncoding.DecodeString(rawPass)
		key3 := sha256.New()
		key3.Write([]byte(username + sid))
		keyHash3 := key3.Sum(nil)

		array6 := array5[:len(array5)-32]
		bytes3 := utils.Decrypt(keyHash3, array6)

		return string(bytes3)
	}

	if strings.HasPrefix(ver, "7") {
		// Reverse the username and concatenate with sid
		reversedUsername := reverseString(username)
		s := reverseString(reversedUsername + sid)

		array7, _ := base64.StdEncoding.DecodeString(rawPass)
		key4 := sha256.New()
		key4.Write([]byte(s))
		keyHash4 := key4.Sum(nil)

		array8 := array7[:len(array7)-32]
		bytes4 := utils.Decrypt(keyHash4, array8)

		return string(bytes4)
	}

	return ""
}

// Helper function to reverse a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func ReadConfigFile(path string) []string {
	var list []string
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return list
	}
	defer file.Close()

	// 创建UTF-16解码器
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	reader := transform.NewReader(file, decoder)

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	var inputStr strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		inputStr.WriteString(line + "\n")
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return list
	}

	input := inputStr.String()

	var username, host, passwd, version string

	re := regexp.MustCompile(`(?m)Version=(.*)`)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		version = match[1]
	}

	re = regexp.MustCompile(`(?m)Host=(.*)`)
	match = re.FindStringSubmatch(input)
	if len(match) > 1 {
		host = match[1]
	}

	re = regexp.MustCompile(`(?m)UserName=(.*)`)
	match = re.FindStringSubmatch(input)
	if len(match) > 1 {
		username = match[1]
	}

	re = regexp.MustCompile(`(?m)Password=(.*)`)
	match = re.FindStringSubmatch(input)
	if len(match) > 1 {
		passwd = match[1]
	}

	list = append(list, version, host, username)
	if len(passwd) > 3 {
		list = append(list, passwd)
	}

	return list
}
func XmanagerSave(path string) {

	rootPath, _ := utils.GetPersonalFolderPath()
	GetAllAccessibleFiles(rootPath + "\\Documents\\NetSarang Computer")
	if len(sessionFiles) != 0 {
		targetPath := filepath.Join(path, XmanagerName)
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}

		sessionstext := DecryptSessions()
		if len(sessionstext) > 1 {
			err := utils.WriteDataToCSV(targetPath+"\\sessions.csv", sessionstext)
			if err != nil {
				return
			}
		}

	}
	fmt.Println("Xmanager 取证结束")

}
