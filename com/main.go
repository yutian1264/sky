package main

import (
	"encoding/base64"
	"fmt"
	"sky/com/utils/rsa"
)

func main() {
	//创建公钥私钥文件
	rsa.CreateRsaFile()
	//获取公钥
	publicKey := rsa.GetPubKey()
	//加密
	b, err := rsa.RsaEncrypt([]byte("测试加密"), []byte(publicKey))
	if err != nil {
		fmt.Println("加密失败")
	}
	//base64加密 使返回的字符串漂亮点
	enc := base64.StdEncoding.EncodeToString(b)

	fmt.Println("加密后的文件是:", enc)
	//获取私钥
	privateKey := rsa.GetPrivKey()
	//先base64解密
	dec, _ := base64.StdEncoding.DecodeString(enc)
	fmt.Println("base64解密:", dec)
	deStr, _ := rsa.RsaDecrypt(dec, []byte(privateKey))

	fmt.Println("解密后铭文:", string(deStr))

}
