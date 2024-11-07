package SoftWares

import (
	"ForensicPro/utils"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var SQLyogtName = "SQLyog"
var keyArray = []byte{
	41, 35, 190, 132, 225, 108, 214, 174, 82, 144,
	73, 241, 201, 187, 33, 143,
}

var ivArray = []byte{
	179, 166, 219, 60, 135, 12, 62, 153, 36, 94,
	13, 28, 6, 183, 71, 222,
}

type IniLine struct {
	Key   string
	Value string
}

type Pixini struct {
	sectionMap map[string][]IniLine
}

func (p *Pixini) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	p.sectionMap = make(map[string][]IniLine)
	currentSection := ""

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
			continue
		}
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			p.sectionMap[currentSection] = append(p.sectionMap[currentSection], IniLine{Key: key, Value: value})
		}
	}
	return nil
}

func (p *Pixini) Set(key, section, value string) {
	for i, line := range p.sectionMap[section] {
		if line.Key == key {
			p.sectionMap[section][i].Value = value
			return
		}
	}
}

func (p *Pixini) ToString() string {
	var buffer bytes.Buffer
	for section, lines := range p.sectionMap {
		buffer.WriteString(fmt.Sprintf("[%s]\n", section))
		for _, line := range lines {
			buffer.WriteString(fmt.Sprintf("%s=%s\n", line.Key, line.Value))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}
func OldDecrypt(text string) string {
	array, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return ""
	}
	for i := range array {
		array[i] = (array[i] << 1) | (array[i] >> 7)
	}
	return string(array)
}
func NewDecrypt(text string) string {
	array, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return ""
	}
	array2 := make([]byte, 128)
	copy(array2, array)

	block, err := aes.NewCipher(keyArray)
	if err != nil {
		return ""
	}

	mode := cipher.NewCFBDecrypter(block, ivArray)
	mode.XORKeyStream(array2, array2)

	return string(array2[:len(array)])
}

func decryptSqlyog(filePath string) string {
	pixini := &Pixini{}
	err := pixini.Load(filePath)
	if err != nil {
		return ""
	}

	for section, lines := range pixini.sectionMap {
		var text string
		var flag bool
		for _, line := range lines {
			if line.Key == "Password" {
				text = line.Value
			}
			if line.Key == "Isencrypted" {
				flag = line.Value == "1"
			}
		}
		if text != "" {
			val := ""
			if flag {
				val = NewDecrypt(text)
			} else {
				val = OldDecrypt(text)
			}
			pixini.Set("Password", section, val)
		}
	}

	return pixini.ToString()
}

func SQLyogSave(path string) {
	appDataPath := utils.GetOperaPath("SQLyog\\sqlyog.ini")
	if _, err := os.Stat(appDataPath); os.IsNotExist(err) {
		return
	}
	targetPath := filepath.Join(path, SQLyogtName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}

	utils.CopyFile(appDataPath, targetPath+"\\sqlyog.ini")

	info := decryptSqlyog(appDataPath)
	if info == "" {
		return
	}
	utils.WriteToFile(info, targetPath+"\\sqlyog_decrypted.ini")

}
