package stores

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"math/rand"
// 	"mime/multipart"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"text/template"
// 	"time"

// 	"crypto/hmac"
// 	"crypto/sha1"
// 	"encoding/base64"
// 	"sync"
// )

// const (
// 	MAX_SINGLE_SIZE         int64 = 8 * 1024 * 1024
// 	UPLOAD_SLICE_BLOCK_SIZE int64 = 1024 * 1024
// )

// /**
//  * CosClient
//  */
// type CosClient struct {
// 	AppID     string
// 	SecretID  string
// 	SecretKey string
// 	Bucket    string
// 	Local     string
// 	UseHttps  bool
// }

// type CosError struct {
// 	Code    int
// 	Message string
// }

// func (e *CosError) Error() string {
// 	return fmt.Sprintf("cos error - %d :%s", e.Code, e.Message)
// }

// type CosBaseResponse struct {
// 	Code    int    `json:"code"`
// 	Message string `json:"message"`
// }

// var (
// 	client = &http.Client{}
// )

// type CosResource struct {
// 	Name string `json:"name"`
// }

// // UploadFile 上传文件
// func (c *CosClient) UploadFile(local io.Reader, remote string, cover bool) {
// 	fileContent, _ := ioutil.ReadAll(local)
// 	fi := bytes.NewBuffer(fileContent)

// 	if int64(len(fileContent)) > MAX_SINGLE_SIZE {
// 		c.UploadLargeFile(fileContent, remote, cover)
// 		return
// 	}

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	writer.WriteField("op", "upload")
// 	//writer.WriteField("sha", base64.StdEncoding.EncodeToString(shaSum[:]))
// 	if cover {
// 		writer.WriteField("insertOnly", "0")
// 	}
// 	writer.WriteField("filecontent", string(fileContent))

// 	request, _ := http.NewRequest("POST", c.buildResourceURL(remote), body)
// 	request.Header.Add("Authorization", c.multiSignature(""))
// 	request.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())

// 	result := CosBaseResponse{}
// 	doRequestAsJson(request, &result)
// 	if result.Code == 0 {
// 		fmt.Printf("[ok   %s]\r\n", remote)
// 	} else {
// 		fmt.Fprintf(os.Stderr, "[failre  %s] - %d:%s\r\n", remote, result.Code, result.Message)
// 	}
// }

// // UploadLargeFile ...
// func (c *CosClient) UploadLargeFile(data []byte, remote string, cover bool) {
// 	length := int64(len(data))

// 	defer func() {
// 		e := recover()
// 		if e != nil {
// 			fmt.Fprintf(os.Stderr, "[failure %s] - %+v\r\n", remote, e)
// 		}
// 	}()

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	writer.WriteField("op", "upload_slice_init")
// 	writer.WriteField("filesize", strconv.FormatInt(length, 10))
// 	writer.WriteField("slice_size", strconv.FormatInt(UPLOAD_SLICE_BLOCK_SIZE, 10))
// 	if cover {
// 		writer.WriteField("insertOnly", "0")
// 	}

// 	url := c.buildResourceURL(remote)
// 	sign := c.multiSignature("")
// 	request, _ := http.NewRequest("POST", url, body)
// 	request.Header.Add("Authorization", sign)
// 	request.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())

// 	var response struct {
// 		Code    int    `json:"code"`
// 		Message string `json:"Message"`
// 		Data    struct {
// 			Session string `json:"session"`
// 		} `json:"data"`
// 	}
// 	doRequestAsJson(request, &response)

// 	if response.Code != 0 {
// 		fmt.Fprintf(os.Stderr, "[failure %s] - %s\r\n", remote, response.Message)
// 	}

// 	session := response.Data.Session
// 	ch := make(chan int, fi.Size()/UPLOAD_SLICE_BLOCK_SIZE)

// 	var offset int64
// 	count := 0

// 	threadPool := make(chan int, 10)
// 	for i := 0; i < 10; i++ {
// 		threadPool <- 1
// 	}

// 	for offset < fi.Size() {

// 		b := make([]byte, UPLOAD_SLICE_BLOCK_SIZE)
// 		length, _ := file.ReadAt(b, offset)
// 		<-threadPool
// 		go func(url, sign, session string, offset int64, bytes []byte, resultCH chan int) {
// 			defer func() {
// 				threadPool <- 1
// 			}()
// 			uploadSlice(url, sign, session, offset, bytes, resultCH)
// 		}(url, sign, session, offset, b[:length], ch)
// 		offset = offset + int64(length)
// 		count++
// 	}

// 	var code int
// 	succ := true
// 	for i := 0; i < count; i++ {

// 		code = <-ch
// 		succ = succ && code == 0
// 	}

// 	if succ {

// 		body := &bytes.Buffer{}
// 		writer := multipart.NewWriter(body)
// 		writer.WriteField("op", "upload_slice_finish")
// 		writer.WriteField("filesize", strconv.FormatInt(fi.Size(), 10))
// 		writer.WriteField("session", session)

