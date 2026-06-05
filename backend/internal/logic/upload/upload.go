package upload

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"novel2script-backend/internal/config"
	"novel2script-backend/internal/dao"
	"novel2script-backend/internal/model/entity"
	"novel2script-backend/internal/service"
	"novel2script-backend/utility/filecheck"
	"novel2script-backend/utility/sanitizer"

	"github.com/gogf/gf/v2/util/grand"
)

type uploadImpl struct{}

func init() {
	service.Upload = &uploadImpl{}
}

// Upload 处理文件上传：校验 → 存储 → 切分章节 → 入库
func (u *uploadImpl) Upload(ctx context.Context, userId, projectId int64, filename string, fileSize int64, mimeType string, reader io.Reader) (*entity.NovelSourceFile, []entity.NovelChapter, error) {
	// 1. 文件校验
	if err := filecheck.CheckFilename(filename); err != nil {
		return nil, nil, err
	}
	if err := filecheck.CheckExtension(filename); err != nil {
		return nil, nil, err
	}
	if err := filecheck.CheckFileSize(fileSize); err != nil {
		return nil, nil, err
	}

	// 2. 生成安全存储路径
	cleanName := filecheck.SanitizeFilename(filename)
	ext := strings.ToLower(filepath.Ext(cleanName))
	uuid := grand.S(16)
	storageDir := filepath.Join(config.C.Upload.StoragePath, "uploads", fmt.Sprintf("%d", userId), fmt.Sprintf("%d", projectId))
	storagePath := filepath.Join(storageDir, uuid+ext)

	// 3. 创建目录并保存文件
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("创建存储目录失败: %w", err)
	}

	dst, err := os.Create(storagePath)
	if err != nil {
		return nil, nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 同时读取内容和计算 hash
	hasher := sha256.New()
	writer := io.MultiWriter(dst, hasher)

	written, err := io.Copy(writer, reader)
	if err != nil {
		os.Remove(storagePath)
		return nil, nil, fmt.Errorf("保存文件失败: %w", err)
	}

	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// 4. 读取文件内容用于章节切分
	content, err := os.ReadFile(storagePath)
	if err != nil {
		return nil, nil, fmt.Errorf("读取文件内容失败: %w", err)
	}

	textContent := string(content)
	// 如果是 docx，提取纯文本（简化处理）
	if ext == ".docx" {
		textContent = extractTextFromDocx(textContent)
	}

	// 清洗文本
	textContent = sanitizer.RemoveControlChars(textContent)
	textContent = sanitizer.CleanText(textContent)

	if len(strings.TrimSpace(textContent)) == 0 {
		os.Remove(storagePath)
		return nil, nil, fmt.Errorf("文件内容为空")
	}

	// 5. 切分章节
	chapters := splitChapters(projectId, textContent)

	// 6. 保存源文件记录
	sourceFile := &entity.NovelSourceFile{
		ProjectId:        projectId,
		OwnerId:          userId,
		OriginalFilename: cleanName,
		StoragePath:      storagePath,
		FileHash:         fileHash,
		FileSize:         written,
		MimeType:         mimeType,
	}

	fileId, err := dao.CreateSourceFile(ctx, sourceFile)
	if err != nil {
		os.Remove(storagePath)
		return nil, nil, fmt.Errorf("保存文件记录失败: %w", err)
	}
	sourceFile.Id = fileId

	// 7. 保存章节到数据库
	if err := dao.BatchCreateChapters(ctx, chapters); err != nil {
		return sourceFile, nil, fmt.Errorf("保存章节失败: %w", err)
	}

	// 8. 更新项目状态
	_ = dao.UpdateProjectStatus(ctx, projectId, "uploaded")

	return sourceFile, chapters, nil
}

// splitChapters 将文本切分为章节
func splitChapters(projectId int64, text string) []entity.NovelChapter {
	lines := strings.Split(text, "\n")

	// 识别章节标题行
	type chapterMark struct {
		index int
		title string
		line  int
	}

	var marks []chapterMark
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if sanitizer.IsChapterTitle(trimmed) {
			marks = append(marks, chapterMark{
				index: len(marks) + 1,
				title: sanitizer.ExtractChapterTitle(trimmed),
				line:  i,
			})
		}
	}

	// 如果没有识别到章节，将整个文本作为单章
	if len(marks) == 0 {
		content := strings.TrimSpace(text)
		if content == "" {
			return nil
		}
		hash := sha256.Sum256([]byte(content))
		return []entity.NovelChapter{
			{
				ProjectId:    projectId,
				ChapterIndex: 1,
				ChapterTitle: "全文",
				Content:      content,
				ContentHash:  hex.EncodeToString(hash[:]),
			},
		}
	}

	// 按章节切分
	var chapters []entity.NovelChapter
	for i, mark := range marks {
		startLine := mark.line
		endLine := len(lines)
		if i+1 < len(marks) {
			endLine = marks[i+1].line
		}

		// 提取章节内容（跳过标题行）
		contentLines := lines[startLine+1 : endLine]
		content := strings.TrimSpace(strings.Join(contentLines, "\n"))
		content = sanitizer.CleanChapterContent(content)

		if content == "" {
			continue
		}

		hash := sha256.Sum256([]byte(content))
		chapters = append(chapters, entity.NovelChapter{
			ProjectId:    projectId,
			ChapterIndex: mark.index,
			ChapterTitle: mark.title,
			Content:      content,
			ContentHash:  hex.EncodeToString(hash[:]),
		})
	}

	return chapters
}

// extractTextFromDocx 从 docx 内容提取纯文本（简化版）
// 完整实现应使用 gooxml 库解析 XML
func extractTextFromDocx(content string) string {
	// 简单提取：<w:t> 标签中的文本
	var result strings.Builder
	lines := strings.Split(content, "<w:t")
	for _, line := range lines[1:] {
		endIdx := strings.Index(line, "</w:t>")
		if endIdx == -1 {
			endIdx = strings.Index(line, "/>")
			if endIdx == -1 {
				continue
			}
		}
		text := line[:endIdx]
		// 去除标签属性
		if idx := strings.Index(text, ">"); idx != -1 {
			text = text[idx+1:]
		}
		result.WriteString(text)
	}

	if result.Len() == 0 {
		// 如果无法解析，返回原文
		return content
	}
	return result.String()
}
