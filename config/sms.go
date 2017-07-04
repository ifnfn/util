package config

import (
	"fmt"

	"github.com/goroom/aliyun_sms"
	"github.com/goroom/logger"
)

// SMSInfo is the details for the SMS server
type SMSInfo struct {
	SignName        string
	TempletCode     string
	AccessKey       string
	AccessKeySecret string
}

// SendMobileMessage sends an email
func (e SMSInfo) SendMobileMessage(to, code string) error {
	TempletCode := e.TempletCode
	AccessKey := e.AccessKey
	AccessKeySecret := e.AccessKeySecret
	SignName := e.SignName

	aliyun_sms, err := aliyun_sms.NewAliyunSms(SignName, TempletCode, AccessKey, AccessKeySecret)
	if err != nil {
		logger.Error(err)
		return err
	}
	cjson := `{"code":"%s"}`
	data := fmt.Sprintf(cjson, code)
	err = aliyun_sms.Send(to, data)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Error("Success")

	return nil
}
