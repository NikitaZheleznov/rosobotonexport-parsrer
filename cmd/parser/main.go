package main

import (
    "log"
    "rosoboronexport-parser/internal/collector"
    "rosoboronexport-parser/internal/config"
    "rosoboronexport-parser/internal/storage"
)

func main() {
    // Загружаем конфигурацию
    cfg := config.DefaultConfig()
    
    // Создаём коллекторы
    matchCollector := collector.NewMatchCollector(cfg.RequestDelay)
    appCollector := collector.NewApplicationCollector(cfg.Parallelism, cfg.RequestDelay)
    
    // Создаём экспортёр
    exporter := storage.NewExcelExporter(cfg.OutputDir)
    
    // Для хранения результатов всех сезонов
    allSeasonsData := make(map[string][]models.Match)
    
    // Обрабатываем каждый сезон
    for _, season := range cfg.Seasons {
        log.Printf("Обработка сезона %s (ID: %s)", season.Name, season.ID)
        
        // Шаг 1: Собираем ID всех матчей за сезон
        matchIDs, err := matchCollector.CollectMatchIDs(season)
        if err != nil {
            log.Printf("Ошибка сбора ID матчей для сезона %s: %v", season.Name, err)
            continue
        }
        
        log.Printf("Найдено %d матчей для сезона %s", len(matchIDs), season.Name)
        
        // Шаг 2: Собираем заявки по каждому матчу
        matches, err := appCollector.CollectApplications(matchIDs, season)
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