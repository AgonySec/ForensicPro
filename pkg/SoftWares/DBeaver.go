package SoftWares

import (
	"ForensicPro/utils"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	//_ "encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var DBeaverName = "DBeaver"

type Configuration struct {
	URL string `json:"url"`
}

type Connection struct {
	Configuration Configuration `json:"configuration"`
}

type Data struct {
	Connections map[string]Connection `json:"connections"`
}

func MatchDataSource(filePath, jdbcKey string) string {
	// 读取 JSON 文件
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return ""
	}

	// 解析 JSON 数据
	var data Data
	if err := json.Unmarshal(bytes, &data); err != nil {
		return ""

	}
	// 查找指定键的 URL
	if conn, exists := data.Connections[jdbcKey]; exists {
		return conn.Configuration.URL
	}

	return ""
}

// ConnectionInfo 提取连接信息
func ConnectionInfo(config string, sources string) [][]string {

	pattern := `"([^\"]+)"\s*:\s*\{\s*"#connection"\s*:\s*\{\s*"user"\s*:\s*"([^\"]+)"\s*,\s*"password"\s*:\s*"([^\"]+)"\s*\}\s*\}`

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(config, -1)
	var CSVData [][]string
	CSVData = append(CSVData, []string{"host", "username", "password"})
	for _, match := range matches {
		key := match[1]
		user := match[2]
		password := match[3]

		CSVData = append(CSVData, []string{MatchDataSource(sources, key), user, password})
	}

	return CSVData
}

func DBeaverDecrypt(filePath string, keyHex string, ivHex string) (string, error) {
	// 读取文件内容
	buffer, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// 将十六进制字符串转换为字节数组
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", err
	}
	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return "", err
	}

	// 创建 AES 解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(iv) != block.BlockSize() {
		return "", fmt.Errorf("IV length must equal block size")
	}

	// 创建 CBC 模式的解密器
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密数据
	decrypted := make([]byte, len(buffer))
	mode.CryptBlocks(decrypted, buffer)

	// 去除 PKCS7 填充
	decrypted = PKCS7Unpad(decrypted, block.BlockSize())

	return string(decrypted), nil
}

// PKCS7Unpad 去除 PKCS7 填充
func PKCS7Unpad(data []byte, blockSize int) []byte {
	padding := int(data[len(data)-1])
	if padding > blockSize || padding == 0 {
		return data
	}
	return data[:len(data)-padding]
}
func DBeaverSave(path string) {

	path1 := utils.GetOperaPath("DBeaverData\\workspace6\\General\\.dbeaver\\data-sources.json")
	path2 := utils.GetOperaPath("DBeaverData\\workspace6\\General\\.dbeaver\\credentials-config.json")
	if _, err := os.Stat(path1); err != nil {
		return
	}
	if _, err := os.Stat(path2); err != nil {
		return
	}
	DBevardecrypt, _ := DBeaverDecrypt(path2, "babb4a9f774ab853c96c2d653dfe544a", "00000000000000000000000000000000")

	content := ConnectionInfo(DBevardecrypt, path1)
	if len(content) > 1 {
		targetPath := filepath.Join(path, DBeaverName)
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}
		err := utils.WriteDataToCSV(targetPath+"\\DBeaver.csv", content)
		if err != nil {
			return
		}
	}

	fmt.Println("DBeaver取证结束")

}
