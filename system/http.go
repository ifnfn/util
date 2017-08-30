package system

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

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
func Fetch(urls, method string, headers map[string]string, data []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, urls, bytes.NewReader(data))

	if proxyURL, exists := headers["Proxy"]; exists {
		delete(headers, "Proxy")
		if proxy, err := url.Parse(proxyURL); err == nil {
			println("use http proxy: ", proxyURL)
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxy),
			}
		}
	}

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

	return nil, fmt.Errorf("http error, %d: %s, %s", resp.StatusCode, http.StatusText(resp.StatusCode), urls)
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

// MD5url 将 URL 转为 MD5
func MD5url(url string) string {
	return strings.ToUpper(GetMD5([]byte(url)))
}

func getCache(url string) ([]byte, error) {
	fileName := "./cache/" + MD5url(url)

	return ioutil.ReadFile(fileName)
}

func saveCache(url string, data []byte) {
	fileName := "./cache/" + MD5url(url)

	ioutil.WriteFile(fileName, data, 0644)
}

// CacheFetch ...
func CacheFetch(url, method string, headers map[string]string, dody []byte, cache bool) ([]byte, error) {
	md5 := MD5url(url)
	if cache {
		if data, err := getCache(url); err == nil {
			// println("cache->", md5, url)
			return data, err
		}
	}

	data, err := Fetch(url, method, headers, dody)

	if err == nil {
		saveCache(url, data)
	}
	// println(string(data))
	println("GET: ", md5, url)

	return data, err
}

// TaobaoGet 淘宝上的得到数据
func TaobaoGet(url string) ([]byte, error) {
	headers := map[string]string{
		"Accept-Language": "zh-CN,zh;q=0.8,en;q=0.6",
		"Referer":         "https://h5.m.taobao.com/app/detail/desc.html?isH5Des=true",
		"User-Agent":      "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
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
