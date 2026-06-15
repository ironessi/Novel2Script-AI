package upload

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestExtractTextFromDocxPreservesParagraphsForChapterSplit(t *testing.T) {
	docx := buildTestDocx(t, []string{
		"《旧楼里的回声》",
		"第一章 雨夜来信",
		"雨下得很急。",
		"林舟站在旧楼门口。",
		"第二章 四楼灯光",
		"苏晚赶到实验楼。",
	})

	text, err := extractTextFromDocx(docx)
	if err != nil {
		t.Fatalf("extractTextFromDocx returned error: %v", err)
	}
	if !strings.Contains(text, "第一章 雨夜来信\n\n雨下得很急。") {
		t.Fatalf("expected paragraph breaks around first chapter, got:\n%s", text)
	}
	if !strings.Contains(text, "第二章 四楼灯光\n\n苏晚赶到实验楼。") {
		t.Fatalf("expected paragraph breaks around second chapter, got:\n%s", text)
	}

	chapters := splitChapters(42, text)
	if len(chapters) != 2 {
		t.Fatalf("expected 2 chapters, got %d: %#v", len(chapters), chapters)
	}
	if chapters[0].ChapterTitle != "第一章 雨夜来信" {
		t.Fatalf("unexpected first chapter title: %q", chapters[0].ChapterTitle)
	}
	if chapters[1].ChapterTitle != "第二章 四楼灯光" {
		t.Fatalf("unexpected second chapter title: %q", chapters[1].ChapterTitle)
	}
}

func TestExtractTextFromDocxRejectsInvalidDocx(t *testing.T) {
	if _, err := extractTextFromDocx([]byte("not a zip")); err == nil {
		t.Fatal("expected invalid docx to return an error")
	}
}

func buildTestDocx(t *testing.T, paragraphs []string) []byte {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatalf("create document.xml: %v", err)
	}

	var xml strings.Builder
	xml.WriteString(`<?xml version="1.0" encoding="UTF-8"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body>`)
	for _, paragraph := range paragraphs {
		xml.WriteString("<w:p><w:r><w:t>")
		xml.WriteString(paragraph)
		xml.WriteString("</w:t></w:r></w:p>")
	}
	xml.WriteString("</w:body></w:document>")

	if _, err := w.Write([]byte(xml.String())); err != nil {
		t.Fatalf("write document.xml: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return buf.Bytes()
}
