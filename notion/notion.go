package notion

import (
	"database/sql"
	"encoding/json"
	"github.com/shinychan95/make-notion-blog/utils"
)

var (
	ApiKey     string
	ApiVersion = "2022-06-28"
	PostDir    string
	ImgDir     string
	db         *sql.DB
)

func Init(apiKey, postDir, imgDir, dbPath string) {
	ApiKey = apiKey
	PostDir = postDir
	ImgDir = imgDir

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	utils.CheckError(err)
}

func Close() {
	err := db.Close()
	utils.CheckError(err)
}

//////////////////////
// Get Data From DB //
//////////////////////

func getCollectionId(rootID string) (colId string) {
	rows, err := db.Query("SELECT collection_id FROM block WHERE id = ? AND type = 'collection_view_page'", rootID)
	utils.CheckError(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&colId)
		utils.CheckError(err)

		if rows.Next() {
			utils.ExecError("more than one row returned")
		}
	}

	if colId == "" {
		utils.ExecError("cannot get collection id")
	}

	return
}

func getCollectionSchema(collectionId string) (schemaMap map[string]Schema) {
	rows, err := db.Query("SELECT schema FROM collection WHERE id = ?", collectionId)
	utils.CheckError(err)
	defer rows.Close()

	var rawSchema string
	for rows.Next() {
		err = rows.Scan(&rawSchema)
		utils.CheckError(err)

		if rows.Next() {
			utils.ExecError("more than one row returned")
		}

	}

	err = json.Unmarshal([]byte(rawSchema), &schemaMap)
	utils.CheckError(err)

	return
}

func getPagesWithProperties(parentId string, schema map[string]Schema) (pages []Page) {
	rows, err := db.Query("SELECT id, properties FROM block WHERE parent_id = ? AND type = 'page' AND is_template IS NULL AND alive = 1", parentId)
	utils.CheckError(err)
	defer rows.Close()

	for rows.Next() {
		var (
			id            string
			rawProperties string
		)
		err = rows.Scan(&id, &rawProperties)
		utils.CheckError(err)

		page := Page{ID: id}
		parsePageProperties(&page, rawProperties, schema)

		pages = append(pages, page)
	}

	return
}

func getRootType(rootID string) (t string) {
	rows, err := db.Query("SELECT type FROM block WHERE id = ?", rootID)
	utils.CheckError(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&t)
		utils.CheckError(err)

		if rows.Next() {
			utils.ExecError("more than one row returned")
			return
		}

		return
	}

	utils.ExecError("root block is not in db")

	return
}

func getBlockData(blockID string) (block Block) {
	rows, err := db.Query("SELECT id, type, content, properties, format FROM block WHERE id = ?", blockID)
	utils.CheckError(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&block.ID, &block.Type, &block.Content, &block.Properties, &block.Format)
		utils.CheckError(err)

		if rows.Next() {
			utils.ExecError("more than one row returned")
			return
		}
	}

	return
}
