package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	studentModel "github.com/yzf120/elysia-backend/model/student"
	teacherModel "github.com/yzf120/elysia-backend/model/teacher"
)

const (
	adminUserStatusEnabled  = "enabled"
	adminUserStatusDisabled = "disabled"

	adminTeacherAuditPending  = "pending"
	adminTeacherAuditApproved = "approved"
	adminTeacherAuditRejected = "rejected"

	studentStatusEnabledValue  int32 = 1
	studentStatusDisabledValue int32 = 0
	teacherStatusEnabledValue  int32 = 1
	teacherStatusDisabledValue int32 = 2
)

type AdminUserManagementService struct {
	studentDAO dao.StudentDAO
	teacherDAO dao.TeacherDAO
}

type AdminStudentListInput struct {
	Page     int
	PageSize int
	Keyword  string
	Major    string
	Grade    string
	Status   string
}

type AdminTeacherListInput struct {
	Page               int
	PageSize           int
	Keyword            string
	Department         string
	VerificationStatus string
	Status             string
}

type AdminBatchStudentStatusInput struct {
	StudentIDs []string `json:"student_ids"`
	Status     string   `json:"status"`
}

type AdminBatchTeacherStatusInput struct {
	TeacherIDs []string `json:"teacher_ids"`
	Status     string   `json:"status"`
}

type AdminStudentDTO struct {
	StudentID        string `json:"student_id"`
	PhoneNumber      string `json:"phone_number"`
	StudentName      string `json:"student_name"`
	StudentNumber    string `json:"student_number"`
	Email            string `json:"email"`
	Major            string `json:"major"`
	Grade            string `json:"grade"`
	ProgrammingLevel string `json:"programming_level"`
	Status           int32  `json:"status"`
	StatusLabel      string `json:"status_label"`
	CreateTime       string `json:"create_time"`
}

type AdminTeacherDTO struct {
	TeacherID               string `json:"teacher_id"`
	PhoneNumber             string `json:"phone_number"`
	TeacherName             string `json:"teacher_name"`
	EmployeeNumber          string `json:"employee_number"`
	SchoolEmail             string `json:"school_email"`
	Department              string `json:"department"`
	Title                   string `json:"title"`
	VerificationStatus      int32  `json:"verification_status"`
	VerificationStatusLabel string `json:"verification_status_label"`
	Status                  int32  `json:"status"`
	StatusLabel             string `json:"status_label"`
	CreateTime              string `json:"create_time"`
}

func NewAdminUserManagementService() *AdminUserManagementService {
	return &AdminUserManagementService{
		studentDAO: dao.NewStudentDAO(),
		teacherDAO: dao.NewTeacherDAO(),
	}
}

func (s *AdminUserManagementService) ListStudents(input AdminStudentListInput) ([]*AdminStudentDTO, int64, error) {
	page, pageSize := normalizePage(input.Page, input.PageSize)
	whereClause, args := buildStudentWhereClause(input)
	total, err := s.studentDAO.CountStudents(whereClause, args)
	if err != nil {
		return nil, 0, errs.NewCommonError(http.StatusInternalServerError, "统计学生数量失败: "+err.Error())
	}
	items, err := s.studentDAO.ListStudents(whereClause, args, int32(pageSize), int32((page-1)*pageSize))
	if err != nil {
		return nil, 0, errs.NewCommonError(http.StatusInternalServerError, "查询学生列表失败: "+err.Error())
	}
	return mapStudentList(items), int64(total), nil
}

func (s *AdminUserManagementService) ListTeachers(input AdminTeacherListInput) ([]*AdminTeacherDTO, int64, error) {
	page, pageSize := normalizePage(input.Page, input.PageSize)
	whereClause, args, err := buildTeacherWhereClause(input)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.teacherDAO.CountTeachers(whereClause, args)
	if err != nil {
		return nil, 0, errs.NewCommonError(http.StatusInternalServerError, "统计教师数量失败: "+err.Error())
	}
	items, err := s.teacherDAO.ListTeachers(whereClause, args, int32(pageSize), int32((page-1)*pageSize))
	if err != nil {
		return nil, 0, errs.NewCommonError(http.StatusInternalServerError, "查询教师列表失败: "+err.Error())
	}
	return mapTeacherList(items), int64(total), nil
}

