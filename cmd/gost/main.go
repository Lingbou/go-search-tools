package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	// "github.com/fatih/color"
	
	"github.com/Lingbou/go-search-tools/internal/config"
	"github.com/Lingbou/go-search-tools/internal/search"
)

var (
	// 配置对象
	cfg = config.NewDefaultConfig()
	
	// 根命令
	rootCmd = &cobra.Command{
		Use:   "gost [command]",
		Short: "文件搜索工具，支持按文件名和内容搜索",
		Long:  "gost 是一个强大的文件搜索工具，可以快速在目录树中查找文件或内容",
	}

	// 搜索文件名的命令
	searchNameCmd = &cobra.Command{
		Use:   "name [flags] <pattern>",
		Short: "按文件名搜索",
		Long:  "按文件名搜索，支持通配符(*和?)",
		Args:  cobra.ExactArgs(1),
		Run:   runSearchName,
	}

	// 搜索文件内容的命令
	searchContentCmd = &cobra.Command{
		Use:   "content [flags] <pattern>",
		Short: "按文件内容搜索",
		Long:  "按文件内容搜索，使用字符串匹配",
		Args:  cobra.ExactArgs(1),
		Run:   runSearchContent,
	}

	// 正则表达式搜索命令
	searchRegexCmd = &cobra.Command{
		Use:   "regex [flags] <pattern>",
		Short: "使用正则表达式搜索文件内容",
		Long:  "使用正则表达式搜索文件内容，支持完整的正则表达式语法",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pattern := args[0]
			
			// 创建正则表达式搜索器
			searcher := search.NewRegexSearcher(cfg, pattern)
			
			// 执行搜索
			if err := searcher.Search(); err != nil {
				os.Exit(1)
			}
		},
	}
)

func init() {
	// 全局参数
	rootCmd.PersistentFlags().StringVarP(&cfg.SearchPath, "path", "p", ".", "搜索路径")
	rootCmd.PersistentFlags().BoolVarP(&cfg.IgnoreCase, "ignore-case", "i", false, "忽略大小写")
	rootCmd.PersistentFlags().BoolVarP(&cfg.ColorOutput, "color", "c", true, "启用颜色输出")
	rootCmd.PersistentFlags().BoolVarP(&cfg.ShowProgress, "progress", "P", false, "显示进度")
	
	// 文件名搜索参数
	searchNameCmd.Flags().BoolVarP(&cfg.Recursive, "recursive", "r", true, "递归搜索子目录")
	searchNameCmd.Flags().IntVarP(&cfg.MaxDepth, "max-depth", "d", -1, "最大递归深度，-1表示不限制")
	searchNameCmd.Flags().StringSliceVarP(&cfg.ExcludeDirs, "exclude-dir", "e", []string{}, "排除的目录")
	searchNameCmd.Flags().StringSliceVarP(&cfg.IncludeExts, "include-ext", "I", []string{}, "只包含的文件扩展名")
	searchNameCmd.Flags().StringSliceVarP(&cfg.ExcludeExts, "exclude-ext", "E", []string{}, "排除的文件扩展名")
	
	// 内容搜索参数
	searchContentCmd.Flags().BoolVarP(&cfg.Recursive, "recursive", "r", true, "递归搜索子目录")
	searchContentCmd.Flags().IntVarP(&cfg.MaxDepth, "max-depth", "d", -1, "最大递归深度，-1表示不限制")
	searchContentCmd.Flags().StringSliceVarP(&cfg.ExcludeDirs, "exclude-dir", "e", []string{}, "排除的目录")
	searchContentCmd.Flags().StringSliceVarP(&cfg.IncludeExts, "include-ext", "I", []string{}, "只包含的文件扩展名")
	searchContentCmd.Flags().StringSliceVarP(&cfg.ExcludeExts, "exclude-ext", "E", []string{}, "排除的文件扩展名")
	searchContentCmd.Flags().IntVarP(&cfg.NumWorkers, "workers", "w", 4, "并行工作线程数")
	searchContentCmd.Flags().DurationVarP(&cfg.Timeout, "timeout", "t", 0, "搜索超时时间，例如10s, 2m等")
	
	// 正则表达式搜索参数
	searchRegexCmd.Flags().BoolVarP(&cfg.Recursive, "recursive", "r", true, "递归搜索子目录")
	searchRegexCmd.Flags().IntVarP(&cfg.MaxDepth, "max-depth", "d", -1, "最大递归深度，-1表示不限制")
	searchRegexCmd.Flags().StringSliceVarP(&cfg.ExcludeDirs, "exclude-dir", "e", []string{}, "排除的目录")
	searchRegexCmd.Flags().StringSliceVarP(&cfg.IncludeExts, "include-ext", "I", []string{}, "只包含的文件扩展名")
	searchRegexCmd.Flags().StringSliceVarP(&cfg.ExcludeExts, "exclude-ext", "E", []string{}, "排除的文件扩展名")
	searchRegexCmd.Flags().IntVarP(&cfg.NumWorkers, "workers", "w", 4, "并行工作线程数")
	searchRegexCmd.Flags().DurationVarP(&cfg.Timeout, "timeout", "t", 0, "搜索超时时间，例如10s, 2m等")
	
	// 将子命令添加到根命令
	rootCmd.AddCommand(searchNameCmd, searchContentCmd, searchRegexCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// 按文件名搜索的执行函数
func runSearchName(cmd *cobra.Command, args []string) {
	pattern := args[0]
	
	// 创建文件名搜索器
	searcher := search.NewNameSearcher(cfg)
	
	// 执行搜索
	if err := searcher.Search(pattern); err != nil {
		os.Exit(1)
	}
}

// 按内容搜索的执行函数
func runSearchContent(cmd *cobra.Command, args []string) {
	pattern := args[0]
	
	// 创建内容搜索器
	searcher := search.NewContentSearcher(cfg, pattern)
	
	// 执行搜索
	if err := searcher.Search(); err != nil {
		os.Exit(1)
	}
}