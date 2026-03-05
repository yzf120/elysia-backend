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




-- =============================================
-- 测试数据批量插入 - problem 表
-- =============================================

INSERT INTO `problem` (
    `title`, `title_slug`, `difficulty`, `tags`,
    `description`, `explanation`, `hint`, `constraints`,
    `advanced_requirement`, `test_cases`
) VALUES

-- 1. 有效的括号（简单）
(
    '有效的括号',
    'valid-parentheses',
    '简单',
    '栈,字符串',
    '给定一个只包括 `(`，`)`，`{`，`}`，`[`，`]` 的字符串 s，判断字符串是否有效。\n\n有效字符串需满足：\n1. 左括号必须用相同类型的右括号闭合。\n2. 左括号必须以正确的顺序闭合。\n3. 每个右括号都有一个对应的相同类型的左括号。\n\n### 输入格式\n第一行输入一个字符串 s。\n\n### 输出格式\n输出 true 或 false。',
    '样例1说明：括号按正确顺序闭合\n样例2说明：括号未按正确顺序闭合',
    '可以使用栈来匹配括号。',
    '1 ≤ s.length ≤ 10^4\ns 仅由括号 ()[]{}  组成',
    NULL,
    '[{"input":"()","expected_output":"true","is_sample":1,"explanation":"括号正确闭合"},{"input":"()[]{}","expected_output":"true","is_sample":1,"explanation":"三对括号均正确闭合"},{"input":"(]","expected_output":"false","is_sample":0,"explanation":"括号类型不匹配"},{"input":"([)]","expected_output":"false","is_sample":0,"explanation":"括号顺序错误"},{"input":"{[]}","expected_output":"true","is_sample":0,"explanation":"嵌套括号正确闭合"}]'
),

-- 2. 最大子数组和（中等）
(
    '最大子数组和',
    'maximum-subarray',
    '中等',
    '数组,动态规划,分治',
    '给你一个整数数组 nums，请你找出一个具有最大和的连续子数组（子数组最少包含一个元素），返回其最大和。\n\n### 输入格式\n第一行输入一个整数 n，表示数组长度。\n第二行输入 n 个整数，表示数组 nums。\n\n### 输出格式\n输出一个整数，表示最大子数组和。',
    '样例1说明：子数组 [4,-1,2,1] 的和最大，为 6\n样例2说明：整个数组即为最大子数组',
    '可以使用动态规划，dp[i] 表示以 nums[i] 结尾的最大子数组和。',
    '1 ≤ n ≤ 10^5\n-10^4 ≤ nums[i] ≤ 10^4',
    '进阶：如果你已经实现复杂度为 O(n) 的解法，尝试使用更为精妙的分治法求解。',
    '[{"input":"9\\n-2 1 -3 4 -1 2 1 -5 4\\n","expected_output":"6","is_sample":1,"explanation":"子数组 [4,-1,2,1] 的和最大为 6"},{"input":"1\\n1\\n","expected_output":"1","is_sample":1,"explanation":"只有一个元素"},{"input":"5\\n5 4 -1 7 8\\n","expected_output":"23","is_sample":0,"explanation":"整个数组之和最大"},{"input":"4\\n-2 -1 -3 -4\\n","expected_output":"-1","is_sample":0,"explanation":"全为负数时取最大单个元素"}]'
),

-- 3. 爬楼梯（简单）
(
    '爬楼梯',
    'climbing-stairs',
    '简单',
    '动态规划,记忆化搜索,数学',
    '假设你正在爬楼梯。需要 n 阶你才能到达楼顶。\n\n每次你可以爬 1 或 2 个台阶。你有多少种不同的方法可以爬到楼顶呢？\n\n### 输入格式\n输入一个整数 n，表示楼梯阶数。\n\n### 输出格式\n输出一个整数，表示爬到楼顶的方法总数。',
    '样例1说明：有 2 种方法（1+1 或 2）\n样例2说明：有 3 种方法（1+1+1、1+2 或 2+1）',
    '这是一个斐波那契数列问题，f(n) = f(n-1) + f(n-2)。',
    '1 ≤ n ≤ 45',
    NULL,
    '[{"input":"2","expected_output":"2","is_sample":1,"explanation":"1+1 或 2，共2种"},{"input":"3","expected_output":"3","is_sample":1,"explanation":"1+1+1、1+2 或 2+1，共3种"},{"input":"1","expected_output":"1","is_sample":0,"explanation":"只有1种方法"},{"input":"10","expected_output":"89","is_sample":0,"explanation":"斐波那契第12项"},{"input":"45","expected_output":"1836311903","is_sample":0,"explanation":"边界最大值"}]'
),