func (s *AdminUserManagementService) BatchUpdateStudentsStatus(input AdminBatchStudentStatusInput) (int, error) {
	studentIDs := sanitizeIDs(input.StudentIDs)
	if len(studentIDs) == 0 {
		return 0, errs.NewCommonError(http.StatusBadRequest, "学生ID不能为空")
	}
	statusValue, err := parseStudentStatus(input.Status)
	if err != nil {
		return 0, err
	}
	if err := s.studentDAO.BatchUpdateStudentStatus(studentIDs, statusValue); err != nil {
		return 0, errs.NewCommonError(http.StatusInternalServerError, "批量更新学生状态失败: "+err.Error())
	}
	return len(studentIDs), nil
}

func (s *AdminUserManagementService) BatchUpdateTeachersStatus(input AdminBatchTeacherStatusInput) (int, error) {
	teacherIDs := sanitizeIDs(input.TeacherIDs)
	if len(teacherIDs) == 0 {
		return 0, errs.NewCommonError(http.StatusBadRequest, "教师ID不能为空")
	}
	statusValue, err := parseTeacherStatus(input.Status)
	if err != nil {
		return 0, err
	}
	if err := s.teacherDAO.BatchUpdateTeacherStatus(teacherIDs, statusValue); err != nil {
		return 0, errs.NewCommonError(http.StatusInternalServerError, "批量更新教师状态失败: "+err.Error())
	}
	return len(teacherIDs), nil
}

func (s *AdminUserManagementService) ExportStudents(input AdminStudentListInput) ([]byte, string, error) {
	whereClause, args := buildStudentWhereClause(input)
	items, err := s.studentDAO.ListStudentsAll(whereClause, args)
	if err != nil {
		return nil, "", errs.NewCommonError(http.StatusInternalServerError, "查询学生导出数据失败: "+err.Error())
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{
			item.PhoneNumber,
			item.StudentName,
			item.StudentNumber,
			item.Email,
			item.Major,
			item.Grade,
			item.ProgrammingLevel,
			studentStatusLabel(item.Status),
			formatTime(item.CreateTime),
		})
	}
	content, err := buildCSV([]string{"手机号", "姓名", "学号", "邮箱", "专业", "年级", "编程基础", "账号状态", "注册时间"}, rows)
	if err != nil {
		return nil, "", errs.NewCommonError(http.StatusInternalServerError, "导出学生列表失败: "+err.Error())
	}
	return content, fmt.Sprintf("学生用户列表_%s.csv", time.Now().Format("20060102150405")), nil
}

func (s *AdminUserManagementService) ExportTeachers(input AdminTeacherListInput) ([]byte, string, error) {
	whereClause, args, err := buildTeacherWhereClause(input)
	if err != nil {
		return nil, "", err
	}
	items, err := s.teacherDAO.ListTeachersAll(whereClause, args)
	if err != nil {
		return nil, "", errs.NewCommonError(http.StatusInternalServerError, "查询教师导出数据失败: "+err.Error())
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{
			item.EmployeeNumber,
			item.TeacherName,
			item.PhoneNumber,
			item.SchoolEmail,
			item.Department,
			item.Title,
			teacherVerificationStatusLabel(item.VerificationStatus),
			teacherStatusLabel(item.Status),
			formatTime(item.CreateTime),
		})
	}
	content, err := buildCSV([]string{"工号", "姓名", "手机号", "学校邮箱", "院系", "职称", "审核状态", "账号状态", "注册时间"}, rows)
	if err != nil {
		return nil, "", errs.NewCommonError(http.StatusInternalServerError, "导出教师列表失败: "+err.Error())
	}
	return content, fmt.Sprintf("教师用户列表_%s.csv", time.Now().Format("20060102150405")), nil
}

func buildStudentWhereClause(input AdminStudentListInput) (string, []interface{}) {
	whereClause := "1=1"
	args := make([]interface{}, 0)
	keyword := strings.TrimSpace(input.Keyword)
	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		whereClause += " AND (student_name LIKE ? OR phone_number LIKE ? OR student_number LIKE ? OR student_id LIKE ?)"
		args = append(args, likeKeyword, likeKeyword, likeKeyword, likeKeyword)
	}
	if major := strings.TrimSpace(input.Major); major != "" {
		whereClause += " AND major LIKE ?"
		args = append(args, "%"+major+"%")
	}
	if grade := strings.TrimSpace(input.Grade); grade != "" {
		whereClause += " AND grade LIKE ?"
		args = append(args, "%"+grade+"%")
	}
	switch strings.ToLower(strings.TrimSpace(input.Status)) {
	case adminUserStatusEnabled:
		whereClause += " AND status = ?"
		args = append(args, studentStatusEnabledValue)
	case adminUserStatusDisabled:
		whereClause += " AND status <> ?"
		args = append(args, studentStatusEnabledValue)
	}
	return whereClause, args
}

