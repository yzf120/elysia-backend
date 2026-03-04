package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/service_impl"
)

var (
	classService *service_impl.ClassServiceImpl
)

// registerClass 注册班级相关路由
func registerClass(publicRouter *mux.Router, protectedRouter *mux.Router) {
	// 增删改：仅教师可操作（受保护路由，教师路由前缀）
	protectedRouter.HandleFunc("/teacher/class/create", createClassHandler).Methods("POST")
	protectedRouter.HandleFunc("/teacher/class/update", updateClassHandler).Methods("POST")
	protectedRouter.HandleFunc("/teacher/class/remove-student", removeStudentHandler).Methods("POST")

	// 查询：学生和教师均可调用（受保护路由，通用路由前缀）
	protectedRouter.HandleFunc("/class/subjects", listSubjectsHandler).Methods("GET")
	protectedRouter.HandleFunc("/class/semesters", listSemestersHandler).Methods("GET")
	protectedRouter.HandleFunc("/class/get-by-code", getClassByCodeHandler).Methods("GET")
	protectedRouter.HandleFunc("/class/members", getClassMembersHandler).Methods("POST")
	protectedRouter.HandleFunc("/class/teacher-classes", getTeacherClassesHandler).Methods("POST")
	protectedRouter.HandleFunc("/class/student-classes", getStudentClassesHandler).Methods("POST")

	// 学生操作：加入/退出班级
	protectedRouter.HandleFunc("/student/class/join", joinClassHandler).Methods("POST")
	protectedRouter.HandleFunc("/student/class/leave", leaveClassHandler).Methods("POST")

	// 公告：教师发布/删除（仅教师），师生均可查询
	protectedRouter.HandleFunc("/teacher/class/announcement/publish", publishAnnouncementHandler).Methods("POST")
	protectedRouter.HandleFunc("/teacher/class/announcement/delete", deleteAnnouncementHandler).Methods("POST")
	protectedRouter.HandleFunc("/class/announcements", getAnnouncementsHandler).Methods("POST")
}

// listSubjectsHandler 查询全量启用科目列表
func listSubjectsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	resp, err := classService.ListSubjects(ctx)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBytes)
		return
	}
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// listSemestersHandler 查询全量启用学期列表
func listSemestersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	resp, err := classService.ListSemesters(ctx)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBytes)
		return
	}
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// createClassHandler 创建班级处理器（仅教师）
func createClassHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	request := &service_impl.CreateClassRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	resp, err := classService.CreateClass(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// updateClassHandler 更新班级信息处理器（仅教师）
func updateClassHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	request := &service_impl.UpdateClassRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	resp, err := classService.UpdateClass(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// removeStudentHandler 教师移除学生处理器（仅教师）
func removeStudentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	request := &service_impl.RemoveStudentRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	resp, err := classService.RemoveStudent(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// getClassByCodeHandler 根据验证码查询班级信息处理器（师生共用）
func getClassByCodeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	classCode := r.URL.Query().Get("class_code")
	if classCode == "" {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, "参数class_code不能为空"),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	request := &service_impl.GetClassByCodeRequest{ClassCode: classCode}
	resp, err := classService.GetClassByCode(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// getClassMembersHandler 获取班级成员列表处理器（师生共用）
func getClassMembersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	request := &service_impl.GetClassMembersRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	resp, err := classService.GetClassMembers(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// getTeacherClassesHandler 获取教师班级列表处理器（师生共用）
func getTeacherClassesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	request := &service_impl.GetTeacherClassesRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	resp, err := classService.GetTeacherClasses(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// getStudentClassesHandler 获取学生班级列表处理器（师生共用）
func getStudentClassesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	request := &service_impl.GetStudentClassesRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	resp, err := classService.GetStudentClasses(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// joinClassHandler 学生加入班级处理器
func joinClassHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	request := &service_impl.JoinClassRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	resp, err := classService.JoinClass(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// leaveClassHandler 学生退出班级处理器
func leaveClassHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	request := &service_impl.LeaveClassRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	resp, err := classService.LeaveClass(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// publishAnnouncementHandler 教师发布公告
func publishAnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)
	req := &service_impl.PublishAnnouncementRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := classService.PublishAnnouncement(ctx, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Code != 0 {
		errResp := &errs.BaseResponse{Data: nil, Error: errs.NewError(int(resp.Code), resp.Message)}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		return
	}
	writeSuccessResponse(w, resp)
}

// deleteAnnouncementHandler 教师删除公告
func deleteAnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)
	req := &service_impl.DeleteAnnouncementRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := classService.DeleteAnnouncement(ctx, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Code != 0 {
		errResp := &errs.BaseResponse{Data: nil, Error: errs.NewError(int(resp.Code), resp.Message)}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		return
	}
	writeSuccessResponse(w, resp)
}

// getAnnouncementsHandler 查询班级公告列表（师生共用）
func getAnnouncementsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)
	req := &service_impl.GetAnnouncementsRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := classService.GetAnnouncements(ctx, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Code != 0 {
		errResp := &errs.BaseResponse{Data: nil, Error: errs.NewError(int(resp.Code), resp.Message)}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		return
	}
	writeSuccessResponse(w, resp)
}
