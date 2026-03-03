-- 代码运行记录表
CREATE TABLE IF NOT EXISTS `code_run` (
    `id`          BIGINT       NOT NULL AUTO_INCREMENT COMMENT '运行记录ID',
    `problem_id`  BIGINT       NOT NULL COMMENT '题目ID',
    `student_id`  VARCHAR(64)  NOT NULL COMMENT '学生ID',
    `language`    VARCHAR(32)  NOT NULL COMMENT '编程语言：python/java/go/cpp/c',
    `code`        LONGTEXT     NOT NULL COMMENT '提交的代码',
    `run_type`    ENUM('test','submit') NOT NULL DEFAULT 'test' COMMENT '运行类型：test=测试样例，submit=提交',
    `status`      ENUM('pending','running','accepted','wrong_answer','time_limit_exceeded','memory_limit_exceeded','compile_error','runtime_error')
                               NOT NULL DEFAULT 'pending' COMMENT '运行状态',
    `output`      TEXT         COMMENT '实际输出',
    `error_msg`   TEXT         COMMENT '错误信息',
    `time_cost`   BIGINT       DEFAULT 0 COMMENT '执行时间（毫秒）',
    `memory_used` BIGINT       DEFAULT 0 COMMENT '内存使用（KB）',
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_student_problem` (`student_id`, `problem_id`),
    INDEX `idx_problem_id` (`problem_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代码运行记录表';
