package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	classModel "github.com/yzf120/elysia-backend/model/class"
)

// ChapterService 章节服务
type ChapterService struct {
	chapterDAO dao.ChapterDAO
	classDAO   dao.ClassDAO
}

// NewChapterService 创建章节服务
func NewChapterService() *ChapterService {
	return &ChapterService{
		chapterDAO: dao.NewChapterDAO(),
		classDAO:   dao.NewClassDAO(),
	}
}

// ==================== 章节操作 ====================

// CreateChapter 创建章节（教师操作）
func (s *ChapterService) CreateChapter(teacherId, classId, title, description string) (*classModel.ClassChapter, error) {
	if teacherId == "" || classId == "" || title == "" {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "必填参数不能为空")
	}

	// 校验班级归属
	class, err := s.classDAO.GetClassById(classId)
	if err != nil || class == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "班级不存在")
	}
	if class.TeacherId != teacherId {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "无权限操作该班级")
	}

	// 计算新章节的排序值（当前最大值 + 10，方便后续插入）
	chapters, _ := s.chapterDAO.ListChaptersByClassId(classId)
	sortOrder := int32(10)
	if len(chapters) > 0 {
		sortOrder = chapters[len(chapters)-1].SortOrder + 10
	}

	chapter := &classModel.ClassChapter{
		ChapterId:   fmt.Sprintf("chap_%d", time.Now().UnixNano()),
		ClassId:     classId,
		Title:       title,
		Description: description,
		SortOrder:   sortOrder,
		Status:      1,
	}

	if err := s.chapterDAO.CreateChapter(chapter); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建章节失败: "+err.Error())
	}

	// 同步更新 class.chapter_ids
	s.syncChapterIds(classId)

	return chapter, nil
}

// UpdateChapter 更新章节（教师操作）
func (s *ChapterService) UpdateChapter(teacherId, chapterId, title, description string) error {
	chapter, err := s.chapterDAO.GetChapterById(chapterId)
	if err != nil || chapter == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "章节不存在")
	}

	// 校验班级归属
	class, err := s.classDAO.GetClassById(chapter.ClassId)
	if err != nil || class == nil || class.TeacherId != teacherId {
		return errs.NewCommonError(errs.ErrBadRequest, "无权限操作")
	}

	updates := map[string]interface{}{}
	if title != "" {
		updates["title"] = title
	}
	if description != "" {
		updates["description"] = description
	}
	if len(updates) == 0 {
		return nil
	}

	if err := s.chapterDAO.UpdateChapter(chapterId, updates); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "更新章节失败: "+err.Error())
	}
	return nil
}

// DeleteChapter 删除章节（教师操作，同时删除其下所有小节）
func (s *ChapterService) DeleteChapter(teacherId, chapterId string) error {
	chapter, err := s.chapterDAO.GetChapterById(chapterId)
	if err != nil || chapter == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "章节不存在")
	}

	// 校验班级归属
	class, err := s.classDAO.GetClassById(chapter.ClassId)
	if err != nil || class == nil || class.TeacherId != teacherId {
		return errs.NewCommonError(errs.ErrBadRequest, "无权限操作")
	}

	// 删除章节下所有小节
	if err := s.chapterDAO.DeleteSectionsByChapterId(chapterId); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "删除小节失败: "+err.Error())
	}

	// 删除章节
	if err := s.chapterDAO.DeleteChapter(chapterId); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "删除章节失败: "+err.Error())
	}

	// 同步更新 class.chapter_ids
	s.syncChapterIds(chapter.ClassId)

	return nil
}

// ReorderChapters 调整章节排序（教师操作）
// orders: [{chapter_id, sort_order}, ...]
func (s *ChapterService) ReorderChapters(teacherId, classId string, orders []dao.ChapterOrder) error {
	// 校验班级归属
	class, err := s.classDAO.GetClassById(classId)
	if err != nil || class == nil || class.TeacherId != teacherId {
		return errs.NewCommonError(errs.ErrBadRequest, "无权限操作")
	}

	if err := s.chapterDAO.BatchUpdateChapterOrder(classId, orders); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "更新章节排序失败: "+err.Error())
	}

	// 同步更新 class.chapter_ids（按新排序）
	s.syncChapterIds(classId)

	return nil
}

// GetChaptersByClassId 查询班级下所有章节（含小节，师生共用）
func (s *ChapterService) GetChaptersByClassId(classId string) ([]*classModel.ClassChapter, map[string][]*classModel.ClassSection, error) {
	chapters, err := s.chapterDAO.ListChaptersByClassId(classId)
	if err != nil {
		return nil, nil, errs.NewCommonError(errs.ErrInternal, "查询章节失败: "+err.Error())
	}

	// 查询每个章节的小节
	sectionMap := make(map[string][]*classModel.ClassSection)
	for _, ch := range chapters {
		sections, err := s.chapterDAO.ListSectionsByChapterId(ch.ChapterId)
		if err != nil {
			continue
		}
		sectionMap[ch.ChapterId] = sections
	}

	return chapters, sectionMap, nil
}