// 		request, _ = http.NewRequest("POST", url, body)
// 		request.Header.Add("Authorization", sign)
// 		request.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())

// 		response := CosBaseResponse{}
// 		doRequestAsJson(request, &response)
// 		if response.Code == 0 {
// 			fmt.Printf("[ok     %s]\r\n", remote)
// 		} else {
// 			fmt.Fprintf(os.Stderr, "[fail %s] - %s\r\n", remote, response.Message)
// 		}

// 	} else {
// 		fmt.Fprintf(os.Stderr, "[failre %s]\r\n", remote)
// 	}

// }

// func uploadSlice(url, sign, session string, offset int64, b []byte, ch chan int) {
// 	body := &bytes.Buffer{}

// 	writer := multipart.NewWriter(body)
// 	writer.WriteField("op", "upload_slice_data")
// 	writer.WriteField("session", session)
// 	writer.WriteField("offset", strconv.FormatInt(offset, 10))
// 	field, _ := writer.CreateFormField("filecontent")
// 	field.Write(b)

// 	request, _ := http.NewRequest("POST", url, body)
// 	request.Header.Add("Authorization", sign)
// 	request.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())

// 	response := CosBaseResponse{}
// 	doRequestAsJson(request, &response)
// 	ch <- response.Code

// }

// func (c *CosClient) DownloadStream(remote string, callback func(io.Reader)) {
// 	request, _ := http.NewRequest("GET", c.buildDownloadUrl(remote), nil)
// 	request.Header.Add("Authorization", c.multiSignature(""))
// 	resp, e := client.Do(request)
// 	if e != nil {
// 		fmt.Fprintf(os.Stderr, "error occurred while download %s : %+v", remote, e)
// 		os.Exit(-1)
// 	}
// 	defer resp.Body.Close()
// 	length, _ := strconv.ParseInt(resp.Header.Get("content-length"), 10, 32)
// 	if length > 1024*1024 {
// 		fmt.Fprintf(os.Stderr, "%s is too large , use `gocos pull` instead\n", remote)
// 		os.Exit(-1)
// 	}
// 	callback(resp.Body)
// }

// type StatFileResult struct {
// 	AccessUrl     string                 `json:"access_url,omitempty"`
// 	Authority     string                 `json:"authority,omitempty"`
// 	BizAttr       string                 `json:"biz_attr"`
// 	Ctime         int64                  `json:"ctime"`
// 	CustomHeaders map[string]interface{} `json:"custom_headers"`
// 	FileLen       int64                  `json:"filelen"`
// 	FileSize      int64                  `json:"filesize"`
// 	Forbid        int                    `json:"forbid"`
// 	Mtime         int64                  `json:"mtime"`
// 	PreviewUrl    string                 `json:"preview_url,omitempty"`
// 	Sha           string                 `json:"sha,omitempty"`
// 	SliceSize     int64                  `json:"slicesize,omitempty"`
// 	SourceUrl     string                 `json:"source_url,omitempty"`
// }

// func (c *CosClient) StatFile(path string) *map[string]interface{} {

// 	request, _ := http.NewRequest("GET", c.buildResourceURL(path)+"?op=stat", nil)
// 	request.Header.Add("Authorization", c.multiSignature(""))
// 	response := &map[string]interface{}{}
// 	doRequestAsJson(request, response)
// 	//json.NewEncoder(os.Stdout).Encode(response)
// 	return response

// }

// func (c *CosClient) DeleteResource(path string, recursive, force bool) {
// 	if strings.HasSuffix(path, "/") && !recursive {
// 		fmt.Fprintln(os.Stderr, "use -r for delete directories")
// 		os.Exit(1)
// 	}

// 	data := struct {
// 		Op string `json:"op"`
// 	}{"delete"}
// 	body, _ := json.Marshal(data)

// 	request, _ := http.NewRequest("POST", c.buildResourceURL(path), bytes.NewBuffer(body))
// 	sign := c.onceSignature(path)
// 	request.Header.Add("Authorization", sign)
// 	request.Header.Add("Content-Type", "application/json")
// 	//println(path
// 	result := CosBaseResponse{}
// 	doRequestAsJson(request, &result)
// 	if result.Code == 0 {
// 		fmt.Printf("[Deleted %s]\r\n", path)
// 	} else {
// 		fmt.Fprintf(os.Stderr, "Failure(%s), %s\r\n", result.Message, path)
// 	}
// }

// func (c *CosClient) Move(src, target string, force bool) {
// 	if strings.HasSuffix(src, "/") {
// 		fmt.Fprintln(os.Stderr, "can not move directory !")
// 		os.Exit(1)
// 	}

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	writer.WriteField("op", "move")
// 	writer.WriteField("dest_fileid", target)
// 	if force {
// 		writer.WriteField("to_over_write", "1")
// 	}

