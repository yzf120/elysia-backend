package router

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/authen"
	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/auth"
	"github.com/yzf120/elysia-backend/model/student/req"
	"github.com/yzf120/elysia-backend/rpc"
	"github.com/yzf120/elysia-backend/service"
	"github.com/yzf120/elysia-backend/service_impl"
	agent_session "github.com/yzf120/elysia-session/proto/agent_session"
	"golang.org/x/crypto/bcrypt"
)

var (
	studentService        *service_impl.StudentServiceImpl
	studentAuthService    *service.StudentAuthService
	studentProfileService *service.StudentProfileService
)

// registerStudent 学生相关路由
func registerStudent(publicRouter *mux.Router, protectedRouter *mux.Router) {
	// 学生端注册登录接口
	studentAuthRouter := publicRouter.PathPrefix("/student/auth").Subrouter()
	studentAuthRouter.HandleFunc("/send-code-register", studentSendCodeHandler).Methods("POST")
	studentAuthRouter.HandleFunc("/send-code-login", studentSendCodeHandler).Methods("POST")
	studentAuthRouter.HandleFunc("/register-sms", studentRegisterWithSMSHandler).Methods("POST")
	studentAuthRouter.HandleFunc("/login-sms", studentLoginWithSMSHandler).Methods("POST")
	studentAuthRouter.HandleFunc("/login-password", studentLoginWithPasswordHandler).Methods("POST")

	protectedRouter.HandleFunc("/student/create", createStudentHandler).Methods("POST")
	protectedRouter.HandleFunc("/student/get", getStudentHandler).Methods("GET")
	protectedRouter.HandleFunc("/student/update", updateStudentHandler).Methods("POST")
	protectedRouter.HandleFunc("/student/list", listStudentsHandler).Methods("POST")
	protectedRouter.HandleFunc("/student/update-progress", updateLearningProgressHandler).Methods("POST")

	// 学生待办事项（未完成章节）
	protectedRouter.HandleFunc("/student/pending-chapters", getPendingChaptersHandler).Methods("GET")

	// 学生个人信息相关接口（需要认证）
	protectedRouter.HandleFunc("/student/profile", getStudentProfileHandler).Methods("GET")
	protectedRouter.HandleFunc("/student/profile", updateStudentProfileHandler).Methods("POST")
	protectedRouter.HandleFunc("/student/profile/password", updateStudentPasswordHandler).Methods("POST")
	protectedRouter.HandleFunc("/student/profile/study-stats", getStudentStudyStatsHandler).Methods("GET")

	// 学生做题画像接口
	protectedRouter.HandleFunc("/student/profile/coding-stats", getStudentCodingStatsHandler).Methods("GET")

	// 学生班级完成度
	protectedRouter.HandleFunc("/student/class-progress", getStudentClassProgressHandler).Methods("GET")
}

// createStudentHandler 创建学生信息处理器
func createStudentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &req.CreateStudentRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	// 调用服务
	resp, err := studentService.CreateStudent(ctx, request)
	if err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	// 序列化并返回响应
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// getStudentHandler 获取学生信息处理器
func getStudentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 从查询参数获取学生ID
	studentId := r.URL.Query().Get("student_id")

	request := &req.GetStudentRequest{
		StudentId: studentId,
	}

	// 调用服务
	resp, err := studentService.GetStudent(ctx, request)
	if err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	// 序列化并返回响应
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// updateStudentHandler 更新学生信息处理器
func updateStudentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &req.UpdateStudentRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	// 调用服务
	resp, err := studentService.UpdateStudent(ctx, request)
	if err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	// 序列化并返回响应
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// listStudentsHandler 查询学生列表处理器
func listStudentsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &req.ListStudentsRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	// 调用服务
	resp, err := studentService.ListStudents(ctx, request)
	if err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	// 序列化并返回响应
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// updateLearningProgressHandler 更新学习进度处理器
func updateLearningProgressHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &req.UpdateLearningProgressRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	// 调用服务
	resp, err := studentService.UpdateLearningProgress(ctx, request)
	if err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	// 序列化并返回响应
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// ==================== 学生待办事项处理器函数 ====================

