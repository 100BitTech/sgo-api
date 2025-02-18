package base

import (
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/samber/oops"
)

// 驼峰改为蛇形
func CaseToSnake(name string) string {
	newName := []rune{}

	// 上一个大写的字母索引
	prevUpper := -1

	for i, c := range name {
		if unicode.IsUpper(c) {
			if i-prevUpper == 1 {
				newName = append(newName, unicode.ToLower(c))
			} else {
				newName = append(newName, '_')
				newName = append(newName, unicode.ToLower(c))
			}

			prevUpper = i
		} else {
			newName = append(newName, c)
		}
	}

	return string(newName)
}

// 模板化字符串
func Tprintf(tmpl string, args map[string]any) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", oops.Wrap(err)
	}

	b := &strings.Builder{}
	if err := t.Execute(b, args); err != nil {
		return "", oops.Wrap(err)
	}

	return b.String(), nil
}

func ExecuteTemplate(t *template.Template, name string, args map[string]any) (string, error) {
	b := &strings.Builder{}
	if err := t.ExecuteTemplate(b, name, args); err != nil {
		return "", oops.Wrap(err)
	}
	return b.String(), nil
}

// 压缩字符串，去除换行和减少空格
func CompressString(s string) string {
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}

// 截断符串，去除换行和减少空格后，只保留最大长度的字符串
func TruncateString(s string, maxLength int) string {
	s = CompressString(s)

	if maxLength == 0 {
		maxLength = 25
	}
	if len(s) > maxLength {
		s = s[0:maxLength]
	}
	s = strings.TrimSpace(s)

	return s
}
