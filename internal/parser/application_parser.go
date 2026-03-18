package parser

import (
	"regexp"
	"rosoboronexport-parser/internal/models"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type ApplicationParser struct{}

func NewApplicationParser() *ApplicationParser {
	return &ApplicationParser{}
}

// ParseGameApplication парсит страницу заявки и возвращает список игроков
func (ap *ApplicationParser) ParseGameApplication(e *colly.HTMLElement) ([]models.Player, error) {
	var players []models.Player

	// Находим секцию с нашей командой
	e.ForEach("div.team-section", func(_ int, section *colly.HTMLElement) {
		// Проверяем, что это "Рособоронэкспорт"
		teamName := section.ChildText("h4")
		if !strings.Contains(teamName, "Рособоронэкспорт") {
			return
		}

		// Парсим полевых игроков
		section.ForEach("div.field-player", func(_ int, playerEl *colly.HTMLElement) {
			player := models.Player{
				Number:   ap.extractNumber(playerEl),
				Name:     ap.extractName(playerEl),
				Position: "Полевой игрок", // Уточним позицию
				Status:   ap.extractStatus(playerEl),
			}
			players = append(players, player)
		})

		// Парсим вратарей
		section.ForEach("div.goalkeeper", func(_ int, playerEl *colly.HTMLElement) {
			player := models.Player{
				Number:   ap.extractNumber(playerEl),
				Name:     ap.extractName(playerEl),
				Position: "Вратарь",
				Status:   ap.extractStatus(playerEl),
			}
			players = append(players, player)
		})
	})

	return players, nil
}

// ParseOpponent парсит название команды соперника
func (ap *ApplicationParser) ParseOpponent(e *colly.HTMLElement) string {
	var opponent string
	e.ForEach("div.team-section", func(i int, section *colly.HTMLElement) {
		if i == 1 { // Вторая секция - соперник
			opponent = section.ChildText("h4")
		}
	})
	return opponent
}

func (ap *ApplicationParser) extractNumber(e *colly.HTMLElement) int {
	// В примере заявки номер был в формате "№2" внутри span
	// Ищем элемент с номером
	var numberText string

	// Пробуем разные варианты, которые могут быть на странице
	e.ForEach("span", func(_ int, el *colly.HTMLElement) {
		text := el.Text
		if strings.Contains(text, "№") || strings.Contains(text, "#") {
			numberText = text
		}
	})

	// Если не нашли в span, ищем в div с классом
	if numberText == "" {
		numberText = e.ChildText("div.player-info span:first-child")
	}

	// Если всё ещё не нашли, ищем в тексте перед именем
	if numberText == "" {
		// В некоторых реализациях номер может быть отдельно
		fullText := e.Text
		// Ищем паттерн "| №2 - Нап." или похожий
		re := regexp.MustCompile(`[|]\s*[№#]?(\d+)`)
		if matches := re.FindStringSubmatch(fullText); len(matches) > 1 {
			numberText = matches[1]
		}
	}

	// Парсим найденный текст
	if numberText != "" {
		// Убираем все не-цифры
		re := regexp.MustCompile(`\d+`)
		if numStr := re.FindString(numberText); numStr != "" {
			if num, err := strconv.Atoi(numStr); err == nil {
				return num
			}
		}
	}

	return 0
}

func (ap *ApplicationParser) extractName(e *colly.HTMLElement) string {
	return e.ChildText("span.name")
}

func (ap *ApplicationParser) extractStatus(e *colly.HTMLElement) string {
	statusEl := e.ChildText("span.status")
	if statusEl == "" {
		return "допущен"
	}
	return statusEl
}
