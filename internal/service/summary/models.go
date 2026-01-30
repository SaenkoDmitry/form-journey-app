package summary

type ExerciseSummary struct {
	ExerciseType string
	Workouts     map[string]struct{}
	Sets         int
	MaxWeight    float32
	AvgWeight    float64
	TotalWeight  float64
	TotalReps    int
	TotalMinutes int
}

type DateSummary struct {
	Workouts    int
	Exercises   map[string]struct{}
	Sets        int
	TotalVolume float32
	MaxWeight   float32
}

type WeekSummary struct {
	SumWeight  float32
	SumMinutes int
	SumMeters  int
}

type Progress struct {
	MaxWeight float32
	MaxReps   int
	AvgWeight float32

	MinMinutes int
	MaxMinutes int
	SumMinutes int

	MinMeters int
	MaxMeters int
	SumMeters int
}
