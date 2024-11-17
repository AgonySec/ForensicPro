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

	status := utils.IsAdmin()
	if status == false {
		fmt.Println("请使用管理员权限运行本程序")
		return
	}
	GetAll("result")
	endTime := time.Now()                 // 记录程序结束时间
	elapsedTime := endTime.Sub(startTime) // 计算程序运行时间
	fmt.Println("=================Windows取证结束=================")
	fmt.Printf("程序运行时间: %v\n", elapsedTime)

}

func GetAll(path string) {
	var wg sync.WaitGroup

	// Browsers
	wg.Add(4)
	go func() { Browsers.ChromeSave(path); wg.Done() }()
	go func() { Browsers.FirefoxSave(path); wg.Done() }()
	go func() { Browsers.IESave(path); wg.Done() }()
	go func() { Browsers.SogouSave(path); wg.Done() }()

	// Messengers
	wg.Add(7)
	go func() { Messengers.LineSave(path); wg.Done() }()
	go func() { Messengers.QQSave(path); wg.Done() }()
	go func() { Messengers.DiscordSave(path); wg.Done() }()
	go func() { Messengers.TelegramSave(path); wg.Done() }()
	go func() { Messengers.EnigmaSave(path); wg.Done() }()
	go func() { Messengers.DingTalkSave(path); wg.Done() }()
	go func() { Messengers.SkypeSave(path); wg.Done() }()

	// FTPS
	wg.Add(3)
	go func() { FTPS.FileZillaSave(path); wg.Done() }()
	go func() { FTPS.SnowflakeSave(path); wg.Done() }()
	go func() { FTPS.WinSCPSave(path); wg.Done() }()

	// Mails
	wg.Add(4)
	go func() { Mails.FoxmailSave(path); wg.Done() }()
	go func() { Mails.MailBirdSave(path); wg.Done() }()
	go func() { Mails.OutlookSave(path); wg.Done() }()
	go func() { Mails.MailMasterSave(path); wg.Done() }()

	// SoftWares
	wg.Add(8)
	go func() { SoftWares.NeteaseCloudMusicSave(path); wg.Done() }()
	go func() { SoftWares.NavicatSave(path); wg.Done() }()
	go func() { SoftWares.VSCodeSave(path); wg.Done() }()
	go func() { SoftWares.XmanagerSave(path); wg.Done() }()
	go func() { SoftWares.FinalShellSave(path); wg.Done() }()
	go func() { SoftWares.SQLyogSave(path); wg.Done() }()
	go func() { SoftWares.SecureCRTSave(path); wg.Done() }()
	go func() { SoftWares.DBeaverSave(path); wg.Done() }()

	// SystemInfos
	wg.Add(6)
	go func() { SystemInfos.ScreenShotInfoSave(path); wg.Done() }()
	go func() { SystemInfos.WifiInfoSave(path); wg.Done() }()
	go func() { SystemInfos.InstalledAppSave(path); wg.Done() }()
	go func() { SystemInfos.SystemInfoSave(path); wg.Done() }()
	go func() { SystemInfos.RegEditSave(path); wg.Done() }()
	go func() { SystemInfos.WindowsLogsSave(path); wg.Done() }()

	wg.Wait()
	err := utils.ZipDirectory(path, "ForensicPro_result.zip")
	if err != nil {
		return
	}

}
