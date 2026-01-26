package dto

type CompleteSet struct {
	WorkoutID       int64
	NeedStartTimer  bool
	NeedMoveToNext  bool
	NeedShowCurrent bool
	Seconds         int
}

type AddOneMoreSet struct {
	WorkoutID int64
}

type RemoveLastSet struct {
	WorkoutID int64
}

type NewSet struct {
	NewReps    int64
	NewWeight  float64
	NewMinutes int64
	NewMeters  int64
}

type SetResult struct {
}
