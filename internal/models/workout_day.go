package models

import (
	"fmt"
	"strings"
	"time"
)

type WorkoutDay struct {
	ID               int64 `gorm:"primaryKey;autoIncrement"`
	UserID           int64
	WorkoutDayTypeID int64
	Exercises        []Exercise `gorm:"foreignKey:WorkoutDayID;constraint:OnDelete:CASCADE"`
	StartedAt        time.Time
	EndedAt          *time.Time
	Completed        bool

	WorkoutDayType *WorkoutDayType `gorm:"foreignKey:WorkoutDayTypeID;references:ID"`
}

func (*WorkoutDay) TableName() string {
	return "workout_days"
}

func (w *WorkoutDay) Status() string {
	if !w.Completed {
		return fmt.Sprintf("🟡 Активна")
	}
	if w.EndedAt != nil {
		return fmt.Sprintf("✅ Завершена в %s", w.EndedAt.Add(3*time.Hour).Format("15:04"))
	}

	return fmt.Sprintf("✅ Завершена")
}

func (w *WorkoutDay) String() string {
	var text strings.Builder

	text.WriteString(fmt.Sprintf("<b>Тип:</b> %s \n", w.WorkoutDayType.Name))
	text.WriteString(fmt.Sprintf("<b>Статус:</b> %s\n", w.Status()))
	text.WriteString(fmt.Sprintf("<b>Дата:</b> 📅 %s\n\n", w.StartedAt.Add(3*time.Hour).Format("02.01.2006")))

	if len(w.Exercises) > 0 {
		text.WriteString("<b>Упражнения:</b>\n")
	}

	for i, exercise := range w.Exercises {
		exerciseObj := exercise.ExerciseType
		text.WriteString(fmt.Sprintf("<b>%s %d. %s</b>\n", exercise.Status(), i+1, exerciseObj.Name))

		for _, set := range exercise.Sets {
			text.WriteString(set.String(w.Completed))
		}
		text.WriteString("\n")
	}

	return text.String()
}
