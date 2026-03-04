package service_impl

import (
	"context"

	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/service"
)

// ChapterServiceImpl 章节服务实现（只做出入参处理）
type ChapterServiceImpl struct {
	chapterService *service.ChapterService
}

// NewChapterServiceImpl 创建章节服务实现
func NewChapterServiceImpl() *ChapterServiceImpl {
	return &ChapterServiceImpl{
		chapterService: service.NewChapterService(),
	}
}

// ==================== 章节接口 ====================

// CreateChapterRequest 创建章节请求
type CreateChapterRequest struct {
	TeacherId   string `json:"teacher_id"`  // 教师ID（必填）
	ClassId     string `json:"class_id"`    // 班级ID（必填）
	Title       string `json:"title"`       // 章节标题（必填）
	Description string `json:"description"` // 章节描述（可选）
}

// CreateChapterResponse 创建章节响应
type CreateChapterResponse struct {
	Code      int32  `json:"code"`
	Message   string `json:"message"`
	ChapterId string `json:"chapter_id"`
}

// CreateChapter 创建章节
func (s *ChapterServiceImpl) CreateChapter(ctx context.Context, req *CreateChapterRequest) (*CreateChapterResponse, error) {
	chapter, err := s.chapterService.CreateChapter(req.TeacherId, req.ClassId, req.Title, req.Description)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &CreateChapterResponse{Code: int32(code), Message: msg}, nil
	}
	return &CreateChapterResponse{
		Code:      consts.SuccessCode,
		Message:   "创建章节成功",
		ChapterId: chapter.ChapterId,
	}, nil
}

// UpdateChapterRequest 更新章节请求
type UpdateChapterRequest struct {
	TeacherId   string `json:"teacher_id"`  // 教师ID（必填）
	ChapterId   string `json:"chapter_id"`  // 章节ID（必填）
	Title       string `json:"title"`       // 章节标题（可选）
	Description string `json:"description"` // 章节描述（可选）
}

// UpdateChapterResponse 更新章节响应
type UpdateChapterResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// UpdateChapter 更新章节
func (s *ChapterServiceImpl) UpdateChapter(ctx context.Context, req *UpdateChapterRequest) (*UpdateChapterResponse, error) {
	if err := s.chapterService.UpdateChapter(req.TeacherId, req.ChapterId, req.Title, req.Description); err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &UpdateChapterResponse{Code: int32(code), Message: msg}, nil
	}
	return &UpdateChapterResponse{Code: consts.SuccessCode, Message: "更新章节成功"}, nil
}

// DeleteChapterRequest 删除章节请求
type DeleteChapterRequest struct {
	TeacherId string `json:"teacher_id"` // 教师ID（必填）
	ChapterId string `json:"chapter_id"` // 章节ID（必填）
}

// DeleteChapterResponse 删除章节响应
type DeleteChapterResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// DeleteChapter 删除章节
func (s *ChapterServiceImpl) DeleteChapter(ctx context.Context, req *DeleteChapterRequest) (*DeleteChapterResponse, error) {
	if err := s.chapterService.DeleteChapter(req.TeacherId, req.ChapterId); err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &DeleteChapterResponse{Code: int32(code), Message: msg}, nil
	}
	return &DeleteChapterResponse{Code: consts.SuccessCode, Message: "删除章节成功"}, nil
}

// ReorderChaptersRequest 调整章节排序请求
type ReorderChaptersRequest struct {
	TeacherId string              `json:"teacher_id"` // 教师ID（必填）
	ClassId   string              `json:"class_id"`   // 班级ID（必填）
	Orders    []ChapterOrderItem  `json:"orders"`     // 排序列表
}

// ChapterOrderItem 章节排序项
type ChapterOrderItem struct {
	ChapterId string `json:"chapter_id"` // 章节ID
	SortOrder int32  `json:"sort_order"` // 排序值
}

// ReorderChaptersResponse 调整章节排序响应
type ReorderChaptersResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// ReorderChapters 调整章节排序
func (s *ChapterServiceImpl) ReorderChapters(ctx context.Context, req *ReorderChaptersRequest) (*ReorderChaptersResponse, error) {
	orders := make([]dao.ChapterOrder, 0, len(req.Orders))
	for _, o := range req.Orders {
		orders = append(orders, dao.ChapterOrder{ChapterId: o.ChapterId, SortOrder: o.SortOrder})
	}
	if err := s.chapterService.ReorderChapters(req.TeacherId, req.ClassId, orders); err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &ReorderChaptersResponse{Code: int32(code), Message: msg}, nil
	}
	return &ReorderChaptersResponse{Code: consts.SuccessCode, Message: "排序更新成功"}, nil
}