// getPendingChaptersHandler 获取学生未完成的章节列表（最多3条）
func getPendingChaptersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	// 从上下文获取学生ID
	studentId, ok := ctx.Value(authen.RoleIDKey).(string)
	if !ok || studentId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 1. 获取学生加入的所有班级
	classMemberDAO := dao.NewClassMemberDAO()
	classMembers, err := classMemberDAO.ListClassesByStudentId(studentId, 100, 0)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "查询班级列表失败")
		return
	}

	classDAO := dao.NewClassDAO()
	chapterDAO := dao.NewChapterDAO()
	codeRunDAO := dao.NewCodeRunDAO()

	type PendingChapter struct {
		ChapterId         string `json:"chapter_id"`
		ChapterTitle      string `json:"chapter_title"`
		CourseName        string `json:"course_name"`
		ClassName         string `json:"class_name"`
		ClassId           string `json:"class_id"`
		TotalSections     int    `json:"total_sections"`
		CompletedSections int    `json:"completed_sections"`
	}

	var pendingChapters []PendingChapter

	// 2. 遍历每个班级，找出未完成的章节
	for _, member := range classMembers {
		classId := member.ClassId

		// 获取班级信息（用于获取班级名称）
		classInfo, err := classDAO.GetClassById(classId)
		if err != nil || classInfo == nil {
			continue
		}

		// 获取班级下所有章节
		chapters, err := chapterDAO.ListChaptersByClassId(classId)
		if err != nil {
			continue
		}

		for _, chapter := range chapters {
			// 获取章节下所有小节
			sections, err := chapterDAO.ListSectionsByChapterId(chapter.ChapterId)
			if err != nil || len(sections) == 0 {
				continue
			}

			// 收集该章节下所有算法题的problem_id
			var problemIds []int64
			for _, section := range sections {
				if section.SectionType == 1 && section.ProblemId != "" {
					pid, err := strconv.ParseInt(section.ProblemId, 10, 64)
					if err == nil && pid > 0 {
						problemIds = append(problemIds, pid)
					}
				}
			}

			// 如果章节没有算法题，视为已完成，跳过
			if len(problemIds) == 0 {
				continue
			}

			// 检查算法题完成情况
			acceptedMap, err := codeRunDAO.BatchGetAcceptedProblems(studentId, problemIds)
			if err != nil {
				continue
			}

			completedCount := 0
			for _, pid := range problemIds {
				if acceptedMap[pid] {
					completedCount++
				}
			}

			// 如果不是全部通过，则为未完成章节
			if completedCount < len(problemIds) {
				pendingChapters = append(pendingChapters, PendingChapter{
					ChapterId:         chapter.ChapterId,
					ChapterTitle:      chapter.Title,
					CourseName:        classInfo.Subject,
					ClassName:         classInfo.ClassName,
					ClassId:           classId,
					TotalSections:     len(problemIds),
					CompletedSections: completedCount,
				})

				// 最多返回3条
				if len(pendingChapters) >= 3 {
					break
				}
			}
		}

		if len(pendingChapters) >= 3 {
			break
		}
	}

	allCompleted := len(pendingChapters) == 0 && len(classMembers) > 0

	writeSuccessResponse(w, map[string]interface{}{
		"pending_chapters": pendingChapters,
		"all_completed":    allCompleted,
		"has_courses":      len(classMembers) > 0,
	})
}

// ==================== 学生认证处理器函数 ====================

// studentSendCodeHandler 学生端发送验证码
func studentSendCodeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	req := &auth.SendCodeRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 从URL路径判断验证码类型
	codeType := consts.Register
	if strings.Contains(r.URL.Path, consts.Login) {
		codeType = consts.Login
	}

	if err := smsService.SendVerificationCode(ctx, req.PhoneNumber, consts.RoleStudent, codeType); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "验证码发送成功",
	})
}

