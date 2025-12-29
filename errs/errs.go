package errs

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const (
	SuccessCode = 0
)

const (
	ErrOperationPending = 202 // 操作正在进行中, 需要等待
)

const (
	ErrUserOccupied               = "正在处理您的其他请求，请稍后再试"
	ErrCharacterUserOccupied      = "跟我聊天的人太火爆了，请稍等片刻再跟我发消息"
	ErrModelNotAvailable          = "模型服务正忙……"
	ErrTooManyRequests            = "当前请求人数过多，请稍后再试"
	ErrAnonTooManyRequests        = "当前请求人数过多，可以注册登陆，享受更好的服务"
	ErrTooManyConcurrentConv      = "您进行中的对话过多，请稍后再试"
	ErrPromptExceedLimit          = "您输入的内容过长，请缩减后输入"
	ErrContinueError              = "当前请求不支持继续生成"
	ErrIncorrectRequest           = "请求参数有误"
	ErrInvalidRequest             = "无效的请求"
	ErrText2ImageFineTuningOption = "精调图片参数有误"
	ErrConvsExceedLimit           = "会话过长，请开始新会话"
	ErrConvSensitive              = "此会话因敏感问题终止，请开始新会话"
	ErrServerAbnormal             = "服务异常，请刷新网页重试"
	ErrInternalServer             = "服务内部错误"
	ErrSpeechesExceedLimit        = "重新生成次数超过上限，请提新问题"
	ErrGetQuotaItem               = "查询模型请求限额出错"
	ErrAcquireQuota               = "获取模型请求限额出错"
	ErrInvokeModelFrequently      = "请求模型频繁，请稍后再试"
	ErrInvokeNetworkSlacking      = "网络开小差了，请重试"
	ErrNotEnoughQuota             = "今日生成额度已用完，明天再来试试吧～"
	ErrNotEnoughQuotaV2           = "今天生成视频的次数已达上限，明天再来试试吧～"

	ErrMessagesExceedLimit    = "会话内容超过40, 按对话时间从旧到新在数组中排列"
	ErrRequestIsNone          = "请求参数不能为空"
	ErrRoleAbnormal           = "角色, 'system', 'user'或者'assistant', 除system外在message中必须是user与assistant交替(一问一答),最后一个为user提问, 且content不能为空"
	ErrModelRequest           = "模型名称有误，请选择指定的模型名称"
	ErrorFileTooLarge         = "文件过大，请上传小于10M的文件"
	ErrorInvalidUrl           = "url非法"
	ErrorUserQuotaLimit       = "今日提问次数已达上限，请明日再试"
	ErrorAnonUserQuotaLimit   = "提问次数已达上限，请登录后再试"
	ErrorAgentDayQuotaLimit   = "你今天的任务余额次数已用完，请明天再试"
	ErrorAgentUserSubmitLimit = "你已提交多项任务，请等待任一项任务完成后再提交新任务吧"
	ErrorNoPermissionForAgent = "该智能体已被设置为私密，无法发送消息"
	ErrorAgentDeleted         = "系统繁忙，请稍后再试"
	ErrorAgentAuditFailed     = "审核失败，请重新编辑后再次使用"

	ErrAlreadyEnterScenario = "当前已是新剧情"

	ErrImageInternalServer   = "生成失败，重新生成"
	ErrImageInputSecurity    = "上传失败，重新上传"
	ErrImageFaceDetection    = "未检测到人脸，重新上传"
	ErrImageSearchPluginFail = "未找到相关图片。"

	ErrSecurityCheck      = "抱歉，我无法回答这个问题，让我们换个话题再聊聊吧。"
	ErrSecurityUserBanned = "账号异常，已限制使用。"

	ErrShareExpired     = "分享链接已过期"
	ErrMessageDelFailed = "提问编辑失败，请重试"

	ErrWxImagePromptSensitive = "生成失败，换个描述试试吧"

	ErrGuoKaoUploadMsg = "尊敬的用户，为保障国家公务员考试的公平性，根据相关政策要求，此功能在国考时段无法使用。\n国考时段：11月29日下午14:00—16:00、11月30日上午 9:00—11:00、11月30日下午14:00—17:00"
)