// 	request, _ := http.NewRequest("POST", c.buildResourceURL(src), body)
// 	request.Header.Add("Authorization", c.onceSignature(src))
// 	request.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())

// 	result := CosBaseResponse{}
// 	doRequestAsJson(request, &result)
// 	if result.Code == 0 {
// 		fmt.Printf("[Move %s to %s Success]\r\n", src, target)
// 	} else {
// 		fmt.Printf("[Move %s to %s failure : %s]\r\n", src, target, result.Message)
// 	}
// }

// func (c *CosClient) onceSignature(file string) string {

// 	var data = struct {
// 		AppID    string
// 		SecretID string
// 		Bucket   string
// 		Exprire  int64
// 		Now      int64
// 		Random   int
// 		File     string
// 	}{c.AppID, c.SecretID, c.Bucket, 0, time.Now().Unix(), rand.Intn(9000000000) + 1000000000, "/" + c.AppID + "/" + c.Bucket + file}
// 	t, _ := template.New("signature-once").Parse("a={{.AppID}}&b={{.Bucket}}&k={{.SecretID}}&e={{.Exprire}}&t={{.Now}}&r={{.Random}}&f={{.File}}")
// 	var s bytes.Buffer
// 	t.Execute(&s, data)

// 	hash := hmac.New(sha1.New, []byte(c.SecretKey))
// 	hash.Write([]byte(s.String()))
// 	sum := hash.Sum(nil)
// 	sign := base64.StdEncoding.EncodeToString(append(sum, []byte(s.String())...))
// 	return sign

// }

// var signatureHolder = struct {
// 	once      sync.Once
// 	signature string
// }{sync.Once{}, ""}

// func (c *CosClient) multiSignature(fileName string) string {
// 	signatureHolder.once.Do(func() {
// 		var data = struct {
// 			AppID    string
// 			SecretID string
// 			Bucket   string
// 			Exprire  int64
// 			Now      int64
// 			Random   int
// 			FileName string
// 		}{c.AppID, c.SecretID, c.Bucket,
// 			time.Now().Unix() + 7776000, time.Now().Unix(),
// 			rand.Intn(9000000000) + 1000000000,
// 			fileName,
// 		}
// 		t, _ := template.New("signature-multi").Parse("a={{.AppID}}&b={{.Bucket}}&k={{.SecretID}}&e={{.Exprire}}&t={{.Now}}&r={{.Random}}&f={{.FileName}}")

// 		var s bytes.Buffer
// 		t.Execute(&s, data)
// 		println(s.String())
// 		s = *bytes.NewBufferString("aaaaaaaaaaaaaaaaaa")

// 		hash := hmac.New(sha1.New, []byte(c.SecretKey))
// 		hash.Write(s.Bytes())
// 		sum := hash.Sum(nil)
// 		println(string(sum))
// 		sign := base64.StdEncoding.EncodeToString(append(sum, s.Bytes()...))
// 		println(sign)
// 		signatureHolder.signature = sign
// 	})
// 	return signatureHolder.signature
// }

// func (c *CosClient) buildResourceURL(path string) string {
// 	var buffer bytes.Buffer
// 	if c.UseHttps {
// 		buffer.WriteString("https")
// 	} else {
// 		buffer.WriteString("http")
// 	}

// 	buffer.WriteString("://")
// 	buffer.WriteString(c.Local)
// 	buffer.WriteString(".file.myqcloud.com/files/v2/")
// 	buffer.WriteString(string(c.AppID))
// 	buffer.WriteString("/")
// 	buffer.WriteString(c.Bucket)
// 	if !strings.HasPrefix(path, "/") {
// 		buffer.WriteString("/")
// 	}
// 	buffer.WriteString(path)
// 	return buffer.String()
// }

// func (c *CosClient) buildDownloadUrl(path string) string {

// 	var buffer bytes.Buffer
// 	if c.UseHttps {
// 		buffer.WriteString("https")
// 	} else {
// 		buffer.WriteString("http")
// 	}

// 	buffer.WriteString("://")
// 	buffer.WriteString(c.Bucket)
// 	buffer.WriteString("-")
// 	buffer.WriteString(c.AppID)
// 	buffer.WriteString(".cos")
// 	buffer.WriteString(c.Local)
// 	buffer.WriteString(".myqcloud.com")
// 	if !strings.HasPrefix(path, "/") {
// 		buffer.WriteString("/")
// 	}
// 	buffer.WriteString(path)
// 	return buffer.String()

// }

// func doRequest(request *(http.Request)) *(http.Response) {
// 	resp, err := client.Do(request)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return resp
// }

// func doRequestAsJson(request *http.Request, val interface{}) error {
// 	resp := doRequest(request)
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}
// 	decoder := json.NewDecoder(bytes.NewReader(body))
// 	decoder.UseNumber()
// 	return decoder.Decode(val)
// }

// func panicError(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
