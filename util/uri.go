package util

import (
	"path"
	"strings"
)

// MatchURI
//'**'             匹配后面所有
//'*'              匹配0或多个非/的字符
//'?'              匹配1个非/的字符
func MatchURI(uri string, patterns ...string) bool {
	if len(patterns) == 0 {
		return false
	}
	if uri[len(uri)-1] == '/' {
		uri = uri[0 : len(uri)-1]
	}
	for _, pattern := range patterns {
		if match0(uri, pattern) {
			return true
		}
	}

	return false
}
func match0(uri string, pattern string) bool {
	// 从uri中分离出path和query
	pathURI := strings.Split(uri, "?")[0]
	// **匹配
	if len(pattern) > 2 && pattern[len(pattern)-2:] == "**" {
		return strings.HasPrefix(pathURI, pattern[0:len(pattern)-3])
	}
	matched, _ := path.Match(pattern, pathURI)
	return matched
}