// 插件相关错误
const (
	ErrBusinessDefault = iota + 10000
	ErrPluginCalNotHitKeyword
	ErrPluginCalNotImplemented
	ErrPluginBrowsingNoResult        // 搜索插件返回结果为空/拒答
	ErrPluginCredibleNoResult        // 可信模型插件返回结果为空/拒答
	ErrPluginIdeNoResult             // 代码插件返回结果为空/拒答
	ErrPluginQuotaLimit              // 插件限流
	ErrPluginImageChatModelEmptyResp // 图生文模型空回复
	ErrPromptSensitive               // 输入敏感
	ErrPromptIntentNotMatch          // 输入意图不匹配
	ErrQuestionLowGrade              // 问题质量低
	ErrPluginBrowsingMainModelIntent // 搜索插件主模型意图
)

// 服务工程错误 登陆相关
const (
	ErrAuthenDefault                     = iota + 20000 // 其他鉴权失败
	ErrAuthenTokenInvalid                               // token无效
	ErrAuthenTokenExpired                               // token过期
	ErrAuthenUserNotRegistered                          // 用户未注册
	ErrAuthenUserForbidden                              // 用户被封号
	ErrAuthenUserUnqualified                            // 用户未获得使用条件
	ErrPhoneLoginVerificationCodeInvalid                // 手机号登录-验证码无效
	ErrPhoneLoginVerificationCodeExpired                // 手机号登录-验证码过期
	ErrServerInnerError                                 // 服务内部错误
	ErrParamMissError                                   // 必填参数缺少
	ErrFileAbnormal                                     // 文档异常，请检查后重新上传
	ErrGetOptimizeSuggestionFailed                      // 获取优化建议失败
	ErrNilConv                                          // 空会话
	ErrErrorUid                                         // 用户id错误
	ErrInConsistentPhone                                // 用户手机号非法
	ErrBadRequest
	ErrEmptyContentFile                   // 文件内容为空
	ErrFileContentInvalid                 // 文件内容不合法
	ErrFileParseFailed                    // 解析失败
	ErrFileAcquireQuotaFailed             // 获取配额失败
	ErrPhoneUsed                          // 手机号已被占用
	ErrAnonRequestForbidden               // 游客请求禁止。需要重新请求
	ErrTencentDocBindCodeError            // 腾讯文档帐号绑定失效
	ErrTencentDocUnBindCodeError          // 腾讯文档帐号解绑失败
	ErrTencentDocGetBindInfoNullCodeError // 腾讯文档帐号绑定信息为空（ps：该错误码端有特殊逻辑，不可修改，原值：20024）
	ErrTencentDocGetBindInfoCodeError     // 腾讯文档获取帐号绑定信息失败
	ErrTencentDocGetDocListCodeError      // 腾讯文档获取文档列表失败
	ErrTencentDocGetDocDetailCodeError    // 腾讯文档获取文档信息失败
	ErrTencentDocEditDocCodeError
	ErrTencentDocDeleteDocCodeError                  // 腾讯文档删除失败
	ErrTencentDocGetFolderIdCodeError                // 获取文件夹失败
	ErrTencentDocDeleteScopeInvalidCodeError         // 缺少腾讯文档删除权限
	ErrOfflineActivityChannelDelCodeError            // 线下活动渠道被删除
	ErrOfflineActivityTokenInvalid                   // 活动渠道token失效
	ErrAppleLoginPhoneNotBind                        // 苹果登录-未绑定手机
	ErrSetParentModelOverLimit                       // 设置家长模式超过频率限制
	ErrTencentDocUnauthorized                        // 腾讯文档帐号未授权（ps：该错误码端有特殊逻辑，不可修改，原值：20036）
	ErrLoginStatusRedis                      = 20051 // 鉴权redis错误，这种场景api返回504，不退出登录
)

// 腾讯文档错误码
const (
	ErrTencentDocPathNotFound       = 10131
	ErrTencentDocDeleteScopeInvalid = 10017
	ErrTencentDocCapacityLimit      = 10059
)

// 会话管理相关
const (
	ErrConversationDeleted          = 21000 //  会话删除
	ErrConversationSensitive        = 21001 //  会话敏感
	ErrConversationTitleSensitive   = 21002 // 更改会话标题，标题敏感
	ErrConversationTooMuch          = 21003 // 请求超过次数
	ErrConversationDeviceTooMuch    = 21004 // 单个用户同时请求量超过限制
	ErrConversationUserBanned       = 21005 // 用户被封禁
	ErrConversationMessageDelFailed = 21006 // 消息删除失败
)

// 服务工程错误 智能体相关
const (
	ErrUnsupportedAppVersion = 22001 // app版本不支持
)

