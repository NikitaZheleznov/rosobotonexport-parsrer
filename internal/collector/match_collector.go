package collector

import (
	"fmt"
	"rosoboronexport-parser/internal/models"
	"rosoboronexport-parser/internal/parser"
	"time"

	"github.com/gocolly/colly"
)

type MatchCollector struct {
	collector *colly.Collector
	parser    *parser.MatchParser
}

func NewMatchCollector(delay time.Duration) *MatchCollector {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
		colly.AllowedDomains("hltr.ru"),
	)

	// Настройка задержек
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*hltr.ru*",
		Delay:       delay,
		RandomDelay: 2 * time.Second,
	})

	return &MatchCollector{
		collector: c,
		parser:    parser.NewMatchParser(),
	}
}

// CollectMatchIDs собирает ID всех матчей команды за указанный сезон
func (mc *MatchCollector) CollectMatchIDs(season models.Season) ([]string, error) {
	url := fmt.Sprintf("%s/teams/%s?season_id=%s",
		"https://hltr.ru", "387", season.ID)

	var matchIDs []string
	var err error

	// Обработчик страницы с матчами
	mc.collector.OnHTML("div.team-matches", func(e *colly.HTMLElement) {
		// Парсим все ссылки на матчи
		ids, parseErr := mc.parser.ExtractMatchIDs(e)
		if parseErr != nil {
			err = parseErr
			return
		}
		matchIDs = append(matchIDs, ids...)
	})

	// Обработчик ошибок
	mc.collector.OnError(func(_ *colly.Response, rErr error) {
		err = rErr
	})

	// Запускаем сбор
	mc.collector.Visit(url)
	mc.collector.Wait()

	return matchIDs, err
}
