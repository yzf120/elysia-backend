-- ----------------------------
-- 教师表：teachers（独立的教师账号体系，不再关联用户表）
-- ----------------------------
DROP TABLE IF EXISTS `teachers`;
CREATE TABLE `teachers` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `teacher_id` varchar(64) NOT NULL COMMENT '教师唯一标识ID（如TCH+时间戳）',
  `phone_number` varchar(20) NOT NULL COMMENT '手机号（用于登录）',
  `password` varchar(255) NOT NULL COMMENT '密码（加密存储）',
  `teacher_name` varchar(100) NOT NULL DEFAULT '' COMMENT '教师姓名',
  `employee_number` varchar(32) NOT NULL COMMENT '教师工号（唯一）',
  `school_email` varchar(128) NOT NULL COMMENT '学校邮箱（唯一）',
  `gender` tinyint(4) NOT NULL DEFAULT 0 COMMENT '性别：0-未设置，1-男，2-女',
  `image_url` varchar(255) DEFAULT '' COMMENT '教师头像URL',
  `teaching_subjects` text COMMENT '授课科目（JSON格式，兼容旧数据）',
  `teaching_years` int(11) NOT NULL DEFAULT 0 COMMENT '教龄（年）',
  `department` varchar(128) DEFAULT '' COMMENT '所属院系/部门',
  `title` varchar(64) DEFAULT '' COMMENT '职称：助教、讲师、副教授、教授',
  `verification_status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '认证状态：0-待审核，1-已通过，2-已驳回',
  `verification_time` datetime DEFAULT NULL COMMENT '认证审核时间',
  `verifier_id` varchar(64) DEFAULT '' COMMENT '审核人ID（关联管理员表）',
  `verification_remark` varchar(512) DEFAULT '' COMMENT '审核备注',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态：0-未激活，1-正常，2-禁用',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_teacher_id` (`teacher_id`),
  UNIQUE KEY `uk_phone_number` (`phone_number`) COMMENT '手机号唯一索引（用于登录）',
  UNIQUE KEY `uk_employee_number` (`employee_number`) COMMENT '工号唯一索引',
  UNIQUE KEY `uk_school_email` (`school_email`) COMMENT '学校邮箱唯一索引',
  KEY `idx_verification_status` (`verification_status`) COMMENT '按认证状态检索',
  KEY `idx_department` (`department`) COMMENT '按院系检索'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='教师信息表（独立账号体系）';

-- ----------------------------
-- 测试数据
-- ----------------------------
INSERT INTO `teachers` (`teacher_id`, `phone_number`, `password`, `teacher_name`, `employee_number`, `school_email`, `gender`, `image_url`, `teaching_subjects`, `teaching_years`, `department`, `title`, `verification_status`, `verification_time`, `verifier_id`, `verification_remark`, `status`, `create_time`, `update_time`) VALUES
-- 已通过认证的教师
('TCH001', '13800001001', '$2a$10$rGJlZjjsJwJmnP77pdHf0.FQ/38UhDF0oErnOW1NtM1PnBrjA8a0u', '张伟', 'T2024001', 'zhangwei@school.edu.cn', 1, 'https://example.com/avatar/teacher1.jpg', '["高等数学", "线性代数"]', 8, '数学系', '副教授', 1, '2024-01-15 10:30:00', 'ADMIN001', '资质审核通过', 1, '2024-01-10 09:00:00', '2024-01-15 10:30:00'),
('TCH002', '13800001002', '$2a$10$rGJlZjjsJwJmnP77pdHf0.FQ/38UhDF0oErnOW1NtM1PnBrjA8a0u', '李娜', 'T2024002', 'lina@school.edu.cn', 2, 'https://example.com/avatar/teacher2.jpg', '["大学英语", "英语口语"]', 5, '外语系', '讲师', 1, '2024-01-16 14:20:00', 'ADMIN001', '资质审核通过', 1, '2024-01-12 10:30:00', '2024-01-16 14:20:00'),
('TCH003', '13800001003', '$2a$10$rGJlZjjsJwJmnP77pdHf0.FQ/38UhDF0oErnOW1NtM1PnBrjA8a0u', '王强', 'T2024003', 'wangqiang@school.edu.cn', 1, 'https://example.com/avatar/teacher3.jpg', '["数据结构", "算法设计"]', 12, '计算机系', '教授', 1, '2024-01-17 09:15:00', 'ADMIN001', '资质审核通过', 1, '2024-01-13 11:00:00', '2024-01-17 09:15:00'),
('TCH004', '13800001004', '$2a$10$rGJlZjjsJwJmnP77pdHf0.FQ/38UhDF0oErnOW1NtM1PnBrjA8a0u', '刘芳', 'T2024004', 'liufang@school.edu.cn', 2, 'https://example.com/avatar/teacher4.jpg', '["大学物理", "量子力学"]', 10, '物理系', '副教授', 1, '2024-01-18 16:45:00', 'ADMIN001', '资质审核通过', 1, '2024-01-14 14:20:00', '2024-01-18 16:45:00'),
('TCH005', '13800001005', '$2a$10$rGJlZjjsJwJmnP77pdHf0.FQ/38UhDF0oErnOW1NtM1PnBrjA8a0u', '陈明', 'T2024005', 'chenming@school.edu.cn', 1, 'https://example.com/avatar/teacher5.jpg', '["有机化学", "无机化学"]', 6, '化学系', '讲师', 1, '2024-01-19 11:30:00', 'ADMIN001', '资质审核通过', 1, '2024-01-15 09:45:00', '2024-01-19 11:30:00'),

