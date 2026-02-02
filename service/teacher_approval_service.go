package service

import (
	"fmt"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/teacher"
)

// TeacherApprovalService 教师审批单服务
type TeacherApprovalService struct {
	approvalDAO dao.TeacherApprovalDAO
	teacherDAO  dao.TeacherDAO
	adminDAO    dao.AdminUserDAO
}

// NewTeacherApprovalService 创建教师审批单服务
func NewTeacherApprovalService() *TeacherApprovalService {
	return &TeacherApprovalService{
		approvalDAO: dao.NewTeacherApprovalDAO(),
		teacherDAO:  dao.NewTeacherDAO(),
		adminDAO:    dao.NewAdminUserDAO(),
	}
}

// CreateApproval 创建审批单（教师注册时自动创建）
func (s *TeacherApprovalService) CreateApproval(teacherId, employeeNumber, schoolEmail, teacherName, phone, department, title, teachingSubjects string, teachingYears int32, applyRemark string) (*teacher.TeacherApproval, error) {
	// 检查该教师是否已有审批单
	existingApproval, _ := s.approvalDAO.GetApprovalByTeacherId(teacherId)
	if existingApproval != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "该教师已有审批单")
	}

	// 生成审批单ID
	approvalId := fmt.Sprintf("APV%d", time.Now().UnixNano())

	// 创建审批单
	approval := &teacher.TeacherApproval{
		ApprovalId:       approvalId,
		TeacherId:        teacherId,
		EmployeeNumber:   employeeNumber,
		SchoolEmail:      schoolEmail,
		TeacherName:      teacherName,
		Phone:            phone,
		Department:       department,
		Title:            title,
		TeachingSubjects: teachingSubjects,
		TeachingYears:    teachingYears,
		ApplyRemark:      applyRemark,
		ApprovalStatus:   0, // 待审批
	}

	if err := s.approvalDAO.CreateApproval(approval); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建审批单失败: "+err.Error())
	}

	return approval, nil
}

// GetApprovalById 根据审批单ID获取审批单
func (s *TeacherApprovalService) GetApprovalById(approvalId string) (*teacher.TeacherApproval, error) {
	approval, err := s.approvalDAO.GetApprovalById(approvalId)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询审批单失败: "+err.Error())
	}
	if approval == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "审批单不存在")
	}
	return approval, nil
}

// GetApprovalByTeacherId 根据教师ID获取审批单
func (s *TeacherApprovalService) GetApprovalByTeacherId(teacherId string) (*teacher.TeacherApproval, error) {
	approval, err := s.approvalDAO.GetApprovalByTeacherId(teacherId)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询审批单失败: "+err.Error())
	}
	if approval == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "审批单不存在")
	}
	return approval, nil
}

// ApproveTeacher 审批教师（管理员操作）
func (s *TeacherApprovalService) ApproveTeacher(approvalId, adminId string, approved bool, remark string) error {
	// 验证管理员身份
	admin, err := s.adminDAO.GetAdminUserByAdminId(adminId)
	if err != nil || admin == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "管理员不存在或无权限")
	}

	// 检查管理员状态
	if admin.Status != 1 {
		return errs.NewCommonError(errs.ErrAuthenUserForbidden, "管理员账号已被禁用")
	}

	// 查询审批单
	approval, err := s.approvalDAO.GetApprovalById(approvalId)
	if err != nil || approval == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "审批单不存在")
	}

	// 检查审批状态
	if approval.ApprovalStatus != 0 {
		return errs.NewCommonError(errs.ErrBadRequest, "该审批单已处理")
	}

	// 更新审批单状态
	approvalUpdates := map[string]interface{}{
		"approver_id":     adminId,
		"approver_name":   admin.RealName,
		"approval_remark": remark,
		"approval_time":   time.Now(),
	}

	if approved {
		approvalUpdates["approval_status"] = 1 // 审批通过
	} else {
		approvalUpdates["approval_status"] = 2 // 审批驳回
	}

	if err := s.approvalDAO.UpdateApproval(approvalId, approvalUpdates); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "更新审批单失败: "+err.Error())
	}

	// 如果审批通过，更新教师状态
	if approved {
		teacherUpdates := map[string]interface{}{
			"verification_status": 1, // 已通过
			"status":              1, // 激活账号
			"verification_time":   time.Now(),
			"verifier_id":         adminId,
			"verification_remark": remark,
		}

		if err := s.teacherDAO.UpdateTeacher(approval.TeacherId, teacherUpdates); err != nil {
			return errs.NewCommonError(errs.ErrInternal, "更新教师状态失败: "+err.Error())
		}
	}

	return nil
}

// ListApprovals 查询审批单列表
func (s *TeacherApprovalService) ListApprovals(page, pageSize int32, filters map[string]interface{}) ([]*teacher.TeacherApproval, int32, error) {
	// 参数校验和默认值设置
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 构建查询条件
	whereClause := "1=1"
	var args []interface{}

	if approvalStatus, ok := filters["approval_status"].(int32); ok {
		whereClause += " AND approval_status = ?"
		args = append(args, approvalStatus)
	}
	if department, ok := filters["department"].(string); ok && department != "" {
		whereClause += " AND department = ?"
		args = append(args, department)
	}
	if teacherName, ok := filters["teacher_name"].(string); ok && teacherName != "" {
		whereClause += " AND teacher_name LIKE ?"
		args = append(args, "%"+teacherName+"%")
	}

	// 查询总数
	total, err := s.approvalDAO.CountApprovals(whereClause, args)
	if err != nil {
		return nil, 0, errs.NewCommonError(errs.ErrInternal, "统计审批单数量失败: "+err.Error())
	}

	// 查询列表
	offset := (page - 1) * pageSize
	approvals, err := s.approvalDAO.ListApprovals(whereClause, args, pageSize, offset)
	if err != nil {
		return nil, 0, errs.NewCommonError(errs.ErrInternal, "查询审批单列表失败: "+err.Error())
	}

	return approvals, total, nil
}

// DeleteApproval 删除审批单（仅限待审批状态）
func (s *TeacherApprovalService) DeleteApproval(approvalId, adminId string) error {
	// 验证管理员身份
	admin, err := s.adminDAO.GetAdminUserByAdminId(adminId)
	if err != nil || admin == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "管理员不存在或无权限")
	}

	// 查询审批单
	approval, err := s.approvalDAO.GetApprovalById(approvalId)
	if err != nil || approval == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "审批单不存在")
	}

	// 只能删除待审批状态的审批单
	if approval.ApprovalStatus != 0 {
		return errs.NewCommonError(errs.ErrBadRequest, "只能删除待审批状态的审批单")
	}

	// 删除审批单
	if err := s.approvalDAO.DeleteApproval(approvalId); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "删除审批单失败: "+err.Error())
	}

	return nil
}
