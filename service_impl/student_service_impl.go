package service_impl

import (
	"context"
	"encoding/json"

	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/student/req"
	"github.com/yzf120/elysia-backend/model/student/rsp"
	"github.com/yzf120/elysia-backend/service"
)

// StudentServiceImpl 学生服务实现（只做出入参处理）
type StudentServiceImpl struct {
	studentService *service.StudentService
}

// NewStudentServiceImpl 创建学生服务实现
func NewStudentServiceImpl() *StudentServiceImpl {
	return &StudentServiceImpl{
		studentService: service.NewStudentService(),
	}
}

// CreateStudent 创建学生信息（注册时补充信息）
func (s *StudentServiceImpl) CreateStudent(ctx context.Context, req *req.CreateStudentRequest) (*rsp.CreateStudentResponse, error) {
	// 调用service层处理业务逻辑
	student, err := s.studentService.CreateStudent(
		req.StudentId,
		req.Major,
		req.Grade,
		req.ProgrammingLevel,
		req.Interests,
		req.LearningTags,
	)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.CreateStudentResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &rsp.CreateStudentResponse{
		Code:      consts.SuccessCode,
		Message:   "创建学生信息成功",
		StudentId: student.StudentId,
	}, nil
}

// GetStudent 获取学生信息
func (s *StudentServiceImpl) GetStudent(ctx context.Context, req *req.GetStudentRequest) (*rsp.GetStudentResponse, error) {
	// 调用service层处理业务逻辑
	student, err := s.studentService.GetStudentByStudentId(req.StudentId)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.GetStudentResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	// 解析JSON字段
	var interests []string
	var learningTags []string
	if student.Interests != "" {
		json.Unmarshal([]byte(student.Interests), &interests)
	}
	if student.LearningTags != "" {
		json.Unmarshal([]byte(student.LearningTags), &learningTags)
	}

	studentInfo := &rsp.StudentInfo{
		StudentId:        student.StudentId,
		StudentNumber:    student.StudentNumber,
		Major:            student.Major,
		Grade:            student.Grade,
		ProgrammingLevel: student.ProgrammingLevel,
		Interests:        interests,
		LearningTags:     learningTags,
		LearningProgress: student.LearningProgress,
		Status:           student.Status,
		CreateTime:       student.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:       student.UpdateTime.Format("2006-01-02 15:04:05"),
	}

	return &rsp.GetStudentResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageQuerySuccess,
		Student: studentInfo,
	}, nil
}

// UpdateStudent 更新学生信息
func (s *StudentServiceImpl) UpdateStudent(ctx context.Context, req *req.UpdateStudentRequest) (*rsp.UpdateStudentResponse, error) {
	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Major != "" {
		updates["major"] = req.Major
	}
	if req.Grade != "" {
		updates["grade"] = req.Grade
	}
	if req.ProgrammingLevel != "" {
		updates["programming_level"] = req.ProgrammingLevel
	}
	if len(req.Interests) > 0 {
		interestsJSON, _ := json.Marshal(req.Interests)
		updates["interests"] = string(interestsJSON)
	}
	if len(req.LearningTags) > 0 {
		learningTagsJSON, _ := json.Marshal(req.LearningTags)
		updates["learning_tags"] = string(learningTagsJSON)
	}

	// 调用service层处理业务逻辑
	_, err := s.studentService.UpdateStudent(req.StudentId, updates)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.UpdateStudentResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &rsp.UpdateStudentResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageUpdateSuccess,
	}, nil
}

// UpdateLearningProgress 更新学习进度
func (s *StudentServiceImpl) UpdateLearningProgress(ctx context.Context, req *req.UpdateLearningProgressRequest) (*rsp.UpdateLearningProgressResponse, error) {
	// 调用service层处理业务逻辑
	err := s.studentService.UpdateLearningProgress(req.StudentId, req.Progress)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.UpdateLearningProgressResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &rsp.UpdateLearningProgressResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageUpdateSuccess,
	}, nil
}

// ListStudents 查询学生列表
func (s *StudentServiceImpl) ListStudents(ctx context.Context, req *req.ListStudentsRequest) (*rsp.ListStudentsResponse, error) {
	// 构建筛选条件
	filters := make(map[string]interface{})
	if req.Major != "" {
		filters["major"] = req.Major
	}
	if req.Grade != "" {
		filters["grade"] = req.Grade
	}
	if req.ProgrammingLevel != "" {
		filters["programming_level"] = req.ProgrammingLevel
	}

	// 调用service层处理业务逻辑
	students, total, err := s.studentService.ListStudents(req.Page, req.PageSize, filters)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.ListStudentsResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	// 转换为响应格式
	studentInfos := make([]*rsp.StudentInfo, 0, len(students))
	for _, student := range students {
		var interests []string
		var learningTags []string
		if student.Interests != "" {
			json.Unmarshal([]byte(student.Interests), &interests)
		}
		if student.LearningTags != "" {
			json.Unmarshal([]byte(student.LearningTags), &learningTags)
		}

		studentInfos = append(studentInfos, &rsp.StudentInfo{
			StudentId:        student.StudentId,
			StudentNumber:    student.StudentNumber,
			Major:            student.Major,
			Grade:            student.Grade,
			ProgrammingLevel: student.ProgrammingLevel,
			Interests:        interests,
			LearningTags:     learningTags,
			LearningProgress: student.LearningProgress,
			Status:           student.Status,
			CreateTime:       student.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:       student.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &rsp.ListStudentsResponse{
		Code:     consts.SuccessCode,
		Message:  consts.MessageQuerySuccess,
		Students: studentInfos,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
