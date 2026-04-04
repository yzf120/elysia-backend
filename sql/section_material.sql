-- 小节学习资料表
-- 一个小节（section_type=3 学习资料）可以包含多条学习资料（支持混合上传）
CREATE TABLE IF NOT EXISTS `section_material` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `material_id` varchar(64) NOT NULL DEFAULT '' COMMENT '资料唯一ID',
  `section_id` varchar(64) NOT NULL DEFAULT '' COMMENT '所属小节ID',
  `chapter_id` varchar(64) NOT NULL DEFAULT '' COMMENT '所属章节ID（冗余）',
  `class_id` varchar(64) NOT NULL DEFAULT '' COMMENT '所属班级ID（冗余）',
  `teacher_id` varchar(64) NOT NULL DEFAULT '' COMMENT '上传教师ID',
  `title` varchar(256) NOT NULL DEFAULT '' COMMENT '资料标题',
  `description` text COMMENT '资料描述/正文（material_type=text时使用）',
  `material_type` varchar(32) NOT NULL DEFAULT '' COMMENT '资料类型：pdf, word, text, video',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '原始文件名',
  `file_path` varchar(512) NOT NULL DEFAULT '' COMMENT '服务端存储路径',
  `file_size` bigint NOT NULL DEFAULT '0' COMMENT '文件大小（字节）',
  `mime_type` varchar(128) NOT NULL DEFAULT '' COMMENT 'MIME类型',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '排序值（同一小节内）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_material_id` (`material_id`),
  KEY `idx_section_id` (`section_id`) USING BTREE,
  KEY `idx_chapter_id` (`chapter_id`) USING BTREE,
  KEY `idx_class_id` (`class_id`) USING BTREE,
  KEY `idx_section_sort` (`section_id`, `sort_order`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='小节学习资料表';

-- class_section 新增 section_type=3 表示学习资料
-- ALTER TABLE `class_section` MODIFY COLUMN `section_type` tinyint NOT NULL DEFAULT '1' COMMENT '小节类型：1-算法题，2-讨论话题，3-学习资料';
