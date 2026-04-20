-- ============================================================
-- AI模型配置表
-- 用于管理员启用/禁用AI模型
-- 注意：GORM AutoMigrate 会自动创建此表，此SQL仅作参考
-- ============================================================
CREATE TABLE IF NOT EXISTS `ai_model_config` (
    `id`          INT          NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `model_id`    VARCHAR(100) NOT NULL COMMENT '模型ID（如 doubao-seed-1-6-lite-251015）',
    `model_name`  VARCHAR(100) NOT NULL COMMENT '模型显示名称',
    `provider`    VARCHAR(50)  NOT NULL COMMENT '提供商：doubao, qwen',
    `description` VARCHAR(500) DEFAULT NULL COMMENT '模型描述',
    `is_enabled`  TINYINT      NOT NULL DEFAULT 1 COMMENT '是否启用：1-启用 0-禁用',
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_model_id` (`model_id`),
    KEY `idx_is_enabled` (`is_enabled`),
    KEY `idx_provider` (`provider`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI模型配置表（管理员启用/禁用模型）';

-- 初始化默认模型数据
INSERT IGNORE INTO `ai_model_config` (`model_id`, `model_name`, `provider`, `description`, `is_enabled`) VALUES
('doubao-seed-1-6-lite-251015', 'Doubao-Seed-1.6-lite', 'doubao', '豆包多模态模型，支持深度思考，适合快速响应场景', 1),
('qwen3-omni-flash', 'Qwen3-Omni-Flash', 'qwen', '通义千问全模态模型，Thinker–Talker 架构，支持深度思考', 1);
