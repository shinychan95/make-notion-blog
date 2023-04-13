package notion

import (
	"database/sql"
	"encoding/json"
	"github.com/shinychan95/make-notion-blog/utils"
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

func AssignNumbersToBlocks(blocks *[]Block) {
	var currentNumber uint8 = 1

	for i := range *blocks {
		if (*blocks)[i].Type == "numbered_list" {
			(*blocks)[i].Number = currentNumber
			currentNumber++

			// Children 순회
			AssignNumbersToBlocks(&((*blocks)[i].Children))
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

func GetRootBlockType(db *sql.DB, rootBlockID string) string {
	rows, err := db.Query("SELECT type FROM block WHERE id = ?", rootBlockID)
	utils.CheckError(err)
	defer rows.Close()

	for rows.Next() {
		var t string
		err = rows.Scan(&t)
		utils.CheckError(err)

		return t
	}

	utils.ExecError("root block is not in db")
	return ""
}

func GetBlockData(db *sql.DB, blockID string) []Block {
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
			block.Children = append(block.Children, GetBlockData(db, childID)...)
		}

		blocks = append(blocks, block)
	}

	return blocks
}
