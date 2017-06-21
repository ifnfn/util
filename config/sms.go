package config

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// SMSInfo is the details for the SMS server
type SMSInfo struct {
	Username  string
	Password  string
	ServerUrl string
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

// SendMobileMessage sends an email
func (e SMSInfo) SendMobileMessage(to, code string) error {
	URL := `http://iot.roabay.com/lib/aliyun-php-sdk-sms/sms.php?to=%s&code=%s&key=%s&v=%s`
	URL = e.ServerUrl
	text := fmt.Sprintf("%s%s%s", to, code, "viewmobile")
	key := GetMD5([]byte(text))
	URL = fmt.Sprintf(URL, to, code, key, GetRandomString(5))

	_, err := SmsRequestServer(URL, "", "GET")

	return err
}
