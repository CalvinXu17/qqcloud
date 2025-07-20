package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
	uuid "github.com/satori/go.uuid"
	"github.com/yuchenfw/gocrypt"
	"github.com/yuchenfw/gocrypt/rsa"
)

var (
	Conf Config
)

func init() {
	status, Conf := ReadIni()
	if status == false {
		logs.Error("配置文件读取失败")
	}
	logs.Info("加载配置文件：", Conf)
}

type LoginResult struct {
	Info  interface{} `json:"info"`
	Token string      `json:"token"`
}

type Config struct {
	Cache      string `json:"Cache"`
	RedisIp    string `json:"RedisIp"`
	RedisPort  string `json:"RedisPort"`
	RedisPwd   string `json:"RedisPwd"`
	Sql        string `json:"Sql"`
	SqlIp      string `json:"SqlIp"`
	SqlPort    string `json:"SqlPort"`
	SqlUser    string `json:"SqlUser"`
	SqlPwd     string `json:"SqlPwd"`
	SqlDb      string `json:"SqlDb"`
	Key        string `json:"Key"`
	SqlRebuild bool   `json:"SqlRebuild"`
}

func GetInterfaceToInt(t1 interface{}) int {
	var t2 int
	switch t1.(type) {
	case uint:
		t2 = int(t1.(uint))
		break
	case int8:
		t2 = int(t1.(int8))
		break
	case uint8:
		t2 = int(t1.(uint8))
		break
	case int16:
		t2 = int(t1.(int16))
		break
	case uint16:
		t2 = int(t1.(uint16))
		break
	case int32:
		t2 = int(t1.(int32))
		break
	case uint32:
		t2 = int(t1.(uint32))
		break
	case int64:
		t2 = int(t1.(int64))
		break
	case uint64:
		t2 = int(t1.(uint64))
		break
	case float32:
		t2 = int(t1.(float32))
		break
	case float64:
		t2 = int(t1.(float64))
		break
	case string:
		t2, _ = strconv.Atoi(t1.(string))
		if t2 == 0 && len(t1.(string)) > 0 {
			f, _ := strconv.ParseFloat(t1.(string), 64)
			t2 = int(f)
		}
		break
	case nil:
		t2 = 0
		break
	case json.Number:
		t3, _ := t1.(json.Number).Int64()
		t2 = int(t3)
		break
	default:
		t2 = t1.(int)
		break
	}
	return t2
}

func IsValueNil(b interface{}) bool {
	if b == nil {
		return true
	}
	type eface struct {
		v   int64
		ptr unsafe.Pointer
	}
	efaceptr := (*eface)(unsafe.Pointer(&b))
	if efaceptr == nil {
		return true
	}
	return uintptr(efaceptr.ptr) == 0x0
}

func GetToken() (token string) {
	t := uuid.NewV4()
	re := md5.Sum([]byte(t.String()))
	code := fmt.Sprintf("%x", re)
	return code
}

func GetStringMd5(str string) (res string) {
	re := md5.Sum([]byte(str))
	code := fmt.Sprintf("%x", re)
	return code
}

func ReadIni() (status bool, conf Config) {
	cfg, err := config.NewConfig("ini", "config.conf")
	if err != nil {
		logs.Error(err)
		return false, Config{}
	}
	Cache, _ := cfg.String("app::cache")
	Key, _ := cfg.String("app::key")
	RedisIp, _ := cfg.String("redis::ip")
	if RedisIp == "" {
		RedisIp = "127.0.1"
	}
	RedisPort, _ := cfg.String("redis::port")
	if RedisPort == "" || RedisPort == "0" {
		RedisPort = "6379"
	}
	RedisPwd, _ := cfg.String("redis::pwd")
	Sql, _ := cfg.String("sql::type")
	SqlIp, _ := cfg.String("sql::ip")
	SqlPort, _ := cfg.String("sql::port")
	SqlUser, _ := cfg.String("sql::user")
	SqlPwd, _ := cfg.String("sql::pwd")
	SqlDb, _ := cfg.String("sql::db")
	rebuild, _ := cfg.String("sql::rebuild")
	SqlRebuild := false
	if rebuild == "true" {
		SqlRebuild = true
	}
	Conf = Config{
		Cache:      Cache,
		RedisIp:    RedisIp,
		RedisPort:  RedisPort,
		RedisPwd:   RedisPwd,
		Sql:        Sql,
		SqlIp:      SqlIp,
		SqlPort:    SqlPort,
		SqlUser:    SqlUser,
		SqlPwd:     SqlPwd,
		SqlDb:      SqlDb,
		Key:        Key,
		SqlRebuild: SqlRebuild,
	}
	return true, Conf
}

var (
	globalCacheClient cache.Cache
	cacheOnce         sync.Once
	cacheStatus       bool
)

