-- ============================================================
-- OJ AI助教 - 意图相关表结构
-- ============================================================

-- 1. 意图字典表：存储标准化的意图分类（4大类 + 二级子意图）
CREATE TABLE `oj_intent_dict` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `intent_level1` VARCHAR(50) NOT NULL COMMENT '一级意图（4大类）：解题相关/知识答疑/操作交互控制/无关兜底',
  `intent_level2` VARCHAR(100) NOT NULL COMMENT '二级子意图（细化标签，便于Agent精准处理）',
  `intent_code` VARCHAR(30) NOT NULL COMMENT '意图编码（唯一标识，如SOLVE_PROBLEM、KNOWLEDGE_QA）',
  `description` VARCHAR(500) DEFAULT NULL COMMENT '意图描述（供管理员理解该意图的含义和适用场景）',
  `match_keywords` TEXT COMMENT '匹配关键词（逗号分隔，用于意图识别辅助）',
  `example_queries` TEXT COMMENT '示例用户问题（JSON数组格式，用于LLM few-shot识别）',
  `rewrite_template` TEXT COMMENT '改写模板（用于标准化请求的模板，支持{query}等占位符）',
  `agent_route` VARCHAR(50) NOT NULL COMMENT '路由到的下游Agent名称：解题Agent/知识Agent/操作Agent/兜底Agent',
  `priority` INT NOT NULL DEFAULT 0 COMMENT '优先级（数值越大优先级越高，用于意图冲突时的优先匹配）',
  `is_valid` TINYINT NOT NULL DEFAULT 1 COMMENT '是否有效：1-有效 0-废弃',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_intent_code` (`intent_code`),
  KEY `idx_intent_level1` (`intent_level1`),
  KEY `idx_is_valid` (`is_valid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='OJ AI助教-意图字典表';

-- 2. 意图提示词模板表：独立管理各意图对应的系统提示词和改写提示词
CREATE TABLE `oj_intent_prompt_template` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `intent_code` VARCHAR(30) NOT NULL COMMENT '关联意图编码（关联oj_intent_dict.intent_code）',
  `template_type` VARCHAR(30) NOT NULL COMMENT '模板类型：system_prompt(系统提示词)/rewrite_prompt(改写提示词)/few_shot(示例)',
  `template_name` VARCHAR(100) NOT NULL COMMENT '模板名称（便于管理员识别）',
  `template_content` TEXT NOT NULL COMMENT '模板内容（支持变量占位符：{user_query}、{problem_id}、{intent_code}等）',
  `is_active` TINYINT NOT NULL DEFAULT 1 COMMENT '是否启用：1-启用 0-禁用（同一intent_code+template_type只能有一个启用）',
  `version` INT NOT NULL DEFAULT 1 COMMENT '版本号（支持多版本管理，便于回滚）',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_intent_code` (`intent_code`),
  KEY `idx_template_type` (`template_type`),
  KEY `idx_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='OJ AI助教-意图提示词模板表';

