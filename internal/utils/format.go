package utils

import (
	"fmt"
	"os"
	// "time"

	"github.com/fatih/color"
)

// FormatSize 格式化文件大小
func FormatSize(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	var i int
	size := float64(bytes)
	
	for size >= 1024 && i < len(units)-1 {
		size /= 1024
		i++
	}
	
	return fmt.Sprintf("%.2f %s", size, units[i])
}

// PrintMatch 打印匹配结果
func PrintMatch(path string, info os.FileInfo, useColor bool) {
	if useColor {
		// 使用颜色区分不同部分
		fileType := ""
		if info.IsDir() {
			fileType = color.CyanString("[目录]")
		} else {
			fileType = color.GreenString("[文件]")
		}
		
		size := ""
		if !info.IsDir() {
			size = FormatSize(info.Size())
		}
		
		modified := info.ModTime().Format("2006-01-02 15:04:05")
		
		fmt.Printf("%s %s %s %s\n", 
			fileType, 
			color.MagentaString(info.Name()), 
			color.BlueString(size), 
			color.YellowString(modified))
	} else {
		fmt.Printf("%s %s %d %s\n", 
			info.Mode(), 
			info.Name(), 
			info.Size(), 
			info.ModTime().Format("2006-01-02 15:04:05"))
	}
}