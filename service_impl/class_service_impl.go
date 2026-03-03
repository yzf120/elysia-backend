package service_impl

import (
	"context"

	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/service"
)

// ClassServiceImpl 班级服务实现（只做出入参处理）
type ClassServiceImpl struct {
	classService   *service.ClassService
	teacherService *service.TeacherService
	subjectService service.SubjectService
	studentDAO     dao.StudentDAO
}

// NewClassServiceImpl 创建班级服务实现
func NewClassServiceImpl() *ClassServiceImpl {
	return &ClassServiceImpl{
		classService:   service.NewClassService(),
		teacherService: service.NewTeacherService(),
		subjectService: service.NewSubjectService(),
		studentDAO:     dao.NewStudentDAO(),
	}
}

// CreateClassRequest 创建班级请求
type CreateClassRequest struct {
	TeacherId   string `json:"teacher_id"`   // 教师ID（必填）
	ClassName   string `json:"class_name"`   // 班级名称（必填）
	SubjectId   string `json:"subject_id"`   // 科目ID（必填）
	Semester    string `json:"semester"`     // 学期（必填）
	Description string `json:"description"`  // 班级描述（可选）
	MaxStudents int32  `json:"max_students"` // 学生人数上限（必填）
}

// CreateClassResponse 创建班级响应
type CreateClassResponse struct {
	Code      int32  `json:"code"`       // 响应码 0-成功 其他-失败
	Message   string `json:"message"`    // 响应消息
	ClassId   string `json:"class_id"`   // 班级ID
	ClassCode string `json:"class_code"` // 班级验证码
}

// CreateClass 创建班级（教师操作）
func (s *ClassServiceImpl) CreateClass(ctx context.Context, req *CreateClassRequest) (*CreateClassResponse, error) {
	// 调用service层处理业务逻辑
	class, err := s.classService.CreateClass(
		req.TeacherId,
		req.ClassName,
		req.SubjectId,
		req.Semester,
		req.Description,
		req.MaxStudents,
	)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &CreateClassResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &CreateClassResponse{
		Code:      consts.SuccessCode,
		Message:   "创建班级成功",
		ClassId:   class.ClassId,
		ClassCode: class.ClassCode,
	}, nil
}

// JoinClassRequest 学生加入班级请求
type JoinClassRequest struct {
	StudentId string `json:"student_id"` // 学生ID（必填）
	ClassCode string `json:"class_code"` // 班级验证码（必填）
}

// JoinClassResponse 学生加入班级响应
type JoinClassResponse struct {
	Code    int32  `json:"code"`    // 响应码 0-成功 其他-失败
	Message string `json:"message"` // 响应消息
}

// JoinClass 学生加入班级
func (s *ClassServiceImpl) JoinClass(ctx context.Context, req *JoinClassRequest) (*JoinClassResponse, error) {
	// 调用service层处理业务逻辑
	err := s.classService.JoinClass(req.StudentId, req.ClassCode)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &JoinClassResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &JoinClassResponse{
		Code:    consts.SuccessCode,
		Message: "加入班级成功",
	}, nil
}

// LeaveClassRequest 学生退出班级请求
type LeaveClassRequest struct {
	StudentId string `json:"student_id"` // 学生ID（必填）
	ClassId   string `json:"class_id"`   // 班级ID（必填）
}

// LeaveClassResponse 学生退出班级响应
type LeaveClassResponse struct {
	Code    int32  `json:"code"`    // 响应码 0-成功 其他-失败
	Message string `json:"message"` // 响应消息
}

// LeaveClass 学生退出班级
func (s *ClassServiceImpl) LeaveClass(ctx context.Context, req *LeaveClassRequest) (*LeaveClassResponse, error) {
	// 调用service层处理业务逻辑
	err := s.classService.LeaveClass(req.StudentId, req.ClassId)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &LeaveClassResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &LeaveClassResponse{
		Code:    consts.SuccessCode,
		Message: "退出班级成功",
	}, nil
}

// RemoveStudentRequest 教师移除学生请求
type RemoveStudentRequest struct {
	TeacherId string `json:"teacher_id"` // 教师ID（必填）
	ClassId   string `json:"class_id"`   // 班级ID（必填）
	StudentId string `json:"student_id"` // 学生ID（必填）
}

