package sanitizer

import (
	"regexp"
	"strings"
	"unicode"
)

// CleanText 清洗小说文本
func CleanText(text string) string {
	// 替换全角空格为半角
	text = strings.ReplaceAll(text, "　", " ")
	// 替换连续多个空格为单个空格（保留换行）
	re := regexp.MustCompile(`[^\S\n]+`)
	text = re.ReplaceAllString(text, " ")
	// 替换连续多个换行为两个换行
	re = regexp.MustCompile(`\n{3,}`)
	text = re.ReplaceAllString(text, "\n\n")
	// 去除首尾空白
	text = strings.TrimSpace(text)
	return text
}

// CleanChapterContent 清洗章节内容
func CleanChapterContent(content string) string {
	content = CleanText(content)
	// 去除章节标题行（如果内容开头包含）
	lines := strings.Split(content, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		if IsChapterTitle(firstLine) {
			lines = lines[1:]
			content = strings.Join(lines, "\n")
			content = strings.TrimSpace(content)
		}
	}
	return content
}

// chapterTitlePattern 章节标题正则
var chapterTitlePattern = regexp.MustCompile(`(?i)^(第[一二三四五六七八九十百千万零\d]+[章节回卷集部篇]|Chapter\s*\d+|CHAPTER\s*\d+|序章|楔子|尾声|番外)`)

// IsChapterTitle 判断是否为章节标题
func IsChapterTitle(line string) bool {
	line = strings.TrimSpace(line)
	if line == "" {
		return false
	}
	return chapterTitlePattern.MatchString(line)
}

// ExtractChapterTitle 提取章节标题
func ExtractChapterTitle(line string) string {
	return strings.TrimSpace(line)
}

// IsControlChar 检查是否为控制字符（保留换行和制表符）
func IsControlChar(r rune) bool {
	return unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t'
}

// RemoveControlChars 去除控制字符
func RemoveControlChars(text string) string {
	return strings.Map(func(r rune) rune {
		if IsControlChar(r) {
			return -1
		}
		return r
	}, text)
}