// studentRegisterWithSMSHandler 学生注册
func studentRegisterWithSMSHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	req := &auth.RegisterWithSMSRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	student, err := studentAuthService.RegisterWithSMS(ctx, req.PhoneNumber, req.Code, req.StudentNumber, req.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":   "注册成功",
		"user_info": student,
	})
}

// studentLoginWithSMSHandler 学生验证码登录
func studentLoginWithSMSHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	req := &auth.LoginWithSMSRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	student, token, err := studentAuthService.LoginWithSMS(ctx, req.PhoneNumber, req.Code)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":   "登录成功",
		"user_info": student,
		"token":     token,
	})
}

// studentLoginWithPasswordHandler 学生密码登录
func studentLoginWithPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	req := &auth.LoginWithPasswordRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	student, token, err := studentAuthService.LoginWithPassword(ctx, req.StudentNumber, req.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":   "登录成功",
		"user_info": student,
		"token":     token,
	})
}

// ==================== 学生个人信息处理器函数 ====================

// getStudentProfileHandler 获取当前登录学生个人信息
func getStudentProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	// 从上下文获取学生ID
	studentId, ok := ctx.Value(authen.RoleIDKey).(string)
	if !ok || studentId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权")
		return
	}

	request := &req.GetStudentRequest{
		StudentId: studentId,
	}

	resp, err := studentService.GetStudent(ctx, request)
	if err != nil || resp.Code != consts.SuccessCode {
		msg := "查询失败"
		if resp != nil {
			msg = resp.Message
		}
		writeErrorResponse(w, http.StatusInternalServerError, msg)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "查询成功",
		"student": resp.Student,
	})
}

