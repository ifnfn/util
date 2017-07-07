package config

import (
	"fmt"

	"crypto/md5"
	"net/http"
	"strings"

	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"time"

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
	Php := false
	if Php {
		URL := `http://iot.roabay.com/lib/aliyun-php-sdk-sms/sms.php?to=%s&code=%s&key=%s&v=%s`
		Param := "?to=%s&code=%s&key=%s&v=%s"
		text := fmt.Sprintf("%s%s%s", to, code, "viewmobile")
		key := GetMD5([]byte(text))
		Param = fmt.Sprintf(Param, to, code, key, GetRandomString(5))
		URL = fmt.Sprintf("%s%s", URL, Param)
		_, err := SmsRequestServer(URL, "", "GET")
		fmt.Println("Url:", URL)
		return err
	} else {
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
	}
	return nil
}

func SmsRequestServer(url, data, method string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	req.Header.Add("User-Agent", "Aliyun sms")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println("Send sms http request error:", err)
	}

	return string(body), err
}

// GetRandomString 生成随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetMD5(text []byte) string {
	sum := md5.Sum(text)

	return hex.EncodeToString(sum[:])
}
