package system

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// HTTPSend ...
func HTTPSend(url, data, method string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	req.Header.Add("Accept", "*/*")

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// HTTPGet Get http
func HTTPGet(url string, headers map[string]string) ([]byte, error) {
	return HTTPSend(url, "", "GET", headers)
}

// HTTPPost Post http
func HTTPPost(url, data string, headers map[string]string) ([]byte, error) {
	return HTTPSend(url, data, "POST", headers)
}

// HTTPDelete Delete http
func HTTPDelete(url string, headers map[string]string) ([]byte, error) {
	return HTTPSend(url, "", "DELETE", headers)
}

// HTTPPatch Patch http
func HTTPPatch(url, data string, headers map[string]string) ([]byte, error) {
	return HTTPSend(url, data, "PATCH", headers)
}

// HTTPGetJSON 得到 JSON 数据
func HTTPGetJSON(url string, headers map[string]string) (interface{}, error) {
	var (
		jsonData interface{}
		err      error
		body     []byte
	)

	if body, err = HTTPGet(url, headers); err == nil {
		err = json.Unmarshal(body, &jsonData)
	}

	return jsonData, err
}

// HTTPPostJSON 得到 JSON 数据
func HTTPPostJSON(url, data string, headers map[string]string) (interface{}, error) {
	var (
		jsonData interface{}
		err      error
		body     []byte
	)

	if body, err = HTTPPost(url, data, headers); err == nil {
		err = json.Unmarshal(body, &jsonData)
	}

	return jsonData, err
}

// HTTPPatchJSON 得到 JSON 数据
func HTTPPatchJSON(url, data string, headers map[string]string) (interface{}, error) {
	var (
		jsonData interface{}
		err      error
		body     []byte
	)

	if body, err = HTTPPatch(url, data, headers); err == nil {
		err = json.Unmarshal(body, &jsonData)
	}

	return jsonData, err
}

// HTTPDeleteJSON 得到 JSON 数据
func HTTPDeleteJSON(url string, headers map[string]string) (interface{}, error) {
	var (
		jsonData interface{}
		err      error
		body     []byte
	)

	if body, err = HTTPDelete(url, headers); err == nil {
		err = json.Unmarshal(body, &jsonData)
	}

	return jsonData, err
}

// TaobaoGet 淘宝上的得到数据
func TaobaoGet(url string) ([]byte, error) {
	headers := map[string]string{
		"Accept-Language": "zh-CN,zh;q=0.8,en;q=0.6",
		"Referer":         "https://h5.m.taobao.com/app/detail/desc.html?isH5Des=true",
		"User-Agent":      "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.96 Mobile Safari/537.36",
	}
	return HTTPSend(url, "", "GET", headers)
}

// TaobaoGetJSON 淘宝上的得到数据
func TaobaoGetJSON(url string) (jsonData map[string]interface{}, err error) {
	if body, e := TaobaoGet(url); e == nil {
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			return nil, err
		}
	}

	return
}
