package workouts

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisegrouptypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
)

type StatsUseCase struct {
	workoutsRepo           workouts.Repo
	dayTypesRepo           daytypes.Repo
	exerciseTypesRepo      exercisetypes.Repo
	exerciseGroupTypesRepo exercisegrouptypes.Repo
}

func NewStatsUseCase(
	workoutsRepo workouts.Repo,
	dayTypesRepo daytypes.Repo,
	exerciseTypesRepo exercisetypes.Repo,
	exerciseGroupTypesRepo exercisegrouptypes.Repo,
) *StatsUseCase {
	return &StatsUseCase{
		workoutsRepo:           workoutsRepo,
		dayTypesRepo:           dayTypesRepo,
		exerciseTypesRepo:      exerciseTypesRepo,
		exerciseGroupTypesRepo: exerciseGroupTypesRepo,
	}
}

func (uc *StatsUseCase) Name() string {
	return "Показать статистику тренировки"
}

func (uc *StatsUseCase) Execute(workoutID int64) (*dto.WorkoutStatistic, error) {
	workoutDay, err := uc.workoutsRepo.Get(workoutID)
	if err != nil {
		return nil, err
	}

	dayType, err := uc.dayTypesRepo.Get(workoutDay.WorkoutDayTypeID)
	if err != nil {
		return nil, err
	}

	totalWeight := 0.0
	completedExercises := 0
	cardioTime := 0

	exercisesMap := make(map[int64]*dto.FormattedExercise)
	exerciseWeightMap := make(map[int64]float64)
	exerciseTimeMap := make(map[int64]int)

	groups, err := uc.exerciseGroupTypesRepo.GetAll()
	if err != nil {
		return nil, err
	}

	groupsMap := make(map[string]string)
	for _, v := range groups {
		groupsMap[v.Code] = v.Name
	}

	for _, exercise := range workoutDay.Exercises {
		exercisesMap[exercise.ID] = dto.MapToFormattedExercise(exercise, groupsMap)

		completedExercises++
		exerciseTime := 0
		exerciseWeight := 0.0
		maxWeight := 0.0

		for _, set := range exercise.Sets {
			if !set.Completed {
				continue
			}
			exerciseWeight += float64(set.GetRealWeight()) * float64(set.GetRealReps())
			exerciseTime += set.GetRealMinutes()
			maxWeight = max(maxWeight, float64(set.GetRealWeight()))
		}
		exerciseWeightMap[exercise.ID] = exerciseWeight
		exerciseTimeMap[exercise.ID] = exerciseTime
		totalWeight += exerciseWeight
		cardioTime += exerciseTime
	}

	return &dto.WorkoutStatistic{
		WorkoutDay:         workoutDay,
		DayType:            dayType,
		ExerciseMap:        exercisesMap,
		ExerciseWeightMap:  exerciseWeightMap,
		ExerciseTimeMap:    exerciseTimeMap,
		TotalWeight:        totalWeight,
		CompletedExercises: completedExercises,
		CardioTime:         cardioTime,
	}, nil
}
