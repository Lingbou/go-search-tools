package search

import (
	// "fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fatih/color"
	
	"github.com/Lingbou/go-search-tools/internal/config"
	"github.com/Lingbou/go-search-tools/internal/matcher"
	"github.com/Lingbou/go-search-tools/internal/utils"
	"github.com/Lingbou/go-search-tools/pkg/filter"
)

// NameSearcher 文件名搜索器
type NameSearcher struct {
	Config *config.SearchConfig
	Filter filter.FileFilter
}

// NewNameSearcher 创建一个新的文件名搜索器
func NewNameSearcher(cfg *config.SearchConfig) *NameSearcher {
	// 创建过滤器
	dirFilter := filter.NewDirectoryFilter(cfg.SearchPath, cfg.ExcludeDirs, cfg.MaxDepth)
	extFilter := filter.NewExtensionFilter(cfg.IncludeExts, cfg.ExcludeExts)
	compositeFilter := filter.NewCompositeFilter(dirFilter, extFilter)
	
	return &NameSearcher{
		Config: cfg,
		Filter: compositeFilter,
	}
}

// Search 执行文件名搜索
func (s *NameSearcher) Search(pattern string) error {
	// 检查路径是否存在
	if _, err := os.Stat(s.Config.SearchPath); os.IsNotExist(err) {
		color.Red("错误: 搜索路径不存在: %s", s.Config.SearchPath)
		return err
	}
	
	// 用于存储匹配的文件
	var matches []string
	var mu sync.Mutex
	
	// 创建进度跟踪器
	progress := utils.NewProgressTracker(s.Config.ShowProgress, "搜索中")
	
	// 计算文件总数用于进度条
	if s.Config.ShowProgress {
		totalFiles := utils.CountFiles(s.Config.SearchPath, s.Config.IncludeExts, s.Config.ExcludeExts)
		if totalFiles == 0 {
			color.Yellow("没有找到文件")
			return nil
		}
		progress.SetTotal(totalFiles)
	}
	
	// 递归搜索文件
	err := filepath.Walk(s.Config.SearchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// 更新进度条
		if s.Config.ShowProgress {
			progress.Increment()
		}
		
		// 应用过滤器
		if !s.Filter.ShouldInclude(path, info) {
			if info.IsDir() && path != s.Config.SearchPath {
				return filepath.SkipDir
			}
			return nil
		}
		
		// 对于目录，只检查过滤条件
		if info.IsDir() {
			return nil
		}
		
		// 文件名匹配
		if matcher.MatchPattern(info.Name(), pattern, s.Config.IgnoreCase) {
			mu.Lock()
			matches = append(matches, path)
			mu.Unlock()
			
			// 打印匹配结果
			utils.PrintMatch(path, info, s.Config.ColorOutput)
		}
		
		return nil
	})
	
	if err != nil {
		color.Red("搜索过程中出错: %v", err)
		return err
	}
	
	// 打印结果摘要
	if len(matches) == 0 {
		color.Yellow("没有找到匹配的文件")
	} else {
		color.Green("共找到 %d 个匹配的文件", len(matches))
	}
	
	return nil
}