package system

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

func ReadExcel(path string) {
	excelFileName := path
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error != nil {
		fmt.Printf("Open excel error:%v", error)
	}

	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				fmt.Printf("%s\n", cell.GetNumberFormat())
			}
		}
	}
}

func WriteExcel() {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file, _ = xlsx.OpenFile("MyXLSXFile.xlsx")
	sheet = file.Sheet["Sheet1"]
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "000101"
	cell = row.AddCell()
	cell.Value = "中文1"
	err = file.Save("MyXLSXFile1.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
