package Browsers

import (
	"ForensicPro/utils"
	"database/sql"
	"fmt"
	_ "github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var profiles []string
var BrowserName string
var BrowserPath string
var MasterKey []byte

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

// 定义浏览器路径字典
var browserOnChromium = map[string]string{
	"Chrome":               getBrowserPath("Google\\Chrome\\User Data"),
	"Chrome Beta":          getBrowserPath("Google\\Chrome Beta\\User Data"),
	"Chromium":             getBrowserPath("Chromium\\User Data"),
	"Chrome SxS":           getBrowserPath("Google\\Chrome SxS\\User Data"),
	"Edge":                 getBrowserPath("Microsoft\\Edge\\User Data"),
	"Brave-Browser":        getBrowserPath("BraveSoftware\\Brave-Browser\\User Data"),
	"QQBrowser":            getBrowserPath("Tencent\\QQBrowser\\User Data"),
	"SogouExplorer":        getBrowserPath("Sogou\\SogouExplorer\\User Data"),
	"360ChromeX":           getBrowserPath("360ChromeX\\Chrome\\User Data"),
	"360Chrome":            getBrowserPath("360Chrome\\Chrome\\User Data"),
	"Vivaldi":              getBrowserPath("Vivaldi\\User Data"),
	"CocCoc":               getBrowserPath("CocCoc\\Browser\\User Data"),
	"Torch":                getBrowserPath("Torch\\User Data"),
	"Kometa":               getBrowserPath("Kometa\\User Data"),
	"Orbitum":              getBrowserPath("Orbitum\\User Data"),
	"CentBrowser":          getBrowserPath("CentBrowser\\User Data"),
	"7Star":                getBrowserPath("7Star\\7Star\\User Data"),
	"Sputnik":              getBrowserPath("Sputnik\\Sputnik\\User Data"),
	"Epic Privacy Browser": getBrowserPath("Epic Privacy Browser\\User Data"),
	"Uran":                 getBrowserPath("uCozMedia\\Uran\\User Data"),
	"Yandex":               getBrowserPath("Yandex\\YandexBrowser\\User Data"),
	"Iridium":              getBrowserPath("Iridium\\User Data"),
	"Opera":                GetOperaPath("Opera Software\\Opera Stable"),
	"Opera GX":             GetOperaPath("Opera Software\\Opera GX Stable"),
}

// 获取浏览器路径
func getBrowserPath(relativePath string) string {
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
func readSQLiteDB(dbPath string, query string, CSVData [][]string) ([][]string, error) {
	// 创建一个 strings.Builder 对象，用于构建最终的字符串结果
	//var builder strings.Builder

	// 创建一个临时文件
	tempFile, err := os.CreateTemp("", "chrome-passwords-*.db")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name()) // 确保临时文件在函数结束时被删除

	// 将数据库文件复制到临时文件
	if err := utils.CopyFile(dbPath, tempFile.Name()); err != nil {
		return nil, err
	}

	// 打开 SQLite 数据库
	db, err := sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 执行 SQL 查询
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 获取查询结果的列数
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 预分配切片以存储每一行的数据
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	//var CSVData [][]string

	// 遍历查询结果
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// 将结果添加到 strings.Builder 中
		var URL, USERNAME, PASSWORD string
		for i, value := range values {
			var fieldValue string
			switch v := value.(type) {
			case []byte:
				fieldValue = string(v)
			case string:
				fieldValue = v
			default:
				fieldValue = fmt.Sprintf("%v", v)
			}
			// 如果当前字段是 password_value，则进行解密
			if columns[i] == "password_value" {
				//masterKey, err := utils.GetMasterKey(BrowserPath)
				if err != nil {
					return nil, fmt.Errorf("failed to get master key: %w", err)
				}
				//decryptedValue, err := decryptPassword([]byte(fieldValue), masterKey)
				decryptedValue, err := utils.DecryptAESGCM([]byte(fieldValue), MasterKey)
				if err != nil {
					fmt.Errorf("failed to decrypt password: %w", err)
				}
				fieldValue = string(decryptedValue)
			}

			// 添加标签和值
			switch columns[i] {
			case "signon_realm":
				URL = fieldValue
				//builder.WriteString("[URL] -> {" + fieldValue + "}\n")
			case "username_value":
				USERNAME = fieldValue
				//builder.WriteString("[USERNAME] -> {" + fieldValue + "}\n")
			case "password_value":
				PASSWORD = fieldValue
				//builder.WriteString("[PASSWORD] -> {" + fieldValue + "}\n")
			}
		}
		CSVData = append(CSVData, []string{URL, USERNAME, PASSWORD})
		//builder.WriteString("\n")
	}

	// 返回构建的字符串结果
	return CSVData, nil
	//return builder.String(), nil
}

