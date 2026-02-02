-- 认证系统重构：删除User表，直接使用Student、Teacher、AdminUser表
-- 执行日期：2026-02-01

-- 1. 修改student表
ALTER TABLE `student` 
  DROP COLUMN `user_id`,
  ADD COLUMN `phone_number` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号' AFTER `student_id`,
  ADD COLUMN `password` varchar(255) DEFAULT NULL COMMENT '密码（可选）' AFTER `phone_number`,
  ADD COLUMN `student_name` varchar(100) DEFAULT NULL COMMENT '学生姓名' AFTER `password`,
  ADD COLUMN `email` varchar(100) DEFAULT NULL COMMENT '邮箱' AFTER `student_number`,
  ADD COLUMN `gender` tinyint DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女' AFTER `email`,
  ADD COLUMN `image_url` varchar(255) DEFAULT NULL COMMENT '头像URL' AFTER `gender`,
  ADD UNIQUE KEY `uk_phone_number` (`phone_number`);

-- 2. 修改teacher表
ALTER TABLE `teacher`
  DROP COLUMN `user_id`,
  ADD COLUMN `phone_number` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号' AFTER `teacher_id`,
  ADD COLUMN `password` varchar(255) DEFAULT NULL COMMENT '密码（可选）' AFTER `phone_number`,
  ADD COLUMN `teacher_name` varchar(100) DEFAULT NULL COMMENT '教师姓名' AFTER `password`,
  ADD COLUMN `gender` tinyint DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女' AFTER `school_email`,
  ADD COLUMN `image_url` varchar(255) DEFAULT NULL COMMENT '头像URL' AFTER `gender`,
  ADD UNIQUE KEY `uk_phone_number` (`phone_number`);

-- 3. 修改admin_user表
ALTER TABLE `admin_user`
  ADD COLUMN `phone_number` varchar(20) DEFAULT NULL COMMENT '手机号' AFTER `username`,
  ADD UNIQUE KEY `uk_phone_number` (`phone_number`);

-- 4. 备份user表（可选，建议先备份再删除）
-- CREATE TABLE `user_backup` AS SELECT * FROM `user`;

-- 5. 删除user表（谨慎操作！）
-- DROP TABLE IF EXISTS `user`;

-- 注意事项：
-- 1. 执行前请先备份数据库
-- 2. 如果有现有数据，需要先迁移数据再执行此脚本
-- 3. 删除user表前请确认所有依赖都已清理
