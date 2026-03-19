package dto

import "time"

type PeriodStats struct {
	AvgTime           time.Duration
	SumTime           time.Duration
	CompletedWorkouts int
	CardioTime        int
	IsWeek            bool
	IsMonth           bool
}

type ExercisesStats struct {
	Items []*ExerciseStat `json:"items"`
	Total int64           `json:"total"`
}

type ExerciseStat struct {
	ID   int64           `json:"id"`
	Date string          `json:"date"`
	Sets []*FormattedSet `json:"sets"`
}
