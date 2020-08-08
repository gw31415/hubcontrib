package image

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"unsafe"
)

func Svg(username string) (string, error) {
	res, err := http.Get("https://github.com/users/" + username + "/contributions")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	html := *(*string)(unsafe.Pointer(&body))
	return ext_elem("svg", html)
}

func ext_elem(tag string, source string) (string, error) {
	last_tag := "</" + tag + ">"
	begin := strings.Index(source, "<"+tag)
	end := strings.Index(source, last_tag)
	if begin == -1 || end == -1 {
		return "", errors.New("target element does not exist.")
	}
	return source[begin : end+len(last_tag)], nil
}
