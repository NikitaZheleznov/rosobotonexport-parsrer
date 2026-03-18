package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"rosoboronexport-parser/internal/models"

	"github.com/PuerkitoBio/goquery"
)

// MatchCollector собирает данные о матчах через API
type MatchCollector struct {
	client  *http.Client
	baseURL string
	teamID  int
	matches []models.Match
}

// NewMatchCollector создает новый экземпляр MatchCollector
func NewMatchCollector(timeout time.Duration) *MatchCollector {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &MatchCollector{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://mtgame.ru/api/v1",
		teamID:  387,
		matches: make([]models.Match, 0),
	}
}

// fetchGames - внутренняя функция для запроса к API
func (mc *MatchCollector) FetchGames(seasonID string) ([]models.APIGame, error) {
	url := fmt.Sprintf("%s/tournament_season/%s/games/?team_id=%d", mc.baseURL, seasonID, mc.teamID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := mc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API вернул статус %d: %s", resp.StatusCode, string(body))
	}

	var games []models.APIGame

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &games); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return games, nil
}

func (mc *MatchCollector) CollectApplications(games []models.APIGame, season models.Season) ([]models.Match, error) {
	mc.matches = make([]models.Match, 0)
	for _, game := range games {
		match, err := mc.GetGameApplicationByID(game)
		if err != nil {
			return nil, err
		}
		match.Team = game.Team.Name
		match.Date = extractDate(game.Datetime)
		match.Opponent = game.CompetitorTeam.Name
		match.GameID = game.ID
		match.SeasonName = season.Name
		match.CreatedAt = time.Now()

		if match != nil {
			mc.matches = append(mc.matches, *match)
		}
	}
	return mc.matches, nil
}

func extractDate(dateTimeStr string) string {
	// Парсим строку в time.Time
	t, err := time.Parse(time.RFC3339, dateTimeStr)
	if err != nil {
		return ""
	}

	// Возвращаем только дату в нужном формате
	return t.Format("02.01.2006")
}

func (mc *MatchCollector) GetGameApplicationByID(game models.APIGame) (*models.Match, error) {
	url := fmt.Sprintf("%s/tournament_hockey_game/%d/game_application_file/?file_type=html&application_type=tournament_user", mc.baseURL, game.ID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	// Устанавливаем ВСЕ заголовки из браузера
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")

	resp, err := mc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API вернул статус %d: %s", resp.StatusCode, string(body))
	}

	// Парсим HTML с помощью goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	match := &models.Match{
		ID: strconv.Itoa(game.ID),
	}

	match.Players, err = ExtractSimplePlayers(doc)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func ExtractSimplePlayers(doc *goquery.Document) ([]models.Player, error) {
	var players []models.Player

	// Находим секцию Рособоронэкспорт
	doc.Find("div.page").Each(func(i int, page *goquery.Selection) {
		if strings.Contains(page.Text(), "Рособоронэкспорт") {
			// Парсим всех игроков в этой секции
			page.Find("td[style*='position: relative']").Each(func(j int, td *goquery.Selection) {
				player := models.Player{}
				text := td.Text()

				// Извлекаем номер
				reNum := regexp.MustCompile(`№(\d+)`)
				if matches := reNum.FindStringSubmatch(text); len(matches) >= 2 {
					player.Number, _ = strconv.ParseInt(matches[1], 10, 64)
				}

				// Извлекаем позицию
				if strings.Contains(text, "Нап.") {
					player.Position = "Нападающий"
				} else if strings.Contains(text, "Защ.") {
					player.Position = "Защитник"
				} else if strings.Contains(text, "Вр.") {
					player.Position = "Вратарь"
				} else {
					player.Position = ""
				}

				// Извлекаем имя
				html, _ := td.Html()
				reName := regexp.MustCompile(`<br\s*/?>\s*([^<]+?)\s*<br\s*/?>`)
				if matches := reName.FindStringSubmatch(html); len(matches) >= 2 {
					player.Name = strings.TrimSpace(matches[1])
				}

				if player.Number > 0 && player.Name != "" && player.Position != "" {
					players = append(players, player)
				}
			})
		}
	})

	return players, nil
}