func ChromePasswords(targetDir string) ([][]string, error) {
	//var builder strings.Builder
	var CSVData [][]string
	// 获取所有浏览器配置文件的数组
	array := profiles //

	if MasterKey == nil {
		return nil, nil
	}
	CSVData = append(CSVData, []string{"URL", "USERNAME", "PASSWORD"})

	// 遍历每个配置文件
	for _, profile := range array {
		LoginDataPath := filepath.Join(BrowserPath, profile, "Login Data")
		// 判断文件是否存在
		if _, err := os.Stat(LoginDataPath); os.IsNotExist(err) {
			return nil, nil
		}
		utils.CopyFile(LoginDataPath, targetDir+"\\Login Data")
		CSVData, _ = readSQLiteDB(LoginDataPath, "SELECT signon_realm, username_value, password_value FROM logins", CSVData)

	}
	return CSVData, nil
}

func readSQLiteDB2(dbPath string, query string, CSVData [][]string) ([][]string, interface{}) {

	// 创建一个 strings.Builder 对象，用于构建最终的字符串结果
	//var builder strings.Builder

	// 创建一个临时文件
	tempFile, err := os.CreateTemp("", "chrome-cookies-*.db")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name()) // 确保临时文件在函数结束时被删除

	// 将数据库文件复制到临时文件
	if err := utils.CopyFile(dbPath, tempFile.Name()); err != nil {
		return nil, err
	}

	// 打开 SQLite 数据库
	db, err := sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 执行 SQL 查询
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 获取查询结果的列数
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 预分配切片以存储每一行的数据
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	// 遍历查询结果
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		var domain string
		var name string
		var path string
		var num float64
		var encryptedValue string
		// 将结果添加到 strings.Builder 中
		for i, value := range values {
			var fieldValue string
			switch v := value.(type) {
			case []byte:
				fieldValue = string(v)
			case string:
				fieldValue = v
			default:
				fieldValue = fmt.Sprintf("%v", v)
			}
			// 如果当前字段是 password_value，则进行解密
			if columns[i] == "encrypted_value" {
				decryptedValue, err := utils.DecryptAESGCM([]byte(fieldValue), MasterKey)
				if err != nil {
					fmt.Errorf("failed to decrypt cookie: %w", err)
				}
				fieldValue = string(decryptedValue)
			}
			if columns[i] == "expires_utc" {
				result, err := strconv.ParseInt(fieldValue, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to parse expires_utc: %w", err)
				}
				// 计算 num
				if float64(result)/1000000.0-11644473600.0 > 0.0 {
					num = float64(result)/1000000.0 - 11644473600.0
				}
			}

			// 添加标签和值
			switch columns[i] {
			case "host_key":
				domain = fieldValue
			case "name":
				name = fieldValue
			case "encrypted_value":
				encryptedValue = fieldValue
			case "path":
				path = fieldValue
			}

		}
		// 构建 JSON 字符串
		//jsonStr := fmt.Sprintf(`{
		//	"domain": "%s",
		//	"expirationDate": %d,
		//	"hostOnly": false,
		//	"name": "%s",
		//	"path": "%s",
		//	"session": true,
		//	"storeId": null,
		//	"value": "%s"
		//}`, domain, num, name, path, encryptedValue)
		//builder.WriteString(jsonStr)
		CSVData = append(CSVData, []string{domain, fmt.Sprintf("%d", num), "false", name, path, "true", "null", encryptedValue})

		//builder.WriteString("\n")
	}

	// 返回构建的字符串结果
	return CSVData, nil
}

