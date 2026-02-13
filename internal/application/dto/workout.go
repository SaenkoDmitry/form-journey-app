package dto

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	"time"
)

type WorkoutItem struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	StartedAt string `json:"started_at"`
	Duration  string `json:"duration"`
	Completed bool   `json:"completed"`
	Status    string `json:"status"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type ShowMyWorkoutsResult struct {
	Items      []WorkoutItem `json:"items"`
	Pagination Pagination    `json:"pagination"`
}

type ConfirmDeleteWorkout struct {
	WorkoutID   int64
	DayTypeName string
}

type DeleteWorkout struct {
}

type ConfirmFinishWorkout struct {
	DayType models.WorkoutDayType
}

type FinishWorkout struct {
	WorkoutID int64
}

type CreateWorkout struct {
	WorkoutID int64
}

type StartWorkout struct {
}

type FormattedWorkout struct {
	ID          int64                `json:"id"`
	UserID      int64                `json:"user_id"`
	Status      string               `json:"status"`
	StartedAt   string               `json:"started_at"`
	Duration    string               `json:"duration"`
	EndedAt     string               `json:"ended_at"`
	DayTypeName string               `json:"day_type_name"`
	Completed   bool                 `json:"completed"`
	Exercises   []*FormattedExercise `json:"exercises"`
}

func MapToFormattedWorkout(w models.WorkoutDay, groupsMap map[string]string) *FormattedWorkout {
	res := &FormattedWorkout{
		ID:          w.ID,
		UserID:      w.UserID,
		StartedAt:   "📆️ " + utils.FormatDateTimeWithDayOfWeek(w.StartedAt),
		Status:      w.Status(),
		Duration:    utils.BetweenTimes(w.StartedAt, w.EndedAt),
		DayTypeName: w.WorkoutDayType.Name,
		Completed:   w.Completed,
	}
	for _, ex := range w.Exercises {
		res.Exercises = append(res.Exercises, MapToFormattedExercise(ex, groupsMap))
	}
	if w.EndedAt != nil {
		res.EndedAt = utils.FormatDate(*w.EndedAt)
	}
	return res
}

func MapToFormattedExercise(ex models.Exercise, groupsMap map[string]string) *FormattedExercise {
	sets := make([]*FormattedSet, 0, len(ex.Sets))
	sumWeight := float32(0)
	for _, s := range ex.Sets {
		if s.Completed {
			sumWeight += s.GetRealWeight() * float32(s.GetRealReps())
		}
		newSet := &FormattedSet{
			ID:              s.ID,
			FormattedString: s.String(ex.WorkoutDay.Completed),
			Completed:       s.Completed,
			Index:           s.Index,
		}
		if s.CompletedAt != nil {
			newSet.CompletedAt = s.CompletedAt.Add(3 * time.Hour).Format("15:04:05")
		}
		sets = append(sets, newSet)
	}
	return &FormattedExercise{
		ID:            ex.ID,
		Name:          ex.ExerciseType.Name,
		Units:         ex.ExerciseType.Units,
		GroupName:     groupsMap[ex.ExerciseType.ExerciseGroupTypeCode],
		RestInSeconds: ex.ExerciseType.RestInSeconds,
		Accent:        ex.ExerciseType.Accent,
		Description:   ex.ExerciseType.Description,
		SumWeight:     sumWeight,
		Index:         ex.Index,
		Sets:          sets,
	}
}

type FormattedExercise struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	Url           string          `json:"url"`
	GroupName     string          `json:"group_name"`
	RestInSeconds int             `json:"rest_in_seconds"`
	Accent        string          `json:"accent"`
	Units         string          `json:"units"`
	Description   string          `json:"description"`
	Index         int             `json:"index"`
	Sets          []*FormattedSet `json:"sets"`
	SumWeight     float32         `json:"sum_weight"`
}

type FormattedSet struct {
	ID              int64  `json:"id"`
	FormattedString string `json:"formatted_string"`
	Completed       bool   `json:"completed"`
	CompletedAt     string `json:"completed_at"`
	Index           int    `json:"index"`
}

type WorkoutProgress struct {
	Workout *FormattedWorkout `json:"workout"`

	TotalExercises     int
	CompletedExercises int

	TotalSets     int
	CompletedSets int

	ProgressPercent int
	RemainingMin    *int

	SessionStarted bool
}

type WorkoutStatistic struct {
	DayType            models.WorkoutDayType
	WorkoutDay         models.WorkoutDay
	TotalWeight        float64
	CompletedExercises int
	CardioTime         int
	ExerciseTypesMap   map[int64]models.ExerciseType
	ExerciseWeightMap  map[int64]float64
	ExerciseTimeMap    map[int64]int
}

type ShowWorkoutByUserID struct {
	Workouts []models.WorkoutDay
	User     *models.User
}
