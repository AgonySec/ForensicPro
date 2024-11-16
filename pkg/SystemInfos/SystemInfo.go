package SystemInfos

import (
	"ForensicPro/utils"
	"bytes"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"golang.design/x/clipboard"
	"golang.org/x/sys/windows/registry"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var SystemInfoName = "Systeminfo"

// 系统信息
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

// USB历史记录,以及MountedDevices挂载过的设备
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

// 自定义注册表项
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

// 安装程序
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

// NetworkList无线信息
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

// RecentDocs最近打开文件
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

// 用户接口的 IP 地址
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

// 系统启动项
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

// 最近一次正常关机时间
func GetShutdownTime(path string) {
	//HKEY_LOCAL_MACHINE\SYSTEM\ControlSet001\Control\Windows
	var result strings.Builder

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Windows`, registry.READ)
	if err != nil {
		return
	}
	defer key.Close()

	ShutdownTime, _, err := key.GetBinaryValue("ShutdownTime")
	if err != nil {
		fmt.Println("读取 ShutdownTime 值失败:", err)
		return
	}
	// 将二进制数据转换为十六进制字符串
	hexString := fmt.Sprintf("%x", ShutdownTime)
	result.WriteString(hexString + "\n")

	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(result.String(), targetPath+"\\"+"ShutdownTime.txt")
	fmt.Println("ShutdownTime信息取证结束")

}

// Prefetch预读取文件
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

// Windows 资源管理器地址栏历史记录
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

// Recent最近打开的文件
func GetRecent(path string) {
	var result strings.Builder

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
		// 输出文件或文件夹名称

		result2 := utils.GetShortcutTargetPath(file)
		if result2 != "" {
			// 获取目标路径的文件名
			result.WriteString(result2)
			result.WriteString("\n")

		}
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(result.String(), targetPath+"\\"+"Recent.txt")
	fmt.Println("最近打开文件信息取证结束")

}

// StartUp启动项
func GetStartUp(path string) {
	startUpPaths := []string{
		utils.GetOperaPath("Microsoft\\Windows\\Start Menu\\Programs\\Startup"),
		"C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs\\Startup",
	}
	var result strings.Builder
	for _, startUpPath := range startUpPaths {
		files, err := filepath.Glob(filepath.Join(startUpPath, "*.lnk"))
		if err != nil {
			fmt.Errorf("获取文件列表失败: %w", err)
			return
		}
		for _, file := range files {
			result2 := utils.GetShortcutTargetPath(file)
			if result2 != "" {
				// 获取目标路径的文件名
				result.WriteString(result2)
				result.WriteString("\n")
			}
		}
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(result.String(), targetPath+"\\"+"StartUp.txt")
	fmt.Println("StartUp启动项信息取证结束")
}

// 获取回收站中的所有文件名
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

// 获取计划任务信息
func GetScheduledJobs(path string) {
	// 执行命令获取计划任务信息
	cmd := exec.Command("schtasks", "/query", "/fo", "LIST", "/v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("获取计划任务信息失败: %v", err)
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(string(output), targetPath+"\\"+"ScheduledJobs.txt")
	fmt.Println("计划任务信息取证结束")

}

// 获取剪切板信息
func GetClipboard(path string) {
	result := string(clipboard.Read(clipboard.FmtText))
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(result, targetPath+"\\"+"Clipboard.txt")
	fmt.Println("剪切板信息取证结束")

}
func GetProcessesOpenedFiles(path string) {
	var builder strings.Builder
	// 获取当前系统中运行的所有进程
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 遍历所有进程
	for _, p := range processes {
		// 获取每个进程打开的文件列表
		openFiles, err := p.OpenFiles()
		if err != nil {
			//fmt.Printf("Error getting open files for process %d: %v\n", p.Pid, err)
			continue
		}
		// 打印每个进程的PID和打开的文件信息
		fmt.Printf("Process ID: %d\n", p.Pid)
		for _, file := range openFiles {
			builder.WriteString(fmt.Sprintf("  File: %s, FD: %d\n", file.Path, file.Fd))
		}
		fmt.Println() // 打印空行以分隔不同进程的信息
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(builder.String(), targetPath+"\\"+"Processes_Opened_Files.txt")
	fmt.Println("Processes_Opened_Files信息取证结束")
}

func SystemInfoSave(path string) {

	//GetSystemInfo(path)
	//GetUSBHistory(path)
	//Get_custom_registry_keys(path)
	//GetRecentDocs(path)

}
