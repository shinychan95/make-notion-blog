package utils

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// SliceToString Convert a string slice to a string, where each element is separated by a comma.
func SliceToString(s []string, transformFn func(string) string) string {
	var sb strings.Builder
	for i, element := range s {
		if i > 0 {
			sb.WriteString(", ")
		}
		if transformFn != nil {
			element = transformFn(element)
		}
		sb.WriteString(element)
	}
	return sb.String()
}

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

func CheckUUIDv4Format(id string) (idWithHyphen string, err error) {
	uuidPattern := regexp.MustCompile(`^([0-9a-fA-F]{8})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{12})$`)
	matches := uuidPattern.FindStringSubmatch(id)

	if matches == nil {
		err = fmt.Errorf("input is not a valid UUID v4")
		return
	}

	idWithHyphen = fmt.Sprintf("%s-%s-%s-%s-%s", matches[1], matches[2], matches[3], matches[4], matches[5])

	return idWithHyphen, err
}

func FindNotionDBPath() (dbPath string) {
	cmd := exec.Command("lsof", "-c", "Notion")
	output, err := cmd.StdoutPipe()
	CheckError(err)

	err = cmd.Start()
	CheckError(err)

	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "notion.db") {
			dbPath = strings.Fields(line)[8]
			return
		}
	}

	err = cmd.Wait()
	CheckError(err)

	// 만약 값을 찾지 못하여도 프로세스 종료
	if dbPath == "" {
		ExecError("not exist notion db")
	}

	return
}
