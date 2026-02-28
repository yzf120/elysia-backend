-- ----------------------------
-- 教师审批单表：teacher_approvals（配套审批流程，与教师表强关联）
-- ----------------------------
DROP TABLE IF EXISTS `teacher_approvals`;
CREATE TABLE `teacher_approvals` (
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

-- ----------------------------
-- 测试数据
-- ----------------------------
INSERT INTO `teacher_approvals` (`approval_id`, `teacher_id`, `employee_number`, `school_email`, `teacher_name`, `phone`, `department`, `title`, `teaching_subjects`, `teaching_years`, `apply_remark`, `approval_status`, `approver_id`, `approver_name`, `approval_remark`, `approval_time`, `create_time`, `update_time`) VALUES
-- 待审批的申请单
('APV001', 'TCH006', 'T2024006', 'zhaomin@school.edu.cn', '赵敏', '13800001006', '中文系', '助教', '["中国文学", "古代汉语"]', 3, '本人毕业于北京大学中文系，希望能加入贵校任教', 0, '', '', '', NULL, '2024-02-01 10:00:00', '2024-02-01 10:00:00'),
('APV002', 'TCH007', 'T2024007', 'sunhao@school.edu.cn', '孙浩', '13800001007', '经济系', '讲师', '["微观经济学", "宏观经济学"]', 4, '拥有经济学博士学位，有丰富的教学经验', 0, '', '', '', NULL, '2024-02-02 15:30:00', '2024-02-02 15:30:00'),
('APV003', 'TCH011', 'T2024011', 'huanglin@school.edu.cn', '黄琳', '13800001011', '艺术系', '讲师', '["美术基础", "油画"]', 5, '毕业于中央美术学院，擅长油画教学', 0, '', '', '', NULL, '2024-02-06 09:30:00', '2024-02-06 09:30:00'),

-- 已通过的申请单
('APV004', 'TCH001', 'T2024001', 'zhangwei@school.edu.cn', '张伟', '13800001001', '数学系', '副教授', '["高等数学", "线性代数"]', 8, '本人在数学领域有多年教学经验', 1, 'ADMIN001', '管理员', '资质审核通过，欢迎加入', '2024-01-15 10:30:00', '2024-01-10 09:00:00', '2024-01-15 10:30:00'),
('APV005', 'TCH002', 'T2024002', 'lina@school.edu.cn', '李娜', '13800001002', '外语系', '讲师', '["大学英语", "英语口语"]', 5, '英语专业八级，有海外留学经历', 1, 'ADMIN001', '管理员', '资质审核通过', '2024-01-16 14:20:00', '2024-01-12 10:30:00', '2024-01-16 14:20:00'),
('APV006', 'TCH003', 'T2024003', 'wangqiang@school.edu.cn', '王强', '13800001003', '计算机系', '教授', '["数据结构", "算法设计"]', 12, '计算机科学博士，发表多篇SCI论文', 1, 'ADMIN001', '管理员', '资质审核通过，欢迎加入', '2024-01-17 09:15:00', '2024-01-13 11:00:00', '2024-01-17 09:15:00'),
('APV007', 'TCH004', 'T2024004', 'liufang@school.edu.cn', '刘芳', '13800001004', '物理系', '副教授', '["大学物理", "量子力学"]', 10, '物理学博士，专注量子力学研究', 1, 'ADMIN001', '管理员', '资质审核通过', '2024-01-18 16:45:00', '2024-01-14 14:20:00', '2024-01-18 16:45:00'),
('APV008', 'TCH005', 'T2024005', 'chenming@school.edu.cn', '陈明', '13800001005', '化学系', '讲师', '["有机化学", "无机化学"]', 6, '化学专业硕士，有企业研发经验', 1, 'ADMIN001', '管理员', '资质审核通过', '2024-01-19 11:30:00', '2024-01-15 09:45:00', '2024-01-19 11:30:00'),

-- 已驳回的申请单
('APV009', 'TCH008', 'T2024008', 'zhoujing@school.edu.cn', '周静', '13800001008', '管理系', '助教', '["市场营销"]', 2, '本科毕业，希望从事教学工作', 2, 'ADMIN001', '管理员', '工号信息有误，请核实后重新提交。另外，教龄不满3年，建议积累更多经验后再申请', '2024-02-05 09:20:00', '2024-02-03 11:00:00', '2024-02-05 09:20:00'),
('APV010', 'TCH012', 'T2024012', 'liuyang@school.edu.cn', '刘洋', '13800001012', '音乐系', '助教', '["钢琴"]', 1, '音乐学院毕业，擅长钢琴教学', 2, 'ADMIN001', '管理员', '学校邮箱格式不正确，请使用学校统一邮箱格式重新申请', '2024-02-07 14:30:00', '2024-02-06 16:00:00', '2024-02-07 14:30:00');
