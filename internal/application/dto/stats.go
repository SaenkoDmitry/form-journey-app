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
