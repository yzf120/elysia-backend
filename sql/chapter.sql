-- 班级章节表
CREATE TABLE IF NOT EXISTS `class_chapter` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `chapter_id` varchar(64) NOT NULL DEFAULT '' COMMENT '章节id',
  `class_id` varchar(64) NOT NULL DEFAULT '' COMMENT '所属班级id',
  `title` varchar(256) NOT NULL DEFAULT '' COMMENT '章节标题',
  `description` text COMMENT '章节描述',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '章节排序（同一班级内，值越小越靠前）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_chapter_id` (`chapter_id`),
  KEY `idx_class_id` (`class_id`) USING BTREE,
  KEY `idx_class_sort` (`class_id`, `sort_order`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='班级章节表';

-- 章节小节表
CREATE TABLE IF NOT EXISTS `class_section` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `section_id` varchar(64) NOT NULL DEFAULT '' COMMENT '小节id',
  `chapter_id` varchar(64) NOT NULL DEFAULT '' COMMENT '所属章节id',
  `class_id` varchar(64) NOT NULL DEFAULT '' COMMENT '所属班级id（冗余，方便查询）',
  `title` varchar(256) NOT NULL DEFAULT '' COMMENT '小节标题',
  `description` text COMMENT '小节描述',
  `section_type` tinyint NOT NULL DEFAULT '1' COMMENT '小节类型：1-算法题，2-讨论话题',
  -- 算法题关联字段（section_type=1 时使用，关联题库）
  `problem_id` varchar(64) NOT NULL DEFAULT '' COMMENT '关联题库的题目ID（section_type=1 时使用）',
  -- 讨论内容字段（section_type=2 时使用）
  `discussion_title` varchar(256) NOT NULL DEFAULT '' COMMENT '讨论话题标题',
  `discussion_content` text COMMENT '讨论话题描述/背景',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '小节排序（同一章节内，值越小越靠前）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_section_id` (`section_id`),
  KEY `idx_chapter_id` (`chapter_id`) USING BTREE,
  KEY `idx_class_id` (`class_id`) USING BTREE,
  KEY `idx_chapter_sort` (`chapter_id`, `sort_order`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='章节小节表';

-- 旧表升级迁移：删除 problem 详情字段，新增 problem_id 字段（若已是新表结构请跳过）
-- 注意：以下语句执行前请确认字段是否存在，已是新表结构则无需执行
ALTER TABLE `class_section` DROP COLUMN `problem_title`;
ALTER TABLE `class_section` DROP COLUMN `problem_difficulty`;
ALTER TABLE `class_section` DROP COLUMN `problem_content`;
ALTER TABLE `class_section` DROP COLUMN `problem_input_desc`;
ALTER TABLE `class_section` DROP COLUMN `problem_output_desc`;
ALTER TABLE `class_section` DROP COLUMN `problem_test_cases`;
ALTER TABLE `class_section` ADD COLUMN `problem_id` varchar(64) NOT NULL DEFAULT '' COMMENT '关联题库的题目ID' AFTER `section_type`;

-- class 表新增 chapter_ids 字段（JSON 数组，存放有序章节id列表）
ALTER TABLE `class` ADD COLUMN `chapter_ids` json DEFAULT NULL COMMENT '章节id列表（有序JSON数组）';

-- ============================================================
-- 示例数据：为班级 cls_1772546178955777000 插入章节和算法题小节
-- ============================================================

-- 插入章节（第一章：数组与哈希表）
INSERT INTO `class_chapter` (
  `chapter_id`,
  `class_id`,
  `title`,
  `description`,
  `sort_order`,
  `status`
) VALUES (
  'chap_1741060000000000001',
  'cls_1772546178955777000',
  '第一章：数组与哈希表',
  '本章介绍数组的基本操作以及哈希表的应用，通过经典题目掌握 O(n) 查找技巧。',
  10,
  1
);

-- 插入小节（算法题：两数之和，关联题库中的题目 id=5）
INSERT INTO `class_section` (
  `section_id`,
  `chapter_id`,
  `class_id`,
  `title`,
  `description`,
  `section_type`,
  `problem_id`,
  `discussion_title`,
  `discussion_content`,
  `sort_order`,
  `status`
) VALUES (
  'sec_1741060000000000002',
  'chap_1741060000000000001',
  'cls_1772546178955777000',
  '两数之和',
  '经典哈希表入门题，掌握用空间换时间的思路。',
  1,
  '5',
  '',
  '',
  10,
  1
);

-- 插入小节（算法题：三数之和，关联题库中的题目）
INSERT INTO `class_section` (
  `section_id`,
  `chapter_id`,
  `class_id`,
  `title`,
  `description`,
  `section_type`,
  `problem_id`,
  `discussion_title`,
  `discussion_content`,
  `sort_order`,
  `status`
) VALUES (
  'sec_1741060000000000003',
  'chap_1741060000000000001',
  'cls_1772546178955777000',
  '三数之和',
  '双指针经典题，在排序数组中高效找出所有不重复的三元组。',
  1,
  'prob_three_sum_001',
  '',
  '',
  20,
  1
);

-- 同步更新 class 表的 chapter_ids 字段
UPDATE `class`
SET `chapter_ids` = JSON_ARRAY('chap_1741060000000000001')
WHERE `class_id` = 'cls_1772546178955777000';