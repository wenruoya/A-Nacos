package utils

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func FindNacos(target string) (string, string, bool) {
	targetUrl, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	responseV1, _ := http.Get(target)
	if responseV1 != nil {
		defer responseV1.Body.Close()
		str_byteV1, _ := ioutil.ReadAll(responseV1.Body)
		if strings.Contains(string(str_byteV1), "Nacos") {
			return target, "[+]" + target + "发现Nacos", true
		}
	}

	url2 := targetUrl.Scheme + "://" + strings.Split(targetUrl.Host, ":")[0] + ":8848/nacos"
	responseV2, _ := http.Get(url2)
	if responseV2 != nil {
		defer responseV2.Body.Close()
		str_byteV2, _ := ioutil.ReadAll(responseV2.Body)
		if strings.Contains(string(str_byteV2), "Nacos") {
			return targetUrl.Scheme + "://" + strings.Split(targetUrl.Host, ":")[0] + ":8848", "[+]" + url2 + "发现Nacos", true
		}
	}
	return target, "[-]" + target + "未发现Nacos", false
}
