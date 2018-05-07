package rsa

import (
"crypto/rsa"
"crypto/x509"
"encoding/pem"
"crypto/rand"
"flag"
"log"
"os"
	"io/ioutil"
	"errors"
	"sky/com/utils/file"
)


func checkFile() bool{

	 result:=false

	blpub,err:=file.IsPathExists("public.pem")
	if err!=nil{
		log.Fatal("读取公钥文件失败")
	}

	blpriv,err1:=file.IsPathExists("private.pem")
	if err1!=nil{
		log.Fatal("读取私钥文件失败")
	}

	if blpub && blpriv{
		result=true;
	}
	return result
}

func CreateRsaFile() {

	if checkFile(){
		return
	}


	var bits int
	flag.IntVar(&bits, "b", 2048, "密钥长度，默认为1024位")
	if err := genRsaKey(bits); err != nil {
		log.Fatal("密钥文件生成失败！")
	}
	log.Println("密钥文件生成成功！")
}


func GetPubKey() string{

	publicKey, err := ioutil.ReadFile("public.pem")
	if err != nil {
		os.Exit(-1)
	}
	return string(publicKey)
}


func GetPrivKey() string{

	publicKey, err := ioutil.ReadFile("private.pem")
	if err != nil {
		os.Exit(-1)
	}
	return string(publicKey)
}
// 加密
func RsaEncrypt(origData,publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext,privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func genRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "私钥",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "公钥",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}
