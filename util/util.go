package util

import (
	"fmt"
)

// 路径部分排序做目录
func SortPath(str []byte) string {

	strLen := len(str)
	for i := 0; i < strLen; i++ {
		for j := 1 + i; j < strLen; j++ {
			if str[i] > str[j] {
				str[i], str[j] = str[j], str[i]
			}
		}
	}

	ret := ""

	for i := 0; i < strLen; i++ {
		ret += fmt.Sprintf("%d", str[i])
	}

	return ret
}

// 组合文件目录路径
func JoinPath(md5_str string) string {

	sortPath := SortPath([]byte(md5_str[:5]))
	return "/data/file/img" + sortPath + "/" + md5_str

}