-- 4. 二叉树的最大深度（简单）
(
    '二叉树的最大深度',
    'maximum-depth-of-binary-tree',
    '简单',
    '树,深度优先搜索,广度优先搜索,二叉树',
    '给定一个二叉树 root，返回其最大深度。\n\n二叉树的最大深度是指从根节点到最远叶子节点的最长路径上的节点数。\n\n二叉树以层序遍历序列输入，空节点用 -1 表示。\n\n### 输入格式\n第一行输入一个整数 n，表示节点总数（含空节点）。\n第二行输入 n 个整数，表示层序遍历序列（-1 表示空节点）。\n\n### 输出格式\n输出一个整数，表示二叉树的最大深度。',
    '样例1说明：树高为3\n样例2说明：树高为2',
    '可以使用递归（DFS）或队列（BFS）来求解。',
    '树中节点的数量在 [0, 10^4] 范围内\n-100 ≤ Node.val ≤ 100',
    NULL,
    '[{"input":"7\\n3 9 20 -1 -1 15 7","expected_output":"3","is_sample":1,"explanation":"最长路径为 3->20->15 或 3->20->7"},{"input":"2\\n1 -1 2","expected_output":"2","is_sample":1,"explanation":"路径为 1->2"},{"input":"0\\n","expected_output":"0","is_sample":0,"explanation":"空树深度为0"},{"input":"1\\n1","expected_output":"1","is_sample":0,"explanation":"只有根节点"}]'
),

-- 5. 反转链表（简单）
(
    '反转链表',
    'reverse-linked-list',
    '简单',
    '链表,递归',
    '给你单链表的头节点 head，请你反转链表，并返回反转后的链表头节点。\n\n链表以空格分隔的整数序列输入，-1 表示链表结束。\n\n### 输入格式\n第一行输入若干个整数（空格分隔），最后以 -1 结尾，表示链表各节点的值。\n\n### 输出格式\n输出反转后链表各节点的值，空格分隔。',
    '样例1说明：1->2->3->4->5 反转为 5->4->3->2->1\n样例2说明：1->2 反转为 2->1',
    '可以使用迭代或递归两种方式实现。',
    '链表中节点的数目范围是 [0, 5000]\n-5000 ≤ Node.val ≤ 5000',
    '进阶：链表可以选用迭代或递归方式完成反转。你能否用两种方法解决这道题？',
    '[{"input":"1 2 3 4 5 -1","expected_output":"5 4 3 2 1","is_sample":1,"explanation":"链表完全反转"},{"input":"1 2 -1","expected_output":"2 1","is_sample":1,"explanation":"两节点链表反转"},{"input":"-1","expected_output":"","is_sample":0,"explanation":"空链表"},{"input":"1 -1","expected_output":"1","is_sample":0,"explanation":"单节点链表不变"}]'
),

