package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var (
	crypt32            = syscall.NewLazyDLL("Crypt32.dll")
	cryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
)

type dataBlob struct {
	cbData uint32
	pbData *byte
}

func DecryptData(encryptedData []byte) ([]byte, error) {
	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("encrypted data is empty")
	}

	var outBlob dataBlob

	// Prepare input data blob
	inBlob := dataBlob{
		cbData: uint32(len(encryptedData)),
		pbData: &encryptedData[0],
	}

	// Call CryptUnprotectData
	ret, _, err := cryptUnprotectData.Call(
		uintptr(unsafe.Pointer(&inBlob)),
		0,
		0,
		0,
		0,
		0,
		uintptr(unsafe.Pointer(&outBlob)),
	)

	if ret == 0 {
		return nil, err
	}

	// Copy decrypted data to a Go slice
	defer func(hmem syscall.Handle) {
		_, err := syscall.LocalFree(hmem)
		if err != nil {

		}
	}(syscall.Handle(unsafe.Pointer(outBlob.pbData)))
	decryptedData := make([]byte, outBlob.cbData)
	copy(decryptedData, (*[1 << 30]byte)(unsafe.Pointer(outBlob.pbData))[:outBlob.cbData:outBlob.cbData])
	// Convert decrypted data to the desired format
	//formattedData2 := formatBytes(decryptedData)

	return decryptedData, nil
}

// 获取解密后的主密钥
func GetMasterKey(BrowserPath string) ([]byte, error) {
	// 	定义Local State 文件路径
	LocalState := filepath.Join(BrowserPath, "Local State")
	// 判断文件是否存在
	if _, err := os.Stat(LocalState); os.IsNotExist(err) {
		return nil, errors.New("Local State 文件不存在")
	}
	data, err := os.ReadFile(LocalState)
	if err != nil {
		return nil, errors.Wrap(err, "读取 Local State 文件失败")
	}

	// 解析 JSON
	var localState map[string]interface{}
	if err := json.Unmarshal(data, &localState); err != nil {
		return nil, errors.Wrap(err, "解析 JSON 失败")
	}

	encryptedKey, ok := localState["os_crypt"].(map[string]interface{})["encrypted_key"].(string)
	if !ok {
		return nil, fmt.Errorf("无法找到加密的主密钥")
	}
	//fmt.Println(encryptedKey)
	decodedKey, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return nil, errors.Wrap(err, "解码 Base64 失败")
	}
	decodedKey = decodedKey[5:]
	decryptedData, err := DecryptData(decodedKey)
	if err != nil {
		fmt.Println("Error decrypting data:", err)
		return nil, nil
	}
	return decryptedData, nil
}

// 解密密码密文
func DecryptAESGCM(encryptedPassword []byte, key []byte) ([]byte, error) {
	if len(encryptedPassword) == 0 {
		return nil, fmt.Errorf("encrypted password is empty")
	}
	str := string(encryptedPassword)
	var plaintext []byte
	// 检查字符串是否以 "v10" 或 "v11" 开头
	if strings.HasPrefix(str, "v10") || strings.HasPrefix(str, "v11") {

		iv := encryptedPassword[3:15]
		payload := encryptedPassword[15:]
		// 创建 AES GCM 解密器
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, fmt.Errorf("failed to create cipher: %v", err)
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, fmt.Errorf("failed to create GCM: %v", err)
		}

		// 解密
		plaintext, err = gcm.Open(nil, iv, payload, nil)
		if err != nil {
			return nil, fmt.Errorf("decryption failed: %v", err)
		}

	} else {
		//fmt.Println("字符串不以 v10 或 v11 开头")
		plaintext, _ = DecryptData(encryptedPassword)
	}

	return plaintext, nil
}

// 辅助函数：复制文件
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// CopyDirectory 递归复制目录
func CopyDirectory(src, dst string) error {
	// 读取源目录的内容
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", src, err)
	}

	// 创建目标目录
	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dst, err)
	}

	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		dstPath := filepath.Join(dst, file.Name())

		if file.IsDir() {
			// 递归复制子目录
			err = CopyDirectory(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// 复制文件
			input, err := ioutil.ReadFile(srcPath)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", srcPath, err)
			}
			err = ioutil.WriteFile(dstPath, input, file.Mode())
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", dstPath, err)
			}
		}
	}

	return nil
}

// 获取本地应用数据文件夹的路径
func getLocalAppData() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir + "\\AppData\\Local", nil
}

// 获取应用数据文件夹的路径
func getAppData() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir + "\\AppData\\Roaming", nil
}

// 获取浏览器路径
func GetLocalAppDataPath(relativePath string) string {
	localAppData, err := getLocalAppData()
	if err != nil {
		panic(err)
	}
	return filepath.Join(localAppData, relativePath)
}

// 获取 Opera 路径
func GetOperaPath(relativePath string) string {
	appData, err := getAppData()
	if err != nil {
		panic(err)
	}
	return filepath.Join(appData, relativePath)
}
