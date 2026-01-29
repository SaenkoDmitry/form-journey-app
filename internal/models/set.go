package models

import (
	"fmt"
	"strings"
	"time"
)

type Set struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	ExerciseID  int64
	Reps        int
	FactReps    int
	Weight      float32
	FactWeight  float32
	Minutes     int
	FactMinutes int
	Meters      int
	FactMeters  int
	Completed   bool
	CompletedAt *time.Time
	Index       int
}

func (*Set) TableName() string {
	return "sets"
}

func (s *Set) String(done bool) string {
	var text strings.Builder

	if s.Meters > 0 {
		text.WriteString(fmt.Sprintf("• %s метров: ", s.FormatMeters()))
	}
	if s.Minutes > 0 {
		text.WriteString(fmt.Sprintf("• %s минут: ", s.FormatMinutes()))
	}
	if s.Reps > 0 {
		text.WriteString(fmt.Sprintf("• %s повторений по %s кг: ", s.FormatReps(), s.FormatWeight()))
	}

	if s.Completed {
		text.WriteString(fmt.Sprintf("✅, %s", s.CompletedAt.Add(3*time.Hour).Format("15:04:05")))
	} else {
		if done {
			text.WriteString("💔")
		} else {
			text.WriteString("🚀")
		}
	}
	text.WriteString("\n")
	return text.String()
}

func (s *Set) FormatReps() string {
	if s.FactReps != 0 {
		return fmt.Sprintf("<strike>%d</strike> <b>%d</b>", s.Reps, s.FactReps)
	}
	return fmt.Sprintf("%d", s.Reps)
}

func (s *Set) FormatWeight() string {
	if s.FactWeight != float32(0) {
		return fmt.Sprintf("<strike>%.0f</strike> <b>%.0f</b>", s.Weight, s.FactWeight)
	}
	return fmt.Sprintf("%.0f", s.Weight)
}

func (s *Set) FormatMinutes() string {
	if s.FactMinutes != 0 {
		return fmt.Sprintf("<strike>%d</strike> <b>%d</b>", s.Minutes, s.FactMinutes)
	}
	return fmt.Sprintf("%d", s.Minutes)
}

func (s *Set) FormatMeters() string {
	if s.FactMeters != 0 {
		return fmt.Sprintf("<strike>%d</strike> <b>%d</b>", s.Meters, s.FactMeters)
	}
	return fmt.Sprintf("%d", s.Meters)
}

func (s *Set) GetRealReps() int {
	if s == nil {
		return 0
	}
	if s.FactReps > 0 {
		return s.FactReps
	}
	return s.Reps
}

func (s *Set) GetRealWeight() float32 {
	if s == nil {
		return 0
	}
	if s.FactWeight > 0 {
		return s.FactWeight
	}
	return s.Weight
}

func (s *Set) GetRealMinutes() int {
	if s == nil {
		return 0
	}
	if s.FactMinutes > 0 {
		return s.FactMinutes
	}
	return s.Minutes
}

func (s *Set) GetRealMeters() int {
	if s == nil {
		return 0
	}
	if s.FactMeters > 0 {
		return s.FactMinutes
	}
	return s.Meters
}
