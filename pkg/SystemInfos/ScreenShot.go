package SystemInfos

import (
	"fmt"
	"github.com/kbinani/screenshot"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

var ScreenShotName = "ScreenShot"

func ScreenShotInfoSave(path string) {
	fmt.Println("ScreenShot")
	// 获取屏幕数量
	n := screenshot.NumActiveDisplays()
	if n == 0 {
		fmt.Println("No screens found.")
		return
	}
	targetPath := filepath.Join(path, ScreenShotName)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}
	// 遍历所有屏幕并保存截图
	for i := 0; i < n; i++ {
		// 获取当前屏幕的边界
		bounds := screenshot.GetDisplayBounds(i)

		// 捕获当前屏幕的截图
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			fmt.Printf("Failed to capture screen %d: %v\n", i, err)
			continue
		}

		// 创建文件并保存截图
		fileName := fmt.Sprintf("screen_%d.jpg", i)
		filePath := filepath.Join(targetPath, fileName)

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Failed to create file for screen %d: %v\n", i, err)
			continue
		}
		defer file.Close()

		// 将图片以JPEG格式写入文件
		err = jpeg.Encode(file, img, nil)
		if err != nil {
			fmt.Printf("Failed to save screenshot for screen %d: %v\n", i, err)
			continue
		}

		fmt.Printf("Screenshot for screen %d saved as %s\n", i, fileName)
	}
}
