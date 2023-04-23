package utils

import (
	"bufio"
	"fmt"
	"os/exec"
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

func ConvertToUUIDv4(input string) (string, error) {
	uuidPattern := regexp.MustCompile(`^([0-9a-fA-F]{8})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{4})-?([0-9a-fA-F]{12})$`)
	matches := uuidPattern.FindStringSubmatch(input)

	if matches == nil {
		return "", fmt.Errorf("input is not a valid UUID v4")
	}

	return fmt.Sprintf("%s-%s-%s-%s-%s", matches[1], matches[2], matches[3], matches[4], matches[5]), nil
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

func RemoveCommonPrefix(firstPath, secondPath string) string {
	firstPathParts := strings.Split(firstPath, "/")
	secondPathParts := strings.Split(secondPath, "/")

	shortestLength := len(firstPathParts)
	if len(secondPathParts) < shortestLength {
		shortestLength = len(secondPathParts)
	}

	lastCommonIndex := -1
	for i := 0; i < shortestLength; i++ {
		if firstPathParts[i] == secondPathParts[i] {
			lastCommonIndex = i
		} else {
			break
		}
	}

	return "/" + strings.Join(secondPathParts[lastCommonIndex+1:], "/")
}
