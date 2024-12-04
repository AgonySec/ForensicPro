package SystemInfos

import (
	"ForensicPro/utils"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/net"
	"golang.design/x/clipboard"
	"golang.org/x/sys/windows/registry"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var SystemInfoName = "Systeminfo"

// GetSystemInfo 系统信息systeminfo
func GetSystemInfo(path string) {
	// 执行systeminfo命令
	cmd := exec.Command("systeminfo")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}
	// 写入文件
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	// 将GBK编码的输出转换为UTF-8编码
	outputStr, err := utils.ConvertGBKToUTF8(output)
	utils.WriteToFile(outputStr, targetPath+"\\"+SystemInfoName+".txt")
	fmt.Println("systeminfo信息取证结束")

}

// GetUSBHistory USB历史记录,以及MountedDevices挂载过的设备，读注册表
func GetUSBHistory(path string) {
	// 定义需要读取的注册表路径
	registryPaths := []string{
		"HKLM\\System\\currentcontrolset\\enum\\usbstor",
		"HKLM\\System\\currentcontrolset\\enum\\usb",
		"HKLM\\System\\MountedDevices",
	}

	var outputBuffer bytes.Buffer

	// 循环读取每个注册表路径的内容
	for _, registryPath := range registryPaths {
		cmd := exec.Command("reg", "query", registryPath, "/s")
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error running command for %s: %v\n", registryPath, err)
			continue
		}

		// 将GBK编码的输出转换为UTF-8编码
		outputStr, err := utils.ConvertGBKToUTF8(output)
		if err != nil {
			fmt.Printf("Error converting encoding for %s: %v\n", registryPath, err)
			continue
		}

		// 将内容追加到缓冲区
		outputBuffer.WriteString(fmt.Sprintf("Registry Path: %s\n", registryPath))
		outputBuffer.WriteString(outputStr)
		outputBuffer.WriteString("\n\n")
	}

	// 写入文件
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}

	filePath := filepath.Join(targetPath, "USBHistory.txt")
	if err := utils.WriteToFile(outputBuffer.String(), filePath); err != nil {
		log.Fatalf("写入文件失败: %v", err)
	}

	fmt.Println("USBHistory信息取证结束")

}

// GetCustomRegistryKeys 自定义注册表项内容
func GetCustomRegistryKeys(path string) {
	//cmd := exec.Command("reg query HKLM\\System\\currentcontrolset\\enum\\usbstor /s")
	cmd := exec.Command("reg", "query", "HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run", "/s")

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}
	// 写入文件
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	// 将GBK编码的输出转换为UTF-8编码
	outputStr, err := utils.ConvertGBKToUTF8(output)
	utils.WriteToFile(outputStr, targetPath+"\\"+"custom_registry_keys.txt")
	fmt.Println("custom_registry_keys信息取证结束")
}

// GetInstalledPrograms 安装程序 注册表 LOCAL_MACHINE\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall
func GetInstalledPrograms(path string) {
	var result strings.Builder

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, registry.READ)
	if err != nil {
		return
	}
	defer key.Close()

	subKeyNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return
	}

	for _, name := range subKeyNames {
		subKey, err := registry.OpenKey(key, name, registry.READ)
		if err != nil {
			continue
		}
		defer subKey.Close()

		displayName, _, err := subKey.GetStringValue("DisplayName")
		if err == nil && displayName != "" {
			result.WriteString(displayName + "\n")
		}
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}

	utils.WriteToFile(result.String(), targetPath+"\\"+"installed_programs.txt")
	fmt.Println("installed_programs信息取证结束")

}

// GetNetworkList NetworkList无线信息 注册表 HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\NetworkList\Profiles
func GetNetworkList(path string) {
	//无线信息
	cmd := exec.Command("reg", "query", "HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\NetworkList\\Profiles", "/s")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}
	// 将GBK编码的输出转换为UTF-8编码
	outputStr, err := utils.ConvertGBKToUTF8(output)
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}

	utils.WriteToFile(outputStr, targetPath+"\\"+"NetworkList.txt")
	fmt.Println("NetworkList无线信息取证结束")
}

