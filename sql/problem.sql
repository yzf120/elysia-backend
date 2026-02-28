-- 创建适配LeetCode风格的OJ平台题目表（新增标签字段，难度为中文）
CREATE TABLE `problem` (
                                    `id` INT NOT NULL AUTO_INCREMENT COMMENT '题目唯一ID（主键，对应LeetCode题号）',
                                    `title` VARCHAR(255) NOT NULL COMMENT '题目标题（如：两数之和）',
                                    `title_slug` VARCHAR(255) NOT NULL COMMENT 'URL友好标题（用于路由，如：two-sum）',
                                    `difficulty` ENUM('简单', '中等', '困难') NOT NULL DEFAULT '简单' COMMENT '题目难度（中文分级）',
                                    `tags` VARCHAR(500) NOT NULL COMMENT '题目标签/知识点（多个标签用英文逗号分隔）',
                                    `description` TEXT NOT NULL COMMENT '题目完整描述（支持Markdown/换行）',
                                    `input_sample` TEXT COMMENT '输入样例（多个样例用换行分隔）',
                                    `output_sample` TEXT COMMENT '输出样例（与输入样例一一对应）',
                                    `explanation` TEXT COMMENT '样例解释（说明样例的推导逻辑）',
                                    `hint` TEXT COMMENT '题目提示信息（LeetCode的「提示」模块）',
                                    `constraints` TEXT COMMENT '数据范围与约束条件（如数组长度、数值范围）',
                                    `advanced_requirement` TEXT COMMENT '进阶要求（如时间/空间复杂度优化要求）',
                                    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '题目创建时间',
                                    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '题目更新时间',
                                    PRIMARY KEY (`id`),
                                    UNIQUE KEY `uk_title_slug` (`title_slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='LeetCode风格OJ平台题目核心表';

-- 插入「两数之和」示例数据（标签为哈希表，难度为简单）
INSERT INTO `problem` (
    `title`,
    `title_slug`,
    `difficulty`,
    `tags`,
    `description`,
    `input_sample`,
    `output_sample`,
    `explanation`,
    `hint`,
    `constraints`,
    `advanced_requirement`
) VALUES (
             '两数之和',
             'two-sum',
             '简单',
             '哈希表,数组',  -- 补充数组标签，更贴合两数之和的知识点，也可仅保留「哈希表」
             '给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target 的那两个整数，并返回它们的数组下标。\n你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。\n你可以按任意顺序返回答案。',
             '示例 1：\n输入：nums = [2,7,11,15], target = 9\n示例 2：\n输入：nums = [3,2,4], target = 6\n示例 3：\n输入：nums = [3,3], target = 6',
             '示例 1：\n输出：[0,1]\n示例 2：\n输出：[1,2]\n示例 3：\n输出：[0,1]',
             '示例 1 解释：因为 nums[0] + nums[1] == 9，返回 [0, 1]。\n示例 2 解释：因为 nums[1] + nums[2] == 6，返回 [1, 2]。\n示例 3 解释：因为 nums[0] + nums[1] == 6，返回 [0, 1]。',
             '可以尝试使用哈希表来降低时间复杂度',
             '2 <= nums.length <= 10^4\n-10^9 <= nums[i] <= 10^9\n-10^9 <= target <= 10^9\n只会存在一个有效答案',
             '进阶：你可以想出一个时间复杂度小于 O(n²) 的算法吗？（推荐使用哈希表实现 O(n) 时间复杂度）'
         );


ALTER TABLE `problem`
    ADD COLUMN `test_cases` JSON NOT NULL COMMENT '测试用例数组：[{"input":"输入文本","expected_output":"输出文本","is_sample":1/0,"explanation":"用例说明"}, ...]' AFTER `advanced_requirement`,
DROP COLUMN `input_sample`,
    DROP COLUMN `output_sample`;

-- 修正后可直接执行的UPDATE语句（双重转义\n，保证JSON格式合法）
UPDATE `problem`
SET `test_cases` = '[
  {"input":"4\\n2 7 11 15\\n9","expected_output":"0 1","is_sample":1,"explanation":"nums[0] + nums[1] = 2 + 7 = 9"},
  {"input":"3\\n3 2 4\\n6","expected_output":"1 2","is_sample":1,"explanation":"nums[1] + nums[2] = 2 + 4 = 6"},
  {"input":"2\\n3 3\\n6","expected_output":"0 1","is_sample":1,"explanation":"nums[0] + nums[1] = 3 + 3 = 6"}
]'
WHERE `title_slug` = 'two-sum';