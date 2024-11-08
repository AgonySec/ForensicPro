package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func ZipDirectory(srcDir, zipFile string) error {
	// 创建一个新的 zip 文件
	zf, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer zf.Close()

	// 创建一个 zip 写入器
	zw := zip.NewWriter(zf)
	defer zw.Close()

	// 遍历目录
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身
		if path == srcDir {
			return nil
		}

		// 获取文件相对于源目录的路径
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 设置文件路径
		header.Name = filepath.Join(filepath.Base(srcDir), path[len(srcDir):])

		// 如果是目录，设置 header 的类型
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		// 创建新的 zip 文件头
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		// 如果是文件，写入文件内容
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
	// 删除源目录
	err = os.RemoveAll(srcDir)
	if err != nil {
		return err
	}
	return err
}
