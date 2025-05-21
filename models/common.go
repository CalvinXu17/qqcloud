package models

import (
	rand2 "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/beego/beego/v2/core/logs"
	"math/rand"
	"os"
	"sort"
	"time"
)

var cacheIndex = []string{
	"cache-project-%s",
	"cache-project-login-%s",
	"cache-project-version-%s",
}

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func RandUpperStr(length int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func RandLowerStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func RandNumStr(length int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func GetRsaKey() (status bool, PublicKey string, PrivateKey string) {
	bits := 1024
	privateKey, err := rsa.GenerateKey(rand2.Reader, bits)
	if err != nil {
		return false, "", ""
	}
	_, err = os.Stat("./key")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("./key", os.ModePerm)
			if err != nil {
				logs.Error("创建key文件夹失败")
				return false, "", ""
			}
		}
	}

	privateKeyStream := x509.MarshalPKCS1PrivateKey(privateKey)
	//file, err := os.Create("./key/private.pem")
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyStream,
	}
	//_ = pem.Encode(file, block)
	privateKeyStr := string(pem.EncodeToMemory(block))
	publicKey := &privateKey.PublicKey
	publicKeyStream, _ := x509.MarshalPKIXPublicKey(publicKey)
	//publicKeyStream := x509.MarshalPKCS1PublicKey(publicKey)
	//file, err = os.Create("./key/public.pem")
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyStream,
	}
	//_ = pem.Encode(file, block)
	publicKeyStr := string(pem.EncodeToMemory(block))
	return true, publicKeyStr, privateKeyStr
}

func In(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}