const (
	ErrTokenForcedExpiration = 23000 // token强制过期 踢登
)

// 模型相关错误
const (
	ErrInternalDefault                    = iota + 30000 // 内部错误
	ErrInternalModelFailed                               // 主模型请求失败
	ErrInternalPluginDrawFailed                          // 文生图插件服务失败
	ErrInternalPluginImageChatFailed                     // 图生文插件服务失败
	ErrInternalPluginSecurityFailed                      // 安全审核服务失败
	ErrInternalCredibleModelFailed                       // 可信模型服务失败/超时
	ErrInternalPluginPdfFailed                           // PDF插件服务失败
	ErrInternalModelInferFailed                          // 主模型推理失败
	ErrInternalOutputCheckFailed                         // 输出检查失败
	ErrExternalModelFailedDefault                        // 默认外部模型比对请求失败
	ErrErnieBotModelFailed                               // 文心一言请求失败
	ErrGPT35ModelFailed                                  // GPT3.5请求失败
	ErrGPT4ModelFailed                                   // GPT4请求失败
	ErrSparkDeskModelFailed                              // 星火模型请求失败
	ErrDoubaoModelFailed                                 // 豆包模型请求失败
	ErrInternalPluginCredibleSearchFailed                // 可信搜索插件失败
	ErrInternalThirdPluginFailed                         // 第三方插件失败
	ErrInternalPluginDocQAFailed                         // 文档解析插件服务失败
	ErrNewMasterAgentInitFailed                          // 新闻哥智能体服务失败
	ErrNewMasterAgentRemoteNewRequestFailed
	ErrNewMasterAgentRemoteRequestFailed    // 新闻哥远程请求失败
	ErrNewMasterAgentRemoteRespStatusFailed // 新闻哥远程响应状态码错误
	ErrNewMasterAgentRemoteRespReadFailed   // 新闻哥远程响应数据流读取失败
	ErrNewMasterAgentRemoteTimeout          // 新闻哥远程api超时
	ErrNewMasterAgentRemoteEventUnOrder     // 新闻哥远程api数据流顺序错误
	ErrNewMasterAgentRemoteRespDataError    // 新闻哥远程api数据格式错误
	ErrInternalMongodbFindFailed            // mongodb查询失败
	ErrorInternalOpenAPIFailed              // OpenAPI调用失败
	ErrInternalMongodbUpdateFailed          // mongodb更新失败
	ErrInternalMongodbInsertFailed          // mongodb插入失败
	ErrInternalPluginImageHelperFailed      // 文生图插件服务失败
	ErrInternalInfoSecuritySugFailed        // sug安全审核失败或不通过
)

const (
	ErrThirdPluginDefault        = iota + 40000 // 第三方插件失败
	ErrThirdPluginParamError                    // 参数检查错误
	ErrThirdPluginInnerError                    // 内部处理错误
	ErrInconsistentCoolingPeriod                // 冷静期边界错误，即：用户登录时收到冷静期提示、继续操作又过了冷静期
)

const (
	ErrParamCheck                           = iota + 50000 // 参数检查错误
	ErrParseURL                                            // 解析URL错误
	ErrConvertResourceURL                                  // 转换资源URL错误
	ErrImageSensitive                                      // 图片敏感
	ErrInternal                                            // 内部错误
	ErrPassonURL                                           // 透传URL错误
	ErrReviseInputSensitive                                // 修正输入敏感
	ErrOutputImageSensitive                                // 输出图片敏感
	ErrImagePromptSensitive                                // 图片提示敏感
	ErrExternalException                                   // 第三方接口调用异常
	ErrExternalSugFail                                     // sug接口调用返回错误
	ErrExternalSdTaskFail                                  // sd task防作弊安全检测接口返回失败
	ErrSessionTitleEmptyModelResp                          // session title模型服务错误，用prompt兜底
	ErrSessionTitleFail                                    // session title内部错误(db查询或写入错误)，title更新失败
	ErrSessionTitlePush                                    // session title生成成功，但推送端侧失败
	ErrIdCardNumberInvalid                                 // 身份证号码不合法
	ErrGenerateShareImageFail                              // 生成分享图内容过长
	ErrGenerateShareLenLimit                               // 生成分享图轮数限制
	ErrSpeechNotFound                                      // 消息已删除
	ErrAsrVoiceShareFileNotAllow                           // asr音频文件用户未授权
	ErrAsrVoiceShareFileInnerError                         // asr音频文件获取内部错误
	ErrGeneralShareVoiceRecordUidIsEmpty                   // 录音笔通用分享的用户ID为空
	ErrGeneralShareVoiceRecordUidIsNotMatch                // 录音笔通用分享的用户ID不匹配
)

