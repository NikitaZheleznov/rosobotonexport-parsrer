package config

import (
	"rosoboronexport-parser/internal/models"
	"time"
)

type Config struct {
	BaseURL      string          // Базовый URL сайта
	TeamID       string          // ID команды (387)
	Seasons      []models.Season // Список сезонов для парсинга
	RequestDelay time.Duration   // Задержка между запросами
	Parallelism  int             // Количество параллельных запросов
	OutputDir    string          // Директория для сохранения Excel
	UserAgent    string          // User-Agent для запросов
}

func DefaultConfig() *Config {
	return &Config{
		BaseURL: "https://hltr.ru",
		TeamID:  "387",
		Seasons: []models.Season{
			{Name: "2021-2022", ID: "14"}, // ID нужно уточнить
			{Name: "2022-2023", ID: "15"},
			{Name: "2023-2024", ID: "16"},
			{Name: "2024-2025", ID: "17"},
			{Name: "2025-2026", ID: "18"},
		},
		RequestDelay: 1 * time.Second,
		Parallelism:  5,
		OutputDir:    "./output",
		UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}
}
