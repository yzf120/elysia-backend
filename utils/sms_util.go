package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// TencentSMSClient 腾讯云短信客户端
type TencentSMSClient struct {
	client *sms.Client
}

// NewTencentSMSClient 创建腾讯云短信客户端
func NewTencentSMSClient() *TencentSMSClient {
	secretId := os.Getenv("TENCENT_SMS_SECRET_ID")
	secretKey := os.Getenv("TENCENT_SMS_SECRET_KEY")
	region := os.Getenv("TENCENT_SMS_REGION")

	if region == "" {
		region = "ap-guangzhou"
	}

	credential := common.NewCredential(secretId, secretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	client, _ := sms.NewClient(credential, region, cpf)

	return &TencentSMSClient{
		client: client,
	}
}

// SendVerificationCode 发送验证码短信
func (c *TencentSMSClient) SendVerificationCode(phoneNumber, code, templateId string) error {
	sdkAppId := os.Getenv("TENCENT_SMS_SDK_APP_ID")
	signName := os.Getenv("TENCENT_SMS_SIGN_NAME")

	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(sdkAppId)
	request.SignName = common.StringPtr(signName)
	request.TemplateId = common.StringPtr(templateId)
	request.TemplateParamSet = common.StringPtrs([]string{code, "5"})
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + phoneNumber})

	response, err := c.client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("腾讯云SDK错误: %v", err)
	}
	if err != nil {
		return fmt.Errorf("发送短信失败: %v", err)
	}

	if len(response.Response.SendStatusSet) > 0 {
		status := response.Response.SendStatusSet[0]
		if *status.Code != "Ok" {
			return fmt.Errorf("短信发送失败: %s", *status.Message)
		}
	}

	return nil
}

// GenerateVerificationCode 生成6位数字验证码
func GenerateVerificationCode() string {
	max := big.NewInt(1000000)
	n, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", n.Int64())
}
