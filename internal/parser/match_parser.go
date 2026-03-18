package parser

import (
    "regexp"
    "github.com/gocolly/colly"
)

type MatchParser struct {
    matchIDRegex *regexp.Regexp
}

func NewMatchParser() *MatchParser {
    return &MatchParser{
        // Регулярное выражение для извлечения ID матча из ссылки
        matchIDRegex: regexp.MustCompile(`/game/(\d+)`),
    }
}

// ExtractMatchIDs извлекает ID матчей из HTML
func (mp *MatchParser) ExtractMatchIDs(e *colly.HTMLElement) ([]string, error) {
    var ids []string
    
    // Ищем все ссылки на страницы матчей
    e.ForEach("a.match-link", func(_ int, el *colly.HTMLElement) {
        href := el.Attr("href")
        if matches := mp.matchIDRegex.FindStringSubmatch(href); len(matches) > 1 {
            ids = append(ids, matches[1])
        }
    })
    
    return ids, nil
}

// ExtractMatchDate извлекает дату матча
func (mp *MatchParser) ExtractMatchDate(e *colly.HTMLElement) string {
    // Ищем элемент с датой
    dateEl := e.ChildText("div.match-date")
    if dateEl == "" {
        // Альтернативный селектор
        dateEl = e.ChildText("span.date")
    }
    return dateEl
}

// ExtractOpponent извлекает название команды соперника
func (mp *MatchParser) ExtractOpponent(e *colly.HTMLElement) string {
    // Ищем название команды соперника
    // Обычно это второй элемент в списке команд
    var opponents []string
    e.ForEach("div.team-name", func(i int, el *colly.HTMLElement) {
        if i == 1 { // Вторая команда - соперник
            opponents = append(opponents, el.Text)
        }
    })
    
    if len(opponents) > 0 {
        return opponents[0]
    }
    return ""
}