// updateStudentProfileHandler 更新当前登录学生个人信息
func updateStudentProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	// 从上下文获取学生ID
	studentId, ok := ctx.Value(authen.RoleIDKey).(string)
	if !ok || studentId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 解析请求体
	var reqBody struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Major     string `json:"major"`
		Grade     string `json:"grade"`
		Interests string `json:"interests"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 构建更新请求
	updateReq := &req.UpdateStudentRequest{
		StudentId: studentId,
		Name:      reqBody.Name,
		Email:     reqBody.Email,
		Major:     reqBody.Major,
		Grade:     reqBody.Grade,
		Interests: []string{reqBody.Interests},
	}
	// 如果兴趣为空字符串，传空slice
	if reqBody.Interests == "" {
		updateReq.Interests = nil
	}

	updateResp, err := studentService.UpdateStudent(ctx, updateReq)
	if err != nil || updateResp.Code != consts.SuccessCode {
		msg := "更新失败"
		if updateResp != nil {
			msg = updateResp.Message
		}
		writeErrorResponse(w, http.StatusBadRequest, msg)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "保存成功",
	})
}

// updateStudentPasswordHandler 修改当前登录学生密码
func updateStudentPasswordHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	// 从上下文获取学生ID
	studentId, ok := r.Context().Value(authen.RoleIDKey).(string)
	if !ok || studentId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 解析请求体
	var reqBody struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if reqBody.OldPassword == "" || reqBody.NewPassword == "" {
		writeErrorResponse(w, http.StatusBadRequest, "旧密码和新密码不能为空")
		return
	}

	if len(reqBody.NewPassword) < 6 {
		writeErrorResponse(w, http.StatusBadRequest, "新密码长度不能少于6位")
		return
	}

	// 查询学生信息验证旧密码
	studentDAO := dao.NewStudentDAO()
	studentModel, err := studentDAO.GetStudentById(studentId)
	if err != nil || studentModel == nil {
		writeErrorResponse(w, http.StatusBadRequest, "学生不存在")
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(studentModel.Password), []byte(reqBody.OldPassword)); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "当前密码错误")
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "密码加密失败")
		return
	}

	// 更新密码
	if err := studentDAO.UpdateStudent(studentId, map[string]interface{}{
		"password": string(hashedPassword),
	}); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "密码修改失败")
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "密码修改成功",
	})
}

// getStudentStudyStatsHandler 获取学生学习统计数据
func getStudentStudyStatsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	// 从上下文获取学生ID
	studentId, ok := ctx.Value(authen.RoleIDKey).(string)
	if !ok || studentId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 1. 获取学生加入的所有班级
	classMemberDAO := dao.NewClassMemberDAO()
	classMembers, err := classMemberDAO.ListClassesByStudentId(studentId, 100, 0)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "查询班级列表失败")
		return
	}

	chapterDAO := dao.NewChapterDAO()
	codeRunDAO := dao.NewCodeRunDAO()

	totalChapters := 0
	completedChapters := 0
	completedCourses := 0
	totalCourses := len(classMembers)

	// 2. 遍历每个班级，统计章节完成情况
	for _, member := range classMembers {
		classId := member.ClassId
		allChaptersCompleted := true

		// 获取班级下所有章节
		chapters, err := chapterDAO.ListChaptersByClassId(classId)
		if err != nil {
			continue
		}

		if len(chapters) == 0 {
			allChaptersCompleted = false
		}

		for _, chapter := range chapters {
			totalChapters++

			// 获取章节下所有小节
			sections, err := chapterDAO.ListSectionsByChapterId(chapter.ChapterId)
			if err != nil {
				allChaptersCompleted = false
				continue
			}

			if len(sections) == 0 {
				// 没有小节的章节视为未完成
				allChaptersCompleted = false
				continue
			}

			// 收集该章节下所有算法题的problem_id
			var problemIds []int64
			for _, section := range sections {
				if section.SectionType == 1 && section.ProblemId != "" {
					pid, err := strconv.ParseInt(section.ProblemId, 10, 64)
					if err == nil && pid > 0 {
						problemIds = append(problemIds, pid)
					}
				}
			}

			// 如果章节有算法题，检查是否全部通过
			chapterCompleted := true
			if len(problemIds) > 0 {
				acceptedMap, err := codeRunDAO.BatchGetAcceptedProblems(studentId, problemIds)
				if err != nil {
					chapterCompleted = false
				} else {
					for _, pid := range problemIds {
						if !acceptedMap[pid] {
							chapterCompleted = false
							break
						}
					}
				}
			} else {
				// 没有算法题的章节（如纯讨论）视为已完成
				chapterCompleted = true
			}

			if chapterCompleted {
				completedChapters++
			} else {
				allChaptersCompleted = false
			}
		}

		// 如果该课程所有章节都完成，课程完成+1
		if allChaptersCompleted && len(chapters) > 0 {
			completedCourses++
		}
	}

	// 3. 查询AI对话次数（通过会话数统计）
	aiSessions := 0
	sessionClient := rpc.GetSessionClient()
	if sessionClient != nil {
		svcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		sessionResp, err := sessionClient.ListSessionsByUser(svcCtx, &agent_session.ListSessionsByUserRequest{
			UserId:   studentId,
			Page:     1,
			PageSize: 1,
		})
		if err == nil && sessionResp != nil {
			aiSessions = int(sessionResp.Total)
		}
	}

	writeSuccessResponse(w, map[string]interface{}{
		"completed_courses":  completedCourses,
		"total_courses":      totalCourses,
		"completed_chapters": completedChapters,
		"total_chapters":     totalChapters,
		"ai_sessions":        aiSessions,
	})
}

// getStudentClassProgressHandler 获取学生各班级的章节完成度
func getStudentClassProgressHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	studentId, ok := ctx.Value(authen.RoleIDKey).(string)
	if !ok || studentId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权")
		return
	}

	classMemberDAO := dao.NewClassMemberDAO()
	classMembers, err := classMemberDAO.ListClassesByStudentId(studentId, 100, 0)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "查询班级列表失败")
		return
	}

	chapterDAO := dao.NewChapterDAO()
	codeRunDAO := dao.NewCodeRunDAO()

	type ClassProgress struct {
		ClassId           string `json:"class_id"`
		TotalChapters     int    `json:"total_chapters"`
		CompletedChapters int    `json:"completed_chapters"`
		ProgressPercent   int    `json:"progress_percent"`
	}

	result := make([]ClassProgress, 0, len(classMembers))

	for _, member := range classMembers {
		classId := member.ClassId

		chapters, err := chapterDAO.ListChaptersByClassId(classId)
		if err != nil {
			result = append(result, ClassProgress{ClassId: classId})
			continue
		}

		totalChapters := 0
		completedChapters := 0

		for _, chapter := range chapters {
			sections, err := chapterDAO.ListSectionsByChapterId(chapter.ChapterId)
			if err != nil {
				continue
			}

			// 收集该章节下所有算法题的problem_id
			var problemIds []int64
			for _, section := range sections {
				if section.SectionType == 1 && section.ProblemId != "" {
					pid, err := strconv.ParseInt(section.ProblemId, 10, 64)
					if err == nil && pid > 0 {
						problemIds = append(problemIds, pid)
					}
				}
			}

			if len(problemIds) == 0 {
				// 没有算法题的章节（纯讨论）视为已完成
				totalChapters++
				completedChapters++
				continue
			}

			totalChapters++

			acceptedMap, err := codeRunDAO.BatchGetAcceptedProblems(studentId, problemIds)
			if err != nil {
				continue
			}

			allDone := true
			for _, pid := range problemIds {
				if !acceptedMap[pid] {
					allDone = false
					break
				}
			}
			if allDone {
				completedChapters++
			}
		}

		percent := 0
		if totalChapters > 0 {
			percent = completedChapters * 100 / totalChapters
		}

		result = append(result, ClassProgress{
			ClassId:           classId,
			TotalChapters:     totalChapters,
			CompletedChapters: completedChapters,
			ProgressPercent:   percent,
		})
	}

	writeSuccessResponse(w, map[string]interface{}{
		"class_progress": result,
	})
}

// getStudentCodingStatsHandler 获取学生做题画像数据（供个人页展示）
func getStudentCodingStatsHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	// 从上下文获取学生ID
	studentId, ok := r.Context().Value(authen.RoleIDKey).(string)
	if !ok || studentId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 延迟初始化 profileService
	if studentProfileService == nil {
		studentProfileService = service.NewStudentProfileService()
	}

	profile, err := studentProfileService.GetProfile(studentId)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "查询做题画像失败")
		return
	}

	if profile == nil {
		// 画像不存在，尝试主动生成一次（可能学生已有提交记录但画像尚未生成）
		studentProfileService.UpdateProfileAfterSubmit(studentId)
		// 重新查询
		profile, err = studentProfileService.GetProfile(studentId)
		if err != nil || profile == nil {
			// 仍然没有画像（说明该学生确实没有提交记录）
			writeSuccessResponse(w, map[string]interface{}{
				"has_profile": false,
			})
			return
		}
	}

	// 解析 JSON 字段
	var commonErrors []string
	if profile.CommonErrors != "" {
		_ = json.Unmarshal([]byte(profile.CommonErrors), &commonErrors)
	}

	var languageStats map[string]int
	if profile.LanguageStats != "" {
		_ = json.Unmarshal([]byte(profile.LanguageStats), &languageStats)
	}

	// 返回前端需要展示的关键字段
	writeSuccessResponse(w, map[string]interface{}{
		"has_profile":             true,
		"total_submissions":       profile.TotalSubmissions,
		"accepted_count":          profile.AcceptedCount,
		"accept_rate":             profile.AcceptRate,
		"solved_problem_count":    profile.SolvedProblemCount,
		"attempted_problem_count": profile.AttemptedProblemCount,
		"preferred_language":      profile.PreferredLanguage,
		"language_stats":          languageStats,
		"common_errors":           commonErrors,
		"difficulty_level":        profile.DifficultyLevel,
		"avg_time_cost":           profile.AvgTimeCost,
		"last_submit_time":        profile.LastSubmitTime,
	})
}
