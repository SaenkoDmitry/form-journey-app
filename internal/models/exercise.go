package models

import "github.com/SaenkoDmitry/training-tg-bot/internal/constants"

type Exercise struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`

	WorkoutDayID int64
	WorkoutDay   *WorkoutDay `gorm:"foreignKey:WorkoutDayID;references:ID"` // join

	ExerciseTypeID int64
	ExerciseType   *ExerciseType `gorm:"foreignKey:ExerciseTypeID;references:ID"` // join

	Sets  []Set `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE"`
	Index int
}

func (*Exercise) TableName() string {
	return "exercises"
}

func (e *Exercise) CloneSets() []Set {
	sets := make([]Set, 0, len(e.Sets))
	for _, set := range e.Sets {
		newSet := Set{
			Reps:    set.GetRealReps(),
			Weight:  set.GetRealWeight(),
			Minutes: set.GetRealMinutes(),
			Meters:  set.GetRealMeters(),
			Index:   set.Index,
		}
		switch {
		case newSet.Reps == 0 && e.ExerciseType.ContainsReps():
			newSet.Reps = constants.DefaultReps
		case newSet.Weight == 0 && e.ExerciseType.ContainsWeight():
			newSet.Weight = constants.DefaultWeight
		case newSet.Minutes == 0 && e.ExerciseType.ContainsMinutes():
			newSet.Minutes = constants.DefaultMinutes
		case newSet.Meters == 0 && e.ExerciseType.ContainsMeters():
			newSet.Meters = constants.DefaultMeters
		}
		sets = append(sets, newSet)
	}
	return sets
}

func (e *Exercise) GetExerciseType() *ExerciseType {
	if e == nil {
		return nil
	}
	return e.ExerciseType
}

func (e *Exercise) Status() string {
	completedExerciseSets := e.CompletedSets()
	allSets := len(e.Sets)

	status := "🔴"
	if completedExerciseSets >= allSets {
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

func (e *Exercise) NextSet() Set {
	for _, set := range e.Sets {
		if !set.Completed {
			return set
		}
	}
	return Set{}
}

func (e *Exercise) LastSet() Set {
	if len(e.Sets) == 0 {
		return Set{}
	}
	return e.Sets[len(e.Sets)-1]
}
