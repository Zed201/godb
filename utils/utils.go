package utils

import (
	"strings"
)

func StartWith(s, p string) bool {
	return strings.HasPrefix(s, p)
}

func JoinS(ss []string, i int, j int) string {
	if j <= 0 {
		j = len(ss) - (-1 * j)
	}
	return strings.Join(ss[i:j], "")
}
