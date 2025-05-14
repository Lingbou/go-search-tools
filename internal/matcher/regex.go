package matcher

import (
	"bufio"
	"context"
	"os"
	"regexp"
	// "strings"
)

// RegexMatcher 正则表达式匹配器
type RegexMatcher struct {
	Pattern     string
	IgnoreCase  bool
	CompiledReg *regexp.Regexp
}

// NewRegexMatcher 创建一个新的正则表达式匹配器
func NewRegexMatcher(pattern string, ignoreCase bool) *RegexMatcher {
	// 处理忽略大小写
	flags := ""
	if ignoreCase {
		flags = "(?i)"
	}
	
	// 编译正则表达式
	reg := regexp.MustCompile(flags + pattern)
	
	return &RegexMatcher{
		Pattern:     pattern,
		IgnoreCase:  ignoreCase,
		CompiledReg: reg,
	}
}

// MatchFile 检查文件内容是否匹配正则表达式
func (m *RegexMatcher) MatchFile(ctx context.Context, filePath string) (bool, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()
	
	// 创建扫描器
	scanner := bufio.NewScanner(file)
	
	// 逐行扫描文件
	for scanner.Scan() {
		// 检查是否超时
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			line := scanner.Text()
			
			// 使用正则表达式匹配
			if m.CompiledReg.MatchString(line) {
				return true, nil
			}
		}
	}
	
	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		return false, err
	}
	
	// 没有找到匹配
	return false, nil
}