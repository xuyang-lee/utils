package internal

import (
	"strings"
)

const (
	Tag = "mapper"
)

// ParamsTag 解析tag, bool控制是否无视tag
func ParamsTag(tag string) (res string, ignore bool) {
	if tag == "" {
		return "", false
	}

	if tag == "-" {
		return "", true
	}

	var sb strings.Builder // 使用strings.Builder提高性能
	runeTag := []rune(tag)
	var jump bool
	for i, r := range runeTag { // 使用 range 循环处理 Unicode 字符
		if jump {
			jump = false
			continue
		}
		if r == '\\' {
			if i+1 < len(runeTag) {
				switch runeTag[i+1] {
				case '\\':
					sb.WriteRune('\\')
					jump = true
				case '-':
					sb.WriteRune('-')
					jump = true
				default:
					sb.WriteRune(r)
				}
			}
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String(), false

}