-- 3. 用户意图记录表：记录用户每次请求的意图识别结果
CREATE TABLE `oj_user_intent_record` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` VARCHAR(64) NOT NULL COMMENT '用户唯一标识（OJ平台用户ID）',
  `session_id` VARCHAR(128) DEFAULT NULL COMMENT '会话ID（关联对话上下文）',
  `question_id` VARCHAR(64) DEFAULT NULL COMMENT '关联的OJ题目ID（NULL表示无关联题目）',
  `original_request` TEXT NOT NULL COMMENT '用户原始请求文本',
  `intent_code` VARCHAR(30) NOT NULL COMMENT '匹配的意图编码（关联oj_intent_dict.intent_code）',
  `intent_level1` VARCHAR(50) NOT NULL COMMENT '一级意图分类（冗余字段，便于统计查询）',
  `rewritten_request` TEXT COMMENT '意图Agent改写后的标准化请求',
  `intent_confidence` DECIMAL(5,2) DEFAULT NULL COMMENT '意图识别置信度（0.00-100.00）',
  `response_time_ms` INT DEFAULT NULL COMMENT '意图识别耗时（毫秒）',
  `recognize_status` TINYINT NOT NULL DEFAULT 1 COMMENT '识别状态：1-成功 0-失败 2-降级（LLM超时时使用关键词匹配）',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '请求时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_intent_code` (`intent_code`),
  KEY `idx_intent_level1` (`intent_level1`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='OJ AI助教-用户意图记录表';

-- ============================================================
-- 初始化数据：意图字典
-- ============================================================
INSERT INTO `oj_intent_dict` (`intent_level1`, `intent_level2`, `intent_code`, `description`, `match_keywords`, `example_queries`, `rewrite_template`, `agent_route`, `priority`) VALUES
-- 解题相关
('解题相关', '题目解题思路', 'SOLVE_THINK', '用户请求某道题的解题思路、算法分析', '怎么做,解法,思路,解题步骤,怎么解,算法思路', '["这道题怎么做","两数之和的解法","帮我分析一下思路"]', '请求分析OJ平台题目{problem_id}的解题思路和算法方案', '解题Agent', 100),
('解题相关', '代码BUG排查', 'SOLVE_BUG', '用户提交的代码存在BUG或结果错误，需要排查', '看bug,代码错了,测试用例没过,WA,RE,运行错误,答案错误', '["帮我看看代码哪里有问题","为什么判题显示WA","运行出错了"]', '请求排查用户提交代码的BUG，判题结果为{judge_result}', '解题Agent', 90),
('解题相关', '代码优化', 'SOLVE_OPTIMIZE', '用户代码功能正确但需要性能优化（TLE/MLE）', '优化代码,超时,运行慢,时间复杂度优化,TLE,MLE,内存超限', '["运行超时怎么优化","代码太慢了","怎么降低时间复杂度"]', '请求优化用户代码的性能，当前问题为{performance_issue}', '解题Agent', 80),
-- 知识答疑
('知识答疑', '算法概念解释', 'KNOWLEDGE_ALGO', '用户询问算法/数据结构相关概念', '动态规划,贪心算法,时间复杂度,空间复杂度,二分查找,DFS,BFS,排序', '["动态规划是什么","时间复杂度O(n²)怎么算","贪心和DP的区别"]', '请求解释{concept}的核心概念及OJ场景中的应用', '知识Agent', 70),
('知识答疑', '错误原因解释', 'KNOWLEDGE_ERROR', '用户询问编译错误、运行错误等原因', '编译错误,运行错误,段错误,栈溢出,什么意思,什么原因', '["编译错误是什么意思","段错误怎么回事","RE是什么"]', '请求解释{error_type}的含义、常见原因及解决方法', '知识Agent', 60),
-- 操作/交互控制
('操作交互控制', '平台功能操作', 'OPERATE_PLATFORM', '用户询问OJ平台的功能使用方法', '提交代码,查看记录,切换语言,导出记录,怎么提交,在哪里看', '["怎么提交代码","查看历史提交记录","怎么切换编程语言"]', '请求讲解OJ平台{operation}的操作步骤', '操作Agent', 50),
('操作交互控制', '对话节奏控制', 'OPERATE_DIALOG', '用户控制对话的节奏和方式', '重新说,换个解法,说简单点,回到刚才,详细一点,总结一下', '["重新说一遍","换个解法讲","说简单点","回到刚才的问题"]', '{dialog_control_action}', '操作Agent', 40),
-- 无关/兜底
('无关兜底', '闲聊/无诉求', 'OTHER_CHAT', '用户无解题/答疑/操作诉求，仅闲聊或反馈', '你好,谢谢,天气,不好用,再见,哈哈', '["你好","谢谢","今天天气怎么样","这平台不好用"]', '用户无解题/答疑/操作诉求，触发闲聊兜底回复', '兜底Agent', 10);

-- ============================================================
-- 初始化数据：意图提示词模板
-- ============================================================
INSERT INTO `oj_intent_prompt_template` (`intent_code`, `template_type`, `template_name`, `template_content`, `is_active`, `version`) VALUES
-- 解题思路的系统提示词
('SOLVE_THINK', 'system_prompt', '解题思路-系统提示词v1', '你是一位专业的OJ编程助教。学生正在寻求解题思路帮助。请注意：\n1. 引导学生思考，给出思路提示而非直接给完整代码\n2. 可以分步骤讲解算法思路\n3. 如果学生提供了题目信息，请结合题目分析\n4. 鼓励学生自己尝试实现', 1, 1),
-- 代码BUG排查的系统提示词
('SOLVE_BUG', 'system_prompt', 'BUG排查-系统提示词v1', '你是一位专业的OJ编程助教。学生的代码存在问题需要帮助排查。请注意：\n1. 仔细分析学生的代码逻辑\n2. 指出具体的问题所在和原因\n3. 给出修改建议而非直接给完整修正代码\n4. 解释为什么这样修改', 1, 1),
-- 代码优化的系统提示词
('SOLVE_OPTIMIZE', 'system_prompt', '代码优化-系统提示词v1', '你是一位专业的OJ编程助教。学生的代码需要性能优化。请注意：\n1. 分析当前代码的时间/空间复杂度\n2. 指出性能瓶颈所在\n3. 给出优化方向和思路\n4. 如果可能，给出复杂度更优的算法建议', 1, 1),
-- 算法概念的系统提示词
('KNOWLEDGE_ALGO', 'system_prompt', '算法概念-系统提示词v1', '你是一位专业的OJ编程助教。学生正在学习算法知识。请注意：\n1. 用通俗易懂的语言解释概念\n2. 结合OJ做题场景给出实际应用示例\n3. 可以对比相关概念帮助理解\n4. 适当给出复杂度分析', 1, 1),
-- 错误解释的系统提示词
('KNOWLEDGE_ERROR', 'system_prompt', '错误解释-系统提示词v1', '你是一位专业的OJ编程助教。学生遇到了编程错误需要帮助理解。请注意：\n1. 清晰解释错误的含义\n2. 列举常见的触发原因\n3. 给出排查和修复的方法\n4. 如果可能，给出预防此类错误的建议', 1, 1),
-- 意图改写提示词（通用）
('SOLVE_THINK', 'rewrite_prompt', '解题思路-改写模板v1', '请将用户的原始问题改写为标准化的解题思路请求。\n原始问题：{user_query}\n题目ID：{problem_id}\n改写要求：明确指出用户需要的是解题思路分析，补充题目ID信息。', 1, 1),
('SOLVE_BUG', 'rewrite_prompt', 'BUG排查-改写模板v1', '请将用户的原始问题改写为标准化的BUG排查请求。\n原始问题：{user_query}\n判题结果：{judge_result}\n改写要求：明确指出用户需要排查代码问题，补充判题结果信息。', 1, 1),
('SOLVE_OPTIMIZE', 'rewrite_prompt', '代码优化-改写模板v1', '请将用户的原始问题改写为标准化的代码优化请求。\n原始问题：{user_query}\n性能问题：{performance_issue}\n改写要求：明确指出用户需要优化代码性能，补充性能问题类型。', 1, 1);