-- 待审核的教师
('TCH006', '13800001006', '$2a$10$rGJlZjjsJwJmnP77pdHf0.FQ/38UhDF0oErnOW1NtM1PnBrjA8a0u', '赵敏', 'T2024006', 'zhaomin@school.edu.cn', 2, '', '["中国文学", "古代汉语"]', 3, '中文系', '助教', 0, NULL, '', '', 0, '2024-02-01 10:00:00', '2024-02-01 10:00:00'),
('TCH007', '13800001007', '$2a$10$rGJlZjjsJwJmnP77pdHf0.FQ/38UhDF0oErnOW1NtM1PnBrjA8a0u', '孙浩', 'T2024007', 'sunhao@school.edu.cn', 1, '', '["微观经济学", "宏观经济学"]', 4, '经济系', '讲师', 0, NULL, '', '', 0, '2024-02-02 15:30:00', '2024-02-02 15:30:00'),

-- 已驳回的教师
('TCH008', '13800001008', '$2a$10$rGJlZjjsJwJmnP77pdHf0.FQ/38UhDF0oErnOW1NtM1PnBrjA8a0u', '周静', 'T2024008', 'zhoujing@school.edu.cn', 2, '', '["市场营销"]', 2, '管理系', '助教', 2, '2024-02-05 09:20:00', 'ADMIN001', '工号信息有误，请重新提交', 2, '2024-02-03 11:00:00', '2024-02-05 09:20:00'),

-- 禁用状态的教师
('TCH009', '13800001009', '$2a$10$rGJlZjjsJwJmnP77pdHf0.FQ/38UhDF0oErnOW1NtM1PnBrjA8a0u', '吴刚', 'T2024009', 'wugang@school.edu.cn', 1, 'https://example.com/avatar/teacher9.jpg', '["体育"]', 7, '体育系', '讲师', 1, '2024-01-20 10:00:00', 'ADMIN001', '资质审核通过', 2, '2024-01-16 08:30:00', '2024-02-10 14:00:00'),

-- 更多正常教师
('TCH010', '13800001010', '$2a$10$rGJlZjjsJwJmnP77pdHf0.FQ/38UhDF0oErnOW1NtM1PnBrjA8a0u', '郑雪', 'T2024010', 'zhengxue@school.edu.cn', 2, 'https://example.com/avatar/teacher10.jpg', '["心理学基础", "发展心理学"]', 9, '心理系', '副教授', 1, '2024-01-21 13:40:00', 'ADMIN001', '资质审核通过', 1, '2024-01-17 10:15:00', '2024-01-21 13:40:00');
