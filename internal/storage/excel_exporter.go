package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"rosoboronexport-parser/internal/models"

	"github.com/xuri/excelize/v2"
)

type ExcelExporter struct {
	outputDir string
}

func NewExcelExporter(outputDir string) *ExcelExporter {
	os.MkdirAll(outputDir, 0755)
	return &ExcelExporter{outputDir: outputDir}
}

func (e *ExcelExporter) ExportAllSeasons(seasonsData map[string][]models.Match) error {
	for seasonName, matches := range seasonsData {
		if err := e.exportSeason(seasonName, matches); err != nil {
			return fmt.Errorf("ошибка экспорта сезона %s: %w", seasonName, err)
		}
	}
	return nil
}

func (e *ExcelExporter) exportSeason(seasonName string, matches []models.Match) error {
	filename := filepath.Join(e.outputDir, fmt.Sprintf("rosoboronexport_%s.xlsx", seasonName))

	f := excelize.NewFile()
	defer f.Close()

	// Создаем стили
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	sheetName := "Заявки"
	f.SetSheetName("Sheet1", sheetName)

	// Заголовки
	headers := []string{"№ матча", "Команда", "Соперник", "Дата матча", "№ игрока", "Имя игрока", "Позиция"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Данные
	rowNum := 2
	for _, match := range matches {
		for _, player := range match.Players {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), match.ID)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), match.Team)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), match.Opponent)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), match.Date)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), player.Number)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), player.Name)
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), player.Position)

			// Применяем стиль к строке
			for col := 0; col < len(headers); col++ {
				cell := fmt.Sprintf("%s%d", string(rune('A'+col)), rowNum)
				f.SetCellStyle(sheetName, cell, cell, dataStyle)
			}

			rowNum++
		}
	}

	// Автоширина колонок
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		width := e.calculateColumnWidth(f, sheetName, col, rowNum-1)
		f.SetColWidth(sheetName, col, col, width)
	}

	return f.SaveAs(filename)
}

func (e *ExcelExporter) calculateColumnWidth(f *excelize.File, sheet, col string, maxRow int) float64 {
	maxWidth := float64(10) // минимальная ширина

	for row := 1; row <= maxRow; row++ {
		cell := fmt.Sprintf("%s%d", col, row)
		value, _ := f.GetCellValue(sheet, cell)
		if width := float64(len(value) + 2); width > maxWidth {
			maxWidth = width
		}
	}

	if maxWidth > 50 {
		return 50
	}
	return maxWidth
}