// todo 如果谷歌浏览器运行中，cookie文件将会被使用中，数据库文件结构会改变!，导致无法读取到内容,目前没有解决方案，待后续研究
func ChromeCookies(targetDir string) ([][]string, error) {

	// 获取所有浏览器配置文件的数组
	array := profiles //
	var CSVData [][]string
	CSVData = append(CSVData, []string{"domain", "expirationDate", "hostOnly", "name", "path", "session", "storeId", "value"})

	if MasterKey == nil {
		return nil, nil
	}
	// 遍历每个配置文件
	for _, profile := range array {
		CookiesPath1 := filepath.Join(BrowserPath, profile, "Cookies")
		CookiesPath2 := filepath.Join(BrowserPath, profile, "Network", "Cookies")
		if _, err := os.Stat(CookiesPath1); os.IsNotExist(err) {
			CookiesPath1 = CookiesPath2
		}

		if _, err := os.Stat(CookiesPath1); os.IsNotExist(err) {
			return nil, fmt.Errorf("both paths do not exist: %s and %s", CookiesPath1, CookiesPath2)
		}
		utils.CopyFile(CookiesPath1, targetDir+"\\Cookies")
		CSVData, _ = readSQLiteDB2(CookiesPath1, "SELECT host_key, name, encrypted_value,path,expires_utc FROM cookies", CSVData)
		//builder.WriteString(result)
	}
	return CSVData, nil
}
func ChromeBooks(targetDir string) (string, error) {
	// 实现获取 bookmarks 的逻辑
	var builder strings.Builder
	// 获取所有浏览器配置文件的数组
	array := profiles //
	// 遍历每个配置文件
	for _, profile := range array {
		// 根据浏览器类型确定历史记录文件的路径
		path := ""
		if strings.Contains(BrowserName, "360") {
			path = filepath.Join(BrowserPath, profile, "360Bookmarks")
		} else {
			path = filepath.Join(BrowserPath, profile, "Bookmarks")
		}
		// 打开文件
		result, _ := utils.ReadFileContent(path)
		utils.CopyFile(path, targetDir+"\\Bookmarks")
		builder.WriteString(result)
	}
	return builder.String(), nil
}

