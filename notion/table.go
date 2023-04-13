package notion

import (
	"encoding/json"
	"fmt"

	"strings"
)

type Table struct {
	ColumnOrder  []string
	ColumnHeader bool
}

func parseTable(tableBlock *Block) {
	var format map[string]interface{}
	if err := json.Unmarshal([]byte(tableBlock.Format.String), &format); err != nil {
		panic(err)
	}

	columnOrder := format["table_block_column_order"].([]interface{})
	columnHeader := format["table_block_column_header"].(bool)

	columnOrderStr := make([]string, len(columnOrder))
	for i, v := range columnOrder {
		columnOrderStr[i] = v.(string)
	}

	tableBlock.Table = &Table{
		ColumnOrder:  columnOrderStr,
		ColumnHeader: columnHeader,
	}
}

func parseTableRow(properties string, columnOrder []string) []string {
	var props map[string]interface{}
	if properties != "" {
		if err := json.Unmarshal([]byte(properties), &props); err != nil {
			panic(err)
		}
	}

	var row []string
	for _, colID := range columnOrder {
		if cellData, ok := props[colID]; ok {
			cell := ParseText(cellData)
			row = append(row, cell)
		} else {
			row = append(row, "")
		}
	}

	return row
}

func createTableMarkdown(tableBlock *Block, tableRowBlocks []Block) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("패닉 복구:", err)
		}
	}()

	parseTable(tableBlock)

	var markdown strings.Builder

	// 맨 첫 번째 행이 헤더가 되도록 설정
	headerRow := tableRowBlocks[0]
	headerCells := parseTableRow(headerRow.Properties.String, tableBlock.Table.ColumnOrder)
	markdown.WriteString("| ")
	for _, cell := range headerCells {
		markdown.WriteString(cell + " | ")
	}
	markdown.WriteString("\n")

	// 헤더와 데이터 사이에 구분선 추가
	markdown.WriteString("| ")
	for range headerCells {
		markdown.WriteString("--- | ")
	}
	markdown.WriteString("\n")

	// 데이터 행 작성
	for _, rowBlock := range tableRowBlocks[1:] { // 첫 번째 행을 건너뛰고 시작
		cells := parseTableRow(rowBlock.Properties.String, tableBlock.Table.ColumnOrder)
		markdown.WriteString("| ")
		for _, cell := range cells {
			markdown.WriteString(cell + " | ")
		}
		markdown.WriteString("\n")
	}

	return markdown.String()
}
