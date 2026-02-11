[//]: # (# 科目管理系统重构说明)

[//]: # ()
[//]: # (## 概述)

[//]: # ()
[//]: # (本次重构引入了独立的科目表（Subject），重新设计了教师、科目和班级之间的关联关系，使系统更加灵活和规范。)

[//]: # ()
[//]: # (## 数据模型设计)

[//]: # ()
[//]: # (### 1. 核心表结构)

[//]: # ()
[//]: # (#### subjects（科目表）)

[//]: # (- **subject_id**: 科目唯一标识)

[//]: # (- **subject_name**: 科目名称)

[//]: # (- **subject_code**: 科目代码（唯一）)

[//]: # (- **category**: 科目分类（理科、文科、艺术等）)

[//]: # (- **description**: 科目描述)

[//]: # (- **credits**: 学分)

[//]: # (- **status**: 状态（1-启用，0-禁用）)

[//]: # ()
[//]: # (#### teacher_subjects（教师-科目关联表）)

[//]: # (- **teacher_id**: 教师ID)

[//]: # (- **subject_id**: 科目ID)

[//]: # (- **start_date**: 开始教授日期)

[//]: # (- **end_date**: 结束教授日期)

[//]: # (- **status**: 状态（1-在教，0-已停止）)

[//]: # (- **remark**: 备注)

[//]: # ()
[//]: # (### 2. 修改的表)

[//]: # ()
[//]: # (#### teachers（教师表）)

[//]: # (- **删除**: teaching_subjects 字段（改用关联表管理）)

[//]: # ()
[//]: # (#### classes（班级表）)

[//]: # (- **修改**: subject 字段改为 subject_id（关联到科目表）)

[//]: # (- **新增**: 索引 idx_teacher, idx_subject)

[//]: # ()
[//]: # (## 业务关系)

[//]: # ()
[//]: # (```)

[//]: # (教师 &#40;Teacher&#41; 1:N 教师-科目关联 &#40;TeacherSubject&#41; N:1 科目 &#40;Subject&#41;)

[//]: # (教师 &#40;Teacher&#41; 1:N 班级 &#40;Class&#41; N:1 科目 &#40;Subject&#41;)

[//]: # (```)

[//]: # ()
[//]: # (### 关系说明)

[//]: # (1. **教师可以开设多个科目**：通过 teacher_subjects 关联表实现)

[//]: # (2. **一个班级对应一个科目**：classes.subject_id 关联到 subjects.subject_id)

[//]: # (3. **教师可以同时开设多个班级**：一个教师可以创建多个班级)

[//]: # (4. **教师只能开设已分配的科目的班级**：创建班级时会验证教师是否有权限教授该科目)

[//]: # ()
[//]: # (## 代码结构)

[//]: # ()
[//]: # (### Model层)

[//]: # (- `/model/subject/subject.go`: 科目和教师-科目关联模型)

[//]: # (- `/model/class/class.go`: 班级模型（已修改）)

[//]: # (- `/model/teacher/teacher.go`: 教师模型（已修改）)

[//]: # ()
[//]: # (### DAO层)

[//]: # (- `/dao/subject_dao.go`: 科目数据访问层)

[//]: # (- `/dao/teacher_subject_dao.go`: 教师-科目关联数据访问层)

[//]: # (- `/dao/class_dao.go`: 班级数据访问层（已增强）)

[//]: # (- `/dao/teacher_dao.go`: 教师数据访问层)

[//]: # ()
[//]: # (### Service层)

[//]: # (- `/service/subject_service.go`: 科目业务逻辑)

[//]: # (- `/service/teacher_subject_service.go`: 教师-科目关联业务逻辑)

[//]: # (- `/service/class_service.go`: 班级业务逻辑（已修改）)

[//]: # ()
[//]: # (## 主要功能)

[//]: # ()
[//]: # (### 科目管理（SubjectService）)

[//]: # (- `CreateSubject`: 创建科目)

[//]: # (- `GetSubjectById`: 根据ID获取科目)

[//]: # (- `GetSubjectByCode`: 根据代码获取科目)

[//]: # (- `UpdateSubject`: 更新科目信息)

[//]: # (- `DeleteSubject`: 删除科目)

[//]: # (- `ListSubjects`: 查询科目列表)

[//]: # (- `EnableSubject`: 启用科目)

[//]: # (- `DisableSubject`: 禁用科目)

[//]: # ()
[//]: # (### 教师-科目关联管理（TeacherSubjectService）)

[//]: # (- `AssignSubjectToTeacher`: 为教师分配科目)

[//]: # (- `RemoveSubjectFromTeacher`: 移除教师的科目)

[//]: # (- `UpdateTeacherSubject`: 更新关联信息)

[//]: # (- `GetTeacherSubjects`: 获取教师的所有科目)

[//]: # (- `GetSubjectTeachers`: 获取科目的所有教师)

[//]: # (- `ListTeacherSubjectRelations`: 查询关联列表)

[//]: # (- `StopTeachingSubject`: 停止教授某科目)

[//]: # (- `ResumeTeachingSubject`: 恢复教授某科目)

[//]: # ()
[//]: # (### 班级管理（ClassService）修改)

[//]: # (- `CreateClass`: 创建班级时需要传入 subjectId，并验证教师是否有权限教授该科目)

[//]: # ()
[//]: # (## 数据库迁移)

[//]: # ()
[//]: # (执行迁移脚本：)

[//]: # (```bash)

[//]: # (mysql -u username -p database_name < migrations/subjects.sql)

[//]: # (```)

[//]: # ()
[//]: # (### 迁移步骤)

[//]: # (1. 创建 subjects 表)

[//]: # (2. 创建 teacher_subjects 表)

[//]: # (3. 修改 teachers 表（删除 teaching_subjects 字段）)

[//]: # (4. 修改 classes 表（添加 subject_id 字段和索引）)

[//]: # (5. 插入示例科目数据（可选）)

[//]: # ()
[//]: # (### 数据迁移注意事项)

[//]: # (如果现有 classes 表中已有数据，需要手动迁移：)

[//]: # (1. 在 subjects 表中创建对应的科目记录)

[//]: # (2. 更新 classes 表，将 subject 名称映射到 subject_id)

[//]: # (3. 确认数据迁移完成后，可删除 classes 表的 subject 字段)

[//]: # ()
[//]: # (示例SQL：)

[//]: # (```sql)

[//]: # (-- 将现有班级的科目名称映射到科目ID)

[//]: # (UPDATE classes c )

[//]: # (INNER JOIN subjects s ON c.subject = s.subject_name )

[//]: # (SET c.subject_id = s.subject_id )

[//]: # (WHERE c.subject_id IS NULL;)

[//]: # ()
[//]: # (-- 确认迁移完成后删除旧字段)

[//]: # (ALTER TABLE classes DROP COLUMN subject;)

[//]: # (```)

[//]: # ()
[//]: # (## 使用示例)

[//]: # ()
[//]: # (### 1. 创建科目)

[//]: # (```go)

[//]: # (subjectService := service.NewSubjectService&#40;&#41;)

[//]: # (subject, err := subjectService.CreateSubject&#40;)

[//]: # (    "高等数学",      // 科目名称)

[//]: # (    "MATH101",      // 科目代码)

[//]: # (    "理科",         // 分类)

[//]: # (    "高等数学基础课程", // 描述)

[//]: # (    4,              // 学分)

[//]: # (&#41;)

[//]: # (```)

[//]: # ()
[//]: # (### 2. 为教师分配科目)

[//]: # (```go)

[//]: # (teacherSubjectService := service.NewTeacherSubjectService&#40;&#41;)

[//]: # (err := teacherSubjectService.AssignSubjectToTeacher&#40;)

[//]: # (    "teacher_123",           // 教师ID)

[//]: # (    "subj_math_001",        // 科目ID)

[//]: # (    time.Now&#40;&#41;,             // 开始日期)

[//]: # (    "主讲教师",              // 备注)

[//]: # (&#41;)

[//]: # (```)

[//]: # ()
[//]: # (### 3. 创建班级（需要科目ID）)

[//]: # (```go)

[//]: # (classService := service.NewClassService&#40;&#41;)

[//]: # (class, err := classService.CreateClass&#40;)

[//]: # (    "teacher_123",          // 教师ID)

[//]: # (    "高数A班",              // 班级名称)

[//]: # (    "subj_math_001",       // 科目ID（不再是科目名称）)

[//]: # (    "2024春季",            // 学期)

[//]: # (    "面向计算机专业",       // 描述)

[//]: # (    100,                   // 最大学生数)

[//]: # (&#41;)

[//]: # (```)

[//]: # ()
[//]: # (### 4. 查询教师的所有科目)

[//]: # (```go)

[//]: # (teacherSubjectService := service.NewTeacherSubjectService&#40;&#41;)

[//]: # (subjects, err := teacherSubjectService.GetTeacherSubjects&#40;"teacher_123"&#41;)

[//]: # (```)

[//]: # ()
[//]: # (### 5. 查询某科目的所有教师)

[//]: # (```go)

[//]: # (teacherSubjectService := service.NewTeacherSubjectService&#40;&#41;)

[//]: # (teacherIds, err := teacherSubjectService.GetSubjectTeachers&#40;"subj_math_001"&#41;)

[//]: # (```)

[//]: # ()
[//]: # (## API接口建议)

[//]: # ()
[//]: # (基于新的数据模型，建议添加以下API接口：)

[//]: # ()
[//]: # (### 科目管理接口)

[//]: # (- `POST /api/subjects` - 创建科目)

[//]: # (- `GET /api/subjects/:id` - 获取科目详情)

[//]: # (- `PUT /api/subjects/:id` - 更新科目信息)

[//]: # (- `DELETE /api/subjects/:id` - 删除科目)

[//]: # (- `GET /api/subjects` - 查询科目列表)

[//]: # ()
[//]: # (### 教师-科目关联接口)

[//]: # (- `POST /api/teachers/:teacherId/subjects` - 为教师分配科目)

[//]: # (- `DELETE /api/teachers/:teacherId/subjects/:subjectId` - 移除教师的科目)

[//]: # (- `GET /api/teachers/:teacherId/subjects` - 获取教师的所有科目)

[//]: # (- `GET /api/subjects/:subjectId/teachers` - 获取科目的所有教师)

[//]: # (- `PUT /api/teacher-subjects/:id/stop` - 停止教授某科目)

[//]: # (- `PUT /api/teacher-subjects/:id/resume` - 恢复教授某科目)

[//]: # ()
[//]: # (### 班级管理接口（修改）)

[//]: # (- `POST /api/classes` - 创建班级（参数改为 subjectId）)

[//]: # (- `GET /api/classes/by-subject/:subjectId` - 按科目查询班级)

[//]: # (- `GET /api/classes/by-teacher-subject` - 按教师和科目查询班级)

[//]: # ()
[//]: # (## 优势)

[//]: # ()
[//]: # (1. **数据规范化**：科目作为独立实体，避免数据冗余)

[//]: # (2. **灵活的权限控制**：通过关联表精确控制教师的科目权限)

[//]: # (3. **便于统计分析**：可以方便地统计每个科目的班级数、教师数等)

[//]: # (4. **支持科目生命周期管理**：可以记录教师教授科目的起止时间)

[//]: # (5. **扩展性强**：未来可以轻松添加科目相关的其他属性和关系)

[//]: # ()
[//]: # (## 注意事项)

[//]: # ()
[//]: # (1. 创建班级前，必须先为教师分配相应的科目权限)

[//]: # (2. 删除科目前，需要检查是否有关联的班级或教师)

[//]: # (3. 教师停止教授某科目后，不影响已创建的班级，但不能创建新班级)

[//]: # (4. 建议定期清理已结束的教师-科目关联记录)

[//]: # ()
[//]: # (## 后续优化建议)

[//]: # ()
[//]: # (1. 添加科目的前置课程关系)

[//]: # (2. 添加科目的学时、考核方式等详细信息)

[//]: # (3. 支持科目的版本管理（课程大纲变更）)

[//]: # (4. 添加科目的资源管理（教材、课件等）)

[//]: # (5. 支持跨学期的科目连续性管理)
