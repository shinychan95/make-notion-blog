package utils

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

func FindNotionDBPath() string {
	cmd := exec.Command("lsof", "-c", "Notion")
	output, err := cmd.StdoutPipe()
	CheckError(err)

	err = cmd.Start()
	CheckError(err)

	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "notion.db") {
			dbPath := strings.Fields(line)[8]
			return dbPath
		}
	}

	err = cmd.Wait()
	CheckError(err)

	return ""
}

// ParsePageTitle 함수 추가
func ParsePageTitle(properties sql.NullString) (pageTitle string, err error) {
	if !properties.Valid {
		return
	}

	var propMap map[string][][]string
	err = json.Unmarshal([]byte(properties.String), &propMap)
	if err != nil {
		return
	}

	titleArr, ok := propMap["title"]
	if !ok || len(titleArr) < 1 || len(titleArr[0]) < 1 {
		return "", fmt.Errorf("title not found in properties")
	}

	return titleArr[0][0], nil
}

func ParsePropTitle(properties string) (text string) {
	var props map[string]interface{}
	if err := json.Unmarshal([]byte(properties), &props); err != nil {
		panic(err)
	}

	for _, value := range props["title"].([]interface{}) {
		values := value.([]interface{})
		v := values[0].(string)

		// 길이가 1보다 큰 경우, text 에 대한 추가 형식 변환이 존재한다.
		if len(values) > 1 {
			for _, format := range values[1].([]interface{}) {
				f := format.([]interface{})
				switch f[0].(string) {
				case "b":
					text += bold(v)
				case "i":
					text += italic(v)
				case "s":
					text += strikethrough(v)
				case "c":
					text += inlineCode(v)
				case "_":
					text += underline(v)
				case "e":
					text += equation(f[1].(string)) // [ "⁍", [["e","x+1"]] ]
				case "a":
					text += link(v, f[1].(string))
				default:
					text += v
					//fmt.Printf("Error: Failed to parse properties. (%v) (%s) type\n", properties, f[0].(string))
				}
			}
		} else {
			text += v
		}

	}

	return
}

func ParseChecked(properties string) bool {
	var propData map[string]interface{}
	err := json.Unmarshal([]byte(properties), &propData)
	CheckError(err)

	checkedData := propData["checked"].([]interface{})
	checkedValue := checkedData[0].([]interface{})[0].(string)

	return checkedValue == "Yes"
}

func ParsePropTitleSimple(properties string) (text string) {
	var props map[string]interface{}
	if err := json.Unmarshal([]byte(properties), &props); err != nil {
		panic(err)
	}

	for _, value := range props["title"].([]interface{}) {
		values := value.([]interface{})
		text += values[0].(string)
	}

	return
}
