package models

import "time"

// Match представляет матч с заявкой
type Match struct {
	ID         string    `json:"id"`          // ID матча
	Date       string    `json:"date"`        // Дата матча (ДД.ММ.ГГГГ)
	Team       string    `json:"team"`        // Всегда "Рособоронэкспорт"
	Opponent   string    `json:"opponent"`    // Команда соперника
	SeasonID   string    `json:"season_id"`   // ID сезона
	SeasonName string    `json:"season_name"` // Название сезона (2021-2022 и т.д.)
	Players    []Player  `json:"players"`     // Список игроков в заявке
	CreatedAt  time.Time `json:"created_at"`
}

// Season представляет сезон с его ID
type Season struct {
	Name       string
	ID         string
	Tournament string
	URL        string
}

// ExcelRow представляет строку для Excel-файла
type ExcelRow struct {
	MatchNumber  string `col:"№ матча"`    // Номер матча/ID
	Team         string `col:"Команда"`    // Рособоронэкспорт
	Opponent     string `col:"Соперник"`   // Название команды соперника
	MatchDate    string `col:"Дата матча"` // 02.02.2021
	PlayerNumber int    `col:"№ игрока"`   // Номер игрока
	PlayerName   string `col:"Имя игрока"` // Фамилия Имя
	Position     string `col:"Позиция"`    // Нападающий/Защитник/Вратарь
}

// APIGame описывает один матч из массива ответа
type APIGame struct {
	ID             int          `json:"id"`
	Datetime       string       `json:"datetime"`
	Team           APITeamBrief `json:"team"`
	CompetitorTeam APITeamBrief `json:"competitor_team"`
}

type APITeamBrief struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MatchApplication представляет полную заявку на матч
// type MatchApplication struct {
// 	GameID   string   `json:"game_id"`
// 	HomeTeam string   `json:"home_team"` // Название команды хозяев
// 	AwayTeam string   `json:"away_team"` // Название команды гостей
// 	Players  []Player `json:"players"`   // Список игроков ВАШЕЙ команды
// 	GameDate string   `json:"game_date"`
// 	GameTime string   `json:"game_time"`
// 	Location string   `json:"location"`
// }

// Player (если у вас ещё нет такой полной структуры)
type Player struct {
	Number   int    `json:"number"`
	Name     string `json:"name"`
	Position string `json:"position"`
}
