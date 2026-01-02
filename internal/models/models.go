package models

import (
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID        int64 `gorm:"primaryKey"`
	Username  string
	ChatID    int64
	CreatedAt time.Time
}

type WorkoutDay struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64
	Name      string
	Exercises []Exercise `gorm:"foreignKey:WorkoutDayID"`
	StartedAt time.Time
	EndedAt   *time.Time
	Completed bool
}

func (w *WorkoutDay) Status() string {
	if !w.Completed {
		return fmt.Sprintf("⏳ *%s* Активна", w.Name)
	}
	if w.EndedAt != nil {
		return fmt.Sprintf("✅ *%s* Завершена в %s", w.Name, w.EndedAt.Add(3*time.Hour).Format("15:04"))
	}

	return fmt.Sprintf("✅ *%s* Завершена", w.Name)
}

func (w *WorkoutDay) String() string {
	var text strings.Builder

	text.WriteString(fmt.Sprintf("Статус: %s\n", w.Status()))
	text.WriteString(fmt.Sprintf("Дата: %s\n\n", w.StartedAt.Add(3*time.Hour).Format("02.01.2006")))
	text.WriteString("*Упражнения:*\n")

	for i, exercise := range w.Exercises {
		text.WriteString(fmt.Sprintf("%s %d. %s: %d/%d подходов\n\n", exercise.Status(), i+1, exercise.Name, exercise.CompletedSets(), len(exercise.Sets)))
	}

	return text.String()
}

type Exercise struct {
	ID            int64 `gorm:"primaryKey"`
	WorkoutDayID  int64
	Name          string
	Sets          []Set `gorm:"foreignKey:ExerciseID"`
	Hint          string
	RestInSeconds int
}

func (e *Exercise) Status() string {
	completedExerciseSets := e.CompletedSets()
	allSets := len(e.Sets)

	status := "🔴"
	if int(completedExerciseSets) >= allSets {
		status = "🟢"
	} else if completedExerciseSets > 0 {
		status = "🟡"
	}
	return status
}

func (e *Exercise) CompletedSets() int {
	completedSets := 0
	for _, set := range e.Sets {
		if set.Completed {
			completedSets++
		}
	}
	return completedSets
}

func (e *Exercise) Next() Set {
	for _, set := range e.Sets {
		if !set.Completed {
			return set
		}
	}
	return Set{}
}

type Set struct {
	ID          int64 `gorm:"primaryKey"`
	ExerciseID  int64
	Reps        int
	Weight      float32
	Completed   bool
	CompletedAt *time.Time
}

type WorkoutSession struct {
	ID                   int64 `gorm:"primaryKey"`
	WorkoutDayID         int64
	CurrentExerciseIndex int
	StartedAt            time.Time
	IsActive             bool
}