// ==================== 小节接口 ====================

// CreateSectionRequest 创建小节请求
type CreateSectionRequest struct {
	TeacherId   string `json:"teacher_id"`   // 教师ID（必填）
	ChapterId   string `json:"chapter_id"`   // 章节ID（必填）
	Title       string `json:"title"`        // 小节标题（必填）
	Description string `json:"description"`  // 小节描述（可选）
	SectionType int32  `json:"section_type"` // 小节类型：1-算法题，2-讨论话题（必填）
	// 算法题关联字段（section_type=1 时填写，关联题库中的题目ID）
	ProblemId string `json:"problem_id"` // 题库中的题目ID
	// 讨论内容字段（section_type=2 时填写）
	DiscussionTitle   string `json:"discussion_title"`   // 讨论话题标题
	DiscussionContent string `json:"discussion_content"` // 讨论话题描述
}

// CreateSectionResponse 创建小节响应
type CreateSectionResponse struct {
	Code      int32  `json:"code"`
	Message   string `json:"message"`
	SectionId string `json:"section_id"`
}

// CreateSection 创建小节
func (s *ChapterServiceImpl) CreateSection(ctx context.Context, req *CreateSectionRequest) (*CreateSectionResponse, error) {
	section, err := s.chapterService.CreateSection(req.TeacherId, req.ChapterId, req.Title, req.Description, req.SectionType,
		req.ProblemId,
		req.DiscussionTitle, req.DiscussionContent)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &CreateSectionResponse{Code: int32(code), Message: msg}, nil
	}
	return &CreateSectionResponse{
		Code:      consts.SuccessCode,
		Message:   "创建小节成功",
		SectionId: section.SectionId,
	}, nil
}

// UpdateSectionRequest 更新小节请求
type UpdateSectionRequest struct {
	TeacherId   string `json:"teacher_id"`  // 教师ID（必填）
	SectionId   string `json:"section_id"`  // 小节ID（必填）
	Title       string `json:"title"`       // 小节标题（可选）
	Description string `json:"description"` // 小节描述（可选）
	// 算法题关联字段
	ProblemId string `json:"problem_id"` // 题库中的题目ID
	// 讨论内容字段
	DiscussionTitle   string `json:"discussion_title"`
	DiscussionContent string `json:"discussion_content"`
}

// UpdateSectionResponse 更新小节响应
type UpdateSectionResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// UpdateSection 更新小节
func (s *ChapterServiceImpl) UpdateSection(ctx context.Context, req *UpdateSectionRequest) (*UpdateSectionResponse, error) {
	if err := s.chapterService.UpdateSection(req.TeacherId, req.SectionId, req.Title, req.Description,
		req.ProblemId,
		req.DiscussionTitle, req.DiscussionContent); err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &UpdateSectionResponse{Code: int32(code), Message: msg}, nil
	}
	return &UpdateSectionResponse{Code: consts.SuccessCode, Message: "更新小节成功"}, nil
}

// DeleteSectionRequest 删除小节请求
type DeleteSectionRequest struct {
	TeacherId string `json:"teacher_id"` // 教师ID（必填）
	SectionId string `json:"section_id"` // 小节ID（必填）
}

// DeleteSectionResponse 删除小节响应
type DeleteSectionResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// DeleteSection 删除小节
func (s *ChapterServiceImpl) DeleteSection(ctx context.Context, req *DeleteSectionRequest) (*DeleteSectionResponse, error) {
	if err := s.chapterService.DeleteSection(req.TeacherId, req.SectionId); err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &DeleteSectionResponse{Code: int32(code), Message: msg}, nil
	}
	return &DeleteSectionResponse{Code: consts.SuccessCode, Message: "删除小节成功"}, nil
}

// ReorderSectionsRequest 调整小节排序请求
type ReorderSectionsRequest struct {
	TeacherId string             `json:"teacher_id"` // 教师ID（必填）
	ChapterId string             `json:"chapter_id"` // 章节ID（必填）
	Orders    []SectionOrderItem `json:"orders"`     // 排序列表
}

