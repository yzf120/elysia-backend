package router

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
	agentpb "github.com/yzf120/elysia-backend/proto/agent"
	"github.com/yzf120/elysia-backend/rpc"
)

// ==================== 知识库管理路由（管理员） ====================

// RegisterKnowledgeRoutes 注册知识库管理路由
func RegisterKnowledgeRoutes(protectedRouter *mux.Router) {
	protectedRouter.HandleFunc("/admin/knowledge/upload", knowledgeUploadHandler).Methods("POST")
	protectedRouter.HandleFunc("/admin/knowledge/import", knowledgeImportHandler).Methods("POST")
	protectedRouter.HandleFunc("/admin/knowledge/list", knowledgeListHandler).Methods("GET")
	protectedRouter.HandleFunc("/admin/knowledge/delete", knowledgeDeleteHandler).Methods("POST")
	protectedRouter.HandleFunc("/admin/knowledge/search", knowledgeSearchHandler).Methods("POST")
}

// knowledgeUploadHandler 上传知识库文件（保存到临时目录）
func knowledgeUploadHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	// 限制上传大小 50MB
	r.ParseMultipartForm(50 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "文件上传失败: "+err.Error())
		return
	}
	defer file.Close()

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowedExts := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true, ".txt": true,
		".md": true, ".csv": true, ".html": true, ".htm": true,
	}
	if !allowedExts[ext] {
		writeErrorResponse(w, http.StatusBadRequest, "不支持的文件格式，支持 PDF、Word、TXT、Markdown、CSV、HTML")
		return
	}

	// 创建临时目录
	uploadDir := "/tmp/elysia_knowledge_uploads"
	os.MkdirAll(uploadDir, 0755)

	// 生成唯一文件名
	timestamp := time.Now().UnixNano()
	savedName := fmt.Sprintf("%d_%s", timestamp, handler.Filename)
	savedPath := filepath.Join(uploadDir, savedName)

	// 保存文件
	dst, err := os.Create(savedPath)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "保存文件失败: "+err.Error())
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "保存文件失败: "+err.Error())
		return
	}

	log.Printf("[Knowledge] 文件上传成功: %s, 大小: %d bytes", handler.Filename, handler.Size)

	writeSuccessResponse(w, map[string]interface{}{
		"message":   "上传成功",
		"file_name": handler.Filename,
		"file_size": handler.Size,
		"file_path": savedPath,
		"file_type": ext,
	})
}

// knowledgeImportHandler 导入知识库文件（解析文件内容 + 通过RPC调用chat-agent存储）
func knowledgeImportHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	// 解析请求体
	var reqBody struct {
		FilePath string `json:"file_path"` // 上传后的临时文件路径
		FileName string `json:"file_name"` // 原始文件名
		FileSize int64  `json:"file_size"` // 文件大小
		FileType string `json:"file_type"` // 文件扩展名
		Tags     string `json:"tags"`      // 标签
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if reqBody.FilePath == "" || reqBody.FileName == "" {
		writeErrorResponse(w, http.StatusBadRequest, "文件路径和文件名不能为空")
		return
	}

	// 解析文件内容
	content, err := parseFileContent(reqBody.FilePath, reqBody.FileType)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "文件解析失败: "+err.Error())
		return
	}

	if strings.TrimSpace(content) == "" {
		writeErrorResponse(w, http.StatusBadRequest, "文件内容为空，无法导入")
		return
	}

	// 生成文档ID
	docID := fmt.Sprintf("doc_%d", time.Now().UnixNano())

	// 通过 RPC 调用 chat-agent 的 StoreKnowledge 方法
	agentClient := rpc.GetAgentClient().GetProxy()
	resp, err := agentClient.StoreKnowledge(ctx, &agentpb.StoreKnowledgeRequest{
		ID:         docID,
		Content:    content,
		SourceType: "knowledge_base",
		SourceID:   reqBody.FileName,
		Tags:       reqBody.Tags,
		FileName:   reqBody.FileName,
		FileSize:   reqBody.FileSize,
		FileType:   reqBody.FileType,
	})
	if err != nil {
		log.Printf("[Knowledge] RPC调用StoreKnowledge失败: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "导入失败: "+err.Error())
		return
	}

	if !resp.Success {
		writeErrorResponse(w, http.StatusInternalServerError, "导入失败: "+resp.Message)
		return
	}

	// 清理临时文件
	os.Remove(reqBody.FilePath)

	log.Printf("[Knowledge] 文件导入成功: %s -> doc_id: %s, 预计索引构建时间: %ds", reqBody.FileName, resp.ID, resp.EstimatedSeconds)

	writeSuccessResponse(w, map[string]interface{}{
		"message":           "文档已接收，正在后台构建索引",
		"doc_id":            resp.ID,
		"estimated_seconds": resp.EstimatedSeconds,
	})
}