func buildTeacherWhereClause(input AdminTeacherListInput) (string, []interface{}, error) {
	whereClause := "1=1"
	args := make([]interface{}, 0)
	keyword := strings.TrimSpace(input.Keyword)
	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		whereClause += " AND (teacher_name LIKE ? OR employee_number LIKE ? OR school_email LIKE ? OR phone_number LIKE ? OR teacher_id LIKE ?)"
		args = append(args, likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword)
	}
	if department := strings.TrimSpace(input.Department); department != "" {
		whereClause += " AND department LIKE ?"
		args = append(args, "%"+department+"%")
	}
	switch strings.ToLower(strings.TrimSpace(input.VerificationStatus)) {
	case "", "all":
	case adminTeacherAuditPending:
		whereClause += " AND verification_status = ?"
		args = append(args, int32(0))
	case adminTeacherAuditApproved:
		whereClause += " AND verification_status = ?"
		args = append(args, int32(1))
	case adminTeacherAuditRejected:
		whereClause += " AND verification_status = ?"
		args = append(args, int32(2))
	default:
		return "", nil, errs.NewCommonError(http.StatusBadRequest, "教师审核状态参数不合法")
	}
	switch strings.ToLower(strings.TrimSpace(input.Status)) {
	case "", "all":
	case adminUserStatusEnabled:
		whereClause += " AND status <> ?"
		args = append(args, teacherStatusDisabledValue)
	case adminUserStatusDisabled:
		whereClause += " AND status = ?"
		args = append(args, teacherStatusDisabledValue)
	default:
		return "", nil, errs.NewCommonError(http.StatusBadRequest, "教师账号状态参数不合法")
	}
	return whereClause, args, nil
}

func parseStudentStatus(status string) (int32, error) {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case adminUserStatusEnabled:
		return studentStatusEnabledValue, nil
	case adminUserStatusDisabled:
		return studentStatusDisabledValue, nil
	default:
		return 0, errs.NewCommonError(http.StatusBadRequest, "学生账号状态参数不合法")
	}
}

func parseTeacherStatus(status string) (int32, error) {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case adminUserStatusEnabled:
		return teacherStatusEnabledValue, nil
	case adminUserStatusDisabled:
		return teacherStatusDisabledValue, nil
	default:
		return 0, errs.NewCommonError(http.StatusBadRequest, "教师账号状态参数不合法")
	}
}

func sanitizeIDs(ids []string) []string {
	result := make([]string, 0, len(ids))
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}

func mapStudentList(items []*studentModel.Student) []*AdminStudentDTO {
	result := make([]*AdminStudentDTO, 0, len(items))
	for _, item := range items {
		result = append(result, &AdminStudentDTO{
			StudentID:        item.StudentId,
			PhoneNumber:      item.PhoneNumber,
			StudentName:      item.StudentName,
			StudentNumber:    item.StudentNumber,
			Email:            item.Email,
			Major:            item.Major,
			Grade:            item.Grade,
			ProgrammingLevel: item.ProgrammingLevel,
			Status:           item.Status,
			StatusLabel:      studentStatusLabel(item.Status),
			CreateTime:       formatTime(item.CreateTime),
		})
	}
	return result
}

func mapTeacherList(items []*teacherModel.Teacher) []*AdminTeacherDTO {
	result := make([]*AdminTeacherDTO, 0, len(items))
	for _, item := range items {
		result = append(result, &AdminTeacherDTO{
			TeacherID:               item.TeacherId,
			PhoneNumber:             item.PhoneNumber,
			TeacherName:             item.TeacherName,
			EmployeeNumber:          item.EmployeeNumber,
			SchoolEmail:             item.SchoolEmail,
			Department:              item.Department,
			Title:                   item.Title,
			VerificationStatus:      item.VerificationStatus,
			VerificationStatusLabel: teacherVerificationStatusLabel(item.VerificationStatus),
			Status:                  item.Status,
			StatusLabel:             teacherStatusLabel(item.Status),
			CreateTime:              formatTime(item.CreateTime),
		})
	}
	return result
}

func studentStatusLabel(status int32) string {
	if status == studentStatusEnabledValue {
		return "正常"
	}
	return "已禁用"
}

func teacherStatusLabel(status int32) string {
	if status == teacherStatusDisabledValue {
		return "已禁用"
	}
	return "正常"
}

func teacherVerificationStatusLabel(status int32) string {
	switch status {
	case 1:
		return "已通过"
	case 2:
		return "已驳回"
	default:
		return "待审核"
	}
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format("2006-01-02 15:04:05")
}

func buildCSV(header []string, rows [][]string) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(buffer)
	if err := writer.Write(header); err != nil {
		return nil, err
	}
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