// SectionOrderItem 小节排序项
type SectionOrderItem struct {
	SectionId string `json:"section_id"` // 小节ID
	SortOrder int32  `json:"sort_order"` // 排序值
}

// ReorderSectionsResponse 调整小节排序响应
type ReorderSectionsResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// ReorderSections 调整小节排序
func (s *ChapterServiceImpl) ReorderSections(ctx context.Context, req *ReorderSectionsRequest) (*ReorderSectionsResponse, error) {
	orders := make([]dao.SectionOrder, 0, len(req.Orders))
	for _, o := range req.Orders {
		orders = append(orders, dao.SectionOrder{SectionId: o.SectionId, SortOrder: o.SortOrder})
	}
	if err := s.chapterService.ReorderSections(req.TeacherId, req.ChapterId, orders); err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &ReorderSectionsResponse{Code: int32(code), Message: msg}, nil
	}
	return &ReorderSectionsResponse{Code: consts.SuccessCode, Message: "排序更新成功"}, nil
}

// ==================== 查询接口（师生共用） ====================

// SectionInfo 小节信息
type SectionInfo struct {
	SectionId   string `json:"section_id"`
	ChapterId   string `json:"chapter_id"`
	ClassId     string `json:"class_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SectionType int32  `json:"section_type"` // 1-算法题，2-讨论话题
	// 算法题关联
	ProblemId string `json:"problem_id"`
	// 讨论内容
	DiscussionTitle   string `json:"discussion_title"`
	DiscussionContent string `json:"discussion_content"`
	SortOrder   int32  `json:"sort_order"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
}

// ChapterInfo 章节信息（含小节列表）
type ChapterInfo struct {
	ChapterId   string         `json:"chapter_id"`
	ClassId     string         `json:"class_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	SortOrder   int32          `json:"sort_order"`
	Sections    []*SectionInfo `json:"sections"` // 小节列表（按 sort_order 升序）
	CreateTime  string         `json:"create_time"`
	UpdateTime  string         `json:"update_time"`
}

// GetClassChaptersRequest 查询班级章节列表请求
type GetClassChaptersRequest struct {
	ClassId string `json:"class_id"` // 班级ID（必填）
}

// GetClassChaptersResponse 查询班级章节列表响应
type GetClassChaptersResponse struct {
	Code     int32          `json:"code"`
	Message  string         `json:"message"`
	Chapters []*ChapterInfo `json:"chapters"` // 章节列表（按 sort_order 升序）
}

// GetClassChapters 查询班级章节列表（师生共用）
func (s *ChapterServiceImpl) GetClassChapters(ctx context.Context, req *GetClassChaptersRequest) (*GetClassChaptersResponse, error) {
	chapters, sectionMap, err := s.chapterService.GetChaptersByClassId(req.ClassId)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &GetClassChaptersResponse{Code: int32(code), Message: msg}, nil
	}

	chapterInfos := make([]*ChapterInfo, 0, len(chapters))
	for _, ch := range chapters {
		info := &ChapterInfo{
			ChapterId:   ch.ChapterId,
			ClassId:     ch.ClassId,
			Title:       ch.Title,
			Description: ch.Description,
			SortOrder:   ch.SortOrder,
			CreateTime:  ch.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:  ch.UpdateTime.Format("2006-01-02 15:04:05"),
			Sections:    make([]*SectionInfo, 0),
		}
		if sections, ok := sectionMap[ch.ChapterId]; ok {
			for _, sec := range sections {
				info.Sections = append(info.Sections, &SectionInfo{
					SectionId:         sec.SectionId,
					ChapterId:         sec.ChapterId,
					ClassId:           sec.ClassId,
					Title:             sec.Title,
					Description:       sec.Description,
					SectionType:       sec.SectionType,
					ProblemId:         sec.ProblemId,
					DiscussionTitle:   sec.DiscussionTitle,
					DiscussionContent: sec.DiscussionContent,
					SortOrder:         sec.SortOrder,
					CreateTime:        sec.CreateTime.Format("2006-01-02 15:04:05"),
					UpdateTime:        sec.UpdateTime.Format("2006-01-02 15:04:05"),
				})
			}
		}
		chapterInfos = append(chapterInfos, info)
	}

	return &GetClassChaptersResponse{
		Code:     consts.SuccessCode,
		Message:  consts.MessageQuerySuccess,
		Chapters: chapterInfos,
	}, nil
}