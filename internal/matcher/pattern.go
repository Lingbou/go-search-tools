package matcher

import (
	"strings"
)

// MatchPattern 实现简单的通配符匹配
// 支持 * 和 ? 通配符
func MatchPattern(s, pattern string, ignoreCase bool) bool {
	if ignoreCase {
		s = strings.ToLower(s)
		pattern = strings.ToLower(pattern)
	}

	return matchPatternInternal(s, pattern)
}

// matchPatternInternal 是实际的匹配逻辑
func matchPatternInternal(s, pattern string) bool {
	if pattern == "" {
		return s == ""
	}
	
	if pattern[0] == '*' {
		if len(pattern) == 1 {
			return true
		}
		
		for i := 0; i <= len(s); i++ {
			if matchPatternInternal(s[i:], pattern[1:]) {
				return true
			}
		}
		return false
	}
	
	return len(s) > 0 && (pattern[0] == '?' || pattern[0] == s[0]) && matchPatternInternal(s[1:], pattern[1:])
}