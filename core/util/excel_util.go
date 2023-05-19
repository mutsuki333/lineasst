/*
	excel.go
	Purpose: excel helpers.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/04/21  v1.0.0 Evan Chen   Initial release
*/

package util

import (
	"strings"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

// Cellname converts coordinates to CellName. This is a wrapper of excelize.CoordinatesToCellName
//
// ex. (1, 1) => "A1"
func Cellname(col int, row int) string {
	name, _ := excelize.CoordinatesToCellName(col, row)
	return name
}

// AutoFitColWidthWithRatio is based on https://github.com/qax-os/excelize/pull/1386
//
// AutoFitColWidth provides a function to autofit columns according to
// their text content with default font size and font.
// Note: this only works on the column with cells which not contains
// formula cell and style with a number format.
//
// For example set column of column H on Sheet1:
//
// err = f.AutoFitColWidth("Sheet1", "H")
//
// Set style of columns C:F on Sheet1:
//
// err = f.AutoFitColWidth("Sheet1", "C:F")
func AutoFitColWidthWithRatio(f *excelize.File, sheetName, columns string, ratio float64) error {
	startColIdx, endColIdx, err := parseColRange(f, columns)
	if err != nil {
		return err
	}

	cols, err := f.Cols(sheetName)
	if err != nil {
		return err
	}

	colIdx := 1
	for cols.Next() {
		if colIdx >= startColIdx && colIdx <= endColIdx {
			rowCells, _ := cols.Rows()
			var max int
			for i := range rowCells {
				rowCell := rowCells[i]

				// Single Byte Character Set(SBCS) is 1 holds
				// Multi-Byte Character System(MBCS) is 2 holds
				var cellLenSBCS, cellLenMBCS int
				for ii := range rowCell {
					if rowCell[ii] < 0x80 {
						cellLenSBCS++
					}
				}

				runeLen := utf8.RuneCountInString(rowCell)
				cellLenMBCS = runeLen - cellLenSBCS

				cellWidth := cellLenSBCS + cellLenMBCS*2
				if cellWidth > max {
					max = cellWidth
				}
			}

			// The ratio of 1.123 is the best approximation I tried my best to
			// find.
			actualMax := float64(max) * ratio

			if actualMax < 9.140625 {
				actualMax = 9.140625
			} else if actualMax >= excelize.MaxColumnWidth {
				actualMax = excelize.MaxColumnWidth
			}
			colName, _ := excelize.ColumnNumberToName(colIdx)

			if err := f.SetColWidth(sheetName, colName, colName, actualMax); err != nil {
				return err
			}
		}

		// fast go away.
		if colIdx == endColIdx {
			break
		}

		colIdx++
	}

	return nil
}

// AutoFitColWidth is based on https://github.com/qax-os/excelize/pull/1386
//
// AutoFitColWidth provides a function to autofit columns according to
// their text content with default font size and font.
// Note: this only works on the column with cells which not contains
// formula cell and style with a number format.
//
// For example set column of column H on Sheet1:
//
// err = f.AutoFitColWidth("Sheet1", "H")
//
// Set style of columns C:F on Sheet1:
//
// err = f.AutoFitColWidth("Sheet1", "C:F")
func AutoFitColWidth(f *excelize.File, sheetName, columns string) error {
	return AutoFitColWidthWithRatio(f, sheetName, columns, 1.123)
}

// parseColRange parse and convert column range with column name to the column number.
func parseColRange(f *excelize.File, columns string) (min, max int, err error) {
	colsTab := strings.Split(columns, ":")
	min, err = excelize.ColumnNameToNumber(colsTab[0])
	if err != nil {
		return
	}
	max = min
	if len(colsTab) == 2 {
		if max, err = excelize.ColumnNameToNumber(colsTab[1]); err != nil {
			return
		}
	}
	if max < min {
		min, max = max, min
	}
	return
}

func FreezeHeader(file *excelize.File, sheet string) error {
	return file.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})
}
