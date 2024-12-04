# ForensicPro - 一款Windows自动取证工具

​																											![](https://badgen.net/github/stars/AgonySec/ForensicPro)![GitHub download counter](https://img.shields.io/github/downloads/AgonySec/ForensicPro/total)

## 介绍

基于golang实现的一款Windows自动取证工具，如果对你有帮助的话，不妨点个star🐟

　　取证内容包括：

　　1．     浏览器取证

　　2．     FTP工具取证

　　3．     邮箱工具取证

　　4．     通讯工具取证

　　5．     常用软件取证

　　6．     系统信息取证，包括WIFI、安装软件目录、截屏、systeminfo、注册表、硬件信息等等。。。



## 编译

GitHub直接clone源码：

```
git clone https://github.com/AgonySec/ForensicPro.git
```

编译：

```
go build -ldflags="-w -s" -o ForensicPro.exe -trimpath .\cmd\forensicpro\main.go
```

如果你不想编译，也可以直接去[release](https://github.com/AgonySec/ForensicPro/releases)页面下载已经打包好的EXE，下载最新版即可

## 快速使用

　　管理员方式运行ForensicPro.exe即可，也可以在管理员命令行窗口运行，无需联网，运行结束之后会生成一个ForensicPro_result.zip，取证结果都保存在里面压缩包中

![image-20241108171619145](assets/image-20241108171619145.png)

## 注意

开发环境：

- go 1.20.1
- Windows 11

## 许可证

　　本项目采用 [MIT License](LICENSE) 许可证。
