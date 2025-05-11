package utils

import (
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

// ProgressTracker 进度跟踪器
type ProgressTracker struct {
	Bar        *progressbar.ProgressBar
	TotalFiles int
	Enabled    bool
}

// NewProgressTracker 创建一个新的进度跟踪器
func NewProgressTracker(enabled bool, description string) *ProgressTracker {
	if !enabled {
		return &ProgressTracker{
			Enabled: false,
		}
	}
	
	return &ProgressTracker{
		Enabled: true,
		Bar:     progressbar.Default(-1, description),
	}
}

// SetTotal 设置总文件数
func (p *ProgressTracker) SetTotal(total int) {
	if !p.Enabled {
		return
	}
	
	p.TotalFiles = total
	p.Bar = progressbar.Default(int64(total), "搜索中")
}

// Increment 增加进度
func (p *ProgressTracker) Increment() {
	if !p.Enabled {
		return
	}
	
	p.Bar.Add(1)
}

// CountFiles 计算目录中的文件总数
func CountFiles(path string, includeExts, excludeExts []string) int {
	count := 0
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			// 扩展名过滤
			ext := filepath.Ext(info.Name())
			if len(includeExts) > 0 {
				found := false
				for _, includeExt := range includeExts {
					if ext == includeExt {
						found = true
						break
					}
				}
				if !found {
					return nil
				}
			}
			
			if len(excludeExts) > 0 {
				for _, excludeExt := range excludeExts {
					if ext == excludeExt {
						return nil
					}
				}
			}
			
			count++
		}
		return nil
	})
	return count
}