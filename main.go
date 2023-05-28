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
var wg sync.WaitGroup
var errCh = make(chan error)

func savePageBlockAsMarkdown(rootID, header, outputDir, path string, published time.Time) {
	blocks := notion.GetBlockData(db, rootID)
	if blocks[0].Type != "page" {
		log.Fatal("root block id is not page")
	}

	fmt.Println("Page title:", path)

	notion.AssignNumbersToBlocks(&blocks)

	var markdownOutput string
	markdownOutput += header + "\n"
	for _, block := range blocks[0].Children {
		markdownOutput += notion.ParseBlock(rootID, block, 0, &wg, errCh)
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, os.ModePerm)
	}

	datePrefix := published.Format("2006-01-02")
	markdownFileName := fmt.Sprintf("%s-%s.md", datePrefix, utils.SanitizeFileName(path))
	markdownFilePath := filepath.Join(outputDir, "", markdownFileName)

	err := ioutil.WriteFile(markdownFilePath, []byte(markdownOutput), 0644)
	utils.CheckError(err)

	fmt.Printf("Markdown file saved: %s\n", markdownFilePath)
}

func saveDatabaseBlockAsMarkdown(rootId, postDir string) {
	// collection_id 값 구하고, 해당 값을 parent_id 로 하는 페이지들을 구한다. (alive 값이 1인 block 페이지만)
	collectionId, _ := notion.GetCollectionId(db, rootId)
	collectionSchema := notion.GetCollectionSchema(db, collectionId)

	// template is NULL, alive 값이 1인 collection 내 페이지들을 가져온다.
	pages, _ := notion.GetPagesWithProperties(db, collectionId, collectionSchema)

	// property 내 Status 가 Drafting 인 글들만 프로세스를 실행한다.
	for _, page := range pages {
		if page.Status == "Drafting" {
			wg.Add(1)
			go func(rootId, header, postDir string) {
				savePageBlockAsMarkdown(rootId, header, postDir, page.Path, page.Published)
				wg.Done()
			}(page.ID, page.ToString(), postDir)
		}
	}

	wg.Wait()

}

func main() {
	// flag 를 사용하여 실행 시 설정 파일 경로 입력 (default: ./config.json)
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

	rootBlockID, err := utils.ConvertToUUIDv4(config.RootId)
	utils.CheckError(err)

	// notion 이미지 상대 경로 저장
	markdownImgDir := utils.RemoveCommonPrefix(config.PostDir, config.ImgDir)

	// notion 이미지 저장을 위한 key 및 경로 저장
	notion.Init(config.ApiKey, config.PostDir, config.ImgDir, markdownImgDir)

	rootBlockType := notion.GetRootBlockType(db, rootBlockID)

	// root block 의 타입에 따라 로직 다르게 동작
	switch rootBlockType {
	case "page":
		fileName := "test"
		savePageBlockAsMarkdown(rootBlockID, "", config.PostDir, fileName, time.Now())
	case "collection_view_page":
		saveDatabaseBlockAsMarkdown(rootBlockID, config.PostDir)
	default:
		utils.ExecError("not possible root block type")
	}

	// 이미지 다운로드 go routine 대기
	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err = range errCh {
		if err != nil {
			log.Fatalf("Error occurred while downloading image: %v", err)
		}
	}

	log.Println("finish make notion into blog")
}
