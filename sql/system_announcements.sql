CREATE TABLE IF NOT EXISTS `system_announcements` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `announcement_id` VARCHAR(64) NOT NULL COMMENT '公告ID',
    `title` VARCHAR(200) NOT NULL COMMENT '公告标题',
    `content` TEXT NOT NULL COMMENT '公告内容',
    `priority` TINYINT NOT NULL DEFAULT 2 COMMENT '优先级：1-低 2-普通 3-高',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态：0-草稿 1-已发布',
    `publisher_admin_id` VARCHAR(64) NOT NULL COMMENT '发布管理员ID',
    `publisher_name` VARCHAR(128) NOT NULL COMMENT '发布管理员姓名',
    `view_count` BIGINT NOT NULL DEFAULT 0 COMMENT '浏览次数',
    `publish_time` DATETIME DEFAULT NULL COMMENT '发布时间',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_announcement_id` (`announcement_id`),
    KEY `idx_status_publish_time` (`status`, `publish_time`),
    KEY `idx_publisher_admin_id` (`publisher_admin_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='平台系统公告表';
