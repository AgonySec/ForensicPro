package main

import (
	"ForensicPro/pkg/Browsers"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("欢迎使用ForensicPro V1.0 by:Agony")
	fmt.Println("下面开始进行Windows取证")
	//path := "C:\\Users\\Administrator\\Desktop\\test.txt"
	//utils.ReadSQLiteDB(path, "select * from cookies")
	//Browsers.ChromeSave("ForensicPro")
	Browsers.FirefoxSave("ForensicPro")
	//Browsers.IESave("ForensicPro")
	//fmt.Println(Browsers.GetMasterKey())
	//Browsers.ChromeCookies()
	fmt.Println("取证结束")
}
