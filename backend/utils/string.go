package utils

import (
	"bytes"
	"regexp"
)

func ToSnake(str string) string {
	reg := regexp.MustCompile(`(\w)([A-Z])`)
	res := bytes.ToLower(reg.ReplaceAll([]byte(str), []byte("${1}_${2}")))
	return string(res)
}