// RemoveStudentResponse 教师移除学生响应
type RemoveStudentResponse struct {
	Code    int32  `json:"code"`    // 响应码 0-成功 其他-失败
	Message string `json:"message"` // 响应消息
}

// RemoveStudent 教师移除学生
func (s *ClassServiceImpl) RemoveStudent(ctx context.Context, req *RemoveStudentRequest) (*RemoveStudentResponse, error) {
	// 调用service层处理业务逻辑
	err := s.classService.RemoveStudent(req.TeacherId, req.ClassId, req.StudentId)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &RemoveStudentResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &RemoveStudentResponse{
		Code:    consts.SuccessCode,
		Message: "移除学生成功",
	}, nil
}

// GetClassMembersRequest 获取班级成员列表请求
type GetClassMembersRequest struct {
	ClassId  string `json:"class_id"`  // 班级ID（必填）
	Page     int32  `json:"page"`      // 页码（从1开始）
	PageSize int32  `json:"page_size"` // 每页数量
}

// ClassMemberInfo 班级成员信息
type ClassMemberInfo struct {
	ClassId       string `json:"class_id"`       // 班级ID
	StudentId     string `json:"student_id"`     // 学生ID
	StudentName   string `json:"student_name"`   // 学生姓名
	StudentNumber string `json:"student_number"` // 学号
	Email         string `json:"email"`          // 邮箱
	JoinTime      string `json:"join_time"`      // 加入时间
	Status        int32  `json:"status"`         // 状态
	Remark        string `json:"remark"`         // 备注
	CreateTime    string `json:"create_time"`    // 创建时间
	UpdateTime    string `json:"update_time"`    // 更新时间
}

// GetClassMembersResponse 获取班级成员列表响应
type GetClassMembersResponse struct {
	Code     int32              `json:"code"`      // 响应码 0-成功 其他-失败
	Message  string             `json:"message"`   // 响应消息
	Members  []*ClassMemberInfo `json:"members"`   // 成员列表
	Total    int32              `json:"total"`     // 总数
	Page     int32              `json:"page"`      // 当前页码
	PageSize int32              `json:"page_size"` // 每页数量
}

