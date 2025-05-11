package matcher

import (
	"context"
	"os"
	"strings"
)

// ContentMatcher 提供文件内容匹配功能
type ContentMatcher struct {
	Pattern    string
	IgnoreCase bool
}

// NewContentMatcher 创建一个新的内容匹配器
func NewContentMatcher(pattern string, ignoreCase bool) *ContentMatcher {
	return &ContentMatcher{
		Pattern:    pattern,
		IgnoreCase: ignoreCase,
	}
}

// MatchFile 检查文件内容是否匹配模式
func (m *ContentMatcher) MatchFile(ctx context.Context, filePath string) (bool, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()
	
	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		return false, err
	}
	
	size := info.Size()
	
	// 限制读取大小，避免处理大文件
	if size > 10*1024*1024 { // 10MB
		return false, nil
	}
	
	// 读取文件内容
	content := make([]byte, size)
	_, err = file.Read(content)
	if err != nil {
		return false, err
	}
	
	// 检查内容是否包含模式
	if m.IgnoreCase {
		return strings.Contains(strings.ToLower(string(content)), strings.ToLower(m.Pattern)), nil
	}
	
	return strings.Contains(string(content), m.Pattern), nil
}