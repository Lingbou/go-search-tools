package config

import (
	"time"
)

// SearchConfig 包含所有搜索相关的配置选项
type SearchConfig struct {
	// 通用选项
	SearchPath   string
	IgnoreCase   bool
	ColorOutput  bool
	ShowProgress bool

	// 文件过滤选项
	Recursive   bool
	MaxDepth    int
	ExcludeDirs []string
	IncludeExts []string
	ExcludeExts []string

	// 内容搜索选项
	NumWorkers int
	Timeout    time.Duration
}

// NewDefaultConfig 返回默认配置
func NewDefaultConfig() *SearchConfig {
	return &SearchConfig{
		SearchPath:   ".",
		IgnoreCase:   false,
		ColorOutput:  true,
		ShowProgress: false,
		Recursive:    true,
		MaxDepth:     -1,
		ExcludeDirs:  []string{},
		IncludeExts:  []string{},
		ExcludeExts:  []string{},
		NumWorkers:   4,
		Timeout:      0,
	}
}