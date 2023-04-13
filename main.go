package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinychan95/make-notion-blog/notion"
	"github.com/shinychan95/make-notion-blog/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var db *sql.DB

func savePageBlockAsMarkdown(rootBlockID, outputDir string) {
	blocks := notion.GetBlockData(db, rootBlockID)
	if blocks[0].Type != "page" {
		log.Fatal("root block id is not page")
	}

	pageTitle := notion.ParsePropTitle(blocks[0].Properties.String)
	fmt.Println("Page title:", pageTitle)

	notion.AssignNumbersToBlocks(&blocks)

	var wg sync.WaitGroup
	errCh := make(chan error)

	var markdownOutput string
	for _, block := range blocks[0].Children {
		markdownOutput += notion.ParseBlock(block, 0, &wg, errCh)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			log.Fatalf("Error occurred while downloading image: %v", err)
		}
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, os.ModePerm)
	}

	datePrefix := time.Now().Format("2006-01-02")
	markdownFileName := fmt.Sprintf("%s-%s.md", datePrefix, utils.SanitizeFileName(pageTitle))
	markdownFilePath := filepath.Join(outputDir, "", markdownFileName)

	err := ioutil.WriteFile(markdownFilePath, []byte(markdownOutput), 0644)
	utils.CheckError(err)

	fmt.Printf("Markdown file saved: %s\n", markdownFilePath)
}

func saveDatabaseBlockAsMarkdown(rootBlockID, outputDir string) {

}

func main() {
	// flag 를 사용하여 실행 시 설정 파일 경로 입력 (default: ./config.j
	configFilePath := flag.String("config", "config.json", "Path to the config.json file")
	flag.Parse()

	// 입력받은 config.json 파일 경로를 사용하여 설정을 읽어옴
	config, err := utils.ReadConfig(*configFilePath)
	utils.CheckError(err)

	// 만약 notion db 경로값 없을 경우, 동적으로 파악
	if config.DBPath == "" {
		config.DBPath = utils.FindNotionDBPath()
	}

	// sqlite3 DB open
	db, err = sql.Open("sqlite3", config.DBPath)
	utils.CheckError(err)
	defer db.Close()

	rootBlockID, err := utils.ConvertToUUIDv4(config.RootBlockID)
	utils.CheckError(err)

	rootBlockType := notion.GetRootBlockType(db, rootBlockID)

	// root block 의 타입에 따라 로직 다르게 동작
	switch rootBlockType {
	case "page":
		savePageBlockAsMarkdown(rootBlockID, config.OutputDir)
	case "collection_view_page":
		saveDatabaseBlockAsMarkdown(rootBlockID, config.OutputDir)
	default:
		utils.ExecError("not possible root block type")
	}

	log.Println("finish make notion into blog")
}
