package consts

// 用户角色常量
const (
	// 角色类型
	RoleStudent = "student" // 学生
	RoleTeacher = "teacher" // 教师
	RoleAdmin   = "admin"   // 管理员
	RoleSuperAdmin = "super_admin" // 超级管理员
)

// 用户状态常量
const (
	UserStatusPending  = 1 // 待审批
	UserStatusActive   = 2 // 可用
	UserStatusRejected = 3 // 驳回
	UserStatusBanned   = 4 // 封号
)

// 教师认证状态常量
const (
	TeacherVerificationPending  = 0 // 待审核
	TeacherVerificationApproved = 1 // 已通过
	TeacherVerificationRejected = 2 // 已驳回
)

// 教师状态常量
const (
	TeacherStatusInactive = 0 // 未激活
	TeacherStatusActive   = 1 // 正常
	TeacherStatusDisabled = 2 // 禁用
)

// 学生状态常量
const (
	StudentStatusDisabled = 0 // 禁用
	StudentStatusActive   = 1 // 正常
)

// 班级状态常量
const (
	ClassStatusEnded    = 0 // 已结束
	ClassStatusOngoing  = 1 // 进行中
	ClassStatusArchived = 2 // 已归档
)

// 班级成员状态常量
const (
	ClassMemberStatusLeft   = 0 // 已退出
	ClassMemberStatusActive = 1 // 正常
)

// 编程水平常量
const (
	ProgrammingLevelBeginner     = "beginner"     // 初学者
	ProgrammingLevelIntermediate = "intermediate" // 中级
	ProgrammingLevelAdvanced     = "advanced"     // 高级
)

// 注册来源常量
const (
	RegisterSourcePhone      = "phone"      // 手机号注册
	RegisterSourceMiniProgram = "miniprogram" // 小程序注册
	RegisterSourceManual     = "manual"     // 手动创建
	RegisterSourceSMS        = "sms"        // 短信验证码注册
	RegisterSourceTeacher    = "teacher"    // 教师注册
)