const (
	ErrLoginSafeHit          = 57000
	ErrRegisterSafeHit       = 57001
	ErrLoginSafeHitVerify    = 57002
	ErrLoginSafeVerifyExpire = 57003
)

const (
	ErrResearchAssistantRemoteSvrFail = 600001 // ai研究助手调用qb研究记录服务返回错误
)

// ErrorMessages 错误码，对应的展示信息
var ErrorMessages = map[int]string{
	ErrAuthenDefault:                        "未知错误",
	ErrAuthenTokenInvalid:                   "token无效",
	ErrAuthenTokenExpired:                   "登录态已过期，请重新登录",
	ErrAuthenUserNotRegistered:              "用户未注册",
	ErrAuthenUserForbidden:                  "用户被封号",
	ErrAuthenUserUnqualified:                "用户未获得使用条件",
	ErrPhoneLoginVerificationCodeInvalid:    "验证码错误",
	ErrPhoneLoginVerificationCodeExpired:    "验证码过期",
	ErrServerInnerError:                     "服务内部错误",
	ErrParamMissError:                       "参数缺失",
	ErrFileAbnormal:                         "文档异常，请检查后重新上传。",
	ErrGetOptimizeSuggestionFailed:          "获取优化建议失败。",
	ErrNilConv:                              "会话不存在。",
	ErrErrorUid:                             "用户ID不匹配。",
	ErrBadRequest:                           "请求参数有误",
	ErrEmptyContentFile:                     "文件内容为空",
	ErrFileContentInvalid:                   "文件解析失败",
	ErrFileParseFailed:                      "文件解析失败",
	ErrFileAcquireQuotaFailed:               "请求频繁，请稍后再试",
	ErrConversationDeleted:                  "换个话题再聊吧",
	ErrConversationSensitive:                "会话敏感",
	ErrTencentDocBindCodeError:              "绑定失败，请稍后再试",
	ErrExternalSdTaskFail:                   "账号安全检测失败",
	ErrGeneralShareVoiceRecordUidIsNotMatch: "通用分享的用户ID不匹配",
	ErrGeneralShareVoiceRecordUidIsEmpty:    "通用分享的用户ID为空",
}

// BatchHistoryListClear 错误码定义
const (
	BatchHistoryListClearSuccess          = 0
	BatchHistoryListClearParaError        = 10000
	BatchHistoryListClearUidIsEmpty       = 100001
	BatchHistoryListClearFailed           = 100002
	BatchHistoryListClearDecompressFailed = 100003
)

const (
	ErrMoveConvProjectDeleted = 300001
	ErrBatchMoveConvFailed    = 300002
)

// CommonError 通用错误信息
type CommonError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Msg     string `json:"msg,omitempty"` // deprecated
	RealErr error  `json:"realErr,omitempty"`
}

type CodeError struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

type CommonErrorData[T any] struct {
	CodeError
	Data *T `json:"data,omitempty"`
}

type VerifyCodeError struct {
	CommonError
	SafeVerifyTypes []string `json:"safeVerifyTypes,omitempty"`
	SafeVerifyCode  string   `json:"safeVerifyCode,omitempty"`
}

type CommonResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error CommonError `json:"error,omitempty"`
}

type CommonErrorResponse struct {
	Error CommonError `json:"error,omitempty"`
}

type VerifyCodeErrorResponse struct {
	Error           CommonError `json:"error,omitempty"`
	SafeVerifyTypes []string    `json:"safeVerifyTypes,omitempty"`
	SafeVerifyCode  string      `json:"safeVerifyCode,omitempty"`
	SafeVerifyUrl   string      `json:"safeVerifyUrl,omitempty"`
	WaterProofAppid string      `json:"waterProofAppid,omitempty"`
}

type TransferErrorResponse struct {
	Error        CommonError `json:"error,omitempty"`
	TransferData string      `json:"transferData,omitempty"` //  鉴权透传data
}

func NewAuthErrorResponse(code int, msg, transferData string) string {
	errMsg := msg
	// 优先使用这里预设的错误码
	if codeMsg, isExist := ErrorMessages[code]; isExist {
		errMsg = codeMsg
	}
	err := TransferErrorResponse{
		Error: CommonError{
			Code:    strconv.Itoa(code),
			Message: errMsg,
		},
		TransferData: transferData,
	}
	jsonStr, _ := json.Marshal(&err)
	return string(jsonStr)
}