// GetRecentDocs RecentDocs最近打开文件 注册表 HKCU\Software\Microsoft\Windows\CurrentVersion\Explorer\RecentDocs
func GetRecentDocs(path string) {
	cmd := exec.Command("reg", "query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\RecentDocs", "/s")

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}
	// 写入文件
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	// 将GBK编码的输出转换为UTF-8编码
	outputStr, err := utils.ConvertGBKToUTF8(output)
	utils.WriteToFile(outputStr, targetPath+"\\"+"RecentDocs.txt")
	fmt.Println("RecentDocs信息取证结束")

}

// GetInterfaces 用户接口的 IP 地址
func GetInterfaces(path string) {
	cmd := exec.Command("reg", "query", "HKLM\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters\\Interfaces", "/s")

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}
	// 写入文件
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	// 将GBK编码的输出转换为UTF-8编码
	outputStr, err := utils.ConvertGBKToUTF8(output)
	utils.WriteToFile(outputStr, targetPath+"\\"+"Interfaces.txt")
	fmt.Println("Interfaces信息取证结束")
}

// GetSystemStartup 系统启动项
func GetSystemStartup(path string) {
	/**
	系统启动项:HKEY_LOCAL_MACHINE\Software\Microsoft\Windows\CurrentVersion\Run
	启动时运行一次:HKEY_LOCAL_MACHINE\Software\Microsoft\Windows\CurrentVersion\RunOnce
	自启动服务:HKEY_LOCAL_MACHINE\System\CurrentControlSet\Services
	特定用户登录时启动:HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run
	HKEY_LOCAL_MACHINE\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Run
	HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon
	HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Terminal Server\Wds\rdpwd
	HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\App Paths
	*/
	// 定义需要读取的注册表路径
	registryPaths := []string{
		"HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run",
		"HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\RunOnce",
		"HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Winlogon",
		"HKLM\\SOFTWARE\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Run",
		"HKLM\\System\\CurrentControlSet\\Services",
		"HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run",
		"HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\App Paths",
		"HKLM\\SYSTEM\\CurrentControlSet\\Control\\Terminal Server\\Wds",
	}

	var outputBuffer bytes.Buffer

	// 循环读取每个注册表路径的内容
	for _, registryPath := range registryPaths {
		cmd := exec.Command("reg", "query", registryPath, "/s")
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error running command for %s: %v\n", registryPath, err)
			continue
		}

		// 将GBK编码的输出转换为UTF-8编码
		outputStr, err := utils.ConvertGBKToUTF8(output)
		if err != nil {
			fmt.Printf("Error converting encoding for %s: %v\n", registryPath, err)
			continue
		}

		// 将内容追加到缓冲区
		outputBuffer.WriteString(fmt.Sprintf("Registry Path: %s\n", registryPath))
		outputBuffer.WriteString(outputStr)
		outputBuffer.WriteString("\n\n")

	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(outputBuffer.String(), targetPath+"\\"+"SystemStartup.txt")
	fmt.Println("SystemStartup信息取证结束")

}

// GetShutdownTime 最近一次正常关机时间
func GetShutdownTime(path string) {
	// 使用 PowerShell 获取事件日志中的关机时间（ID 1074 表示正常关机）
	cmd := exec.Command("powershell", "Get-WinEvent -LogName System | Where-Object {$_.Id -eq 1074} | Select-Object -First 1 TimeCreated")
	result, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(string(result), targetPath+"\\"+"ShutdownTime.txt")
	fmt.Println("ShutdownTime信息取证结束")

}

// GetPrefetch 预读取文件
func GetPrefetch(path string) {
	// 指定文件夹路径
	prefetchFolderPath := "C:\\Windows\\Prefetch"
	// 写入文件
	targetPath := filepath.Join(path, SystemInfoName, "prefetch")
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}

	// 遍历文件夹中的所有文件
	err := filepath.Walk(prefetchFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 如果是文件，将文件内容复制到目标路径
		if !info.IsDir() {
			relativePath, err := filepath.Rel(prefetchFolderPath, path)
			if err != nil {
				return err
			}
			targetFilePath := filepath.Join(targetPath, relativePath)

			// 确保目标文件夹存在
			if err := os.MkdirAll(filepath.Dir(targetFilePath), os.ModePerm); err != nil {
				return err
			}

			// 打开源文件
			sourceFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer sourceFile.Close()

			// 创建目标文件
			targetFile, err := os.Create(targetFilePath)
			if err != nil {
				return err
			}
			defer targetFile.Close()

			// 将文件内容复制到目标文件
			_, err = io.Copy(targetFile, sourceFile)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("遍历预读取文件文件夹时出错:", err)
		return
	}
	fmt.Println("预读取文件已成功全部导出")

}