// 初始化缓存客户端，仅执行一次
func initCache() {
	var ac cache.Cache
	var err error

	if Conf.Cache == "file" {
		ac, err = cache.NewCache("file", `{"CachePath":"./cache","FileSuffix":".cache","DirectoryLevel":"2","EmbedExpiry":"3600"}`)
	} else if Conf.Cache == "redis" {
		configStr := fmt.Sprintf(`{"key":"qqyCache","conn":"%s:%s","dbNum":"0","password":"%s"}`, Conf.RedisIp, Conf.RedisPort, Conf.RedisPwd)
		ac, err = cache.NewCache("redis", configStr)
	} else {
		logs.Error("未知缓存类型:", Conf.Cache)
		return
	}

	if err != nil {
		logs.Error("缓存初始化失败:", err)
		return
	}

	globalCacheClient = ac
	cacheStatus = true
}

func GetCacheAC() (bool, cache.Cache) {
	cacheOnce.Do(initCache)
	return cacheStatus, globalCacheClient
}

func GetManagerId(id interface{}) int {
	if id == "" {
		return 0
	}
	return GetInterfaceToInt(id)
}

func GetTokenString(Ctx *context.Context) string {
	token := Ctx.Request.Header["Token"]
	if len(token) == 1 {
		return token[0]
	}
	return ""
}

func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

type RSACrypt struct {
	PublicKey  string
	PrivateKey string
	Type       gocrypt.Encode
}

func (r *RSACrypt) RSASign(s string, hashType int) string {
	var hash = []gocrypt.Hash{gocrypt.MD5, gocrypt.SHA1, gocrypt.SHA224, gocrypt.SHA256, gocrypt.SHA384, gocrypt.SHA512}
	r.PublicKey = strings.Replace(r.PublicKey, "-----BEGIN RSA PUBLIC KEY-----", "", -1)
	r.PublicKey = strings.Replace(r.PublicKey, "-----END RSA PUBLIC KEY-----", "", -1)
	r.PrivateKey = strings.Replace(r.PrivateKey, "-----BEGIN RSA PRIVATE KEY-----", "", -1)
	r.PrivateKey = strings.Replace(r.PrivateKey, "-----END RSA PRIVATE KEY-----", "", -1)
	var secretInfo = rsa.RSASecret{
		PublicKey:          r.PublicKey,
		PublicKeyDataType:  r.Type,
		PrivateKey:         r.PrivateKey,
		PrivateKeyDataType: r.Type,
		PrivateKeyType:     gocrypt.PKCS1,
	}
	handle := rsa.NewRSACrypt(secretInfo)
	sign, err := handle.Sign(s, hash[hashType], gocrypt.Base64)
	if err != nil {
		return ""

	}
	return sign
}

func GetFloatLen(a float64) int {
	n := fmt.Sprint(a)
	t := strings.Split(n, ".")
	if len(t) <= 1 {
		return 0
	}
	return len(t[1])
}

func GetAddrIp(a string) string {
	t := strings.Split(a, ":")
	if len(t) <= 1 {
		return ""
	}
	return t[0]
}

func GetVersionString(a float64) string {
	l := GetFloatLen(a)
	var s strings.Builder
	n := fmt.Sprint(a)
	s.WriteString(n)
	if l == 0 {
		s.WriteString(".00")
	}
	if l == 1 {
		s.WriteString("0")
	}
	if l == 2 {
		s.WriteString("")
	}
	return s.String()
}

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func AesEncrypt(text []byte, key []byte, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("错误 -" + err.Error())
		return "", err
	}
	blockSize := block.BlockSize()
	originData := pad(text, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesDecrypt(text string, key []byte, iv []byte) (string, error) {
	decode_data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", nil
	}
	block, _ := aes.NewCipher(key)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origin_data := make([]byte, len(decode_data))
	blockMode.CryptBlocks(origin_data, decode_data)
	return string(unpad(origin_data)), nil
}

func unpad(ciphertext []byte) []byte {
	length := len(ciphertext)
	unpadding := int(ciphertext[length-1])
	return ciphertext[:(length - unpadding)]
}

func PostReq(url string, params string) (string, error) {
	var err error
	// 1. 创建http客户端实例
	client := &http.Client{}
	// 2. 创建请求实例
	req, err := http.NewRequest("POST", url, strings.NewReader(params))
	if err != nil {
		return "", err
	}
	// 3. 设置请求头，可以设置多个
	req.Header.Set("Host", " ")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 4. 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 5. 一次性读取请求到的数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body), err
}

type DateObj struct {
	Start     int64
	End       int64
	StartDate string
	EndDate   string
	Date      string
}

func GetDateRange(days int) []DateObj {
	var dateRange []DateObj
	i := 0
	for i < days {
		year, month, day := time.Now().AddDate(0, 0, -6).AddDate(0, 0, i).Date()
		t := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		startDate, _ := time.Parse("2006-01-02 15:04:05", t.Format("2006-01-02 15:04:05"))
		t = time.Date(year, month, day, 23, 59, 59, 0, time.Local)
		endDate, _ := time.Parse("2006-01-02 15:04:05", t.Format("2006-01-02 15:04:05"))
		dateRange = append(dateRange, DateObj{
			Start:     startDate.Unix(),
			End:       endDate.Unix(),
			StartDate: startDate.Format("2006-01-02 15:04:05"),
			EndDate:   endDate.Format("2006-01-02 15:04:05"),
			Date:      startDate.Format("2006-01-02"),
		})
		i++
	}
	return dateRange
}
