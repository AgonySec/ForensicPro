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
	fmt.Println(`         _____                        _      ____            
  /\    |  ___|__  _ __ ___ _ __  ___(_) ___|  _ \ _ __ ___  
 /  \   | |_ / _ \| '__/ _ \ '_ \/ __| |/ __| |_) | '__/ _ \ 
/ /\ \  |  _| (_) | | |  __/ | | \__ \ | (__|  __/| | | (_) |
\/  \/  |_|  \___/|_|  \___|_| |_|___/_|\___|_|   |_|  \___/ `)
	fmt.Println("欢迎使用 ^ ForensicPro V1.0 一款Windows自动取证工具 by:Agony")
	fmt.Println("===================================================")
	fmt.Println("下面开始进行Windows取证")
	GetAll()
	endTime := time.Now()                 // 记录程序结束时间
	elapsedTime := endTime.Sub(startTime) // 计算程序运行时间
	fmt.Println("Windows取证结束")
	fmt.Printf("程序运行时间: %v\n", elapsedTime)

}

func GetAll() {
	var wg sync.WaitGroup

	// Browsers
	wg.Add(4)
	go func() { Browsers.ChromeSave("ForensicPro"); wg.Done() }()
	go func() { Browsers.FirefoxSave("ForensicPro"); wg.Done() }()
	go func() { Browsers.IESave("ForensicPro"); wg.Done() }()
	go func() { Browsers.SogouSave("ForensicPro"); wg.Done() }()

	// Messengers
	wg.Add(7)
	go func() { Messengers.LineSave("ForensicPro"); wg.Done() }()
	go func() { Messengers.QQSave("ForensicPro"); wg.Done() }()
	go func() { Messengers.DiscordSave("ForensicPro"); wg.Done() }()
	go func() { Messengers.TelegramSave("ForensicPro"); wg.Done() }()
	go func() { Messengers.EnigmaSave("ForensicPro"); wg.Done() }()
	go func() { Messengers.DingTalkSave("ForensicPro"); wg.Done() }()
	go func() { Messengers.SkypeSave("ForensicPro"); wg.Done() }()

	// FTPS
	wg.Add(3)
	go func() { FTPS.FileZillaSave("ForensicPro"); wg.Done() }()
	go func() { FTPS.SnowflakeSave("ForensicPro"); wg.Done() }()
	go func() { FTPS.WinSCPSave("ForensicPro"); wg.Done() }()

	// Mails
	wg.Add(4)
	go func() { Mails.FoxmailSave("ForensicPro"); wg.Done() }()
	go func() { Mails.MailBirdSave("ForensicPro"); wg.Done() }()
	go func() { Mails.OutlookSave("ForensicPro"); wg.Done() }()
	go func() { Mails.MailMasterSave("ForensicPro"); wg.Done() }()

	// SoftWares
	wg.Add(8)
	go func() { SoftWares.NeteaseCloudMusicSave("ForensicPro"); wg.Done() }()
	go func() { SoftWares.NavicatSave("ForensicPro"); wg.Done() }()
	go func() { SoftWares.VSCodeSave("ForensicPro"); wg.Done() }()
	go func() { SoftWares.XmanagerSave("ForensicPro"); wg.Done() }()
	go func() { SoftWares.FinalShellSave("ForensicPro"); wg.Done() }()
	go func() { SoftWares.SQLyogSave("ForensicPro"); wg.Done() }()
	go func() { SoftWares.SecureCRTSave("ForensicPro"); wg.Done() }()
	go func() { SoftWares.DBeaverSave("ForensicPro"); wg.Done() }()

	// SystemInfos
	wg.Add(3)
	go func() { SystemInfos.ScreenShotInfoSave("ForensicPro"); wg.Done() }()
	go func() { SystemInfos.WifiInfoSave("ForensicPro"); wg.Done() }()
	go func() { SystemInfos.InstalledAppSave("ForensicPro"); wg.Done() }()

	wg.Wait()

	err := utils.ZipDirectory("ForensicPro", "ForensicPro_result.zip")
	if err != nil {
		return
	}

}