// knowledgeListHandler 获取知识库文档列表
func knowledgeListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	// 解析分页参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	keyword := r.URL.Query().Get("keyword")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	_ = keyword // 关键词搜索在 chat-agent 端实现

	// 通过 RPC 调用 chat-agent 的 ListKnowledge 方法
	agentClient := rpc.GetAgentClient().GetProxy()
	resp, err := agentClient.ListKnowledge(ctx, &agentpb.ListKnowledgeRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		log.Printf("[Knowledge] RPC调用ListKnowledge失败: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"total": resp.Total,
		"items": resp.Items,
		"page":  page,
	})
}

// knowledgeDeleteHandler 删除知识库文档
func knowledgeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	var reqBody struct {
		DocID string   `json:"doc_id"` // 单个删除
		IDs   []string `json:"ids"`    // 批量删除
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 收集要删除的 ID 列表
	var deleteIDs []string
	if reqBody.DocID != "" {
		deleteIDs = append(deleteIDs, reqBody.DocID)
	}
	deleteIDs = append(deleteIDs, reqBody.IDs...)

	if len(deleteIDs) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, "请指定要删除的文档ID")
		return
	}

	agentClient := rpc.GetAgentClient().GetProxy()
	var failedIDs []string

	for _, id := range deleteIDs {
		resp, err := agentClient.DeleteKnowledge(ctx, &agentpb.DeleteKnowledgeRequest{
			ID: id,
		})
		if err != nil || !resp.Success {
			failedIDs = append(failedIDs, id)
			log.Printf("[Knowledge] 删除文档失败: id=%s, err=%v", id, err)
		}
	}

	if len(failedIDs) > 0 {
		writeSuccessResponse(w, map[string]interface{}{
			"message":    fmt.Sprintf("部分删除成功，%d 个失败", len(failedIDs)),
			"failed_ids": failedIDs,
		})
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": fmt.Sprintf("成功删除 %d 个文档，索引清理在后台进行", len(deleteIDs)),
		"async":   true,
	})
}

