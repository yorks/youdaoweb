package youdaoweb

import (
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// 映射 HTML 标签到 fatih/color 样式
var tagStyles = map[string]*color.Color{
	"b":      color.New(color.Bold),
	"strong": color.New(color.Bold),
	"i":      color.New(color.Italic),
	"em":     color.New(color.Italic),
	"u":      color.New(color.Underline),
}

// 正则：匹配 <tag>...</tag>，不使用反向引用
var tagRe = regexp.MustCompile(`(?i)<([a-z]+)>(.*?)</[a-z]+>`)

// 递归解析 HTML 样式
func RenderHTML(input string) string {
	for {
		m := tagRe.FindStringSubmatch(input)
		if m == nil {
			break
		}

		fullMatch := m[0]
		tag := strings.ToLower(m[1])
		content := m[2]

		// 递归处理嵌套标签
		content = RenderHTML(content)

		var replacement string
		if style, ok := tagStyles[tag]; ok {
			replacement = style.Sprint(content)
		} else {
			replacement = content // 未定义样式的标签直接去掉
		}

		input = strings.Replace(input, fullMatch, replacement, 1)
	}

	return input
}

