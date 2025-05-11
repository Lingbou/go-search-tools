package search

import (
	"context"
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

// ContentSearcher 文件内容搜索器
type ContentSearcher struct {
	Config  *config.SearchConfig
	Filter  filter.FileFilter
	Matcher *matcher.ContentMatcher
}

// NewContentSearcher 创建一个新的内容搜索器
func NewContentSearcher(cfg *config.SearchConfig, pattern string) *ContentSearcher {
	// 创建过滤器
	dirFilter := filter.NewDirectoryFilter(cfg.SearchPath, cfg.ExcludeDirs, cfg.MaxDepth)
	extFilter := filter.NewExtensionFilter(cfg.IncludeExts, cfg.ExcludeExts)
	compositeFilter := filter.NewCompositeFilter(dirFilter, extFilter)
	
	// 创建内容匹配器
	contentMatcher := matcher.NewContentMatcher(pattern, cfg.IgnoreCase)
	
	return &ContentSearcher{
		Config:  cfg,
		Filter:  compositeFilter,
		Matcher: contentMatcher,
	}
}

// Search 执行内容搜索
func (s *ContentSearcher) Search() error {
	// 检查路径是否存在
	if _, err := os.Stat(s.Config.SearchPath); os.IsNotExist(err) {
		color.Red("错误: 搜索路径不存在: %s", s.Config.SearchPath)
		return err
	}
	
	// 创建上下文用于超时控制
	var ctx context.Context
	var cancel context.CancelFunc
	if s.Config.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), s.Config.Timeout)
		defer cancel()
	} else {
		ctx = context.Background()
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
	
	// 创建文件通道
	filesCh := make(chan string)
	resultsCh := make(chan string)
	
	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < s.Config.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range filesCh {
				// 检查是否超时
				select {
				case <-ctx.Done():
					return
				default:
					// 搜索文件内容
					matched, err := s.Matcher.MatchFile(ctx, filePath)
					if err == nil && matched {
						resultsCh <- filePath
					}
					
					// 更新进度条
					if s.Config.ShowProgress {
						progress.Increment()
					}
				}
			}
		}()
	}
	
	// 收集结果的协程
	go func() {
		wg.Wait()
		close(resultsCh)
	}()
	
	// 遍历文件并发送到通道
	go func() {
		defer close(filesCh)
		
		err := filepath.Walk(s.Config.SearchPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			// 检查是否超时
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
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
				
				// 发送文件路径到通道
				filesCh <- path
				
				return nil
			}
		})
		
		if err != nil && err != ctx.Err() {
			color.Red("搜索过程中出错: %v", err)
		}
	}()
	
	// 处理结果
	for filePath := range resultsCh {
		// 获取文件信息
		info, err := os.Stat(filePath)
		if err != nil {
			color.Red("获取文件信息失败: %s - %v", filePath, err)
			continue
		}
		
		mu.Lock()
		matches = append(matches, filePath)
		mu.Unlock()
		
		// 打印匹配结果
		utils.PrintMatch(filePath, info, s.Config.ColorOutput)
	}
	
	// 检查是否超时
	if ctx.Err() != nil {
		color.Yellow("搜索超时，已找到 %d 个匹配的文件", len(matches))
	} else if len(matches) == 0 {
		color.Yellow("没有找到匹配的文件")
	} else {
		color.Green("共找到 %d 个匹配的文件", len(matches))
	}
	
	return nil
}