func NewErrorResponse(code int, msg string) string {
	errMsg := msg
	// 优先使用这里预设的错误码
	if codeMsg, isExist := ErrorMessages[code]; isExist {
		errMsg = codeMsg
	}
	err := CommonErrorResponse{
		Error: CommonError{
			Code:    strconv.Itoa(code),
			Message: errMsg,
		},
	}
	jsonStr, _ := json.Marshal(&err)
	return string(jsonStr)
}

func NewCommonError(code int, msg string) *CommonError {
	return &CommonError{
		Code:    strconv.Itoa(code),
		Message: msg,
		Msg:     msg,
	}
}

func NewCommonErrorData[T any](code int, msg string, data *T) *CommonErrorData[T] {
	err := CommonErrorData[T]{
		Data: data,
		CodeError: CodeError{
			Code:    code,
			Message: msg,
		},
	}
	return &err
}

// DefaultInnerServerError 缺省内部错误控制
var DefaultInnerServerError = CommonError{
	Code:    "500",
	Message: "服务开小差拉，可以稍后再试下",
}

// DefaultBadRequestError ...
var DefaultBadRequestError = CommonError{
	Code:    "400",
	Message: "请求参数有误",
}

// DefaultInnerServerCodeError 缺省内部错误控制
var DefaultInnerServerCodeError = CodeError{
	Code:    500,
	Message: "服务开小差拉，可以稍后再试下",
}

// DefaultBadRequestCodeError ...
var DefaultBadRequestCodeError = CodeError{
	Code:    400,
	Message: "请求参数有误",
}

// GetBadRequestErr ...
func GetBadRequestErr(msg string) CommonError {
	return CommonError{
		Code:    "400",
		Message: msg,
	}
}

// GetInternalErr ...
func GetInternalErr(msg string) CommonError {
	return CommonError{
		Code:    "500",
		Message: msg,
	}
}

// GetBadRequestErr ...
func GetBadRequestCodeErr(msg string) CodeError {
	return CodeError{
		Code:    400,
		Message: msg,
	}
}

// GetInternalErr ...
func GetInternalCodeErr(msg string) CodeError {
	return CodeError{
		Code:    500,
		Message: msg,
	}
}

func GetSuccessResp() *CommonResponse {
	return &CommonResponse{}
}

func GetErr(code int, msg string) *CommonResponse {
	return &CommonResponse{
		Error: CommonError{
			Code:    strconv.Itoa(code),
			Message: msg,
		},
	}
}

// IsErrorResponse response判断
func IsErrorResponse(response *CommonResponse) bool {
	return IsErrResponseCode(response.Error.Code)
}

// IsErrResponseCode response code判断
func IsErrResponseCode(errcode string) bool {
	return errcode != "" && errcode != "0" && errcode != "200"
}

// IsClientErrResponseCode response code判断
func IsClientErrResponseCode(errcode string) bool {
	return errcode == "400"
}

// IsSuccessResponseCode 判断错误码
func IsSuccessResponseCode(code string) bool {
	return code == "" || code == "0" || code == "200"
}

// IsSuccessCode 成功
func IsSuccessCode(code int) bool {
	return code == 0 || code == 200
}

// GetHttpStatusCode code码转换为httpStatus
func GetHttpStatusCode(code string) int {
	i, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		return http.StatusInternalServerError
	}
	if i < 100 || i > 999 {
		return http.StatusInternalServerError
	}
	return int(i)
}

// CommonRspV2 通用错误
type CommonRspV2 struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *CommonRspV2) Serialize() []byte {
	body, err := json.Marshal(r)
	if err != nil {
		return []byte{}
	}
	return body
}

// GetCommonSuccessResp 通用成功返回
func GetCommonSuccessResp(data interface{}) *CommonRspV2 {
	return &CommonRspV2{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func NewCommonErrRspV2(code int, message string) *CommonRspV2 {
	return &CommonRspV2{
		Code:    int32(code),
		Message: message,
		Data:    struct{}{},
	}
}

func NewCommonRspV2(code int, message string, data interface{}) *CommonRspV2 {
	return &CommonRspV2{
		Code:    int32(code),
		Message: message,
		Data:    data,
	}
}
