package task

import (
	"fmt"
	"time"
)

// RecurrenceGenerator генерирует даты выполнения задач на основе конфигурации
type RecurrenceGenerator struct {
	config RecurrenceConfig
}

// NewRecurrenceGenerator создает новый генератор
func NewRecurrenceGenerator(config RecurrenceConfig) (*RecurrenceGenerator, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &RecurrenceGenerator{config: config}, nil
}

// GenerateNextOccurrences генерирует даты следующих N выполнений задачи
// startDate - дата начала периодичности
// limit - количество дат для генерирования
func (rg *RecurrenceGenerator) GenerateNextOccurrences(startDate time.Time, limit int) []time.Time {
	if rg.config.RecurrenceType == RecurrenceTypeNone {
		return nil
	}

	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Проверка max_occurrences
	if rg.config.MaxOccurrences != nil && limit > *rg.config.MaxOccurrences {
		limit = *rg.config.MaxOccurrences
	}

	var occurrences []time.Time

	switch rg.config.RecurrenceType {
	case RecurrenceTypeDaily:
		occurrences = rg.generateDaily(startDate, limit)
	case RecurrenceTypeMonthly:
		occurrences = rg.generateMonthly(startDate, limit)
	case RecurrenceTypeWeekly:
		occurrences = rg.generateWeekly(startDate, limit)
	case RecurrenceTypeOddEven:
		occurrences = rg.generateOddEven(startDate, limit)
	case RecurrenceTypeSpecific:
		occurrences = rg.generateSpecific(startDate, limit)
	}

	// Фильтруем по endDate если указана
	if rg.config.EndDate != nil {
		filtered := make([]time.Time, 0)
		for _, t := range occurrences {
			if t.Before(*rg.config.EndDate) || t.Equal(*rg.config.EndDate) {
				filtered = append(filtered, t)
			}
		}
		occurrences = filtered
	}

	return occurrences
}

// generateDaily генерирует ежедневные даты
func (rg *RecurrenceGenerator) generateDaily(startDate time.Time, limit int) []time.Time {
	var occurrences []time.Time
	interval := 1
	if rg.config.IntervalDays != nil {
		interval = *rg.config.IntervalDays
	}

	current := startDate
	for len(occurrences) < limit {
		occurrences = append(occurrences, current)
		current = current.AddDate(0, 0, interval)
	}

	return occurrences
}

// generateMonthly генерирует ежемесячные даты на определенное число
func (rg *RecurrenceGenerator) generateMonthly(startDate time.Time, limit int) []time.Time {
	var occurrences []time.Time

	dayOfMonth := 1
	if rg.config.MonthlyDayOfMonth != nil {
		dayOfMonth = *rg.config.MonthlyDayOfMonth
	}

	// Начинаем с первого дня месяца startDate
	current := time.Date(startDate.Year(), startDate.Month(), 1, startDate.Hour(), startDate.Minute(), startDate.Second(), startDate.Nanosecond(), startDate.Location())

	// Если startDate уже прошел в текущем месяце, начинаем со следующего
	if startDate.Day() > dayOfMonth {
		current = current.AddDate(0, 1, 0)
	}

	for len(occurrences) < limit {
		// Устанавливаем день месяца
		year, month, _ := current.Date()
		nextDate := time.Date(year, month, dayOfMonth, startDate.Hour(), startDate.Minute(), startDate.Second(), startDate.Nanosecond(), startDate.Location())

		// Если день больше чем дней в месяце (например 31 число в феврале), используем последний день
		daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, startDate.Location()).Day()
		if dayOfMonth > daysInMonth {
			nextDate = time.Date(year, month+1, 0, startDate.Hour(), startDate.Minute(), startDate.Second(), startDate.Nanosecond(), startDate.Location())
		}

		if nextDate.After(startDate) || nextDate.Equal(startDate) {
			occurrences = append(occurrences, nextDate)
		}

		current = current.AddDate(0, 1, 0)
	}

	return occurrences
}

// generateWeekly генерирует еженедельные даты на конкретные дни
func (rg *RecurrenceGenerator) generateWeekly(startDate time.Time, limit int) []time.Time {
	var occurrences []time.Time

	if rg.config.WeeklyDays == nil || len(*rg.config.WeeklyDays) == 0 {
		return occurrences
	}

	current := startDate
	for len(occurrences) < limit {
		// Проверяем день недели (0 = Monday в ISO 8601)
		weekDay := int(current.Weekday())
		if weekDay == 0 {
			weekDay = 7 // Sunday
		}
		weekDay-- // Convert to 0-6 where 0 is Monday

		for _, targetDay := range *rg.config.WeeklyDays {
			if weekDay == targetDay {
				occurrences = append(occurrences, current)
				break
			}
		}

		if len(occurrences) >= limit {
			break
		}

		current = current.AddDate(0, 0, 1)
	}

	return occurrences
}

// generateOddEven генерирует даты четных/нечетных дней
func (rg *RecurrenceGenerator) generateOddEven(startDate time.Time, limit int) []time.Time {
	var occurrences []time.Time

	oddEvenType := DayTypeOdd
	if rg.config.OddEvenType != nil {
		oddEvenType = *rg.config.OddEvenType
	}

	current := startDate
	for len(occurrences) < limit {
		day := current.Day()
		isOdd := day%2 == 1

		if (oddEvenType == DayTypeOdd && isOdd) || (oddEvenType == DayTypeEven && !isOdd) {
			occurrences = append(occurrences, current)
		}

		current = current.AddDate(0, 0, 1)
	}

	return occurrences
}

// generateSpecific генерирует даты для конкретных дней
func (rg *RecurrenceGenerator) generateSpecific(startDate time.Time, limit int) []time.Time {
	var occurrences []time.Time

	if rg.config.SpecificDates == nil {
		return occurrences
	}

	// Парсим все даты и сортируем их
	dateMap := make(map[string]time.Time)
	for _, dateStr := range *rg.config.SpecificDates {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		// Установим время из startDate
		t = time.Date(t.Year(), t.Month(), t.Day(), startDate.Hour(), startDate.Minute(), startDate.Second(), startDate.Nanosecond(), startDate.Location())
		if t.After(startDate) || t.Equal(startDate) {
			dateMap[dateStr] = t
		}
	}

	// Рекурсивно вызываем в течение нескольких лет для обработки годичного цикла
	for year := startDate.Year(); year <= startDate.Year()+10 && len(occurrences) < limit; year++ {
		for _, dateStr := range *rg.config.SpecificDates {
			if len(occurrences) >= limit {
				break
			}
			t, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				continue
			}

			// Переносим дату на нужный год
			t = time.Date(year, t.Month(), t.Day(), startDate.Hour(), startDate.Minute(), startDate.Second(), startDate.Nanosecond(), startDate.Location())

			if t.After(startDate) || t.Equal(startDate) {
				// Проверяем, что мы еще не добавили эту дату
				found := false
				for _, existing := range occurrences {
					if existing.Equal(t) {
						found = true
						break
					}
				}
				if !found {
					occurrences = append(occurrences, t)
				}
			}
		}
	}

	return occurrences
}

// GetNextOccurrence получает следующую дату выполнения после заданной
func (rg *RecurrenceGenerator) GetNextOccurrence(afterDate time.Time) *time.Time {
	occurrences := rg.GenerateNextOccurrences(afterDate.AddDate(0, 0, 1), 1)
	if len(occurrences) > 0 {
		return &occurrences[0]
	}
	return nil
}

func (rg *RecurrenceGenerator) String() string {
	return fmt.Sprintf("RecurrenceGenerator{type: %s}", rg.config.RecurrenceType)
}
