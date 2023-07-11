package notion

import (
	"encoding/json"
	"fmt"
	"github.com/shinychan95/make-notion-blog/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Page struct {
	ID         string
	Title      string
	Status     string
	Path       string
	Author     string
	Categories []string
	Tags       []string
	Published  time.Time
}

func (pg *Page) GetMetaString() string {
	var sb strings.Builder

	sb.WriteString("---\n")

	sb.WriteString(fmt.Sprintf("title: %s\n", pg.Title))
	sb.WriteString(fmt.Sprintf("author: %s\n", pg.Author))
	sb.WriteString(fmt.Sprintf("date: %s\n", pg.Published.Format("2006-01-02 15:04:05 -0700")))

	sb.WriteString("categories: [" + utils.SliceToString(pg.Categories, nil) + "]\n")
	sb.WriteString("tags: [" + utils.SliceToString(pg.Tags, strings.ToLower) + "]\n")

	sb.WriteString("---\n")

	return sb.String()
}

func handlePage(page Page, wg *sync.WaitGroup, errCh chan error) {
	fmt.Println("Page title:", page.Path)

	if page.ID == "1519a0a9-70f1-444e-95b4-f6e6fac46131" {
		fmt.Println("")
	}

	// page block 하위 모든 block parsing
	pageBlock := getBlockData(page.ID)
	parseChildBlocks(&pageBlock)
	setNumberedListValue(&pageBlock.Children)

	//////////////////////
	// markdown 결과 출력 //
	//////////////////////

	var markdownOutput string

	// 내부 헤더
	markdownOutput += page.GetMetaString() + "\n"

	// 내부 컨텐츠
	for _, block := range pageBlock.Children {
		markdownOutput += ParseBlock(page.ID, block, 0, wg, errCh)
	}

	if _, err := os.Stat(PostDir); os.IsNotExist(err) {
		os.MkdirAll(PostDir, os.ModePerm)
	}

	datePrefix := page.Published.Format("2006-01-02")
	markdownFileName := fmt.Sprintf("%s-%s.md", datePrefix, utils.SanitizeFileName(page.Path))
	markdownFilePath := filepath.Join(PostDir, "", markdownFileName)

	err := ioutil.WriteFile(markdownFilePath, []byte(markdownOutput), 0644)
	utils.CheckError(err)

	fmt.Printf("Markdown file saved: %s\n", markdownFilePath)
}

func parsePageProperties(page *Page, rawProperties string, schema map[string]Schema) {
	var propertiesMap map[string][][]interface{}
	err := json.Unmarshal([]byte(rawProperties), &propertiesMap)
	utils.CheckError(err)

	// INFO - Author 의 경우, static 하게 입력한다.
	page.Author = "chanyoung.kim"

	for key, value := range propertiesMap {
		schemaValue := schema[key]
		propertyValue := value[0]

		switch schemaValue.Name {
		case "Categories":
			page.Categories = strings.Split(propertyValue[0].(string), ",")
		case "Tags":
			page.Tags = strings.Split(propertyValue[0].(string), ",")
		case "Status":
			page.Status = propertyValue[0].(string)
		case "Title":
			page.Title = propertyValue[0].(string)
		case "Path":
			page.Path = propertyValue[0].(string)
		case "Published":
			dateProperty := propertyValue[1].([]interface{})[0].([]interface{})[1].(map[string]interface{})
			dateString := dateProperty["start_date"].(string)
			timeString, ok := dateProperty["start_time"]
			if !ok {
				timeString = "00:00"
			}
			dateTime := dateString + "T" + timeString.(string) + ":00"
			location, _ := time.LoadLocation("Asia/Seoul")
			page.Published, err = time.ParseInLocation("2006-01-02T15:04:05", dateTime, location)
			utils.CheckError(err)
		}
	}

	return
}
