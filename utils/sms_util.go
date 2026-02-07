package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dypnsapi20170525 "github.com/alibabacloud-go/dypnsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
)

// AliyunSMSClient 阿里云短信客户端
type AliyunSMSClient struct {
	client *dypnsapi20170525.Client
}

// NewAliyunSMSClient 创建阿里云短信客户端
func NewAliyunSMSClient() *AliyunSMSClient {
	client, err := createAliyunClient()
	if err != nil {
		fmt.Printf("创建阿里云短信客户端失败: %v\n", err)
		return nil
	}

	return &AliyunSMSClient{
		client: client,
	}
}

// createAliyunClient 使用凭据初始化账号Client
func createAliyunClient() (*dypnsapi20170525.Client, error) {
	// 工程代码建议使用更安全的无AK方式，凭据配置方式请参见：https://help.aliyun.com/document_detail/378661.html
	// 这里会自动从环境变量读取 ALIBABA_CLOUD_ACCESS_KEY_ID 和 ALIBABA_CLOUD_ACCESS_KEY_SECRET
	config := new(credential.Config).
		SetType("access_key").
		SetAccessKeyId(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")).
		SetAccessKeySecret(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"))
	credential, err := credential.NewCredential(config)
	if err != nil {
		return nil, err
	}

	openapiConfig := &openapi.Config{
		Credential: credential,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dypnsapi
	openapiConfig.Endpoint = tea.String("dypnsapi.aliyuncs.com")

	client, err := dypnsapi20170525.NewClient(openapiConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// SendVerificationCode 发送验证码短信
func (c *AliyunSMSClient) SendVerificationCode(phoneNumber, code, templateCode string) error {
	if c.client == nil {
		return fmt.Errorf("阿里云短信客户端未初始化")
	}

	signName := os.Getenv("ALIYUN_SMS_SIGN_NAME")
	if signName == "" {
		signName = "速通互联验证码" // 默认签名
	}

	// 构建模板参数 {"code":"123456","min":"5"}
	templateParam := fmt.Sprintf(`{"code":"%s","min":"5"}`, code)

	sendSmsVerifyCodeRequest := &dypnsapi20170525.SendSmsVerifyCodeRequest{
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		PhoneNumber:   tea.String(phoneNumber),
		TemplateParam: tea.String(templateParam),
	}

	runtime := &util.RuntimeOptions{}

	var sendErr error
	tryErr := func() error {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				sendErr = r
			}
		}()

		resp, err := c.client.SendSmsVerifyCodeWithOptions(sendSmsVerifyCodeRequest, runtime)
		if err != nil {
			return err
		}

		// 检查响应状态
		if resp.Body != nil && resp.Body.Code != nil && *resp.Body.Code != "OK" {
			message := "未知错误"
			if resp.Body.Message != nil {
				message = *resp.Body.Message
			}
			return fmt.Errorf("短信发送失败: %s", message)
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}

		// 错误处理
		errorMsg := tea.StringValue(error.Message)

		// 尝试解析诊断信息
		var data interface{}
		if error.Data != nil {
			d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
			d.Decode(&data)
			if m, ok := data.(map[string]interface{}); ok {
				if recommend, ok := m["Recommend"]; ok {
					errorMsg = fmt.Sprintf("%s (建议: %v)", errorMsg, recommend)
				}
			}
		}

		return fmt.Errorf("阿里云短信发送失败: %s", errorMsg)
	}

	if sendErr != nil {
		return fmt.Errorf("短信发送异常: %v", sendErr)
	}

	return nil
}

// GenerateVerificationCode 生成6位数字验证码
func GenerateVerificationCode() string {
	max := big.NewInt(1000000)
	n, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", n.Int64())
}
