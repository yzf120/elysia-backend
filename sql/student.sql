-- 学生表
CREATE TABLE IF NOT EXISTS `students` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `student_id` VARCHAR(64) NOT NULL COMMENT '学生唯一标识ID',
  `phone_number` VARCHAR(20) NOT NULL COMMENT '手机号',
  `password` VARCHAR(255) DEFAULT NULL COMMENT '密码（加密存储）',
  `student_name` VARCHAR(100) DEFAULT NULL COMMENT '学生姓名',
  `student_number` VARCHAR(32) DEFAULT NULL COMMENT '学号',
  `email` VARCHAR(100) DEFAULT NULL COMMENT '邮箱',
  `gender` TINYINT DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女',
  `image_url` VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
  `major` VARCHAR(128) DEFAULT NULL COMMENT '专业',
  `grade` VARCHAR(32) DEFAULT NULL COMMENT '年级',
  `programming_level` VARCHAR(32) DEFAULT NULL COMMENT '编程基础：beginner, intermediate, advanced',
  `interests` TEXT COMMENT '兴趣爱好（JSON格式）',
  `learning_tags` TEXT COMMENT '学习标签（JSON格式）',
  `learning_progress` TEXT COMMENT '学习进度（JSON格式）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-正常',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_student_id` (`student_id`),
  UNIQUE KEY `uk_phone_number` (`phone_number`),
  KEY `idx_student_number` (`student_number`),
  KEY `idx_email` (`email`),
  KEY `idx_major` (`major`),
  KEY `idx_grade` (`grade`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学生信息表';

-- ============================================
-- 测试数据插入
-- ============================================

-- 清空测试数据（可选，谨慎使用）
-- TRUNCATE TABLE `students`;

-- 插入测试学生数据
INSERT INTO `students` (
  `student_id`, 
  `phone_number`, 
  `password`, 
  `student_name`, 
  `student_number`, 
  `email`, 
  `gender`, 
  `image_url`, 
  `major`, 
  `grade`, 
  `programming_level`, 
  `interests`, 
  `learning_tags`, 
  `learning_progress`, 
  `status`
) VALUES 
-- 学生1：计算机专业，高级水平
(
  'STU001', 
  '13800138001', 
  '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbmxNsZjv.Ij6Oc4i', -- 密码: password123
  '张三', 
  '2021001001', 
  'zhangsan@example.com', 
  1, 
  'https://example.com/avatars/zhangsan.jpg', 
  '计算机科学与技术', 
  '2021级', 
  'advanced', 
  '["编程", "算法", "开源项目"]', 
  '["Go", "Python", "数据结构", "算法"]', 
  '{"completed": ["Go基础", "数据结构"], "in_progress": ["分布式系统"], "planned": ["云原生"]}', 
  1
),

-- 学生2：软件工程专业，中级水平
(
  'STU002', 
  '13800138002', 
  '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbmxNsZjv.Ij6Oc4i', 
  '李四', 
  '2021001002', 
  'lisi@example.com', 
  1, 
  'https://example.com/avatars/lisi.jpg', 
  '软件工程', 
  '2021级', 
  'intermediate', 
  '["Web开发", "移动开发", "UI设计"]', 
  '["JavaScript", "React", "Node.js"]', 
  '{"completed": ["HTML/CSS", "JavaScript基础"], "in_progress": ["React进阶"], "planned": ["TypeScript"]}', 
  1
),

-- 学生3：数据科学专业，初级水平
(
  'STU003', 
  '13800138003', 
  '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbmxNsZjv.Ij6Oc4i', 
  '王五', 
  '2022001001', 
  'wangwu@example.com', 
  2, 
  'https://example.com/avatars/wangwu.jpg', 
  '数据科学与大数据技术', 
  '2022级', 
  'beginner', 
  '["数据分析", "机器学习", "可视化"]', 
  '["Python", "数据分析", "机器学习"]', 
  '{"completed": ["Python基础"], "in_progress": ["数据分析"], "planned": ["机器学习入门"]}', 
  1
),

-- 学生4：人工智能专业，中级水平
(
  'STU004', 
  '13800138004', 
  '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbmxNsZjv.Ij6Oc4i', 
  '赵六', 
  '2022001002', 
  'zhaoliu@example.com', 
  1, 
  'https://example.com/avatars/zhaoliu.jpg', 
  '人工智能', 
  '2022级', 
  'intermediate', 
  '["深度学习", "计算机视觉", "NLP"]', 
  '["Python", "TensorFlow", "PyTorch", "深度学习"]', 
  '{"completed": ["Python进阶", "机器学习基础"], "in_progress": ["深度学习"], "planned": ["计算机视觉"]}', 
  1
),

-- 学生5：网络工程专业，高级水平
(
  'STU005', 
  '13800138005', 
  '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbmxNsZjv.Ij6Oc4i', 
  '孙七', 
  '2020001001', 
  'sunqi@example.com', 
  1, 
  'https://example.com/avatars/sunqi.jpg', 
  '网络工程', 
  '2020级', 
  'advanced', 
  '["网络安全", "云计算", "DevOps"]', 
  '["Linux", "Docker", "Kubernetes", "网络协议"]', 
  '{"completed": ["网络基础", "Linux系统管理", "Docker"], "in_progress": ["Kubernetes"], "planned": ["云原生安全"]}', 
  1
),

-- 学生6：信息安全专业，中级水平
(
  'STU006', 
  '13800138006', 
  '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbmxNsZjv.Ij6Oc4i', 
  '周八', 
  '2023001001', 
  'zhouba@example.com', 
  2, 
  'https://example.com/avatars/zhouba.jpg', 
  '信息安全', 
  '2023级', 
  'intermediate', 
  '["渗透测试", "密码学", "安全审计"]', 
  '["Python", "网络安全", "Web安全"]', 
  '{"completed": ["网络基础", "Python编程"], "in_progress": ["Web安全"], "planned": ["渗透测试"]}', 
  1
),

-- 学生7：物联网工程专业，初级水平
(
  'STU007', 
  '13800138007', 
  '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbmxNsZjv.Ij6Oc4i', 
  '吴九', 
  '2023001002', 
  'wujiu@example.com', 
  1, 
  'https://example.com/avatars/wujiu.jpg', 
  '物联网工程', 
  '2023级', 
  'beginner', 
  '["嵌入式开发", "传感器", "智能硬件"]', 
  '["C语言", "Arduino", "嵌入式"]', 
  '{"completed": ["C语言基础"], "in_progress": ["Arduino编程"], "planned": ["嵌入式Linux"]}', 
  1
),

-- 学生8：电子信息工程专业，中级水平（禁用状态）
(
  'STU008', 
  '13800138008', 
  '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbmxNsZjv.Ij6Oc4i', 
  '郑十', 
  '2021001003', 
  'zhengshi@example.com', 
  2, 
  'https://example.com/avatars/zhengshi.jpg', 
  '电子信息工程', 
  '2021级', 
  'intermediate', 
  '["信号处理", "通信技术", "FPGA"]', 
  '["MATLAB", "Verilog", "信号处理"]', 
  '{"completed": ["MATLAB基础", "数字信号处理"], "in_progress": ["FPGA设计"], "planned": ["通信原理"]}', 
  0
),

-- 学生9：数字媒体技术专业，初级水平
(
  'STU009', 
  '13800138009', 
  '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbmxNsZjv.Ij6Oc4i', 
  '陈十一', 
  '2023001003', 
  'chenshiyi@example.com', 
  2, 
  'https://example.com/avatars/chenshiyi.jpg', 
  '数字媒体技术', 
  '2023级', 
  'beginner', 
  '["游戏开发", "3D建模", "动画制作"]', 
  '["Unity", "Blender", "游戏设计"]', 
  '{"completed": ["Unity基础"], "in_progress": ["3D建模"], "planned": ["游戏开发进阶"]}', 
  1
),

-- 学生10：计算机科学专业，高级水平
(
  'STU010', 
  '13800138010', 
  '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbmxNsZjv.Ij6Oc4i', 
  '刘十二', 
  '2020001002', 
  'liushier@example.com', 
  1, 
  'https://example.com/avatars/liushier.jpg', 
  '计算机科学与技术', 
  '2020级', 
  'advanced', 
  '["系统架构", "微服务", "性能优化"]', 
  '["Java", "Spring", "微服务", "分布式系统"]', 
  '{"completed": ["Java高级", "Spring全家桶", "微服务架构"], "in_progress": ["性能优化"], "planned": ["云原生架构"]}', 
  1
);

-- ============================================
-- 查询验证
-- ============================================

-- 查询所有学生
SELECT * FROM `students`;

-- 查询正常状态的学生
SELECT * FROM `students` WHERE `status` = 1;

-- 按专业统计学生数量
SELECT `major`, COUNT(*) as student_count 
FROM `students` 
WHERE `status` = 1 
GROUP BY `major` 
ORDER BY student_count DESC;

-- 按年级统计学生数量
SELECT `grade`, COUNT(*) as student_count 
FROM `students` 
WHERE `status` = 1 
GROUP BY `grade` 
ORDER BY `grade`;

-- 按编程水平统计
SELECT `programming_level`, COUNT(*) as student_count 
FROM `students` 
WHERE `status` = 1 
GROUP BY `programming_level`;

-- 查询特定学号的学生
SELECT * FROM `students` WHERE `student_number` = '2021001001';

-- 查询特定手机号的学生
SELECT * FROM `students` WHERE `phone_number` = '13800138001';

-- fix_password.sql
USE elysia;

UPDATE `students`
SET `password` = '$2a$10$QLYXi0yDbEq0evy90NkyyOzM7FNiQF6XJroaYYYeiksu76xAZvlEu'
WHERE `student_id` IN (
                       'STU001', 'STU002', 'STU003', 'STU004', 'STU005',
                       'STU006', 'STU007', 'STU008', 'STU009', 'STU010'
    );

SELECT '密码更新完成！' as message;
SELECT student_id, student_name, LEFT(password, 30) as password_hash FROM `students`;