// knowledgeSearchHandler 检索知识库
func knowledgeSearchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	var reqBody struct {
		Query string `json:"query"`
		TopK  int32  `json:"top_k"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if reqBody.Query == "" {
		writeErrorResponse(w, http.StatusBadRequest, "查询内容不能为空")
		return
	}

	if reqBody.TopK <= 0 {
		reqBody.TopK = 5
	}

	agentClient := rpc.GetAgentClient().GetProxy()
	resp, err := agentClient.SearchKnowledge(ctx, &agentpb.SearchKnowledgeRequest{
		Query: reqBody.Query,
		TopK:  reqBody.TopK,
	})
	if err != nil {
		log.Printf("[Knowledge] RPC调用SearchKnowledge失败: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "检索失败: "+err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"items": resp.Items,
	})
}

// ==================== 文件解析辅助函数 ====================

// parseFileContent 解析文件内容为纯文本
func parseFileContent(filePath string, fileType string) (string, error) {
	switch fileType {
	case ".txt", ".md":
		// 纯文本和 Markdown 直接按文本读取
		return parseTxtFile(filePath)
	case ".pdf":
		// 使用 ledongthuc/pdf 库解析 PDF 文件
		content, err := parsePdfFile(filePath)
		if err != nil {
			log.Printf("[Knowledge] PDF 解析失败，尝试降级为文本读取: %v", err)
			return parseTxtFile(filePath)
		}
		return content, nil
	case ".docx":
		// 使用 nguyenthenguyen/docx 库解析 DOCX 文件
		content, err := parseDocxFile(filePath)
		if err != nil {
			log.Printf("[Knowledge] DOCX 解析失败: %v", err)
			return "", fmt.Errorf("DOCX 文件解析失败: %w", err)
		}
		return content, nil
	case ".doc":
		// .doc 是旧版二进制格式，Go 生态暂无成熟纯 Go 解析库
		// 尝试提取其中的可读文本（有损提取）
		content, err := parseDocLegacyFile(filePath)
		if err != nil {
			return "", fmt.Errorf(".doc 文件解析失败，建议转换为 .docx 后重新上传: %w", err)
		}
		return content, nil
	case ".csv":
		return parseCsvFile(filePath)
	case ".html", ".htm":
		return parseHtmlFile(filePath)
	default:
		return parseTxtFile(filePath)
	}
}

// parseTxtFile 解析纯文本文件
func parseTxtFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	// 增大缓冲区以支持长行
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	return sb.String(), nil
}

// parsePdfFile 使用 ledongthuc/pdf 库解析 PDF 文件提取文本
func parsePdfFile(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开 PDF 文件失败: %w", err)
	}
	defer f.Close()

	var sb strings.Builder
	totalPages := r.NumPage()

	for i := 1; i <= totalPages; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		text, err := p.GetPlainText(nil)
		if err != nil {
			log.Printf("[Knowledge] PDF 第 %d 页文本提取失败: %v", i, err)
			continue
		}
		sb.WriteString(text)
		sb.WriteString("\n")
	}

	content := sb.String()
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("PDF 文件未提取到有效文本内容（可能是扫描件/图片型 PDF）")
	}

	log.Printf("[Knowledge] PDF 解析成功: %d 页, 提取 %d 字符", totalPages, len(content))
	return content, nil
}

// parseDocxFile 使用 nguyenthenguyen/docx 库解析 DOCX 文件提取文本
func parseDocxFile(filePath string) (string, error) {
	r, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return "", fmt.Errorf("打开 DOCX 文件失败: %w", err)
	}
	defer r.Close()

	doc := r.Editable()
	content := doc.GetContent()

	// GetContent 返回的是 XML 格式，需要提取纯文本
	// 移除 XML 标签，保留文本内容
	content = stripXmlTags(content)

	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("DOCX 文件未提取到有效文本内容")
	}

	log.Printf("[Knowledge] DOCX 解析成功: 提取 %d 字符", len(content))
	return content, nil
}

// parseDocLegacyFile 尝试从旧版 .doc 二进制文件中提取可读文本
func parseDocLegacyFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取 .doc 文件失败: %w", err)
	}

	// 从二进制数据中提取可打印的 UTF-8/ASCII 文本片段
	var sb strings.Builder
	var currentWord strings.Builder

	for _, b := range data {
		// 只保留可打印 ASCII 字符和常见空白字符
		if (b >= 0x20 && b <= 0x7E) || b == '\n' || b == '\r' || b == '\t' {
			currentWord.WriteByte(b)
		} else {
			if currentWord.Len() >= 4 { // 至少4个连续可读字符才保留，过滤噪声
				sb.WriteString(currentWord.String())
				sb.WriteString(" ")
			}
			currentWord.Reset()
		}
	}
	// 处理最后一段
	if currentWord.Len() >= 4 {
		sb.WriteString(currentWord.String())
	}

	content := sb.String()
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf(".doc 文件未提取到有效文本，建议转换为 .docx 格式后重新上传")
	}

	log.Printf("[Knowledge] .doc 文件有损提取完成: 提取 %d 字符（建议使用 .docx 格式以获得更好效果）", len(content))
	return content, nil
}

// parseCsvFile 解析 CSV 文件，将每行转为可读文本
func parseCsvFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开 CSV 文件失败: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true // 兼容不规范的 CSV

	var sb strings.Builder
	lineNum := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("[Knowledge] CSV 第 %d 行解析警告: %v", lineNum+1, err)
			continue
		}
		sb.WriteString(strings.Join(record, " | "))
		sb.WriteString("\n")
		lineNum++
	}

	log.Printf("[Knowledge] CSV 解析成功: %d 行", lineNum)
	return sb.String(), nil
}

// parseHtmlFile 解析 HTML 文件，去除标签保留纯文本
func parseHtmlFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取 HTML 文件失败: %w", err)
	}

	content := stripHtmlTags(string(data))

	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("HTML 文件未提取到有效文本内容")
	}

	log.Printf("[Knowledge] HTML 解析成功: 提取 %d 字符", len(content))
	return content, nil
}

// stripXmlTags 移除 XML/DOCX 标签，提取纯文本
func stripXmlTags(s string) string {
	var sb strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			sb.WriteRune(' ') // 标签边界用空格分隔
			continue
		}
		if !inTag {
			sb.WriteRune(r)
		}
	}
	// 合并多余空白
	result := sb.String()
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}
	return strings.TrimSpace(result)
}

// stripHtmlTags 移除 HTML 标签，保留纯文本，处理常见 HTML 实体
func stripHtmlTags(s string) string {
	// 先处理常见 HTML 实体
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")

	// 将块级标签替换为换行
	blockTags := []string{"</p>", "</div>", "</li>", "</tr>", "<br>", "<br/>", "<br />"}
	for _, tag := range blockTags {
		s = strings.ReplaceAll(s, tag, "\n")
	}

	// 移除所有 HTML 标签
	return stripXmlTags(s)
}
