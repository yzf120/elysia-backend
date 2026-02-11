-- 科目表和教师-科目关联表的数据库迁移脚本

-- 1. 创建科目表
CREATE TABLE IF NOT EXISTS `subjects` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `subject_id` VARCHAR(64) NOT NULL COMMENT '科目ID',
    `subject_name` VARCHAR(128) NOT NULL COMMENT '科目名称',
    `subject_code` VARCHAR(32) NOT NULL COMMENT '科目代码',
    `category` VARCHAR(64) DEFAULT NULL COMMENT '科目分类：理科、文科、艺术等',
    `description` TEXT COMMENT '科目描述',
    `credits` INT DEFAULT 0 COMMENT '学分',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_subject_id` (`subject_id`),
    UNIQUE KEY `uk_subject_code` (`subject_code`),
    KEY `idx_category` (`category`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='科目表';

-- 2. 创建教师-科目关联表
CREATE TABLE IF NOT EXISTS `teacher_subjects` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `teacher_id` VARCHAR(64) NOT NULL COMMENT '教师ID',
    `subject_id` VARCHAR(64) NOT NULL COMMENT '科目ID',
    `start_date` DATE DEFAULT NULL COMMENT '开始教授日期',
    `end_date` DATE DEFAULT NULL COMMENT '结束教授日期',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-在教，0-已停止',
    `remark` VARCHAR(512) DEFAULT NULL COMMENT '备注',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_teacher_subject` (`teacher_id`, `subject_id`),
    KEY `idx_teacher_id` (`teacher_id`),
    KEY `idx_subject_id` (`subject_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教师-科目关联表';

-- 3. 修改教师表：删除teaching_subjects字段
ALTER TABLE `teachers` DROP COLUMN IF EXISTS `teaching_subjects`;

-- 4. 修改班级表：将subject字段改为subject_id，并添加索引
-- 先检查是否存在subject字段，如果存在则进行迁移
ALTER TABLE `classes` 
    ADD COLUMN `subject_id` VARCHAR(64) DEFAULT NULL COMMENT '科目ID' AFTER `teacher_id`;

-- 如果原来有subject字段的数据，可以先保留，后续手动迁移
-- 这里添加索引
ALTER TABLE `classes` 
    ADD INDEX `idx_teacher` (`teacher_id`),
    ADD INDEX `idx_subject` (`subject_id`);

-- 如果确认数据已迁移，可以删除旧的subject字段
-- ALTER TABLE `classes` DROP COLUMN IF EXISTS `subject`;

-- 5. 插入一些示例科目数据（可选）
INSERT INTO `subjects` (`subject_id`, `subject_name`, `subject_code`, `category`, `description`, `credits`, `status`) VALUES
('subj_math_001', '高等数学', 'MATH101', '理科', '高等数学基础课程', 4, 1),
('subj_eng_001', '大学英语', 'ENG101', '文科', '大学英语基础课程', 3, 1),
('subj_phy_001', '大学物理', 'PHY101', '理科', '大学物理基础课程', 4, 1),
('subj_chem_001', '大学化学', 'CHEM101', '理科', '大学化学基础课程', 3, 1),
('subj_cs_001', '计算机科学导论', 'CS101', '理科', '计算机科学入门课程', 3, 1),
('subj_prog_001', '程序设计基础', 'PROG101', '理科', 'Python/Java程序设计', 4, 1),
('subj_hist_001', '中国近代史', 'HIST101', '文科', '中国近代史纲要', 2, 1),
('subj_art_001', '美术鉴赏', 'ART101', '艺术', '美术作品鉴赏课程', 2, 1),
('subj_music_001', '音乐鉴赏', 'MUSIC101', '艺术', '音乐作品鉴赏课程', 2, 1),
('subj_pe_001', '体育', 'PE101', '体育', '体育基础课程', 1, 1)
ON DUPLICATE KEY UPDATE `subject_name` = VALUES(`subject_name`);

-- 6. 数据迁移说明
-- 如果classes表中已有数据，需要手动将subject字段的值迁移到subject_id
-- 可以通过以下步骤：
-- a. 先在subjects表中创建对应的科目记录
-- b. 然后更新classes表，将subject名称对应到subject_id
-- 示例：
-- UPDATE classes c 
-- INNER JOIN subjects s ON c.subject = s.subject_name 
-- SET c.subject_id = s.subject_id 
-- WHERE c.subject_id IS NULL;
