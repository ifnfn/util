package system

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	curl "github.com/andelf/go-curl"
)

func HTTPSPostEx(url, data, token string, rdata *[]byte) error {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	easy.Setopt(curl.OPT_URL, url)
	easy.Setopt(curl.OPT_POST, false)
	easy.Setopt(curl.OPT_SSL_VERIFYHOST, false)
	easy.Setopt(curl.OPT_SSL_VERIFYPEER, false)
	easy.Setopt(curl.OPT_TIMEOUT, 30)
	easy.Setopt(curl.OPT_POSTFIELDS, data)
	easy.Setopt(curl.OPT_POSTFIELDSIZE, len(data))
	easy.Setopt(curl.OPT_HEADER, false)
	easy.Setopt(curl.OPT_TRANSFERTEXT, true)

	easy.Setopt(curl.OPT_HTTPHEADER, []string{
		"Content-type: application/json;charset='utf-8'",
		"Expect:",
		"Accept: */*",
		"Cache-Control: no-cache",
		"Pragma: no-cache",
		fmt.Sprintf("Content-Length: %d", len(data)),
		fmt.Sprintf("authtoken: %s", token),
	})
	IsReturn := false
	ReturnData := func(buf []byte, userdata interface{}) bool {
		num := len(buf)
		*rdata = buf[:num]
		IsReturn = true
		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, ReturnData)

	if err := easy.Perform(); err != nil {
		println("ERROR: ", err.Error())
		return err
	}
	time.Sleep(10000) // wait gorotine
	if IsReturn == false {
		time.Sleep(200000000)
	}

	return nil
}

func HTTPSPost(url, data, token string) ([]byte, error) {
	return HTTPSSend(url, data, token, "POST")
}

func HTTPSSend(url, data, token, method string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	req.Header["authtoken"] = []string{token}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("User-Agent", "Golang https client")
	req.Header.Add("Cache-control", "no-cache")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Accept-Encoding", "gzip,deflate,br")

	fmt.Println("Url:", url, " authtoken:", token)
	fmt.Println("Data:", data)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Fetch Httpclient
func Fetch(url, method string, headers map[string]string, data []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(data))

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified || resp.StatusCode == http.StatusNotFound {
		return ioutil.ReadAll(resp.Body)
	}

	return nil, fmt.Errorf("http error, %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
}

// HTTPSend ...
func HTTPSend(url, data, method string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNotModified || resp.StatusCode == http.StatusNotFound {
		return ioutil.ReadAll(resp.Body)
	}

	return nil, fmt.Errorf("http error, %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
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

func getCache(url string) ([]byte, error) {
	url = strings.ToUpper(url)
	md5 := strings.ToUpper(GetMD5([]byte(url)))
	fileName := "./cache/" + md5

	return ioutil.ReadFile(fileName)
}

func saveCache(url string, data []byte) {
	url = strings.ToUpper(url)
	md5 := strings.ToUpper(GetMD5([]byte(url)))
	fileName := "./cache/" + md5

	ioutil.WriteFile(fileName, data, 0644)
}

// CacheFetch ...
func CacheFetch(url, method string, headers map[string]string, dody []byte, cache bool) ([]byte, error) {
	if cache {
		if data, err := getCache(url); err == nil {
			u := strings.ToUpper(url)
			md5 := strings.ToUpper(GetMD5([]byte(u)))
			println("cache->", md5, url)
			return data, err
		}
	}

	data, err := Fetch(url, method, headers, dody)

	if err == nil {
		saveCache(url, data)
	}

	return data, err
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
