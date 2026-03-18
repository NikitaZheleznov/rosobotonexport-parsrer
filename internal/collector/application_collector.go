package collector

import (
	"fmt"
	"rosoboronexport-parser/internal/models"
	"rosoboronexport-parser/internal/parser"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type ApplicationCollector struct {
	collector *colly.Collector
	parser    *parser.ApplicationParser
	matches   []models.Match
	mu        sync.Mutex
}

func NewApplicationCollector(parallelism int, delay time.Duration) *ApplicationCollector {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0"),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: parallelism,
		Delay:       delay,
	})

	return &ApplicationCollector{
		collector: c,
		parser:    parser.NewApplicationParser(),
		matches:   make([]models.Match, 0),
	}
}

// CollectApplications собирает заявки для списка ID матчей
func (ac *ApplicationCollector) CollectApplications(matchIDs []string, season models.Season) ([]models.Match, error) {
	// Обработчик успешного ответа
	ac.collector.OnHTML("div.game-application", func(e *colly.HTMLElement) {
		// Парсим игроков
		players, err := ac.parser.ParseGameApplication(e)
		if err != nil {
			return
		}

		// Парсим соперника
		opponent := ac.parser.ParseOpponent(e)

		// Создаём объект матча
		match := models.Match{
			ID:         e.Request.URL.Query().Get("id"), // Нужно извлечь из URL
			Team:       "Рособоронэкспорт",
			Opponent:   opponent,
			SeasonID:   season.ID,
			SeasonName: season.Name,
			Players:    players,
			CreatedAt:  time.Now(),
		}

		ac.mu.Lock()
		ac.matches = append(ac.matches, match)
		ac.mu.Unlock()
	})

	// Запускаем сбор для всех ID
	for _, id := range matchIDs {
		url := fmt.Sprintf("https://mtgame.ru/api/v1/tournament_hockey_game/%s/game_application_file/?file_type=html&application_type=tournament_user", id)
		ac.collector.Visit(url)
	}

	ac.collector.Wait()
	return ac.matches, nil
}
