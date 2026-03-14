-- ========================================
-- 教师审核单初始化 SQL
-- 为现有教师生成对应的审核单记录
-- ========================================

-- 清空现有审核单数据（如果存在）
TRUNCATE TABLE teacher_approvals;

-- 为所有现有教师生成审核单记录
INSERT INTO teacher_approvals (
    approval_id,
    teacher_id,
    employee_number,
    school_email,
    teacher_name,
    phone,
    department,
    title,
    teaching_subjects,
    teaching_years,
    apply_remark,
    approval_status,
    approver_id,
    approver_name,
    approval_remark,
    approval_time,
    create_time,
    update_time
)
SELECT
    CONCAT('APV', LPAD(ROW_NUMBER() OVER (ORDER BY t.create_time), 6, '0')) AS approval_id,
    t.teacher_id,
    t.employee_number,
    t.school_email,
    t.teacher_name,
    t.phone_number AS phone,
    IFNULL(t.department, '') AS department,
    IFNULL(t.title, '') AS title,
    IFNULL(t.teaching_subjects, '[]') AS teaching_subjects,
    IFNULL(t.teaching_years, 0) AS teaching_years,
    CONCAT('教师注册申请，教龄', IFNULL(t.teaching_years, 0), '年') AS apply_remark,
    CASE
        WHEN t.verification_status = 1 THEN 1  -- 已通过 -> 审批通过
        WHEN t.verification_status = 2 THEN 2  -- 已驳回 -> 审批驳回
        ELSE 0                                  -- 待审核 -> 待审批
    END AS approval_status,
    IFNULL(t.verifier_id, '') AS approver_id,
    CASE
        WHEN t.verifier_id != '' THEN '管理员'
        ELSE ''
    END AS approver_name,
    IFNULL(t.verification_remark, '') AS approval_remark,
    t.verification_time AS approval_time,
    t.create_time,
    NOW() AS update_time
FROM teachers t
ORDER BY t.create_time;
