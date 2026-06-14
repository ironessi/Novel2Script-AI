package filecheck

import (
	"fmt"
	"path/filepath"
	"strings"
)

// AllowedTypes 允许的文件类型
var AllowedTypes = map[string]bool{
	".txt":  true,
	".md":   true,
	".docx": true,
}

// AllowedMIMEs 允许的 MIME 类型
var AllowedMIMEs = map[string]bool{
	"text/plain":    true,
	"text/markdown": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

// MaxFileSize 最大文件大小（字节），默认 20MB
var MaxFileSize int64 = 20 * 1024 * 1024

// CheckExtension 检查文件扩展名是否合法
func CheckExtension(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if !AllowedTypes[ext] {
		return fmt.Errorf("不支持的文件类型: %s，仅支持 .txt/.md/.docx", ext)
	}
	return nil
}

// CheckMIME 检查 MIME 类型是否合法
func CheckMIME(mimeType string) error {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	if idx := strings.Index(mimeType, ";"); idx >= 0 {
		mimeType = strings.TrimSpace(mimeType[:idx])
	}
	if !AllowedMIMEs[mimeType] {
		return fmt.Errorf("不支持的 MIME 类型: %s", mimeType)
	}
	return nil
}

// CheckFileSize 检查文件大小是否超限
func CheckFileSize(size int64) error {
	if size > MaxFileSize {
		return fmt.Errorf("文件大小 %d 字节超过限制 %d 字节", size, MaxFileSize)
	}
	if size == 0 {
		return fmt.Errorf("文件为空")
	}
	return nil
}

// SanitizeFilename 清洗文件名，防止路径穿越
func SanitizeFilename(filename string) string {
	// 只保留文件名部分，去除路径
	filename = filepath.Base(filename)
	// 替换危险字符
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	return filename
}

// CheckFilename 检查文件名是否安全
func CheckFilename(filename string) error {
	cleaned := SanitizeFilename(filename)
	if cleaned == "" || cleaned == "." {
		return fmt.Errorf("文件名无效")
	}
	// 检查是否包含脚本标签
	lower := strings.ToLower(cleaned)
	if strings.Contains(lower, "<script") || strings.Contains(lower, "javascript:") {
		return fmt.Errorf("文件名包含非法内容")
	}
	return nil
}