// GetExplorerTypedPaths 资源管理器地址栏历史记录
func GetExplorerTypedPaths(path string) {
	// Windows 资源管理器地址栏历史记录 HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Explorer\TypedPaths
	cmd := exec.Command("reg", "query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\TypedPaths", "/s")

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}
	// 写入文件
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	// 将GBK编码的输出转换为UTF-8编码
	outputStr, err := utils.ConvertGBKToUTF8(output)
	utils.WriteToFile(outputStr, targetPath+"\\"+"ExplorerTypedPaths.txt")
	fmt.Println("Windows 资源管理器地址栏历史记录取证结束")
}

// GetRecent 获取最近打开的文件
func GetRecent(path string) {
	var CSVData [][]string
	CSVData = append(CSVData, []string{"文件名", "文件路径"})
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`, registry.QUERY_VALUE)
	if err != nil {
		return
	}
	defer key.Close()

	// 读取 Recent 文件夹路径
	recentPath, _, err := key.GetStringValue("Recent")
	if err != nil {
		return
	}
	// 读取 recentPath 目录下的所有文件和文件夹
	// 获取当前目录下的所有 .lnk 文件 ，不包括子目录
	files, err := filepath.Glob(filepath.Join(recentPath, "*.lnk"))
	if err != nil {
		fmt.Errorf("获取文件列表失败: %w", err)
		return
	}

	for _, file := range files {
		fileName, targetPath := utils.GetShortcutTargetPath(file)
		if fileName != "" || targetPath != "" {
			CSVData = append(CSVData, []string{fileName, targetPath})
		}
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	err = utils.WriteDataToCSV(targetPath+"\\"+"Recent.csv", CSVData)
	if err != nil {
		return
	}
	fmt.Println("最近打开文件信息取证结束")

}

// GetStartUp 获取系统启动项和软件启动项
func GetStartUp(path string) {
	startUpPaths := []string{
		utils.GetOperaPath("Microsoft\\Windows\\Start Menu\\Programs\\Startup"),
		"C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs\\Startup",
	}
	var CSVData [][]string
	CSVData = append(CSVData, []string{"文件名", "文件路径"})

	for _, startUpPath := range startUpPaths {
		files, err := filepath.Glob(filepath.Join(startUpPath, "*.lnk"))
		if err != nil {
			fmt.Errorf("获取文件列表失败: %w", err)
			return
		}
		for _, file := range files {
			fileName, targetPath := utils.GetShortcutTargetPath(file)
			if fileName != "" || targetPath != "" {
				CSVData = append(CSVData, []string{fileName, targetPath})
			}
		}
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	err := utils.WriteDataToCSV(targetPath+"\\"+"StartUp.csv", CSVData)
	if err != nil {
		return
	}
	fmt.Println("StartUp启动项信息取证结束")
}

// GetRecycleBin 获取回收站中的所有文件名
func GetRecycleBin(path string) {
	var result strings.Builder

	recycleBinPath := "C:\\$Recycle.Bin"

	// 获取回收站文件夹中的所有目录
	dirEntries, err := ioutil.ReadDir(recycleBinPath)
	if err != nil {
		fmt.Errorf("无法读取回收站目录: %w", err)
		return
	}

	// 遍历所有子目录，查找回收的文件
	for _, entry := range dirEntries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), "1001") {
			userBinPath := filepath.Join(recycleBinPath, entry.Name())
			fileEntries, err := ioutil.ReadDir(userBinPath)
			if err != nil {
				continue
			}
			// 遍历该目录中的文件
			for _, fileEntry := range fileEntries {
				result.WriteString(fileEntry.Name())
				result.WriteString("\n")
			}
		}
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(result.String(), targetPath+"\\"+"RecycleBin.txt")
	fmt.Println("RecycleBin信息取证结束")

}

// GetScheduledJobs 获取计划任务信息
func GetScheduledJobs(path string) {
	// 执行命令获取计划任务信息
	cmd := exec.Command("schtasks", "/query", "/fo", "LIST")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("获取计划任务信息失败: %v", err)
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}

	csvFilePath := filepath.Join(targetPath, "ScheduledJobs.csv")
	err = writeCSV(csvFilePath, string(output))
	if err != nil {
		log.Fatalf("写入CSV文件失败: %v", err)
	}
	fmt.Println("计划任务信息取证结束")
}
func writeCSV(filePath string, data string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var headers []string
	var records [][]string

	lines := strings.Split(data, "\n")
	currentRecord := make(map[string]string)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if len(currentRecord) > 0 {
				if len(headers) == 0 {
					headers = getKeysSorted(currentRecord)
					writer.Write(headers)
				}
				record := make([]string, len(headers))
				for i, header := range headers {
					record[i] = currentRecord[header]
				}
				records = append(records, record)
				currentRecord = make(map[string]string)
			}
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			currentRecord[key] = value
		}
	}

	// Handle the last record if there is no trailing empty line
	if len(currentRecord) > 0 {
		if len(headers) == 0 {
			headers = getKeysSorted(currentRecord)
			writer.Write(headers)
		}
		record := make([]string, len(headers))
		for i, header := range headers {
			record[i] = currentRecord[header]
		}
		records = append(records, record)
	}

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func getKeysSorted(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// GetClipboard 获取剪切板信息
func GetClipboard(path string) {
	result := string(clipboard.Read(clipboard.FmtText))
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(result, targetPath+"\\"+"Clipboard.txt")
	fmt.Println("剪切板信息取证结束")

}

// 根据协议类型返回协议名称
func protocolName(protocolType uint32) string {
	switch protocolType {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	default:
		return fmt.Sprintf("未知协议 (%d)", protocolType)
	}
}

func getConnections() ([]net.ConnectionStat, error) {
	conns, err := net.Connections("all") // 获取所有类型的连接（TCP/UDP）
	if err != nil {
		return nil, err
	}
	return conns, nil
}

// GetSockets 获取套接字信息
func GetSockets(path string) {
	var CSVData [][]string

	conns, err := getConnections()
	if err != nil {
		fmt.Println("获取网络连接失败:", err)
		return
	}
	CSVData = append(CSVData, []string{"Local Address", "Remote Address", "Status", "Protocol"})

	for _, conn := range conns {
		protocol := protocolName(conn.Type) // 获取协议名称
		CSVData = append(CSVData, []string{
			conn.Laddr.String(),
			conn.Raddr.String(),
			conn.Status,
			protocol})
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	err = utils.WriteDataToCSV(targetPath+"\\"+"Sockets.csv", CSVData)
	if err != nil {
		return
	}
	fmt.Println("Sockets信息取证结束")

}

// GetSessions 获取会话信息
func GetSessions(path string) {
	cmd := exec.Command("quser")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("执行 quser 命令失败: %v", err)
	}

	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	err = utils.WriteToFile(string(output), targetPath+"\\"+"Sessions.txt")
	if err != nil {
		return
	}
	fmt.Println("Sessions信息取证结束")
}

// GetProcesses 获取进程信息
func GetProcesses(path string) {
	var processes []utils.Win32_Process
	var CSVData [][]string
	query := "SELECT ProcessID, Name, CommandLine, CreationDate, ExecutablePath, ParentProcessId, Status, ThreadCount FROM Win32_Process"

	// 执行WMI查询
	err := wmi.Query(query, &processes)
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	if len(processes) == 0 {
		fmt.Println("No processes found.")
		return
	}

	// 创建新的Excel文件
	//f := excelize.NewFile()
	// 写入表头
	headers := []string{"ProcessID", "Name", "CommandLine", "ParentProcessId", "CreationDate", "ExecutablePath", "Status", "ThreadCount"}
	CSVData = append(CSVData, headers)

	// 将用户信息逐行写入Excel
	for _, process := range processes {
		CSVData = append(CSVData, []string{
			fmt.Sprintf("%d", process.ProcessID),
			process.Name,
			process.CommandLine,
			fmt.Sprintf("%d", process.ParentProcessId),
			process.CreationDate,
			process.ExecutablePath,
			process.Status,
			fmt.Sprintf("%d", process.ThreadCount),
		})

	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	// 保存生成的Excel文件
	utils.WriteDataToCSV(targetPath+"\\"+"Processes.csv", CSVData)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Processes信息取证结束")
}

// GetNetworksCards 获取网卡信息
func GetNetworksCards(path string) {
	var CSVData [][]string
	var adapters []utils.Win32_NetworkAdapter
	query := "SELECT Name, AdapterType, DeviceID, Manufacturer, MACAddress, Speed, NetEnabled FROM Win32_NetworkAdapter"

	// 执行 WMI 查询
	err := wmi.Query(query, &adapters)
	if err != nil {
		return
	}

	if len(adapters) == 0 {
		fmt.Println("No network adapters found.")
	} else {
		headers := []string{"Name", "AdapterType", "DeviceID", "Manufacturer", "MACAddress", "Speed", "NetEnabled"}
		CSVData = append(CSVData, headers)
		for _, adapter := range adapters {
			CSVData = append(CSVData, []string{
				adapter.Name,
				adapter.AdapterType,
				adapter.DeviceID,
				adapter.Manufacturer,
				adapter.MacAddress,
				fmt.Sprintf("%d", adapter.Speed),
				fmt.Sprintf("%v", adapter.NetEnabled)})
		}
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteDataToCSV(targetPath+"\\"+"NetworksCards.csv", CSVData)
	fmt.Println("NetworksCards信息取证结束")
}

// GetRoutesTables 获取路由表信息
func GetRoutesTables(path string) {
	cmd := exec.Command("netstat", "-rn")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	err = utils.WriteToFile(string(output), targetPath+"\\"+"RoutesTables.txt")
	if err != nil {
		return
	}
	fmt.Println("RoutesTables信息取证结束")

}

// GetNamedPipes 获取命名管道信息
func GetNamedPipes(path string) {
	// 执行 PowerShell 脚本查询命名管道
	cmd := exec.Command("powershell", "-Command", "Get-ChildItem -Path '\\\\.\\pipe\\'")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}

	err = utils.WriteToFile(string(output), targetPath+"\\"+"NamedPipes.txt")
	if err != nil {
		return
	}
	fmt.Println("NamedPipes信息取证结束")
}

// GetAccounts 获取账户信息
func GetAccounts(path string) {
	var accounts []utils.Win32_UserAccount
	var CSVData [][]string
	query := "SELECT AccountType, Description, Disabled, Domain, " +
		"FullName, LocalAccount, Lockout, Name, " +
		"PasswordChangeable, PasswordExpires, PasswordRequired, SID, " +
		"SIDType, Status FROM Win32_UserAccount"

	// 执行WMI查询
	err := wmi.Query(query, &accounts)
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	if len(accounts) == 0 {
		fmt.Println("No accounts found.")
		return
	}

	// 设置Excel表头
	headers := []string{
		"AccountType", "Description", "Disabled", "Domain", "FullName",
		"LocalAccount", "Lockout", "Name", "PasswordChangeable", "PasswordExpires",
		"PasswordRequired", "SID", "SIDType", "Status",
	}
	CSVData = append(CSVData, headers)
	// 将用户信息逐行写入Excel
	for _, account := range accounts {
		CSVData = append(CSVData, []string{
			strconv.Itoa(int(account.AccountType)),
			account.Description,
			strconv.FormatBool(account.Disabled),
			account.Domain,
			account.FullName,
			strconv.FormatBool(account.LocalAccount),
			strconv.FormatBool(account.Lockout),
			account.Name,
			strconv.FormatBool(account.PasswordChangeable),
			strconv.FormatBool(account.PasswordExpires),
			strconv.FormatBool(account.PasswordRequired),
			account.SID,
			strconv.Itoa(int(account.SIDType)),
			account.Status,
		})
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	// 保存生成的Excel文件
	err = utils.WriteDataToCSV(targetPath+"\\"+"Accounts.csv", CSVData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Accounts信息取证结束")
}

// GetArpTable 获取arp表信息
func GetArpTable(path string) {
	var builder strings.Builder

	// 执行 arp -a 命令获取 ARP 表信息
	cmd := exec.Command("arp", "-a")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	result, _ := utils.ConvertGBKToUTF8(output)

	// 输出 ARP 表
	builder.WriteString("ARP Table Information:\n")
	builder.WriteString(result + "\n")
	// 解析 ARP 表信息
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if strings.Contains(line, "动态") {
			// 打印出动态映射的 IP 和 MAC 地址
			builder.WriteString(line + "\n")
		}
	}

	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}

	utils.WriteToFile(builder.String(), targetPath+"\\"+"ArpTable.txt")
	fmt.Println("ArpTable信息取证结束")
}

// 获取文件系统快照信息
func GetFsSnapshot(path string) {
	// 执行 vssadmin list shadows 获取 VSS 快照信息
	cmd := exec.Command("vssadmin", "List", "Shadows")
	output, err := cmd.Output()
	if err != nil {
		return
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(string(output), targetPath+"\\"+"FsSnapshot.txt")
	fmt.Println("FsSnapshot信息取证结束")
}

// 获取DNS缓存信息
func GetDnsCaches(path string) {
	// 执行 ipconfig /displaydns 获取 DNS 缓存信息
	cmd := exec.Command("ipconfig", "/displaydns")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	result, _ := utils.ConvertGBKToUTF8(output)
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(result, targetPath+"\\"+"DnsCaches.txt")
	fmt.Println("DnsCaches信息取证结束")
}

// 获取共享信息
func GetShares(path string) {
	// 执行 net share 获取共享信息
	cmd := exec.Command("net", "share")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	result, _ := utils.ConvertGBKToUTF8(output)

	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(result, targetPath+"\\"+"Shares.txt")
	fmt.Println("Shares信息取证结束")
}

// GetKB 获取KB补丁信息
func GetKB(path string) {
	var kb []utils.KB
	var Data [][]string
	Data = append(Data, []string{"Caption", "CSName", "HotFixID", "Description", "InstalledBy", "InstalledOn"})

	// 构建WMI查询 wmic qfe
	query := `SELECT Caption,CSName, HotFixID, Description,InstalledBy, InstalledOn FROM Win32_QuickFixEngineering`

	// 执行查询
	err := wmi.Query(query, &kb)
	if err != nil {
		log.Fatal(err)
	}

	// 打印查询结果
	for _, v := range kb {

		Data = append(Data, []string{v.Caption, v.CSName, v.HotFixID, v.Description, v.InstalledBy, v.InstalledOn})
	}

	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	err = utils.WriteDataToCSV(targetPath+"\\"+"KB.csv", Data)
	if err != nil {
		return
	}
	fmt.Println("KB信息取证结束")
}

// GetVolume 获取卷信息
func GetVolume(path string) {
	var volumes []utils.Win32_Volume
	var Data [][]string
	// 获取卷信息
	err := wmi.Query("SELECT * FROM Win32_Volume", &volumes)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(volumes) == 0 {
		fmt.Println("No volumes found.")
		return
	}
	Data = append(Data, []string{"Capacity", "DriveType", "DriveLetter", "FileSystem", "FreeSpace", "Label", "Name"})
	for _, v := range volumes {
		Data = append(Data, []string{
			fmt.Sprintf("%d", v.Capacity),
			fmt.Sprintf("%d", v.DriveType),
			v.DriveLetter,
			v.FileSystem,
			fmt.Sprintf("%d", v.FreeSpace),
			v.Label,
			v.Name})
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	err = utils.WriteDataToCSV(targetPath+"\\"+"Volume.csv", Data)
	if err != nil {
		return
	}
	fmt.Println("Volume信息取证结束")
}

// GetBios 获取BIOS信息
func GetBios(path string) {
	var bios []utils.Win32_BIOS
	var Data [][]string
	err := wmi.Query("SELECT * FROM Win32_BIOS", &bios)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(bios) == 0 {
		fmt.Println("No volumes found.")
		return
	}
	Data = append(Data, []string{"Manufacturer", "Name", "SerialNumber", "SMBIOSBIOSVersion", "Version"})
	for _, v := range bios {
		Data = append(Data, []string{
			v.Manufacturer,
			v.Name,
			v.SerialNumber,
			v.SMBIOSBIOSVersion,
			v.Version,
		})
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	err = utils.WriteDataToCSV(targetPath+"\\"+"Bios.csv", Data)
	if err != nil {
		return
	}
	fmt.Println("Bios信息取证结束")
}

// GetLogicalDisk 获取磁盘信息
func GetLogicalDisk(path string) {
	var logicaldisks []utils.Win32_LogicalDisk
	var Data [][]string
	Data = append(Data, []string{"Caption", "Description", "ProviderName", "SystemName", "DeviceID", "FileSystem", "FreeSpace", "Name", "Size", "VolumeName", "VolumeSerialNumber"})

	// 获取卷信息
	err := wmi.Query("SELECT * FROM Win32_LogicalDisk", &logicaldisks)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(logicaldisks) == 0 {
		return
	}
	for _, v := range logicaldisks {
		Data = append(Data, []string{
			v.Caption,
			v.Description,
			v.ProviderName,
			v.SystemName,
			v.DeviceID,
			v.FileSystem,
			fmt.Sprintf("%d", v.FreeSpace),
			v.Name,
			fmt.Sprintf("%d", v.Size),
			v.VolumeName,
			v.VolumeSerialNumber,
		})
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	err = utils.WriteDataToCSV(targetPath+"\\"+"Logicaldisk.csv", Data)
	if err != nil {
		return
	}
	fmt.Println("Logicaldisk信息取证结束")
}

func GetCPU(path string) {
	var cpu []utils.Win32_CPU
	var Data [][]string
	Data = append(Data, []string{"Caption", "DeviceID", "Manufacturer", "MaxClockSpeed", "SocketDesignation", "Name", "NumberOfCores", "NumberOfLogicalProcessors", "ProcessorId"})
	err := wmi.Query("SELECT * FROM Win32_Processor", &cpu)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(cpu) == 0 {
		return
	}
	for _, v := range cpu {
		Data = append(Data, []string{
			v.Caption,
			v.DeviceID,
			v.Manufacturer,
			fmt.Sprintf("%d", v.MaxClockSpeed),
			v.SocketDesignation,
			v.Name,
			fmt.Sprintf("%d", v.NumberOfCores),
			fmt.Sprintf("%d", v.NumberOfLogicalProcessors),
			v.ProcessorId,
		})
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	err = utils.WriteDataToCSV(targetPath+"\\"+"CPU.csv", Data)
	if err != nil {
		return
	}
	fmt.Println("CPU信息取证结束")
}

// GetPhysicalMemory 获取物理内存信息
func GetPhysicalMemory(path string) {
	var memories []utils.Win32_PhysicalMemory
	var CSVData [][]string
	query := "SELECT * FROM Win32_PhysicalMemory"

	// 执行WMI查询
	err := wmi.Query(query, &memories)
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	if len(memories) == 0 {
		fmt.Println("No physical memory found.")
		return
	}

	// 设置CSV表头
	headers := []string{"Capacity", "DeviceLocator", "Manufacturer", "PartNumber", "SerialNumber", "Speed", "Tag", "MemoryType", "Name", "TotalWidth"}
	CSVData = append(CSVData, headers)

	// 将物理内存信息逐行写入CSV
	for _, memory := range memories {
		CSVData = append(CSVData, []string{
			fmt.Sprintf("%d", memory.Capacity),
			memory.DeviceLocator,
			memory.Manufacturer,
			memory.PartNumber,
			memory.SerialNumber,
			strconv.FormatUint(uint64(memory.Speed), 10),
			memory.Tag,
			fmt.Sprintf("%d", memory.MemoryType),
			memory.Name,
			fmt.Sprintf("%d", memory.TotalWidth),
		})
	}

	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	err = utils.WriteDataToCSV(targetPath+"\\"+"PhysicalMemory.csv", CSVData)
	if err != nil {
		return
	}
	fmt.Println("PhysicalMemory信息取证结束")
}

func SystemInfoSave(path string) {
	//GetProcessesOpenedFiles(path)
	GetPhysicalMemory(path)
	GetSystemInfo(path)
	GetVolume(path)
	GetLogicalDisk(path)
	GetCPU(path)
	GetBios(path)
	GetUSBHistory(path)
	GetRecentDocs(path)
	GetCustomRegistryKeys(path)
	GetInstalledPrograms(path)
	GetProcesses(path)
	GetShares(path)
	GetScheduledJobs(path)
	GetSessions(path)
	GetRecent(path)
	GetClipboard(path)
	GetInterfaces(path)
	GetArpTable(path)
	GetNamedPipes(path)
	GetAccounts(path)
	GetNetworksCards(path)
	GetRoutesTables(path)
	GetFsSnapshot(path)
	GetDnsCaches(path)
	GetStartUp(path)
	GetSystemStartup(path)
	GetExplorerTypedPaths(path)
	GetSockets(path)
	GetNetworkList(path)
	GetPrefetch(path)
	GetShutdownTime(path)
	GetRecycleBin(path)
	GetKB(path)
}
