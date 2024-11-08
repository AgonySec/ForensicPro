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
	Browsers.ChromeSave("ForensicPro")
	Browsers.FirefoxSave("ForensicPro")
	Browsers.IESave("ForensicPro")
	Browsers.SogouSave("ForensicPro")

	Messengers.LineSave("ForensicPro")
	Messengers.QQSave("ForensicPro")
	Messengers.DiscordSave("ForensicPro")
	Messengers.TelegramSave("ForensicPro")
	Messengers.EnigmaSave("ForensicPro")
	Messengers.DingTalkSave("ForensicPro")
	Messengers.SkypeSave("ForensicPro")

	FTPS.FileZillaSave("ForensicPro")
	FTPS.SnowflakeSave("ForensicPro")
	FTPS.WinSCPSave("ForensicPro")

	Mails.FoxmailSave("ForensicPro")
	Mails.MailBirdSave("ForensicPro")
	Mails.OutlookSave("ForensicPro")
	Mails.MailMasterSave("ForensicPro")

	SoftWares.NeteaseCloudMusicSave("ForensicPro")
	SoftWares.NavicatSave("ForensicPro")
	SoftWares.VSCodeSave("ForensicPro")
	SoftWares.XmanagerSave("ForensicPro")
	SoftWares.FinalShellSave("ForensicPro")
	SoftWares.SQLyogSave("ForensicPro")
	SoftWares.SecureCRTSave("ForensicPro")
	SoftWares.DBeaverSave("ForensicPro")

	SystemInfos.ScreenShotInfoSave("ForensicPro")
	SystemInfos.WifiInfoSave("ForensicPro")
	SystemInfos.InstalledAppSave("ForensicPro")

	err := utils.ZipDirectory("ForensicPro", "ForensicPro_result.zip")
	if err != nil {
		return
	}
}
