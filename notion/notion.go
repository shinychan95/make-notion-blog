package notion

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	ApiKey         string
	ApiVersion     = "2022-06-28"
	OutputDir      string
	RelativeImgDir string
)

func init() {
	ApiKey = ""
	OutputDir = ""
}

type Page struct {
	ID         string
	Categories []string
	Tags       []string
	Author     string
	Date       time.Time
	Status     string
	Title      string
}

type Schema struct {
	Name      string              `json:"name"`
	Type      string              `json:"type"`
	Options   []map[string]string `json:"options,omitempty"`
	OptionIds []string            `json:"optionIds,omitempty"`
}

type Schemas map[string]Schema

// Init 추가적인 인자의
func Init(apiKey, outputDir, relativeImgDir string) {
	ApiKey = apiKey
	OutputDir = outputDir
	RelativeImgDir = relativeImgDir
}

func (pg *Page) ToString() string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", pg.Title))
	sb.WriteString(fmt.Sprintf("author: %s\n", pg.Author))
	sb.WriteString(fmt.Sprintf("date: %s\n", pg.Date.Format("2006-01-02 15:04:05 -0700")))
	sb.WriteString("categories: ")

	for i, category := range pg.Categories {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("[%s]", category))
	}
	sb.WriteString("\n")

	sb.WriteString("tags: ")

	for i, tag := range pg.Tags {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(tag)
	}
	sb.WriteString("\n")

	sb.WriteString("---\n")

	return sb.String()
}

func GetPagesWithProperties(db *sql.DB, parentId string, schema map[string]Schema) ([]Page, error) {
	rows, err := db.Query("SELECT id, properties FROM block WHERE parent_id = ? AND type = 'page' AND is_template IS NULL AND alive = 1", parentId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var pages []Page
	for rows.Next() {
		var (
			id            string
			rawProperties string
		)
		err = rows.Scan(&id, &rawProperties)
		if err != nil {
			log.Fatal(err)
		}

		var page Page
		page, err = parsePageProperties(rawProperties, schema)
		page.ID = id
		if err != nil {
			log.Fatal(err)
		}
		pages = append(pages, page)
	}

	if pages == nil {
		return []Page{}, errors.New("no pages found")
	}

	return pages, nil
}

func GetCollectionId(db *sql.DB, rootID string) (string, error) {
	rows, err := db.Query("SELECT collection_id FROM block WHERE id = ? AND type = 'collection_view_page'", rootID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var collectionId string
	for rows.Next() {
		err = rows.Scan(&collectionId)
		if err != nil {
			log.Fatal(err)
		}
	}

	if collectionId == "" {
		return "", errors.New("cannot get collection id")
	}

	return collectionId, nil
}

func GetCollectionSchema(db *sql.DB, collectionId string) map[string]Schema {
	rows, err := db.Query("SELECT schema FROM collection WHERE id = ?", collectionId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var rawSchema string
	for rows.Next() {
		err = rows.Scan(&rawSchema)
		if err != nil {
			log.Fatal(err)
		}
	}

	var schemaMap map[string]Schema
	err = json.Unmarshal([]byte(rawSchema), &schemaMap)
	if err != nil {
		log.Fatal(err)
	}

	return schemaMap
}

func parsePageProperties(rawProperties string, schema map[string]Schema) (Page, error) {
	var propertiesMap map[string][][]interface{}
	err := json.Unmarshal([]byte(rawProperties), &propertiesMap)
	if err != nil {
		return Page{}, err
	}

	page := Page{}

	// Author, Date 의 경우, static 하게 입력한다.
	page.Author = "Chanyoung Kim"
	page.Date = time.Now()

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
		}
	}

	return page, nil
}
