package models

import "time"

// Player представляет игрока в заявке
type Player struct {
	Number      int    `json:"number"`       // Номер игрока
	Name        string `json:"name"`         // Полное имя
	Position    string `json:"position"`     // Позиция (Нападающий/Защитник/Вратарь)
	BirthDate   string `json:"birth_date"`   // Дата рождения (опционально)
	IsCaptain   bool   `json:"is_captain"`   // Капитан?
	IsAssistant bool   `json:"is_assistant"` // Ассистент?
	Index       string `json:"index"`        // Индекс игрока (если есть)
	Status      string `json:"status"`       // Статус (допущен/не допущен)
}

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
	Name string // "2021-2022"
	ID   string // "14", "15" и т.д.
	URL  string // Полный URL страницы команды за сезон
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