// GetClassMembers 获取班级成员列表
func (s *ClassServiceImpl) GetClassMembers(ctx context.Context, req *GetClassMembersRequest) (*GetClassMembersResponse, error) {
	// 调用service层处理业务逻辑
	members, total, err := s.classService.GetClassMembers(req.ClassId, req.Page, req.PageSize)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &GetClassMembersResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	// 转换为响应格式，并查询学生详情
	memberInfos := make([]*ClassMemberInfo, 0, len(members))
	for _, member := range members {
		info := &ClassMemberInfo{
			ClassId:    member.ClassId,
			StudentId:  member.StudentId,
			JoinTime:   member.JoinTime.Format("2006-01-02 15:04:05"),
			Status:     member.Status,
			Remark:     member.Remark,
			CreateTime: member.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime: member.UpdateTime.Format("2006-01-02 15:04:05"),
		}
		// 查询学生详情（姓名、学号、邮箱）
		if student, err := s.studentDAO.GetStudentById(member.StudentId); err == nil && student != nil {
			info.StudentName = student.StudentName
			info.StudentNumber = student.StudentNumber
			info.Email = student.Email
		}
		memberInfos = append(memberInfos, info)
	}

	return &GetClassMembersResponse{
		Code:     consts.SuccessCode,
		Message:  consts.MessageQuerySuccess,
		Members:  memberInfos,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetStudentClassesRequest 获取学生加入的班级列表请求
type GetStudentClassesRequest struct {
	StudentId string `json:"student_id"` // 学生ID（必填）
	Page      int32  `json:"page"`       // 页码（从1开始）
	PageSize  int32  `json:"page_size"`  // 每页数量
}

// GetStudentClassesResponse 获取学生加入的班级列表响应
type GetStudentClassesResponse struct {
	Code     int32        `json:"code"`      // 响应码 0-成功 其他-失败
	Message  string       `json:"message"`   // 响应消息
	Classes  []*ClassInfo `json:"classes"`   // 班级列表（含完整班级信息）
	Total    int32        `json:"total"`     // 总数
	Page     int32        `json:"page"`      // 当前页码
	PageSize int32        `json:"page_size"` // 每页数量
}

// GetStudentClasses 获取学生加入的班级列表
func (s *ClassServiceImpl) GetStudentClasses(ctx context.Context, req *GetStudentClassesRequest) (*GetStudentClassesResponse, error) {
	// 调用service层处理业务逻辑，获取学生加入的班级成员记录
	members, total, err := s.classService.GetStudentClasses(req.StudentId, req.Page, req.PageSize)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &GetStudentClassesResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	// 根据 class_id 查询完整班级信息
	classInfos := make([]*ClassInfo, 0, len(members))
	for _, member := range members {
		class, err := s.classService.GetClassById(member.ClassId)
		if err != nil || class == nil {
			continue
		}
		info := &ClassInfo{
			ClassId:         class.ClassId,
			ClassName:       class.ClassName,
			ClassCode:       class.ClassCode,
			TeacherId:       class.TeacherId,
			SubjectId:       class.SubjectId,
			Subject:         class.Subject,
			Semester:        class.Semester,
			MaxStudents:     class.MaxStudents,
			CurrentStudents: class.CurrentStudents,
			Description:     class.Description,
			Announcement:    class.Announcement,
			QrCodeUrl:       class.QrCodeUrl,
			Status:          class.Status,
			CreateTime:      class.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:      class.UpdateTime.Format("2006-01-02 15:04:05"),
		}
		// 查询教师姓名
		if teacher, err := s.teacherService.GetTeacherById(class.TeacherId); err == nil && teacher != nil {
			info.TeacherName = teacher.TeacherName
		}
		// 查询科目名称（兼容旧数据）
		if class.Subject != "" {
			info.SubjectName = class.Subject
		} else if subj, err := s.subjectService.GetSubjectById(class.SubjectId); err == nil && subj != nil {
			info.SubjectName = subj.SubjectName
		}
		classInfos = append(classInfos, info)
	}

	return &GetStudentClassesResponse{
		Code:     consts.SuccessCode,
		Message:  consts.MessageQuerySuccess,
		Classes:  classInfos,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetTeacherClassesRequest 获取教师创建的班级列表请求
type GetTeacherClassesRequest struct {
	TeacherId string `json:"teacher_id"` // 教师ID（必填）
	Page      int32  `json:"page"`       // 页码（从1开始）
	PageSize  int32  `json:"page_size"`  // 每页数量
}

// ClassInfo 班级信息
type ClassInfo struct {
	ClassId         string `json:"class_id"`         // 班级ID
	ClassName       string `json:"class_name"`       // 班级名称
	ClassCode       string `json:"class_code"`       // 班级验证码
	TeacherId       string `json:"teacher_id"`       // 教师ID
	TeacherName     string `json:"teacher_name"`     // 教师姓名
	SubjectId       string `json:"subject_id"`       // 科目ID
	SubjectName     string `json:"subject_name"`     // 科目名称
	Subject         string `json:"subject"`          // 科目名称（冗余字段，方便展示）
	Semester        string `json:"semester"`         // 学期
	MaxStudents     int32  `json:"max_students"`     // 学生人数上限
	CurrentStudents int32  `json:"current_students"` // 当前学生人数
	Description     string `json:"description"`      // 班级描述
	Announcement    string `json:"announcement"`     // 班级公告
	QrCodeUrl       string `json:"qr_code_url"`      // 二维码URL
	Status          int32  `json:"status"`           // 状态
	CreateTime      string `json:"create_time"`      // 创建时间
	UpdateTime      string `json:"update_time"`      // 更新时间
}

// GetTeacherClassesResponse 获取教师创建的班级列表响应
type GetTeacherClassesResponse struct {
	Code     int32        `json:"code"`      // 响应码 0-成功 其他-失败
	Message  string       `json:"message"`   // 响应消息
	Classes  []*ClassInfo `json:"classes"`   // 班级列表
	Total    int32        `json:"total"`     // 总数
	Page     int32        `json:"page"`      // 当前页码
	PageSize int32        `json:"page_size"` // 每页数量
}

// GetTeacherClasses 获取教师创建的班级列表
func (s *ClassServiceImpl) GetTeacherClasses(ctx context.Context, req *GetTeacherClassesRequest) (*GetTeacherClassesResponse, error) {
	// 调用service层处理业务逻辑
	classes, total, err := s.classService.GetTeacherClasses(req.TeacherId, req.Page, req.PageSize)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &GetTeacherClassesResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	// 转换为响应格式
	classInfos := make([]*ClassInfo, 0, len(classes))
	for _, class := range classes {
		info := &ClassInfo{
			ClassId:         class.ClassId,
			ClassName:       class.ClassName,
			ClassCode:       class.ClassCode,
			TeacherId:       class.TeacherId,
			SubjectId:       class.SubjectId,
			Subject:         class.Subject,
			Semester:        class.Semester,
			MaxStudents:     class.MaxStudents,
			CurrentStudents: class.CurrentStudents,
			Description:     class.Description,
			Announcement:    class.Announcement,
			QrCodeUrl:       class.QrCodeUrl,
			Status:          class.Status,
			CreateTime:      class.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:      class.UpdateTime.Format("2006-01-02 15:04:05"),
		}
		// 查询教师姓名
		if teacher, err := s.teacherService.GetTeacherById(class.TeacherId); err == nil && teacher != nil {
			info.TeacherName = teacher.TeacherName
		}
		// 查询科目名称（兼容旧数据）
		if class.Subject != "" {
			info.SubjectName = class.Subject
		} else if subj, err := s.subjectService.GetSubjectById(class.SubjectId); err == nil && subj != nil {
			info.SubjectName = subj.SubjectName
		}
		classInfos = append(classInfos, info)
	}

	return &GetTeacherClassesResponse{
		Code:     consts.SuccessCode,
		Message:  consts.MessageQuerySuccess,
		Classes:  classInfos,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// UpdateClassRequest 更新班级信息请求
type UpdateClassRequest struct {
	TeacherId    string `json:"teacher_id"`   // 教师ID（必填）
	ClassId      string `json:"class_id"`     // 班级ID（必填）
	ClassName    string `json:"class_name"`   // 班级名称（可选）
	SubjectId    string `json:"subject_id"`   // 科目ID（可选）
	Description  string `json:"description"`  // 班级描述（可选）
	Announcement string `json:"announcement"` // 班级公告（可选）
	MaxStudents  int32  `json:"max_students"` // 学生人数上限（可选）
	Status       int32  `json:"status"`       // 状态（可选，-1表示不更新）
}

// UpdateClassResponse 更新班级信息响应
type UpdateClassResponse struct {
	Code    int32  `json:"code"`    // 响应码 0-成功 其他-失败
	Message string `json:"message"` // 响应消息
}

// UpdateClass 更新班级信息
func (s *ClassServiceImpl) UpdateClass(ctx context.Context, req *UpdateClassRequest) (*UpdateClassResponse, error) {
	// 构建更新字段
	updates := make(map[string]interface{})
	if req.ClassName != "" {
		updates["class_name"] = req.ClassName
	}
	if req.SubjectId != "" {
		updates["subject_id"] = req.SubjectId
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Announcement != "" {
		updates["announcement"] = req.Announcement
	}
	if req.MaxStudents > 0 {
		updates["max_students"] = req.MaxStudents
	}
	if req.Status >= 0 {
		updates["status"] = req.Status
	}

	// 调用service层处理业务逻辑
	_, err := s.classService.UpdateClass(req.TeacherId, req.ClassId, updates)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &UpdateClassResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &UpdateClassResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageUpdateSuccess,
	}, nil
}

// GetClassByCodeRequest 根据验证码获取班级信息请求
type GetClassByCodeRequest struct {
	ClassCode string `json:"class_code"` // 班级验证码（必填）
}

// GetClassByCodeResponse 根据验证码获取班级信息响应
type GetClassByCodeResponse struct {
	Code    int32      `json:"code"`    // 响应码 0-成功 其他-失败
	Message string     `json:"message"` // 响应消息
	Class   *ClassInfo `json:"class"`   // 班级信息
}

// GetClassByCode 根据验证码获取班级信息
func (s *ClassServiceImpl) GetClassByCode(ctx context.Context, req *GetClassByCodeRequest) (*GetClassByCodeResponse, error) {
	// 调用service层处理业务逻辑
	class, err := s.classService.GetClassByCode(req.ClassCode)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &GetClassByCodeResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	classInfo := &ClassInfo{
		ClassId:         class.ClassId,
		ClassName:       class.ClassName,
		ClassCode:       class.ClassCode,
		TeacherId:       class.TeacherId,
		SubjectId:       class.SubjectId,
		Subject:         class.Subject,
		Semester:        class.Semester,
		MaxStudents:     class.MaxStudents,
		CurrentStudents: class.CurrentStudents,
		Description:     class.Description,
		Announcement:    class.Announcement,
		QrCodeUrl:       class.QrCodeUrl,
		Status:          class.Status,
		CreateTime:      class.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:      class.UpdateTime.Format("2006-01-02 15:04:05"),
	}
	// 查询教师姓名
	if teacher, err := s.teacherService.GetTeacherById(class.TeacherId); err == nil && teacher != nil {
		classInfo.TeacherName = teacher.TeacherName
	}
	// 查询科目名称（兼容旧数据）
	if class.Subject != "" {
		classInfo.SubjectName = class.Subject
	} else if subj, err := s.subjectService.GetSubjectById(class.SubjectId); err == nil && subj != nil {
		classInfo.SubjectName = subj.SubjectName
	}

	return &GetClassByCodeResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageQuerySuccess,
		Class:   classInfo,
	}, nil
}

// ==================== 查询科目列表 ====================

// ListSubjectsRequest 查询科目列表请求
type ListSubjectsRequest struct{}

// SubjectItem 科目简要信息
type SubjectItem struct {
	SubjectId   string `json:"subject_id"`
	SubjectName string `json:"subject_name"`
}

// ListSubjectsResponse 查询科目列表响应
type ListSubjectsResponse struct {
	Code     int32          `json:"code"`
	Message  string         `json:"message"`
	Subjects []*SubjectItem `json:"subjects"`
}

// ListSubjects 查询全量启用科目列表（供创建班级时选择）
func (s *ClassServiceImpl) ListSubjects(ctx context.Context) (*ListSubjectsResponse, error) {
	subjectDAO := dao.NewSubjectDAO()
	subjects, err := subjectDAO.ListSubjects("status = ?", []interface{}{1}, 1000, 0)
	if err != nil {
		return &ListSubjectsResponse{Code: errs.ErrInternal, Message: "查询科目列表失败"}, err
	}
	items := make([]*SubjectItem, 0, len(subjects))
	for _, s := range subjects {
		items = append(items, &SubjectItem{
			SubjectId:   s.SubjectId,
			SubjectName: s.SubjectName,
		})
	}
	return &ListSubjectsResponse{
		Code:     consts.SuccessCode,
		Message:  consts.MessageQuerySuccess,
		Subjects: items,
	}, nil
}

// ==================== 查询学期列表 ====================

// SemesterItem 学期简要信息
type SemesterItem struct {
	SemesterId   string `json:"semester_id"`
	SemesterName string `json:"semester_name"`
}

// ListSemestersResponse 查询学期列表响应
type ListSemestersResponse struct {
	Code      int32           `json:"code"`
	Message   string          `json:"message"`
	Semesters []*SemesterItem `json:"semesters"`
}

// ListSemesters 查询全量启用学期列表（供创建班级时选择）
func (s *ClassServiceImpl) ListSemesters(ctx context.Context) (*ListSemestersResponse, error) {
	semesterDAO := dao.NewSemesterDAO()
	semesters, err := semesterDAO.ListSemesters(1)
	if err != nil {
		return &ListSemestersResponse{Code: errs.ErrInternal, Message: "查询学期列表失败"}, err
	}
	items := make([]*SemesterItem, 0, len(semesters))
	for _, sem := range semesters {
		items = append(items, &SemesterItem{
			SemesterId:   sem.SemesterId,
			SemesterName: sem.SemesterName,
		})
	}
	return &ListSemestersResponse{
		Code:      consts.SuccessCode,
		Message:   consts.MessageQuerySuccess,
		Semesters: items,
	}, nil
}
