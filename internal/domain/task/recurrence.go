package task

import (
	"fmt"
	"time"
)

// RecurrenceType определяет тип периодичности задачи
type RecurrenceType string

const (
	RecurrenceTypeNone        RecurrenceType = "none"        // Задача не периодическая
	RecurrenceTypeDaily       RecurrenceType = "daily"       // Ежедневно каждый N-й день
	RecurrenceTypeMonthly     RecurrenceType = "monthly"     // Ежемесячно на определенное число
	RecurrenceTypeWeekly      RecurrenceType = "weekly"      // Еженедельно в определенные дни недели
	RecurrenceTypeOddEven     RecurrenceType = "odd_even"    // Четные/нечетные дни месяца
	RecurrenceTypeSpecific    RecurrenceType = "specific"    // На конкретные даты
)

// Valid проверяет валидность типа
func (rt RecurrenceType) Valid() bool {
	switch rt {
	case RecurrenceTypeNone, RecurrenceTypeDaily, RecurrenceTypeMonthly,
		RecurrenceTypeWeekly, RecurrenceTypeOddEven, RecurrenceTypeSpecific:
		return true
	default:
		return false
	}
}

// DayType для четных/нечетных дней
type DayType string

const (
	DayTypeEven DayType = "even"
	DayTypeOdd  DayType = "odd"
)

// Valid проверяет валидность типа дня
func (dt DayType) Valid() bool {
	return dt == DayTypeEven || dt == DayTypeOdd
}

// RecurrenceConfig хранит конфигурацию периодичности
type RecurrenceConfig struct {
	RecurrenceType RecurrenceType `json:"recurrence_type"`

	// Для RecurrenceTypeDaily: каждый N-й день
	IntervalDays *int `json:"interval_days,omitempty"`

	// Для RecurrenceTypeMonthly: число месяца (1-30)
	MonthlyDayOfMonth *int `json:"monthly_day_of_month,omitempty"`

	// Для RecurrenceTypeWeekly: дни недели (0=Пн, 1=Вт, ... 6=Вс)
	WeeklyDays *[]int `json:"weekly_days,omitempty"`

	// Для RecurrenceTypeOddEven: четные или нечетные дни
	OddEvenType *DayType `json:"odd_even_type,omitempty"`

	// Для RecurrenceTypeSpecific: конкретные даты (YYYY-MM-DD)
	SpecificDates *[]string `json:"specific_dates,omitempty"`

	// Окончание периодичности (если nil - бесконечная)
	EndDate *time.Time `json:"end_date,omitempty"`

	// Количество повторений (если 0 - без ограничений)
	MaxOccurrences *int `json:"max_occurrences,omitempty"`
}

// Validate проверяет валидность конфигурации
func (rc *RecurrenceConfig) Validate() error {
	if !rc.RecurrenceType.Valid() {
		return fmt.Errorf("invalid recurrence type: %s", rc.RecurrenceType)
	}

	switch rc.RecurrenceType {
	case RecurrenceTypeDaily:
		if rc.IntervalDays == nil || *rc.IntervalDays < 1 {
			return fmt.Errorf("interval_days must be at least 1")
		}
	case RecurrenceTypeMonthly:
		if rc.MonthlyDayOfMonth == nil || *rc.MonthlyDayOfMonth < 1 || *rc.MonthlyDayOfMonth > 31 {
			return fmt.Errorf("monthly_day_of_month must be between 1 and 31")
		}
	case RecurrenceTypeWeekly:
		if rc.WeeklyDays == nil || len(*rc.WeeklyDays) == 0 {
			return fmt.Errorf("weekly_days must not be empty")
		}
		for _, day := range *rc.WeeklyDays {
			if day < 0 || day > 6 {
				return fmt.Errorf("weekly day must be between 0 and 6")
			}
		}
	case RecurrenceTypeOddEven:
		if rc.OddEvenType == nil || !rc.OddEvenType.Valid() {
			return fmt.Errorf("invalid odd_even_type")
		}
	case RecurrenceTypeSpecific:
		if rc.SpecificDates == nil || len(*rc.SpecificDates) == 0 {
			return fmt.Errorf("specific_dates must not be empty")
		}
		// Валидация формата дат
		for _, dateStr := range *rc.SpecificDates {
			if _, err := time.Parse("2006-01-02", dateStr); err != nil {
				return fmt.Errorf("invalid date format in specific_dates: %s", dateStr)
			}
		}
	}

	// Проверка окончания периодичности
	if rc.EndDate != nil && rc.MaxOccurrences != nil {
		if *rc.MaxOccurrences < 1 {
			return fmt.Errorf("max_occurrences must be at least 1 if specified")
		}
	}

	return nil
}
