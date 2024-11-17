package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
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

func CopyFile(src, dst string) error {
	// 使用 os.OpenFile 以只读模式打开源文件
	srcFile, err := os.OpenFile(src, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %v", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("复制文件失败: %v", err)
	}

	return nil
}

// 辅助函数：复制文件
func CopyFile_old(src, dst string) error {
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
func CopyLockDirectory(src, dst string) error {
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
			// 尝试以只读模式打开文件
			srcFile, err := os.OpenFile(srcPath, os.O_RDONLY, 0)
			if err != nil {
				fmt.Printf("failed to open file %s: %v, skipping...\n", srcPath, err)
				continue
			}
			defer srcFile.Close()

			// 读取文件内容
			input, err := ioutil.ReadAll(srcFile)
			if err != nil {
				fmt.Printf("failed to read file %s: %v, skipping...\n", srcPath, err)
				continue
			}

			// 写入目标文件
			err = ioutil.WriteFile(dstPath, input, file.Mode())
			if err != nil {
				fmt.Printf("failed to write file %s: %v, skipping...\n", dstPath, err)
				continue
			}
		}
	}

	return nil
}

// 获取Local本地应用数据文件夹的路径
func getLocalAppData() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir + "\\AppData\\Local", nil
}

// 获取Roaming应用数据文件夹的路径
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

// 将GBK编码的字节切片转换为UTF-8编码的字符串
func ConvertGBKToUTF8(gbkBytes []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(gbkBytes), simplifiedchinese.GBK.NewDecoder())
	utf8Bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(utf8Bytes), nil
}

type Navicat11Cipher struct {
	blowfishCipher *blowfish.Cipher
}

// StringToByteArray converts a hex string to a byte array
func StringToByteArray(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

// XorBytes performs XOR between two byte slices
func XorBytes(a, b []byte, len int) {
	for i := 0; i < len; i++ {
		a[i] ^= b[i]
	}
}

// initializes the Navicat11Cipher with a default key
func NewNavicat11Cipher() (*Navicat11Cipher, error) {
	bytes := []byte("3DC5CA39")
	hash := sha1.New()
	hash.Write(bytes)
	key := hash.Sum(nil)

	blowfishCipher, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &Navicat11Cipher{
		blowfishCipher: blowfishCipher,
	}, nil
}

// NewNavicat11CipherWithCustomKey initializes the Navicat11Cipher with a custom user key
func NewNavicat11CipherWithCustomKey(customUserKey string) (*Navicat11Cipher, error) {
	bytes := []byte(customUserKey)
	hash := sha1.New()
	hash.Write(bytes)
	key := hash.Sum(nil)[:8] // Use the first 8 bytes as the key

	blowfishCipher, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &Navicat11Cipher{
		blowfishCipher: blowfishCipher,
	}, nil
}

// DecryptString decrypts the provided ciphertext using the blowfish cipher
func (n *Navicat11Cipher) DecryptString(ciphertext string) (string, error) {
	num := 8
	array, err := StringToByteArray(ciphertext)
	if err != nil {
		return "", err
	}

	array2 := make([]byte, num)
	for i := 0; i < num; i++ {
		array2[i] = byte(0xFF)
	}

	// Initial block encryption
	n.blowfishCipher.Encrypt(array2, array2)

	var result []byte
	num2 := len(array) / num
	num3 := len(array) % num
	var array4, array5 []byte

	for i := 0; i < num2; i++ {
		// Process each 8-byte block
		array4 = array[i*num : (i+1)*num]
		array5 = make([]byte, num)
		copy(array5, array4)

		n.blowfishCipher.Decrypt(array4, array4)
		XorBytes(array4, array2, num)
		result = append(result, array4...)

		XorBytes(array2, array5, num)
	}

	if num3 != 0 {
		// Handle remainder block if any
		array4 = make([]byte, num)
		copy(array4, array[num2*num:])
		n.blowfishCipher.Encrypt(array2, array2)
		XorBytes(array4, array2, num)
		result = append(result, array4[:num3]...)
	}

	return string(result), nil
}

// GetPersonalFolderPath 获取用户的个人文件夹路径
func GetPersonalFolderPath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.HomeDir, nil
}

// GetCurrentUserSID 获取当前用户的SID
func GetCurrentUserSID() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("whoami", "/all")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// 正则表达式匹配SID
	re := regexp.MustCompile(`S-1-\d+-\d+-\d+-\d+-\d+-\d+`)
	matches := re.FindStringSubmatch(out.String())
	if len(matches) < 1 {
		return "", fmt.Errorf("SID not found")
	}

	return matches[0], nil
}

// 检查程序是否以管理员身份运行
func IsAdmin() bool {
	_, err := exec.Command("net", "session").Output()
	if err != nil {
		return false
	}
	return true
}
func PrintBanner() {
	fmt.Println(`         _____                        _      ____            
  /\    |  ___|__  _ __ ___ _ __  ___(_) ___|  _ \ _ __ ___  
 /  \   | |_ / _ \| '__/ _ \ '_ \/ __| |/ __| |_) | '__/ _ \ 
/ /\ \  |  _| (_) | | |  __/ | | \__ \ | (__|  __/| | | (_) |
\/  \/  |_|  \___/|_|  \___|_| |_|___/_|\___|_|   |_|  \___/ `)
	fmt.Println("欢迎使用 ^ ForensicPro V1.0 一款Windows自动取证工具 by:Agony")
	fmt.Println("===================================================")
	fmt.Println("下面开始进行Windows取证")
}