-- 6. 合并两个有序数组（简单）
(
    '合并两个有序数组',
    'merge-sorted-array',
    '简单',
    '数组,双指针,排序',
    '给你两个按非递减顺序排列的整数数组 nums1 和 nums2，另有两个整数 m 和 n，分别表示 nums1 和 nums2 中的元素数目。\n\n请你合并 nums2 到 nums1 中，使合并后的数组同样按非递减顺序排列。\n\n### 输入格式\n第一行输入整数 m，表示 nums1 的有效元素个数。\n第二行输入 m 个整数，表示 nums1。\n第三行输入整数 n，表示 nums2 的元素个数。\n第四行输入 n 个整数，表示 nums2。\n\n### 输出格式\n输出合并后的数组，空格分隔。',
    '样例1说明：合并后为 [1,2,2,3,5,6]\n样例2说明：nums1 为空，结果为 nums2',
    '从后往前合并可以避免覆盖问题，时间复杂度 O(m+n)。',
    '0 ≤ m, n ≤ 200\n1 ≤ m + n ≤ 200\n-10^9 ≤ nums1[i], nums2[j] ≤ 10^9',
    NULL,
    '[{"input":"3\\n1 2 3\\n3\\n2 5 6","expected_output":"1 2 2 3 5 6","is_sample":1,"explanation":"合并两个有序数组"},{"input":"0\\n\\n1\\n1","expected_output":"1","is_sample":1,"explanation":"nums1为空"},{"input":"1\\n2\\n2\\n1 3","expected_output":"1 2 3","is_sample":0,"explanation":"nums2元素插入nums1"},{"input":"3\\n1 2 3\\n0\\n","expected_output":"1 2 3","is_sample":0,"explanation":"nums2为空，结果不变"}]'
),

-- 7. 最长公共前缀（简单）
(
    '最长公共前缀',
    'longest-common-prefix',
    '简单',
    '字符串,字典树',
    '编写一个函数来查找字符串数组中的最长公共前缀。\n\n如果不存在公共前缀，返回空字符串 ""。\n\n### 输入格式\n第一行输入一个整数 n，表示字符串数量。\n接下来 n 行，每行输入一个字符串。\n\n### 输出格式\n输出最长公共前缀字符串，若不存在则输出空行。',
    '样例1说明：三个字符串的公共前缀为 "fl"\n样例2说明：三个字符串无公共前缀',
    '可以纵向扫描，逐列比较每个字符串的同一位置字符。',
    '1 ≤ strs.length ≤ 200\n0 ≤ strs[i].length ≤ 200\nstrs[i] 仅由小写英文字母组成',
    NULL,
    '[{"input":"3\\nflower\\nflow\\nflight","expected_output":"fl","is_sample":1,"explanation":"公共前缀为fl"},{"input":"3\\ndog\\nracecar\\ncar","expected_output":"","is_sample":1,"explanation":"无公共前缀"},{"input":"1\\nabc","expected_output":"abc","is_sample":0,"explanation":"只有一个字符串"},{"input":"2\\nab\\nab","expected_output":"ab","is_sample":0,"explanation":"两个相同字符串"}]'
),

-- 8. 买卖股票的最佳时机（简单）
(
    '买卖股票的最佳时机',
    'best-time-to-buy-and-sell-stock',
    '简单',
    '数组,动态规划',
    '给定一个数组 prices，它的第 i 个元素 prices[i] 表示一支给定股票第 i 天的价格。\n\n你只能选择某一天买入这只股票，并选择在未来的某一个不同的日子卖出该股票。设计一个算法来计算你所能获取的最大利润。\n\n返回你可以从这笔交易中获取的最大利润。如果你不能获取任何利润，返回 0。\n\n### 输入格式\n第一行输入一个整数 n，表示天数。\n第二行输入 n 个整数，表示每天的股票价格。\n\n### 输出格式\n输出一个整数，表示最大利润。',
    '样例1说明：第2天买入（价格=1），第5天卖出（价格=6），利润=5\n样例2说明：价格持续下降，无法获利',
    '记录遍历过程中的最低价格，用当前价格减去最低价格即为当前最大利润。',
    '1 ≤ prices.length ≤ 10^5\n0 ≤ prices[i] ≤ 10^4',
    NULL,
    '[{"input":"6\\n7 1 5 3 6 4","expected_output":"5","is_sample":1,"explanation":"第2天买入价格1，第5天卖出价格6，利润5"},{"input":"5\\n7 6 4 3 1","expected_output":"0","is_sample":1,"explanation":"价格持续下降，不交易利润为0"},{"input":"1\\n5","expected_output":"0","is_sample":0,"explanation":"只有一天无法交易"},{"input":"4\\n2 4 1 7","expected_output":"6","is_sample":0,"explanation":"第3天买入价格1，第4天卖出价格7，利润6"}]'
),

