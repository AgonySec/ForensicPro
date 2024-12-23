package main

import (
	"ForensicPro/pkg/Browsers"
	"ForensicPro/pkg/FTPS"
	"ForensicPro/pkg/Mails"
	"ForensicPro/pkg/Messengers"
	"ForensicPro/pkg/SoftWares"
	"ForensicPro/pkg/SystemInfos"
	"ForensicPro/utils"
	"fmt"
	"sync"
	"time"
)

func main() {
	startTime := time.Now() // 记录程序开始时间
	utils.PrintBanner()     // 打印Banner
	//判断是否为管理员权限
	if utils.IsAdmin() == false {
		fmt.Println("请使用管理员权限运行本程序")
		return
	}
	path := "result"
	GetAll(path)
	endTime := time.Now()                 // 记录程序结束时间
	elapsedTime := endTime.Sub(startTime) // 计算程序运行时间
	fmt.Println("=================Windows取证结束=================")
	fmt.Printf("程序运行时间: %v\n", elapsedTime)
	// 提示用户按任意键继续
	fmt.Println("按任意键继续...")
	fmt.Scanln()
}
func GetAll(path string) {
	var wg sync.WaitGroup

	// 设置超时时间
	timeout := 2 * time.Minute

	// Browsers
	runWithTimeout(&wg, timeout, func() { Browsers.ChromeSave(path) })
	runWithTimeout(&wg, timeout, func() { Browsers.FirefoxSave(path) })
	runWithTimeout(&wg, timeout, func() { Browsers.IESave(path) })
	runWithTimeout(&wg, timeout, func() { Browsers.SogouSave(path) })

	// Messengers
	runWithTimeout(&wg, timeout, func() { Messengers.LineSave(path) })
	runWithTimeout(&wg, timeout, func() { Messengers.QQSave(path) })
	runWithTimeout(&wg, timeout, func() { Messengers.DiscordSave(path) })
	runWithTimeout(&wg, timeout, func() { Messengers.TelegramSave(path) })
	runWithTimeout(&wg, timeout, func() { Messengers.EnigmaSave(path) })
	runWithTimeout(&wg, timeout, func() { Messengers.DingTalkSave(path) })
	runWithTimeout(&wg, timeout, func() { Messengers.SkypeSave(path) })

	// FTPS
	runWithTimeout(&wg, timeout, func() { FTPS.FileZillaSave(path) })
	runWithTimeout(&wg, timeout, func() { FTPS.SnowflakeSave(path) })
	runWithTimeout(&wg, timeout, func() { FTPS.WinSCPSave(path) })

	// Mails
	runWithTimeout(&wg, timeout, func() { Mails.FoxmailSave(path) })
	runWithTimeout(&wg, timeout, func() { Mails.MailBirdSave(path) })
	runWithTimeout(&wg, timeout, func() { Mails.OutlookSave(path) })
	runWithTimeout(&wg, timeout, func() { Mails.MailMasterSave(path) })

	// SoftWares
	runWithTimeout(&wg, timeout, func() { SoftWares.NeteaseCloudMusicSave(path) })
	runWithTimeout(&wg, timeout, func() { SoftWares.NavicatSave(path) })
	runWithTimeout(&wg, timeout, func() { SoftWares.VSCodeSave(path) })
	runWithTimeout(&wg, timeout, func() { SoftWares.XmanagerSave(path) })
	runWithTimeout(&wg, timeout, func() { SoftWares.FinalShellSave(path) })
	runWithTimeout(&wg, timeout, func() { SoftWares.SQLyogSave(path) })
	runWithTimeout(&wg, timeout, func() { SoftWares.SecureCRTSave(path) })
	runWithTimeout(&wg, timeout, func() { SoftWares.DBeaverSave(path) })

	// SystemInfos
	runWithTimeout(&wg, timeout, func() { SystemInfos.InstalledAppSave(path) })
	runWithTimeout(&wg, timeout, func() { SystemInfos.SystemInfoSave(path) })
	runWithTimeout(&wg, timeout, func() { SystemInfos.RegEditSave(path) })
	runWithTimeout(&wg, timeout, func() { SystemInfos.WindowsLogsSave(path) })
	runWithTimeout(&wg, timeout, func() { SystemInfos.ScreenShotInfoSave(path) })
	runWithTimeout(&wg, timeout, func() { SystemInfos.WifiInfoSave(path) })

	wg.Wait()
	fmt.Println("已完成全部取证，压缩文件中：")
	err := utils.ZipDirectory(path, "ForensicPro_result.zip")
	if err != nil {
		return
	}
}

func runWithTimeout(wg *sync.WaitGroup, timeout time.Duration, f func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		done := make(chan bool)
		go func() {
			f()
			done <- true
		}()
		select {
		case <-done:
			// 函数执行完成
		case <-time.After(timeout):
			fmt.Println("函数执行超时")
		}
	}()
}