func ChromeExtensions() (string, error) {
	// 实现获取 extensions 的逻辑
	var builder strings.Builder
	// 获取所有浏览器配置文件的数组
	array := profiles //
	// 遍历每个配置文件
	// 遍历每个用户配置文件
	for _, profile := range array {
		// 构建扩展目录的路径
		path := filepath.Join(BrowserPath, profile, "Extensions")
		// 检查扩展目录是否存在
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue // 如果不存在，跳过该配置文件
		}

		// 读取扩展目录下的所有子目录
		directories, err := ioutil.ReadDir(path)
		if err != nil {
			continue // 如果读取失败，跳过
		}

		// 遍历每个扩展目录
		for _, dir := range directories {
			if dir.IsDir() { // 确保是目录
				// 检查该目录是否存在
				path2 := filepath.Join(path, dir.Name())
				if _, err := os.Stat(path2); os.IsNotExist(err) {
					continue // 如果不存在，跳过该目录
				}

				// 读取该目录下的所有文件
				dir2, err := ioutil.ReadDir(path2)
				if err != nil {
					continue // 如果读取失败，跳过
				}
				if len(dir2) == 0 {
					continue // 如果目录为空，跳过
				}

				// 构建第一个子目录的路径
				path3 := filepath.Join(path2, dir2[0].Name())
				if _, err := os.Stat(path3); os.IsNotExist(err) {
					continue // 如果不存在，跳过该目录
				}

				// 构建 manifest.json 文件的路径
				manifestPath := filepath.Join(path3, "manifest.json")
				if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
					continue // 如果不存在，跳过该目录
				}

				// 读取 manifest.json 文件内容
				manifestData, err := ioutil.ReadFile(manifestPath)
				if err != nil {
					continue // 如果读取失败，跳过
				}

				// 定义正则表达式，匹配 "name": "扩展名称"
				re := regexp.MustCompile(`"name": "(.*?)"`)
				// 查找所有匹配项
				matches := re.FindAllStringSubmatch(string(manifestData), -1)

				// 遍历匹配结果
				for _, match := range matches {
					if len(match) > 1 { // 确保有匹配的内容
						fileName := dir.Name() // 获取扩展目录名
						value := match[1]      // 获取扩展名称
						// 将扩展目录名和名称添加到结果字符串中
						builder.WriteString(fmt.Sprintf("%s    %s\n", fileName, value))
					}
				}
			}
		}
	}
	return builder.String(), nil
}
func ChromeHistory(targetDir string) string {
	//fmt.Println("开始进行chrome浏览器历史记录取证")
	// 创建一个 strings.Builder 对象，用于构建最终的字符串结果
	var builder strings.Builder
	// 获取所有浏览器配置文件的数组
	array := profiles //

	// 遍历每个配置文件
	for _, profile := range array {
		// 根据浏览器类型确定历史记录文件的路径
		historyPath := ""
		if strings.Contains(BrowserName, "360") {
			historyPath = filepath.Join(BrowserPath, profile, "360History")
		} else {
			historyPath = filepath.Join(BrowserPath, profile, "History")
		}

		// 判断文件是否存在
		if _, err := os.Stat(historyPath); os.IsNotExist(err) {
			return ""
		}
		utils.CopyFile(historyPath, targetDir+"\\History")
		result, err := utils.ReadSQLiteDB_url(historyPath, "SELECT url FROM urls")
		if err != nil {
			fmt.Println("Error reading database:", err)
			continue
		}

		// 将结果添加到 strings.Builder 中
		builder.WriteString(result)
	}

	// 返回构建的字符串结果
	return builder.String()
}
func ChromeSave(path string) {
	// 遍历浏览器列表
	for browserName, browserPath := range browserOnChromium {
		// 初始化默认配置文件列表
		profiles = []string{"Default"}
		BrowserName = browserName
		BrowserPath = browserPath
		// 获取主密钥
		MasterKey, _ = utils.GetMasterKey(BrowserPath)

		// 获取所有配置文件目录
		dirEntries, err := ioutil.ReadDir(BrowserPath)
		if err != nil {
			continue
		}
		for i := 1; i < 100; i++ {
			//filepath.Join(browserPath, fmt.Sprintf("Profile %d", i))
			for _, entry := range dirEntries {
				if entry.Name() == fmt.Sprintf("Profile %d", i) {
					profiles = append(profiles, entry.Name())
					break
				}
			}
		}
		// 创建目标目录
		targetDir := filepath.Join(path, BrowserName)
		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			continue
		}
		utils.CopyFile(BrowserPath+"\\Local State", targetDir+"\\Local State")
		//取证历史记录文件
		history := ChromeHistory(targetDir)
		if history != "" {
			// 将历史记录写入到文件
			outputFile := BrowserName + "_history.txt"
			if err := utils.WriteToFile(history, targetDir+"\\"+outputFile); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
		// 取证书签
		books, err := ChromeBooks(targetDir)
		if books != "" {
			outputFile := BrowserName + "_books.txt"
			if err := utils.WriteToFile(books, targetDir+"\\"+outputFile); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
		// 取证密码
		passwords, err := ChromePasswords(targetDir)
		if len(passwords) > 1 {
			outputFile := BrowserName + "_passwords.csv"
			if err := utils.WriteDataToCSV(targetDir+"\\"+outputFile, passwords); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
		// 取证拓展
		extensions, err := ChromeExtensions()
		if extensions != "" {
			outputFile := BrowserName + "_extensions.txt"
			if err := utils.WriteToFile(extensions, targetDir+"\\"+outputFile); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
		// 取证cookie
		cookies, err := ChromeCookies(targetDir)
		if len(cookies) > 1 {
			outputFile := BrowserName + "_cookies.csv"
			if err := utils.WriteDataToCSV(targetDir+"\\"+outputFile, cookies); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}

		}
		array := profiles

		// 遍历每个配置文件
		for _, profile := range array {
			//configPath := filepath.Join(BrowserPath, profile)
			//err := CreateDirectory(configPath)
			if err != nil {
				fmt.Println(err)
				return
			}
			// 复制 Local Storage 目录
			localStorageSrc := filepath.Join(BrowserPath, profile, "Local Storage")
			localStorageDst := filepath.Join(targetDir, profile, "Local Storage")
			if _, err := os.Stat(localStorageSrc); err == nil {
				err = utils.CopyLockDirectory(localStorageSrc, localStorageDst)
				if err != nil {
					fmt.Println(err)
				}
			}
			// 复制 Local Extension Settings 目录
			localExtSettingsSrc := filepath.Join(BrowserPath, profile, "Local Extension Settings")
			localExtSettingsDst := filepath.Join(targetDir, profile, "Local Extension Settings")
			if _, err := os.Stat(localExtSettingsSrc); err == nil {
				err = utils.CopyLockDirectory(localExtSettingsSrc, localExtSettingsDst)
				if err != nil {
					fmt.Println(err)
				}
			}

			// 复制 Sync Extension Settings 目录
			syncExtSettingsSrc := filepath.Join(BrowserPath, profile, "Sync Extension Settings")
			syncExtSettingsDst := filepath.Join(targetDir, profile, "Sync Extension Settings")
			if _, err := os.Stat(syncExtSettingsSrc); err == nil {
				err = utils.CopyLockDirectory(syncExtSettingsSrc, syncExtSettingsDst)
				if err != nil {
					fmt.Println(err)
				}
			}

		}
	}
	fmt.Println("基于Chrome内核的浏览器取证结束")

}