-- 9. 环形链表（简单）
(
    '环形链表',
    'linked-list-cycle',
    '简单',
    '哈希表,链表,双指针',
    '给你一个链表的头节点 head，判断链表中是否有环。\n\n如果链表中有某个节点，可以通过连续跟踪 next 指针再次到达，则链表中存在环。\n\n输入格式：链表节点值序列（空格分隔），以及环的入口位置 pos（-1 表示无环，0 表示尾节点指向头节点，以此类推）。\n\n### 输入格式\n第一行输入若干整数（空格分隔），表示链表节点值，以 -999 结尾（-999 不是节点值）。\n第二行输入整数 pos，表示尾节点指向的节点下标（-1 表示无环）。\n\n### 输出格式\n输出 true 或 false。',
    '样例1说明：尾节点指向下标1的节点，存在环\n样例2说明：pos=-1，无环',
    '使用快慢指针（Floyd 判圈算法），快指针每次走两步，慢指针每次走一步，若相遇则有环。',
    '链表中节点的数目范围是 [0, 10^4]\n-10^5 ≤ Node.val ≤ 10^5\npos 为 -1 或者链表中的一个有效索引',
    '进阶：你能用 O(1)（即常量）内存解决此问题吗？',
    '[{"input":"3 2 0 -4 -999\\n1","expected_output":"true","is_sample":1,"explanation":"尾节点指向下标1，存在环"},{"input":"1 2 -999\\n0","expected_output":"true","is_sample":1,"explanation":"尾节点指向头节点，存在环"},{"input":"1 -999\\n-1","expected_output":"false","is_sample":0,"explanation":"只有一个节点且无环"},{"input":"-999\\n-1","expected_output":"false","is_sample":0,"explanation":"空链表无环"}]'
),

-- 10. 搜索插入位置（简单）
(
    '搜索插入位置',
    'search-insert-position',
    '简单',
    '数组,二分查找',
    '给定一个排序数组和一个目标值，在数组中找到目标值，并返回其索引。如果目标值不存在于数组中，返回它将会被按顺序插入的位置。\n\n请你使用时间复杂度为 O(log n) 的算法。\n\n### 输入格式\n第一行输入一个整数 n，表示数组长度。\n第二行输入 n 个整数，表示有序数组 nums。\n第三行输入目标值 target。\n\n### 输出格式\n输出一个整数，表示目标值的索引或插入位置。',
    '样例1说明：目标值5在数组中，下标为2\n样例2说明：目标值2不在数组中，应插入下标1',
    '使用二分查找，当 left > right 时，left 即为插入位置。',
    '1 ≤ nums.length ≤ 10^4\n-10^4 ≤ nums[i] ≤ 10^4\nnums 为无重复元素的升序排列数组\n-10^4 ≤ target ≤ 10^4',
    NULL,
    '[{"input":"4\\n1 3 5 6\\n5","expected_output":"2","is_sample":1,"explanation":"目标值5在数组中，下标为2"},{"input":"4\\n1 3 5 6\\n2","expected_output":"1","is_sample":1,"explanation":"目标值2应插入下标1"},{"input":"4\\n1 3 5 6\\n7","expected_output":"4","is_sample":0,"explanation":"目标值7应插入末尾"},{"input":"4\\n1 3 5 6\\n0","expected_output":"0","is_sample":0,"explanation":"目标值0应插入开头"}]'
);

-- =============================================
-- 新增时间限制和内存限制字段
-- time_limit: 单位 ms，默认 1000ms
-- memory_limit: 单位 MB，默认 256MB
-- =============================================
ALTER TABLE `problem`
    ADD COLUMN `time_limit`   INT NOT NULL DEFAULT 1000 COMMENT '时间限制（单位：ms，默认1000ms）' AFTER `test_cases`,
    ADD COLUMN `memory_limit` INT NOT NULL DEFAULT 256  COMMENT '内存限制（单位：MB，默认256MB）'  AFTER `time_limit`;

-- 将已有数据统一设置为默认值
UPDATE `problem` SET `time_limit` = 1000, `memory_limit` = 256;
