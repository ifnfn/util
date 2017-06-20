package system

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// SMTPInfo is the details for the SMS server
type SMSInfo struct {
	Username  string
	Password  string
	ServerUrl string
}

func GetHmacSha1Sign(str, key string) string {
	str = fmt.Sprintf(`POST&%s2F&%s`, "%", str)
	hmacSign := HmacSha1(str, fmt.Sprintf("%s&", key))
	return hmacSign
}

func SmsRequestServer(url, data, method string) (string, error) {
	//	fmt.Print(method, " URL:", url)

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

func getParam(to, code string) string {

	smTpl := "SMS_38790045"
	AccessKey := "LTAIl0WFf9fazIQR"
	AccessKeySecret := "PQiy1Ss3yDn6dkZOfvmSyXS25E6fUO"

	SignName := "容贝智城"

	random := fmt.Sprintf("%s-%s-%s", GetRandomNumberString(5), GetRandomNumberString(5), GetRandomNumberString(6))
	timeSign := GetAliyunFormatTime()

	param := make(map[string]interface{})
	param["AccessKeyId"] = AccessKey
	param["Action"] = "SingleSendSms"

	param["Format"] = "XML"
	cjson := `{"code":"%s"}`
	orgParamStr := fmt.Sprintf(cjson, code)

	param["ParamString"] = orgParamStr
	param["RecNum"] = to
	param["RegionId"] = "cn-hangzhou"

	param["SignName"] = SignName
	param["SignatureMethod"] = "HMAC-SHA1"
	param["SignatureNonce"] = random
	param["SignatureVersion"] = "1.0"

	param["TemplateCode"] = smTpl
	param["Timestamp"] = timeSign
	param["Version"] = "2016-09-27"
	param["Format"] = "JSON"

	sorted_keys := make([]string, 0)
	for k, _ := range param {
		sorted_keys = append(sorted_keys, k)
	}
	// sort 'string' key in increasing order
	tmp := ""
	sort.Strings(sorted_keys)
	for _, k := range sorted_keys {
		if len(tmp) > 0 {
			tmp = fmt.Sprintf("%s&", tmp)
		}
		v := param[k]
		tmp = fmt.Sprintf("%s%s=%s", tmp, k, v)
	}
	org := tmp

	tmp = ""
	for _, k := range sorted_keys {
		//		fmt.Printf("k=%v, v=%v\n", k, param[k])
		if len(tmp) > 0 {
			tmp = fmt.Sprintf("%s&", tmp)
		}
		a := fmt.Sprintf("%s", param[k])
		v := url.QueryEscape(a)
		tmp = fmt.Sprintf("%s%s=%s", tmp, k, v)
	}

	hParam := url.QueryEscape(tmp)
	fmt.Println("Wait Encode:", hParam)

	hSign := GetHmacSha1Sign(hParam, AccessKeySecret)
	desParam := fmt.Sprintf("Signature=%s&%s", hSign, org)
	return desParam
}

// SendMobileMessage sends an email
func SendMobileMessage(to, code string) error {
	Php := true
	if Php {
		URL := `http://iot.roabay.com/lib/aliyun-php-sdk-sms/sms.php?to=%s&code=%s&key=%s&v=%s`
		text := fmt.Sprintf("%s%s%s", to, code, "viewmobile")
		key := GetMD5([]byte(text))
		URL = fmt.Sprintf(URL, to, code, key, GetRandomNumberString(5))

		_, err := SmsRequestServer(URL, "", "GET")

		return err
	}

	URL := "https://sms.aliyuncs.com/"
	data := getParam(to, code)
	_, err := SmsRequestServer(URL, data, "POST")

	return err
}
