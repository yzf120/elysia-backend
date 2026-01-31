package service_impl

import (
	"context"
	"encoding/json"

	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/teacher/req"
	"github.com/yzf120/elysia-backend/model/teacher/rsp"
	"github.com/yzf120/elysia-backend/service"
)

// TeacherServiceImpl 教师服务实现（只做出入参处理）
type TeacherServiceImpl struct {
	teacherService *service.TeacherService
}

// NewTeacherServiceImpl 创建教师服务实现
func NewTeacherServiceImpl() *TeacherServiceImpl {
	return &TeacherServiceImpl{
		teacherService: service.NewTeacherService(),
	}
}

// RegisterTeacher 教师注册（工号+学校邮箱双重验证）
func (s *TeacherServiceImpl) RegisterTeacher(ctx context.Context, req *req.RegisterTeacherRequest) (*rsp.RegisterTeacherResponse, error) {
	// 调用service层处理业务逻辑
	teacher, err := s.teacherService.RegisterTeacher(
		ctx,
		req.PhoneNumber,
		req.Password,
		req.EmployeeNumber,
		req.SchoolEmail,
		req.RealName,
		req.Department,
		req.TeachingSubjects,
	)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.RegisterTeacherResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &rsp.RegisterTeacherResponse{
		Code:      consts.SuccessCode,
		Message:   "教师注册成功，等待管理员审核",
		TeacherId: teacher.TeacherId,
	}, nil
}

// LoginTeacher 教师登录
func (s *TeacherServiceImpl) LoginTeacher(ctx context.Context, req *req.LoginTeacherRequest) (*rsp.LoginTeacherResponse, error) {
	// 调用service层处理业务逻辑
	teacher, user, token, err := s.teacherService.LoginTeacher(ctx, req.PhoneNumber, req.Password)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.LoginTeacherResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	// 解析授课科目JSON
	var teachingSubjects []string
	if teacher.TeachingSubjects != "" {
		json.Unmarshal([]byte(teacher.TeachingSubjects), &teachingSubjects)
	}

	teacherInfo := &rsp.TeacherInfo{
		TeacherId:          teacher.TeacherId,
		UserId:             teacher.UserId,
		EmployeeNumber:     teacher.EmployeeNumber,
		SchoolEmail:        teacher.SchoolEmail,
		TeachingSubjects:   teachingSubjects,
		TeachingYears:      teacher.TeachingYears,
		Department:         teacher.Department,
		Title:              teacher.Title,
		VerificationStatus: teacher.VerificationStatus,
		Status:             teacher.Status,
		CreateTime:         teacher.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:         teacher.UpdateTime.Format("2006-01-02 15:04:05"),
	}

	userInfo := &rsp.UserInfo{
		UserId:      user.UserId,
		UserName:    user.UserName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		ChineseName: user.ChineseName,
		UserType:    user.UserType,
		Status:      user.Status,
	}

	return &rsp.LoginTeacherResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageLoginSuccess,
		Teacher: teacherInfo,
		User:    userInfo,
		Token:   token,
	}, nil
}

// VerifyTeacher 审核教师（管理员操作）
func (s *TeacherServiceImpl) VerifyTeacher(ctx context.Context, req *req.VerifyTeacherRequest) (*rsp.VerifyTeacherResponse, error) {
	// 调用service层处理业务逻辑
	err := s.teacherService.VerifyTeacher(req.TeacherId, req.VerifierId, req.Approved, req.Remark)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.VerifyTeacherResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	message := "教师审核通过"
	if !req.Approved {
		message = "教师审核驳回"
	}

	return &rsp.VerifyTeacherResponse{
		Code:    consts.SuccessCode,
		Message: message,
	}, nil
}

// GetTeacher 获取教师信息
func (s *TeacherServiceImpl) GetTeacher(ctx context.Context, req *req.GetTeacherRequest) (*rsp.GetTeacherResponse, error) {
	// 调用service层处理业务逻辑
	teacher, err := s.teacherService.GetTeacherByUserId(req.UserId)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.GetTeacherResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	// 解析授课科目JSON
	var teachingSubjects []string
	if teacher.TeachingSubjects != "" {
		json.Unmarshal([]byte(teacher.TeachingSubjects), &teachingSubjects)
	}

	teacherInfo := &rsp.TeacherInfo{
		TeacherId:          teacher.TeacherId,
		UserId:             teacher.UserId,
		EmployeeNumber:     teacher.EmployeeNumber,
		SchoolEmail:        teacher.SchoolEmail,
		TeachingSubjects:   teachingSubjects,
		TeachingYears:      teacher.TeachingYears,
		Department:         teacher.Department,
		Title:              teacher.Title,
		VerificationStatus: teacher.VerificationStatus,
		Status:             teacher.Status,
		CreateTime:         teacher.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:         teacher.UpdateTime.Format("2006-01-02 15:04:05"),
	}

	return &rsp.GetTeacherResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageQuerySuccess,
		Teacher: teacherInfo,
	}, nil
}

// UpdateTeacher 更新教师信息
func (s *TeacherServiceImpl) UpdateTeacher(ctx context.Context, req *req.UpdateTeacherRequest) (*rsp.UpdateTeacherResponse, error) {
	// 构建更新字段
	updates := make(map[string]interface{})
	if len(req.TeachingSubjects) > 0 {
		teachingSubjectsJSON, _ := json.Marshal(req.TeachingSubjects)
		updates["teaching_subjects"] = string(teachingSubjectsJSON)
	}
	if req.TeachingYears > 0 {
		updates["teaching_years"] = req.TeachingYears
	}
	if req.Department != "" {
		updates["department"] = req.Department
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}

	// 调用service层处理业务逻辑
	_, err := s.teacherService.UpdateTeacher(req.TeacherId, updates)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.UpdateTeacherResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &rsp.UpdateTeacherResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageUpdateSuccess,
	}, nil
}

// ListTeachers 查询教师列表
func (s *TeacherServiceImpl) ListTeachers(ctx context.Context, req *req.ListTeachersRequest) (*rsp.ListTeachersResponse, error) {
	// 构建筛选条件
	filters := make(map[string]interface{})
	if req.Department != "" {
		filters["department"] = req.Department
	}
	if req.VerificationStatus >= 0 {
		filters["verification_status"] = req.VerificationStatus
	}
	if req.Status >= 0 {
		filters["status"] = req.Status
	}

	// 调用service层处理业务逻辑
	teachers, total, err := s.teacherService.ListTeachers(req.Page, req.PageSize, filters)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.ListTeachersResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	// 转换为响应格式
	teacherInfos := make([]*rsp.TeacherInfo, 0, len(teachers))
	for _, teacher := range teachers {
		var teachingSubjects []string
		if teacher.TeachingSubjects != "" {
			json.Unmarshal([]byte(teacher.TeachingSubjects), &teachingSubjects)
		}

		teacherInfos = append(teacherInfos, &rsp.TeacherInfo{
			TeacherId:          teacher.TeacherId,
			UserId:             teacher.UserId,
			EmployeeNumber:     teacher.EmployeeNumber,
			SchoolEmail:        teacher.SchoolEmail,
			TeachingSubjects:   teachingSubjects,
			TeachingYears:      teacher.TeachingYears,
			Department:         teacher.Department,
			Title:              teacher.Title,
			VerificationStatus: teacher.VerificationStatus,
			Status:             teacher.Status,
			CreateTime:         teacher.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:         teacher.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &rsp.ListTeachersResponse{
		Code:     consts.SuccessCode,
		Message:  consts.MessageQuerySuccess,
		Teachers: teacherInfos,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
