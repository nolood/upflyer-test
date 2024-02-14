package xlsx

import (
	"fmt"
	"os"
	"time"

	"github.com/nolood/upflyer-test.git/internal/config"
	"github.com/xuri/excelize/v2"
)

var (
	filename = "table.xlsx"
	sheet    = "Sheet1"
)

var sheetData = Sheet{}

type Sheet struct {
	Link    string
	DOB     string
	Age     int
	Subs    int
	ViewsM7 int
	ER      string
}

// Лучше использовать базу данных, и в дальнейшем от туда выгружать в excel

func AddToSheet(sheet Sheet) {
	sheetData = sheet
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		createAndOperate()
	} else {
		openAndOperate()
	}
}

func createAndOperate() {
	f := excelize.NewFile()

	f.SetCellValue(sheet, "A1", "#")
	f.SetCellValue(sheet, "B1", "Link")
	f.SetCellValue(sheet, "C1", "DOB")
	f.SetCellValue(sheet, "D1", "Today")
	f.SetCellValue(sheet, "E1", "Age")
	f.SetCellValue(sheet, "F1", "Subs")
	f.SetCellValue(sheet, "G1", "ViewsM7")
	f.SetCellValue(sheet, "H1", "ER")

	operate(f)

	if err := f.SaveAs(filename); err != nil {
		config.Logger.Error(err.Error())
		return
	}
}

func openAndOperate() {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		config.Logger.Error(err.Error())
		return
	}

	operate(f)

	defer f.Save()
}

func operate(f *excelize.File) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		config.Logger.Error(err.Error())
		return
	}

	nextRow := len(rows) + 1

	f.SetCellValue(sheet, fmt.Sprintf("A%d", nextRow), nextRow-1)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", nextRow), sheetData.Link)
	f.SetCellValue(sheet, fmt.Sprintf("C%d", nextRow), sheetData.DOB)
	f.SetCellValue(sheet, fmt.Sprintf("D%d", nextRow), time.Now().Format("2006.01.02"))
	f.SetCellValue(sheet, fmt.Sprintf("E%d", nextRow), sheetData.Age)
	f.SetCellValue(sheet, fmt.Sprintf("F%d", nextRow), sheetData.Subs)
	f.SetCellValue(sheet, fmt.Sprintf("G%d", nextRow), sheetData.ViewsM7)
	f.SetCellValue(sheet, fmt.Sprintf("H%d", nextRow), sheetData.ER)

}
