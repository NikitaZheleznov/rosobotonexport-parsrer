package main

import (
	"log"
	"rosoboronexport-parser/internal/collector"
	"rosoboronexport-parser/internal/config"
	"rosoboronexport-parser/internal/models"
	"rosoboronexport-parser/internal/storage"
)

func main() {

	cfg := config.DefaultConfig()
	matchCollector := collector.NewMatchCollector(cfg.RequestDelay)
	exporter := storage.NewExcelExporter(cfg.OutputDir)
	allSeasonsData := make(map[string][]models.Match)

	// Обрабатываем каждый сезон
	for _, season := range cfg.Seasons {
		log.Printf("Обработка сезона %s (ID: %s)", season.Name, season.ID)

		// Шаг 1: Собираем ID всех матчей за сезон
		games, err := matchCollector.FetchGames(season.ID)
		if err != nil {
			log.Printf("Ошибка сбора ID матчей для сезона %s: %v", season.Name, err)
			continue
		}

		log.Printf("Найдено %d матчей для сезона %s", len(games), season.Name)

		// Шаг 2: Собираем заявки по каждому матчу
		matches, err := matchCollector.CollectApplications(games, season)
		if err != nil {
			log.Printf("Ошибка сбора заявок для сезона %s: %v", season.Name, err)
			continue
		}

		allSeasonsData[season.Name] = matches
		log.Printf("Сезон %s обработан, собрано %d матчей", season.Name, len(matches))
	}

	// Шаг 3: Экспортируем все данные в Excel
	if err := exporter.ExportAllSeasons(allSeasonsData); err != nil {
		log.Fatalf("Ошибка экспорта в Excel: %v", err)
	}

	log.Println("Готово! Все файлы сохранены в директории", cfg.OutputDir)
}
