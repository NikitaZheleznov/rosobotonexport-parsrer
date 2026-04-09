package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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

func (e *ExcelExporter) ExportAllSeasons(teamName, seasonName string, matches []models.Match) error {
	if err := e.ExportMatchesToExcel(teamName, seasonName, matches); err != nil {
		return fmt.Errorf("ошибка экспорта сезона %s: %w", seasonName, err)
	}
	return nil
}

func (e *ExcelExporter) ExportMatchesToExcel(teamName, seasonName string, matches []models.Match) error {
	filename := filepath.Join(e.outputDir, fmt.Sprintf("Список_всех_составов_на_игры_сезона_%s_%s.xlsx", seasonName, teamName))

	f := excelize.NewFile()
	defer f.Close()

	for i, match := range matches {
		sheetName := strconv.Itoa(match.GameID)
		// Создаем новый лист
		index, err := f.NewSheet(sheetName)
		if err != nil {
			return fmt.Errorf("ошибка создания листа %s: %w", sheetName, err)
		}

		// Делаем первый лист активным
		if i == 0 {
			f.SetActiveSheet(index)
		}

		// Заполняем лист данными матча
		if err := e.fillMatchSheet(f, sheetName, match); err != nil {
			return fmt.Errorf("ошибка заполнения листа %s: %w", sheetName, err)
		}
	}

	// Удаляем лист по умолчанию "Sheet1"
	f.DeleteSheet("Sheet1")

	// Сохраняем файл
	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("ошибка сохранения Excel: %w", err)
	}

	return nil
}

// fillMatchSheet заполняет один лист данными матча
func (e *ExcelExporter) fillMatchSheet(f *excelize.File, sheetName string, match models.Match) error {
	// Стили
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  14,
			Color: "#1F4E79",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	dataStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "D3D3D3", Style: 1},
		},
	})

	// Заголовок матча с ID матча в скобках
	title := fmt.Sprintf("%s - %s (%s)",
		match.Team, match.Opponent, match.Date)

	// ID матча мелким шрифтом под заголовком
	subtitle := fmt.Sprintf("ID матча: %d", match.GameID)

	// Устанавливаем основной заголовок
	f.SetCellValue(sheetName, "A1", title)
	f.MergeCell(sheetName, "A1", "C1")
	f.SetCellStyle(sheetName, "A1", "C1", titleStyle)
	f.SetRowHeight(sheetName, 1, 30)

	// Добавляем ID матча в отдельной строке
	f.SetCellValue(sheetName, "A2", subtitle)
	f.MergeCell(sheetName, "A2", "C2")

	// Стиль для подзаголовка (более мелкий шрифт)
	subtitleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  10,
			Color: "#666666",
			Bold:  false,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
		},
	})
	f.SetCellStyle(sheetName, "A2", "C2", subtitleStyle)

	// Заголовки таблицы
	headers := []string{"№", "Имя игрока", "Позиция"}
	for i, header := range headers {
		col := string(rune('A' + i))
		cell := fmt.Sprintf("%s%d", col, 3) // Заголовки на 3 строке
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Данные игроков
	for i, player := range match.Players {
		rowNum := 4 + i

		// Номер
		numCell := fmt.Sprintf("A%d", rowNum)
		f.SetCellInt(sheetName, numCell, player.Number)
		f.SetCellStyle(sheetName, numCell, numCell, dataStyle)

		// Имя
		nameCell := fmt.Sprintf("B%d", rowNum)
		f.SetCellValue(sheetName, nameCell, player.Name)
		f.SetCellStyle(sheetName, nameCell, nameCell, dataStyle)

		// Позиция
		posCell := fmt.Sprintf("C%d", rowNum)
		f.SetCellValue(sheetName, posCell, player.Position)
		f.SetCellStyle(sheetName, posCell, posCell, dataStyle)
	}

	// Настройка ширины колонок
	f.SetColWidth(sheetName, "A", "A", 8)  // №
	f.SetColWidth(sheetName, "B", "B", 35) // Имя игрока
	f.SetColWidth(sheetName, "C", "C", 20) // Позиция

	// Добавляем информацию о матче
	infoStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  10,
			Color: "#666666",
		},
	})

	// Последняя строка с количеством игроков
	lastRow := 4 + len(match.Players)
	totalCell := fmt.Sprintf("A%d", lastRow+1)
	f.SetCellValue(sheetName, totalCell, fmt.Sprintf("Всего игроков: %d", len(match.Players)))
	f.MergeCell(sheetName, totalCell, fmt.Sprintf("C%d", lastRow+1))
	f.SetCellStyle(sheetName, totalCell, fmt.Sprintf("C%d", lastRow+1), infoStyle)

	return nil
}
