package models

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
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

	User           *User           `gorm:"foreignKey:UserID;references:ID"`
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

	text.WriteString(fmt.Sprintf("<b>День:</b> <u>%s</u> \n", w.WorkoutDayType.Name))
	text.WriteString(fmt.Sprintf("<b>Начата:</b> 📅 %s\n", utils.FormatDateTimeWithDayOfWeek(w.StartedAt)))
	text.WriteString(fmt.Sprintf("<b>Статус:</b> %s\n", w.Status()))
	if w.Completed {
		text.WriteString(fmt.Sprintf("<b>Длительность:</b> %s\n", utils.BetweenTimes(w.StartedAt, w.EndedAt)))
	}
	text.WriteString("\n")

	if len(w.Exercises) > 0 {
		text.WriteString("<b>УПРАЖНЕНИЯ:</b>\n")
	}

	for i, exercise := range w.Exercises {
		sumWeight := float32(0)
		exerciseObj := exercise.ExerciseType
		text.WriteString(fmt.Sprintf("<b>%d. %s</b>\n", i+1, exerciseObj.Name))

		for _, set := range exercise.Sets {
			if set.Completed {
				sumWeight += set.GetRealWeight() * float32(set.GetRealReps())
			}
			text.WriteString(set.String(w.Completed))
		}
		if sumWeight > 0 {
			text.WriteString(fmt.Sprintf("<u>Общий вес</u>: %.0f кг\n", sumWeight))
		}
		text.WriteString("\n")
	}

	return text.String()
}
