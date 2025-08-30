package excel

import (
	"fmt"
	"learn_bot_admin_panel/tools/sql_null"
	"reflect"
	"time"
	"unicode/utf8"

	"github.com/tealeg/xlsx"
)

type CellType string

const (
	CellTypeString        CellType = "string"
	CellTypeInt           CellType = "int"
	CellTypeDate          CellType = "date"
	CellTypeBool          CellType = "bool"
	CellTypeSqlNullString CellType = "sql_null_string"
	CellTypeSqlNullTime   CellType = "sql_null_time"
)

var SetCellByTypeMap = map[CellType]func(cell *xlsx.Cell, value reflect.Value){
	CellTypeString: func(cell *xlsx.Cell, value reflect.Value) {
		cell.SetString(value.String())
	},

	CellTypeInt: func(cell *xlsx.Cell, value reflect.Value) {
		cell.SetInt(value.Interface().(int))
	},

	CellTypeDate: func(cell *xlsx.Cell, value reflect.Value) {
		cell.SetDate(value.Interface().(time.Time))
	},

	CellTypeBool: func(cell *xlsx.Cell, value reflect.Value) {
		localizedVal := "Нет"
		if value.Interface().(bool) {
			localizedVal = "Да"
		}

		cell.SetString(localizedVal)
	},

	CellTypeSqlNullString: func(cell *xlsx.Cell, value reflect.Value) {
		typeVal := value.Interface().(sql_null.NullString)
		cell.SetString(typeVal.OptionalResult())
	},

	CellTypeSqlNullTime: func(cell *xlsx.Cell, value reflect.Value) {
		typeVal := value.Interface().(sql_null.NullTime)

		if typeVal.Valid {
			cell.SetDate(typeVal.Time)
			return
		}

		cell.SetString("-")
	},
}

func BuildExcelFileFromStruct[T any](data []T, sheetName string) (*xlsx.File, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		return nil, err
	}

	headerRow := sheet.AddRow()

	t := reflect.TypeOf(data[0])

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("excel_head")
		if tag == "" {
			continue
		}

		cell := headerRow.AddCell()
		cell.Value = tag

		headerRow.Sheet.SetColWidth(i, i, columnAutoWidth(tag))
	}

	for _, val := range data {
		val := reflect.ValueOf(val)

		row := sheet.AddRow()
		for i := 0; i < t.NumField(); i++ {
			cType := t.Field(i).Tag.Get("excel_cell")
			if cType == "" {
				continue
			}

			cell := row.AddCell()
			excelCellType := CellType(cType)

			setCellByType(cell, excelCellType, val.Field(i))
		}
	}

	return file, nil
}

func setCellByType(cell *xlsx.Cell, cType CellType, value reflect.Value) {
	set, exsits := SetCellByTypeMap[cType]
	if !exsits {
		cell.SetString(fmt.Sprintf("%v", value.Interface()))
		return
	}

	set(cell, value)
}

func columnAutoWidth(header string) (width float64) {
	length := utf8.RuneCountInString(header)

	if length < 10 {
		length = 11
	} else if length >= 30 {
		length = 15
	}

	return float64(length) * 1.2
}
