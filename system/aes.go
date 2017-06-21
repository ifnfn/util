package system

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"io"

	"roabay.com/util/config"
)

type AesEncrypt struct {
}

func (this *AesEncrypt) getKey() []byte {
	strKey := config.Server.DatabaseKey
	keyLen := len(strKey)
	if keyLen < 16 {
		panic("res key 长度不能小于16")
	}
	arrKey := []byte(strKey)

	if keyLen >= 32 {
		//取前32个字节
		return arrKey[:32]
	}
	if keyLen >= 24 {
		//取前24个字节
		return arrKey[:24]
	}
	//取前16个字节
	return arrKey[:16]
}

// 加密后再用 Base64 编码
func (this *AesEncrypt) EncryptBase64(src string) (string, error) {
	result, err := this.Encrypt(src)
	if err == nil {
		return base64.StdEncoding.EncodeToString(result), nil
	}

	return "", err
}

// 先用 Base64 解码后再解密
func (this *AesEncrypt) DecryptBase64(src string) (string, error) {
	data, _ := base64.StdEncoding.DecodeString(src)
	result, err := this.Decrypt(data)
	if err == nil {
		return result, nil
	}

	return "", err
}

//加密字符串
func (this *AesEncrypt) Encrypt(strMesg string) ([]byte, error) {
	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	encrypted := make([]byte, len(strMesg))
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, []byte(strMesg))

	return encrypted, nil
}

//解密字符串
func (this *AesEncrypt) Decrypt(src []byte) (strDesc string, err error) {
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, src)

	return string(decrypted), nil
}

// AesEncrpyt 加密
func AesEncrpyt(src io.Reader, dst io.Writer, key string, ctp int) error {
	md5Sum := md5.Sum([]byte(key)) // len = 16
	block, err := aes.NewCipher(md5Sum[:])
	if err != nil {
		return err
	}

	iv := md5Sum[:block.BlockSize()]
	var stream cipher.Stream

	switch ctp {
	case 1:
		stream = cipher.NewCFBEncrypter(block, iv)
	case 2:
		stream = cipher.NewCTR(block, iv)
	default:
		stream = cipher.NewOFB(block, iv)
	}

	writer := &cipher.StreamWriter{S: stream, W: dst}
	// Copy the input file to the output file, encrypting as we go.
	if _, err := io.Copy(writer, src); err != nil {
		return err
	}

	return nil
}

// AesDecrypt 解密
func AesDecrypt(src io.Reader, dst io.Writer, key string, ctp int) error {
	md5Sum := md5.Sum([]byte(key)) // len = 16
	block, err := aes.NewCipher(md5Sum[:])
	if err != nil {
		return err
	}

	iv := md5Sum[:block.BlockSize()]
	var stream cipher.Stream

	switch ctp {
	case 1:
		stream = cipher.NewCFBDecrypter(block, iv[:])
	case 2:
		stream = cipher.NewCTR(block, iv[:])
	default:
		stream = cipher.NewOFB(block, iv[:])
	}

	reader := &cipher.StreamReader{S: stream, R: src}
	// Copy the input file to the output file, decrypting as we go.
	if _, err := io.Copy(dst, reader); err != nil {
		return err
	}

	return nil
}
