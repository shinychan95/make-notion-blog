package notion

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/shinychan95/make-notion-blog/markdown"
	"github.com/shinychan95/make-notion-blog/utils"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

type Block struct {
	ID         string
	Type       string
	Number     uint8
	ParsedProp ParsedProp
	Content    sql.NullString
	Children   []Block
	Properties sql.NullString
	Format     sql.NullString
	Table      *Table
}

type ParsedProp struct {
	Title    string
	Language string // for code type
}

func parseChildBlocks(block *Block) {
	childIDs, err := extractChildIDs(block.Content)
	utils.CheckError(err)

	for _, childID := range childIDs {
		childBlock := getBlockData(childID)
		parseChildBlocks(&childBlock)

		block.Children = append(block.Children, childBlock)
	}

	return
}

func setNumberedListValue(blocks *[]Block) {
	var currentNumber uint8 = 1

	for i := range *blocks {
		if (*blocks)[i].Type == "numbered_list" {
			(*blocks)[i].Number = currentNumber
			currentNumber++
			
			setNumberedListValue(&((*blocks)[i].Children))
		} else {
			currentNumber = 1
		}
	}
}

// extractChildIDs 함수 추가
func extractChildIDs(content sql.NullString) (childIDs []string, err error) {
	if !content.Valid {
		return
	}

	err = json.Unmarshal([]byte(content.String), &childIDs)
	utils.CheckError(err)

	return
}

/////////////////////////////////////////
// Block 자체 그리고 Type 에 따른 Parse 로직 //
/////////////////////////////////////////

func ParseBlock(pageID string, block Block, indentLv int, wg *sync.WaitGroup, errCh chan error) string {
	var output string

	if block.Properties.String != "" {
		block.ParsedProp.Title = ParsePropTitle(block.Properties.String)
		block.ParsedProp.Language = ParsePropLanguage(block.Properties.String)
	}

	indent := strings.Repeat("   ", indentLv)
	text := strings.ReplaceAll(block.ParsedProp.Title, "\n", "\n"+indent)

	switch block.Type {
	case "header":
		output = markdown.Header(indent, text)
	case "sub_header":
		output = markdown.SubHeader(indent, text)
	case "sub_sub_header":
		output = markdown.SubSubHeader(indent, text)
	case "text":
		output = markdown.Text(indent, text)
	case "code":
		output = markdown.Code(indent, block.ParsedProp.Language, text)
	case "divider":
		output = markdown.Divider(indent)
	case "bulleted_list":
		output = markdown.BulletedList(indent, text)
	case "numbered_list":
		output = markdown.NumberedList(indent, block.Number, text)
	case "toggle":
		var content string
		for _, child := range block.Children {
			content += ParseBlock(pageID, child, indentLv+1, wg, errCh)
		}
		output = markdown.Toggle(indent, text, content)
		block.Children = nil
	case "quote":
		output = markdown.Quote(indent, text)
	case "callout":
		output = markdown.Callout(indent, text)
	case "image":
		imageFileName := SaveImageIfNotExist(pageID, block.ID, wg, errCh)
		output = markdown.Image(indent, filepath.Join("/assets/pages", pageID, imageFileName))
	case "to_do":
		output = markdown.ToDo(indent, text, ParseChecked(block.Properties.String))
	case "table":
		output = createTableMarkdown(&block, block.Children)
		block.Children = nil
	default:
		log.Printf("Unsupported block type: %s", block.Type)
		output = ""
	}

	for _, child := range block.Children {
		output += ParseBlock(pageID, child, indentLv+1, wg, errCh)
	}

	return output
}

func ParsePropLanguage(properties string) (language string) {
	var props map[string]interface{}
	if err := json.Unmarshal([]byte(properties), &props); err != nil {
		panic(err)
	}

	if langValue, ok := props["language"]; ok {
		langArray := langValue.([]interface{})
		language = langArray[0].([]interface{})[0].(string)
	} else {
		language = ""
	}

	return
}

func ParsePropTitle(properties string) (text string) {
	var props map[string]interface{}
	if err := json.Unmarshal([]byte(properties), &props); err != nil {
		panic(err)
	}

	text = ParseText(props["title"])

	return
}

func ParseText(text interface{}) (parsedText string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("패닉 복구:", err)
		}
	}()

	// INFO - [ ["type",[["b"]]], [" "], ["자체가",[["i"]]], [" "], ["하나의",[["_"]]], [" "], ["변환으로",[["s"]]], ...]
	for _, value := range text.([]interface{}) {
		values := value.([]interface{})
		v := values[0].(string)

		// 길이가 1보다 큰 경우, text 에 대한 추가 형식 변환이 존재한다.
		if len(values) > 1 {
			for _, format := range values[1].([]interface{}) {
				f := format.([]interface{})
				switch f[0].(string) {
				case "b":
					v = markdown.Bold(v)
				case "i":
					v = markdown.Italic(v)
				case "s":
					v = markdown.Strikethrough(v)
				case "c":
					v = markdown.InlineCode(v)
				case "_":
					v = markdown.Underline(v)
				case "e":
					v = markdown.Equation(f[1].(string)) // [ "⁍", [["e","x+1"]] ]
				case "a":
					v = markdown.Link(v, f[1].(string))
				case "h":
					// 배경색이므로 무시
				case "‣":
					// 페이지 혹은 기타 노션 내부 링크이므로 무시
				default:
					//fmt.Printf("Error: Failed to parse properties. (%v) (%s) type\n", properties, f[0].(string))
				}
			}
			parsedText += v
		} else {
			parsedText += v
		}
	}

	return
}

func ParseChecked(properties string) bool {
	var propData map[string]interface{}
	err := json.Unmarshal([]byte(properties), &propData)
	utils.CheckError(err)

	checkedData := propData["checked"].([]interface{})
	checkedValue := checkedData[0].([]interface{})[0].(string)

	return checkedValue == "Yes"
}