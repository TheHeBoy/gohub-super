package sms

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	jsoniter "github.com/json-iterator/go"
	"gohub/pkg/logger"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Aliyun struct{}

// Send 实现Driver中的Send方法
func (a *Aliyun) Send(phone string, message Message, config map[string]string) bool {
	client, err := CreateClient(tea.String(config["access_key_id"]), tea.String(config["access_key_secret"]))
	if err != nil {
		logger.Error("短信[阿里云]", "解析绑定错误", err.Error())
		return false
	}
	logger.Debug("短信[阿里云]", "配置信息", config)

	param, err := json.Marshal(message.Data)
	if err != nil {
		logger.Error("短信[阿里云]", "短信模板参数解析错误", err.Error())
		return false
	}

	// 发送参数
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(config["sign_name"]),
		TemplateCode:  tea.String(message.Template),
		PhoneNumbers:  tea.String(phone),
		TemplateParam: tea.String(string(param)),
	}

	// 其他运行参数
	runtime := &util.RuntimeOptions{}

	_, err = client.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		var errs = &tea.SDKError{}
		if _t, ok := err.(*tea.SDKError); ok {
			errs = _t
		} else {
			errs.Message = tea.String(err.Error())
		}

		var r dysmsapi20170525.SendSmsResponseBody
		err = json.Unmarshal([]byte(*errs.Data), &r)
		logger.ErrorIf(err)

		return false
	}

	return true
}

// CreateClient
/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}
