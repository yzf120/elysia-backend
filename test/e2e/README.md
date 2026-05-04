# 端到端测试用例

## 目录结构

按任务书功能点分类，共22个功能模块、70个测试用例。

| 目录 | 功能点 | 任务书对应 | 用例数 |
|------|--------|-----------|--------|
| F01_student_auth | 学生注册与登录 | 3a | 5 |
| F02_teacher_auth | 教师注册与登录 | 3a | 3 |
| F03_problem_management | 题目管理 | 3a | 5 |
| F04_code_run | 代码运行与判题 | 3a | 7 |
| F05_class_course | 班级与课程管理 | 3a | 5 |
| F06_material | 学习资料管理 | 3a | 2 |
| F07_ai_chat | AI文字对话交互 | 3b | 5 |
| F08_ai_image | AI图片交互 | 3b | 2 |
| F09_code_awareness | 源代码感知 | 3b | 2 |
| F10_compile_error_awareness | 编译器报错感知 | 3b | 2 |
| F11_testcase_awareness | 未通过测试用例感知 | 3b | 2 |
| F12_intent_recognition | 意图识别与路由 | 3b | 6 |
| F13_heuristic_answer | 启发式回答 | 3b | 1 |
| F14_rag | RAG检索增强 | 3b | 1 |
| F15_exercise_recommend | 练习题推荐 | 3b | 1 |
| F16_session_history | 会话历史记录 | 3c | 3 |
| F17_favorite | 会话收藏 | 3c | 4 |
| F18_coding_profile | 学生做题画像 | 3c | 3 |
| F19_qa_profile | 问答行为画像 | 3c | 2 |
| F20_study_stats | 学习统计数据 | 3c | 2 |
| F21_content_security | 内容安全审核 | 系统保障 | 4 |
| F22_admin_management | 管理员后台管理 | 系统管理 | 4 |

## 使用方法

1. 使用 IntelliJ IDEA / VS Code REST Client 插件打开 `.http` 文件
2. 选择 `dev` 环境（基于 `http-client.env.json`）
3. 按顺序执行测试用例
4. 部分测试用例有依赖关系（如需要先登录获取token），请按文件内注释的顺序执行

## 前置条件

- 后端服务运行在 `localhost:8080`
- MySQL、Redis、MongoDB 等依赖服务已启动
- chat-agent、llm-tool、session、access 等微服务已启动
