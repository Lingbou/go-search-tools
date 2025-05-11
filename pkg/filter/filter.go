package filter

import (
	"os"
	"path/filepath"
	"strings"
)

// FileFilter 定义文件过滤器接口
type FileFilter interface {
	ShouldInclude(path string, info os.FileInfo) bool
}

// CompositeFilter 组合多个过滤器
type CompositeFilter struct {
	filters []FileFilter
}

// NewCompositeFilter 创建一个新的组合过滤器
func NewCompositeFilter(filters ...FileFilter) *CompositeFilter {
	return &CompositeFilter{
		filters: filters,
	}
}

// ShouldInclude 检查是否应该包含文件
// 所有过滤器都必须返回 true 才会包含
func (f *CompositeFilter) ShouldInclude(path string, info os.FileInfo) bool {
	for _, filter := range f.filters {
		if !filter.ShouldInclude(path, info) {
			return false
		}
	}
	return true
}

// DirectoryFilter 目录过滤器
type DirectoryFilter struct {
	ExcludeDirs []string
	MaxDepth    int
	BasePath    string
}

// NewDirectoryFilter 创建一个新的目录过滤器
func NewDirectoryFilter(basePath string, excludeDirs []string, maxDepth int) *DirectoryFilter {
	return &DirectoryFilter{
		ExcludeDirs: excludeDirs,
		MaxDepth:    maxDepth,
		BasePath:    basePath,
	}
}

// ShouldInclude 检查是否应该包含目录
func (f *DirectoryFilter) ShouldInclude(path string, info os.FileInfo) bool {
	if !info.IsDir() {
		return true
	}
	
	// 检查是否是要排除的目录
	for _, excludeDir := range f.ExcludeDirs {
		if info.Name() == excludeDir {
			return false
		}
	}
	
	// 检查深度
	if f.MaxDepth >= 0 {
		depth := strings.Count(strings.TrimPrefix(path, f.BasePath), string(os.PathSeparator))
		if depth > f.MaxDepth {
			return false
		}
	}
	
	return true
}

// ExtensionFilter 扩展名过滤器
type ExtensionFilter struct {
	IncludeExts []string
	ExcludeExts []string
}

// NewExtensionFilter 创建一个新的扩展名过滤器
func NewExtensionFilter(includeExts, excludeExts []string) *ExtensionFilter {
	return &ExtensionFilter{
		IncludeExts: includeExts,
		ExcludeExts: excludeExts,
	}
}

// ShouldInclude 检查是否应该包含文件
func (f *ExtensionFilter) ShouldInclude(path string, info os.FileInfo) bool {
	if info.IsDir() {
		return true
	}
	
	ext := filepath.Ext(info.Name())
	
	// 检查包含扩展名
	if len(f.IncludeExts) > 0 {
		found := false
		for _, includeExt := range f.IncludeExts {
			if ext == includeExt {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// 检查排除扩展名
	if len(f.ExcludeExts) > 0 {
		for _, excludeExt := range f.ExcludeExts {
			if ext == excludeExt {
				return false
			}
		}
	}
	
	return true
}