-- 学期表数据库脚本

-- 1. 创建学期表
CREATE TABLE IF NOT EXISTS `semesters` (
    `id`            BIGINT       NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `semester_id`   VARCHAR(64)  NOT NULL COMMENT '学期ID，如 sem_2026_spring',
    `semester_name` VARCHAR(64)  NOT NULL COMMENT '学期名称，如 2026春',
    `year`          SMALLINT     NOT NULL COMMENT '学年，如 2026',
    `term`          TINYINT      NOT NULL COMMENT '学期：1-春季，2-秋季',
    `start_date`    DATE         NOT NULL COMMENT '学期开始日期',
    `end_date`      DATE         NOT NULL COMMENT '学期结束日期',
    `status`        TINYINT      NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    `remark`        VARCHAR(256) DEFAULT NULL COMMENT '备注',
    `create_time`   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_semester_id`   (`semester_id`),
    UNIQUE KEY `uk_year_term`     (`year`, `term`),
    KEY `idx_status`              (`status`),
    KEY `idx_start_date`          (`start_date`),
    KEY `idx_end_date`            (`end_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学期表';

-- 2. 插入测试数据（近3年，共6个学期）
INSERT INTO `semesters` (`semester_id`, `semester_name`, `year`, `term`, `start_date`, `end_date`, `status`, `remark`) VALUES
('sem_2024_spring', '2024春', 2024, 1, '2024-02-26', '2024-07-05', 1, '2024年春季学期'),
('sem_2024_autumn', '2024秋', 2024, 2, '2024-09-02', '2025-01-17', 1, '2024年秋季学期'),
('sem_2025_spring', '2025春', 2025, 1, '2025-02-24', '2025-07-04', 1, '2025年春季学期'),
('sem_2025_autumn', '2025秋', 2025, 2, '2025-09-01', '2026-01-16', 1, '2025年秋季学期'),
('sem_2026_spring', '2026春', 2026, 1, '2026-02-23', '2026-07-03', 1, '2026年春季学期（当前）'),
('sem_2026_autumn', '2026秋', 2026, 2, '2026-09-07', '2027-01-15', 1, '2026年秋季学期')
ON DUPLICATE KEY UPDATE
    `semester_name` = VALUES(`semester_name`),
    `start_date`    = VALUES(`start_date`),
    `end_date`      = VALUES(`end_date`),
    `status`        = VALUES(`status`),
    `remark`        = VALUES(`remark`);
