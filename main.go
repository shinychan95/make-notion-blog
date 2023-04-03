package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/shinychan95/make-notion-blog/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	notionVersion = "2022-06-28"
)

var config Config

type Block struct {
	ID         string
	Type       string
	Number     uint8
	Content    sql.NullString
	Children   []Block
	Properties sql.NullString
	Format     sql.NullString
}

type ImageBlock struct {
	Object         string    `json:"object"`
	ID             string    `json:"id"`
	CreatedTime    time.Time `json:"created_time"`
	LastEditedTime time.Time `json:"last_edited_time"`
	Type           string    `json:"type"`
	Image          struct {
		Type string `json:"type"`
		File struct {
			URL        string    `json:"url"`
			ExpiryTime time.Time `json:"expiry_time"`
		} `json:"file"`
	} `json:"image"`
}

// Add your Notion API Key here
const notionAPIKey = "secret_RrJPUd5a8BLNDZZp6BqGosxNikfmARDK3BcTzydyjBr"

func downloadImage(url, imagePath string) error {
	resp, err := http.Get(url)
	utils.CheckError(err)
	defer resp.Body.Close()

	// 이미지 파일 저장
	out, err := os.Create(imagePath)
	utils.CheckError(err)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	utils.CheckError(err)

	return nil
}

func getImageURL(blockID string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.notion.com/v1/blocks/%s", blockID), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+notionAPIKey)
	req.Header.Add("Notion-Version", notionVersion)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var imageBlock ImageBlock
	err = json.Unmarshal(body, &imageBlock)
	if err != nil {
		return "", err
	}

	return imageBlock.Image.File.URL, nil
}

func assignNumbersToBlocks(blocks *[]Block) {
	var currentNumber uint8 = 1

	for i := range *blocks {
		if (*blocks)[i].Type == "numbered_list" {
			(*blocks)[i].Number = currentNumber
			currentNumber++

			// Children 순회
			assignNumbersToBlocks(&((*blocks)[i].Children))
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

func getBlockData(db *sql.DB, blockID string) []Block {
	rows, err := db.Query("SELECT id, type, content, properties, format FROM block WHERE id = ?", blockID)
	utils.CheckError(err)
	defer rows.Close()

	var blocks []Block

	for rows.Next() {
		var block Block
		err = rows.Scan(&block.ID, &block.Type, &block.Content, &block.Properties, &block.Format)
		utils.CheckError(err)

		childIDs, err := extractChildIDs(block.Content)
		utils.CheckError(err)

		for _, childID := range childIDs {
			block.Children = append(block.Children, getBlockData(db, childID)...)
		}

		blocks = append(blocks, block)
	}

	return blocks
}

func parseBlock(block Block, indentLevel int) string {
	var output string

	var title string
	if block.Properties.String != "" {
		title = utils.ParsePropTitle(block.Properties.String)
	}

	indent := strings.Repeat("  ", indentLevel)

	switch block.Type {
	case "header":
		output = fmt.Sprintf("%s# %s\n", indent, title)
	case "sub_header":
		output = fmt.Sprintf("%s## %s\n", indent, title)
	case "sub_sub_header":
		output = fmt.Sprintf("%s### %s\n", indent, title)
	case "text":
		output = fmt.Sprintf("%s %s\n", indent, title)
	case "paragraph":
		output = fmt.Sprintf("%s%s\n", indent, title)
	case "code":
		output = fmt.Sprintf("%s```yaml\n%s%s\n%s```\n", indent, indent, title, indent)
	case "divider":
		output = fmt.Sprintf("%s---\n", indent)
	case "bulleted_list":
		output = fmt.Sprintf("%s- %s\n", indent, title)
	case "numbered_list":
		output = fmt.Sprintf("%s%d. %s\n", indent, block.Number, title)
	case "toggle":
		output = fmt.Sprintf("%s<details>\n%s<summary>%s</summary>\n", indent, indent, title)
		for _, child := range block.Children {
			output += parseBlock(child, indentLevel+1)
		}
		output += fmt.Sprintf("%s</details>\n", indent)
	case "quote":
		output = fmt.Sprintf("%s> %s\n", indent, title)
	case "callout":
		output = fmt.Sprintf("%s💡 %s\n", indent, title)
		for _, child := range block.Children {
			output += parseBlock(child, indentLevel+1)
		}
	case "image":
		imageURL, err := getImageURL(block.ID)
		utils.CheckError(err)

		imageFileName := fmt.Sprintf("%s.png", block.ID)
		imagePath := filepath.Join(config.OutputDirectory, "assets", imageFileName)

		err = downloadImage(imageURL, imagePath)
		utils.CheckError(err)

		output = fmt.Sprintf("%s![](%s)\n", indent, "/assets/"+imageFileName)
	case "to_do":
		checked := utils.ParseChecked(block.Properties.String)
		if checked {
			output = fmt.Sprintf("%s- [x] %s\n", indent, title)
		} else {
			output = fmt.Sprintf("%s- [ ] %s\n", indent, title)
		}

	default:
		log.Printf("Unsupported block type: %s", block.Type)
		output = ""
	}

	for _, child := range block.Children {
		output += parseBlock(child, indentLevel+1)
	}

	return output
}

func main() {
	var err error

	configFilePath := "config.json"
	config, err = readConfig(configFilePath)
	utils.CheckError(err)

	// 만약 notion db 경로값 없을 경우, 동적으로 파악
	if config.DBPath == "" {
		config.DBPath = utils.FindNotionDBPath()
		if config.DBPath == "" {
			log.Fatal("notion location is missing")
		}
	}

	// sqlite3 DB open
	db, err := sql.Open("sqlite3", config.DBPath)
	utils.CheckError(err)
	defer db.Close()

	rootID := "0e00d47d-bb28-4497-9438-ce3e2dbfda68"
	blocks := getBlockData(db, rootID)

	// TODO -  데이터베이스 내 변경된 페이지 추적
	pageTitle, err := utils.ParsePageTitle(blocks[0].Properties)
	utils.CheckError(err)
	fmt.Println("Page title:", pageTitle)

	assignNumbersToBlocks(&blocks)

	var markdownOutput string
	for _, block := range blocks[0].Children {
		markdownOutput += parseBlock(block, 0)
	}

	outputDir := config.OutputDirectory
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, os.ModePerm)
	}

	datePrefix := time.Now().Format("2006-01-02")
	markdownFileName := fmt.Sprintf("%s-%s.md", datePrefix, utils.SanitizeFileName(pageTitle))
	markdownFilePath := filepath.Join(outputDir, "_posts", markdownFileName)

	err = ioutil.WriteFile(markdownFilePath, []byte(markdownOutput), 0644)
	utils.CheckError(err)

	fmt.Printf("Markdown file saved: %s\n", markdownFilePath)
}
