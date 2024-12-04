package SystemInfos

import (
	"ForensicPro/utils"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 存在bug,若使用管理员方式运行则会报错... 先弃用
func GetProcessesOpenedFiles(path string) {
	var builder strings.Builder
	// 获取当前系统中运行的所有进程
	processes, err := process.Processes()
	if err != nil {
		return
	}
	// 遍历所有进程
	for _, p := range processes {
		if p == nil {
			continue
		}
		openFiles, err := p.OpenFiles()
		if err != nil {
			log.Printf("Error getting open files for process %d: %v\n", p.Pid, err)
			continue
		}
		if openFiles == nil {
			continue
		}
		builder.WriteString(fmt.Sprintf("Process ID: %d\n", p.Pid))
		for _, file := range openFiles {
			builder.WriteString(fmt.Sprintf("  File: %s, FD: %d\n", file.Path, file.Fd))
		}
		builder.WriteString("\n") // 打印空行以分隔不同进程的信息
	}
	targetPath := filepath.Join(path, SystemInfoName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	utils.WriteToFile(builder.String(), targetPath+"\\"+"Processes_Opened_Files.txt")
	fmt.Println("Processes_Opened_Files信息取证结束")
}
