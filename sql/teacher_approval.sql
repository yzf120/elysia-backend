-- ----------------------------
-- 教师审批单表：teacher_approval（配套审批流程，与教师表强关联）
-- ----------------------------
DROP TABLE IF EXISTS `teacher_approval`;
CREATE TABLE `teacher_approval` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '审批单自增主键',
  `approval_id` varchar(64) NOT NULL COMMENT '审批单唯一标识ID（如APV+时间戳）',
  `teacher_id` varchar(64) NOT NULL COMMENT '关联教师表teacher_id（外键逻辑）',
  `employee_number` varchar(32) NOT NULL COMMENT '提交的教师工号（冗余存储，避免关联）',
  `school_email` varchar(128) NOT NULL COMMENT '提交的学校邮箱（冗余存储）',
  `teacher_name` varchar(64) NOT NULL COMMENT '教师姓名（提交的基础信息）',
  `phone` varchar(16) NOT NULL COMMENT '教师联系电话（提交的基础信息）',
  `department` varchar(128) DEFAULT '' COMMENT '提交的所属部门/院系',
  `title` varchar(64) DEFAULT '' COMMENT '提交的教师职称',
  `teaching_subjects` text COMMENT '提交的授课科目',
  `teaching_years` int(11) NOT NULL DEFAULT 0 COMMENT '提交的教龄',
  `apply_remark` varchar(512) DEFAULT '' COMMENT '教师提交的申请备注',
  `approval_status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '审批单状态：0-待审批 1-审批通过 2-审批驳回',
  `approver_id` varchar(64) DEFAULT '' COMMENT '实际审批人ID（关联管理员表）',
  `approver_name` varchar(64) DEFAULT '' COMMENT '实际审批人姓名（冗余存储）',
  `approval_remark` varchar(512) DEFAULT '' COMMENT '审批人填写的审批意见',
  `approval_time` datetime DEFAULT NULL COMMENT '审批完成时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '审批单创建时间（教师提交时间）',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '审批单更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_approval_id` (`approval_id`),
  UNIQUE KEY `uk_teacher_id` (`teacher_id`) COMMENT '一个教师仅对应一个审批单',
  KEY `idx_approval_status` (`approval_status`) COMMENT '按审批状态检索（高频查询）',
  KEY `idx_create_time` (`create_time`) COMMENT '按提交时间检索'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='教师注册审批单表';
