package notion

import (
	"encoding/json"
	"fmt"
	"github.com/shinychan95/make-notion-blog/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

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

const APIKey = "secret_RrJPUd5a8BLNDZZp6BqGosxNikfmARDK3BcTzydyjBr"
const Version = "2022-06-28"

func SaveImage(blockId, outputDir string, wg *sync.WaitGroup, errCh chan error) {
	imageURL, err := getImageURL(blockId)
	utils.CheckError(err)

	imageFileName := fmt.Sprintf("%s.png", blockId)
	imagePath := filepath.Join(outputDir, "assets", imageFileName)

	wg.Add(1)
	downloadImageIfNotExist(imageURL, imagePath, wg, errCh)
}

func downloadImageIfNotExist(imageURL, imagePath string, wg *sync.WaitGroup, errCh chan error) {
	if !checkImageExist(imagePath) {
		go func(url, path string) {
			defer wg.Done()
			fmt.Printf("Downloading image: %s\n", url) // 이미지 다운로드 시작 메시지 출력
			err := downloadImage(url, path)
			if err != nil {
				errCh <- err
			} else {
				fmt.Printf("Image downloaded: %s\n", path) // 이미지 다운로드 완료 메시지 출력
			}
		}(imageURL, imagePath)
	} else {
		wg.Done()
	}
}

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

	req.Header.Add("Authorization", "Bearer "+APIKey)
	req.Header.Add("Notion-Version", Version)

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

func checkImageExist(imagePath string) bool {
	// 파일 존재 여부 확인
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return false
	}

	return true
}
