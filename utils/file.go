package utils

import (
	"regexp"
	"strings"
)

func SanitizeFileName(filename string) string {
	// 파일 이름에 사용할 수 없는 문자를 정규식으로 정의
	reg, err := regexp.Compile(`[<>:"/\\|?*]`)
	if err != nil {
		panic(err)
	}

	// 파일 이름에서 유효하지 않은 문자를 제거하고 공백을 하이픈(-)으로 대체
	sanitized := reg.ReplaceAllString(filename, "")
	sanitized = strings.ReplaceAll(sanitized, " ", "-")

	return sanitized
}
