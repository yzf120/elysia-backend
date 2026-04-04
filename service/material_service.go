package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	classModel "github.com/yzf120/elysia-backend/model/class"
)

const (
	materialUploadDir   = "uploads/section-materials"
	maxMaterialFileSize = 200 * 1024 * 1024 // 200MB（视频可能较大）
)

// MaterialService 学习资料服务
type MaterialService struct {
	materialDAO dao.MaterialDAO
	chapterDAO  dao.ChapterDAO
	classDAO    dao.ClassDAO
}

// NewMaterialService 创建学习资料服务
func NewMaterialService() *MaterialService {
	return &MaterialService{
		materialDAO: dao.NewMaterialDAO(),
		chapterDAO:  dao.NewChapterDAO(),
		classDAO:    dao.NewClassDAO(),
	}
}

// CreateMaterial 创建学习资料（教师操作）
func (s *MaterialService) CreateMaterial(teacherId, sectionId, title, description, materialType string,
	fileHeader *multipart.FileHeader) (*classModel.SectionMaterial, error) {

	if teacherId == "" || sectionId == "" || title == "" || materialType == "" {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "必填参数不能为空")
	}

	// 校验资料类型
	validTypes := map[string]bool{
		classModel.MaterialTypePDF:   true,
		classModel.MaterialTypeWord:  true,
		classModel.MaterialTypeText:  true,
		classModel.MaterialTypeVideo: true,
	}
	if !validTypes[materialType] {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "资料类型不合法（支持：pdf, word, text, video）")
	}

	// 查询小节
	section, err := s.chapterDAO.GetSectionById(sectionId)
	if err != nil || section == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "小节不存在")
	}

	// 校验班级归属
	class, err := s.classDAO.GetClassById(section.ClassId)
	if err != nil || class == nil || class.TeacherId != teacherId {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "无权限操作")
	}

	material := &classModel.SectionMaterial{
		MaterialId:   fmt.Sprintf("mat_%d", time.Now().UnixNano()),
		SectionId:    sectionId,
		ChapterId:    section.ChapterId,
		ClassId:      section.ClassId,
		TeacherId:    teacherId,
		Title:        title,
		Description:  description,
		MaterialType: materialType,
		Status:       1,
	}

	// 计算排序值
	existingMaterials, _ := s.materialDAO.ListMaterialsBySectionId(sectionId)
	sortOrder := int32(10)
	if len(existingMaterials) > 0 {
		sortOrder = existingMaterials[len(existingMaterials)-1].SortOrder + 10
	}
	material.SortOrder = sortOrder

	// 处理文件上传（非text类型需要文件）
	if materialType != classModel.MaterialTypeText {
		if fileHeader == nil {
			return nil, errs.NewCommonError(errs.ErrBadRequest, "非文字类型资料必须上传文件")
		}
		if fileHeader.Size > maxMaterialFileSize {
			return nil, errs.NewCommonError(errs.ErrBadRequest, "文件大小超过限制（最大200MB）")
		}

		// 校验文件类型
		if err := s.validateFileType(materialType, fileHeader); err != nil {
			return nil, err
		}

		// 保存文件
		filePath, err := s.saveUploadedFile(fileHeader, section.ClassId, sectionId)
		if err != nil {
			return nil, errs.NewCommonError(errs.ErrInternal, "文件保存失败: "+err.Error())
		}

		material.FileName = fileHeader.Filename
		material.FilePath = filePath
		material.FileSize = fileHeader.Size
		material.MimeType = detectMimeType(fileHeader)
	}

	if err := s.materialDAO.CreateMaterial(material); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建学习资料失败: "+err.Error())
	}

	return material, nil
}

// GetMaterialsBySection 查询小节下所有学习资料
func (s *MaterialService) GetMaterialsBySection(sectionId string) ([]*classModel.SectionMaterial, error) {
	return s.materialDAO.ListMaterialsBySectionId(sectionId)
}

// GetMaterialFile 获取学习资料文件信息（用于下载/查看）
func (s *MaterialService) GetMaterialFile(materialId string) (*classModel.SectionMaterial, string, error) {
	material, err := s.materialDAO.GetMaterialById(materialId)
	if err != nil || material == nil {
		return nil, "", errs.NewCommonError(http.StatusNotFound, "资料不存在")
	}

	if material.FilePath == "" {
		return nil, "", errs.NewCommonError(http.StatusNotFound, "该资料没有附件")
	}

	absPath, err := filepath.Abs(material.FilePath)
	if err != nil || !fileExists(absPath) {
		return nil, "", errs.NewCommonError(http.StatusNotFound, "文件不存在")
	}

	return material, absPath, nil
}

// DeleteMaterial 删除单条学习资料（教师操作）
func (s *MaterialService) DeleteMaterial(teacherId, materialId string) error {
	material, err := s.materialDAO.GetMaterialById(materialId)
	if err != nil || material == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "资料不存在")
	}

	// 校验权限
	class, err := s.classDAO.GetClassById(material.ClassId)
	if err != nil || class == nil || class.TeacherId != teacherId {
		return errs.NewCommonError(errs.ErrBadRequest, "无权限操作")
	}

	// 删除文件
	if material.FilePath != "" {
		_ = os.Remove(material.FilePath)
	}

	if err := s.materialDAO.DeleteMaterial(materialId); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "删除学习资料失败: "+err.Error())
	}
	return nil
}

// DeleteMaterialsBySection 删除小节下所有学习资料（教师操作，删除小节时调用）
func (s *MaterialService) DeleteMaterialsBySection(sectionId string) error {
	materials, _ := s.materialDAO.ListMaterialsBySectionId(sectionId)
	for _, m := range materials {
		if m.FilePath != "" {
			_ = os.Remove(m.FilePath)
		}
	}
	return s.materialDAO.DeleteMaterialsBySectionId(sectionId)
}

// validateFileType 校验文件类型与资料类型是否匹配
func (s *MaterialService) validateFileType(materialType string, fh *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	switch materialType {
	case classModel.MaterialTypePDF:
		if ext != ".pdf" {
			return errs.NewCommonError(errs.ErrBadRequest, "PDF类型资料仅支持 .pdf 文件")
		}
	case classModel.MaterialTypeWord:
		if ext != ".doc" && ext != ".docx" {
			return errs.NewCommonError(errs.ErrBadRequest, "Word类型资料仅支持 .doc/.docx 文件")
		}
	case classModel.MaterialTypeVideo:
		allowed := map[string]bool{".mp4": true, ".avi": true, ".mkv": true, ".mov": true, ".webm": true}
		if !allowed[ext] {
			return errs.NewCommonError(errs.ErrBadRequest, "视频类型资料仅支持 mp4/avi/mkv/mov/webm 格式")
		}
	}
	return nil
}

// saveUploadedFile 保存上传文件到本地
func (s *MaterialService) saveUploadedFile(fh *multipart.FileHeader, classId, sectionId string) (string, error) {
	dir := filepath.Join(materialUploadDir, classId, sectionId)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	ext := filepath.Ext(fh.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join(dir, fileName)

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return destPath, nil
}

// detectMimeType 检测MIME类型
func detectMimeType(fh *multipart.FileHeader) string {
	ct := fh.Header.Get("Content-Type")
	if ct != "" && ct != "application/octet-stream" {
		return ct
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	mimeMap := map[string]string{
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
		".mkv":  "video/x-matroska",
		".mov":  "video/quicktime",
		".webm": "video/webm",
	}
	if mime, ok := mimeMap[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
