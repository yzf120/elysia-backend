CREATE TABLE `ai_violation_records` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户ID（学生ID或教师ID）',
    `user_role` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户角色：student / teacher',
    `session_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '会话ID',
    `sender_type` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '发送方：user-用户消息 / ai-AI回复',
    `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '违规内容（截取前500字符）',
    `suggestion` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '审核建议：Block-违规 / Review-疑似',
    `label` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '违规标签：Porn-色情 / Abuse-谩骂 / Ad-广告 等',
    `score` int NOT NULL DEFAULT 0 COMMENT '违规置信度 0-100',
    `request_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '腾讯云内容安全审核请求ID',
    `content_type` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'text' COMMENT '内容类型：text-文本 / image-图片',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`) USING BTREE COMMENT '按用户ID查询违规记录',
    KEY `idx_session_id` (`session_id`) USING BTREE COMMENT '按会话ID查询违规记录',
    KEY `idx_create_time` (`create_time`) USING BTREE COMMENT '按创建时间排序查询'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='AI对话违规记录表（用户消息或AI回复触发内容安全审核未通过时写入）';