// ==================== 小节操作 ====================

// CreateSection 创建小节（教师操作）
func (s *ChapterService) CreateSection(teacherId, chapterId, title, description string, sectionType int32,
	problemId string,
	discussionTitle, discussionContent string) (*classModel.ClassSection, error) {
	if teacherId == "" || chapterId == "" || title == "" {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "必填参数不能为空")
	}
	if sectionType != classModel.SectionTypeProblem && sectionType != classModel.SectionTypeDiscussion {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "小节类型不合法（1-算法题，2-讨论话题）")
	}

	chapter, err := s.chapterDAO.GetChapterById(chapterId)
	if err != nil || chapter == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "章节不存在")
	}

	// 校验班级归属
	class, err := s.classDAO.GetClassById(chapter.ClassId)
	if err != nil || class == nil || class.TeacherId != teacherId {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "无权限操作")
	}

	// 计算排序值
	sections, _ := s.chapterDAO.ListSectionsByChapterId(chapterId)
	sortOrder := int32(10)
	if len(sections) > 0 {
		sortOrder = sections[len(sections)-1].SortOrder + 10
	}

	section := &classModel.ClassSection{
		SectionId:         fmt.Sprintf("sec_%d", time.Now().UnixNano()),
		ChapterId:         chapterId,
		ClassId:           chapter.ClassId,
		Title:             title,
		Description:       description,
		SectionType:       sectionType,
		ProblemId:         problemId,
		DiscussionTitle:   discussionTitle,
		DiscussionContent: discussionContent,
		SortOrder:         sortOrder,
		Status:            1,
	}

	if err := s.chapterDAO.CreateSection(section); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建小节失败: "+err.Error())
	}

	return section, nil
}

// UpdateSection 更新小节（教师操作）
func (s *ChapterService) UpdateSection(teacherId, sectionId, title, description string,
	problemId string,
	discussionTitle, discussionContent string) error {
	section, err := s.chapterDAO.GetSectionById(sectionId)
	if err != nil || section == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "小节不存在")
	}

	// 校验班级归属
	class, err := s.classDAO.GetClassById(section.ClassId)
	if err != nil || class == nil || class.TeacherId != teacherId {
		return errs.NewCommonError(errs.ErrBadRequest, "无权限操作")
	}

	updates := map[string]interface{}{}
	if title != "" {
		updates["title"] = title
	}
	if description != "" {
		updates["description"] = description
	}
	if problemId != "" {
		updates["problem_id"] = problemId
	}
	if discussionTitle != "" {
		updates["discussion_title"] = discussionTitle
	}
	if discussionContent != "" {
		updates["discussion_content"] = discussionContent
	}
	if len(updates) == 0 {
		return nil
	}

	if err := s.chapterDAO.UpdateSection(sectionId, updates); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "更新小节失败: "+err.Error())
	}
	return nil
}

// DeleteSection 删除小节（教师操作）
func (s *ChapterService) DeleteSection(teacherId, sectionId string) error {
	section, err := s.chapterDAO.GetSectionById(sectionId)
	if err != nil || section == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "小节不存在")
	}

	// 校验班级归属
	class, err := s.classDAO.GetClassById(section.ClassId)
	if err != nil || class == nil || class.TeacherId != teacherId {
		return errs.NewCommonError(errs.ErrBadRequest, "无权限操作")
	}

	if err := s.chapterDAO.DeleteSection(sectionId); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "删除小节失败: "+err.Error())
	}
	return nil
}

// ReorderSections 调整小节排序（教师操作）
func (s *ChapterService) ReorderSections(teacherId, chapterId string, orders []dao.SectionOrder) error {
	chapter, err := s.chapterDAO.GetChapterById(chapterId)
	if err != nil || chapter == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "章节不存在")
	}

	// 校验班级归属
	class, err := s.classDAO.GetClassById(chapter.ClassId)
	if err != nil || class == nil || class.TeacherId != teacherId {
		return errs.NewCommonError(errs.ErrBadRequest, "无权限操作")
	}

	if err := s.chapterDAO.BatchUpdateSectionOrder(chapterId, orders); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "更新小节排序失败: "+err.Error())
	}
	return nil
}

// ==================== 内部工具 ====================

// syncChapterIds 同步更新 class.chapter_ids（按 sort_order 排序后的 chapter_id 列表）
func (s *ChapterService) syncChapterIds(classId string) {
	chapters, err := s.chapterDAO.ListChaptersByClassId(classId)
	if err != nil {
		return
	}
	ids := make([]string, 0, len(chapters))
	for _, ch := range chapters {
		ids = append(ids, ch.ChapterId)
	}
	idsJSON, _ := json.Marshal(ids)
	s.classDAO.UpdateClass(classId, map[string]interface{}{
		"chapter_ids": string(idsJSON),
	})
}
