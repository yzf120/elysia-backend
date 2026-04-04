-- code_run 表新增 teacher_id 字段（教师代码调试运行记录）
ALTER TABLE `code_run` ADD COLUMN `teacher_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '教师ID（教师调试时使用）' AFTER `student_id`;
ALTER TABLE `code_run` ADD INDEX `idx_teacher_problem` (`teacher_id`, `problem_id`);
