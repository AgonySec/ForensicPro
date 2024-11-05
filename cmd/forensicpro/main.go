package main

import (
	"ForensicPro/pkg/SoftWares"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func main() {
	startTime := time.Now() // 记录程序开始时间

	//fmt.Println("欢迎使用ForensicPro V1.0 by:Agony")
	//fmt.Println("下面开始进行Windows取证")
	//
	//Browsers.ChromeSave("ForensicPro")
	//Browsers.FirefoxSave("ForensicPro")
	//Browsers.IESave("ForensicPro")
	//
	//Messengers.LineSave("ForensicPro")
	//Messengers.QQSave("ForensicPro")
	//Messengers.DiscordSave("ForensicPro")
	//Messengers.TelegramSave("ForensicPro")
	//Messengers.EnigmaSave("ForensicPro")
	//Messengers.DingTalkSave("ForensicPro")
	//
	//FTPS.FileZillaSave("ForensicPro")
	//FTPS.SnowflakeSave("ForensicPro")
	//
	//Mails.FoxmailSave("ForensicPro")
	//
	//SoftWares.NeteaseCloudMusicSave("ForensicPro")
	//SoftWares.NavicatSave("ForensicPro")
	//SoftWares.VSCodeSave("ForensicPro")
	SoftWares.XmanagerSave("ForensicPro")
	//
	//SystemInfos.ScreenShotInfoSave("ForensicPro")
	//SystemInfos.WifiInfoSave("ForensicPro")
	//SystemInfos.InstalledAppSave("ForensicPro")
	endTime := time.Now()                 // 记录程序结束时间
	elapsedTime := endTime.Sub(startTime) // 计算程序运行时间

	//fmt.Println("取证结束")
	fmt.Printf("程序运行时间: %v\n", elapsedTime)

}
