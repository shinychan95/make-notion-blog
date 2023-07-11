package main

import (
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shinychan95/make-notion-blog/notion"
	"github.com/shinychan95/make-notion-blog/utils"
	"log"
	"sync"
)

var wg sync.WaitGroup
var errCh = make(chan error)

func main() {
	// flag
	configPath := flag.String("config", "config.json", "Path to config.json file")
	execType := flag.String("type", "collection_view_page", "type to execute ('page' or 'collection_view_page')")

	flag.Parse()

	if *execType != "page" && *execType != "collection_view_page" {
		log.Fatal("Invalid command. Only 'page' or 'collection' are accepted.")
	}

	//////////////////
	// program init //
	//////////////////
	config, err := utils.ReadConfig(*configPath)
	utils.CheckError(err)

	rootID, err := utils.CheckUUIDv4Format(config.RootID)
	utils.CheckError(err)

	notion.Init(config.ApiKey, config.PostDir, config.ImgDir, config.DBPath) // db open

	//////////////////
	// program exec //
	//////////////////
	switch *execType {
	case "page":
		// single page 는 우선 지원하지 않도록 하자.
	case "collection_view_page":
		notion.HandleCollectionViewPage(rootID, &wg, errCh)
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

	return
}
