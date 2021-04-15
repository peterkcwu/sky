package util

import (
	"image"
	"log"
	"net/url"
	"regexp"
)

var regexpUrlParse *regexp.Regexp

var NoImg *image.RGBA

func init() {

	var err error
	// 初始化正则表达式
	regexpUrlParse, err = regexp.Compile("[a-z0-9]{32}")
	if err != nil {
		log.Fatalln("regexpUrlParse:", err)
	}

	// 创建 RGBA 画板大小 - 用于找不到图片时用
	NoImg = image.NewRGBA(image.Rect(0, 0, 400, 400))

}

func UrlParse(md5_url string) string {

	if md5_url == "" {
		return ""
	}

	if len(md5_url) < 32 {
		return ""
	}

	// 进行 url 解析
	parse, err := url.Parse(md5_url)
	if err != nil {
		return ""
	}

	parsePath := parse.Path

	if len(parse.Path) != 32 {
		return ""
	}

	// 匹配是否是 md5 的长度
	if !IsMD5Path(parsePath) {
		return ""
	}

	// 组合文件完整路径
	return JoinPath(parsePath) + "/" + parsePath

}

// 匹配是否是 md5 的长度
func IsMD5Path(str string) bool {

	return regexpUrlParse.MatchString(str)

}
