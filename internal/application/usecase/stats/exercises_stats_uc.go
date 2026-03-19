package stats

import (
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type GetExercisesStatsUseCase struct {
	usersRepo     users.Repo
	exercisesRepo exercises.Repo
}

func NewGetExercisesStatsUseCase(usersRepo users.Repo, exercisesRepo exercises.Repo) *GetExercisesStatsUseCase {
	return &GetExercisesStatsUseCase{
		usersRepo:     usersRepo,
		exercisesRepo: exercisesRepo,
	}
}

func (uc *GetExercisesStatsUseCase) Name() string {
	return "Статистика пользователя по упражнению"
}

func (uc *GetExercisesStatsUseCase) Execute(userID, exerciseTypeID int64, offset, limit int) (*dto.ExercisesStats, error) {
	exerciseObjs, err := uc.exercisesRepo.FindAllByUserIDAndExTypeID(userID, exerciseTypeID, offset, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.ExerciseStat, 0)
	for _, ex := range exerciseObjs {
		stat := &dto.ExerciseStat{
			ID:   ex.ID,
			Date: ex.WorkoutDay.StartedAt.Add(3 * time.Hour).Format("02.01.2006 15:04"),
			Sets: make([]*dto.FormattedSet, 0),
		}
		for _, s := range ex.Sets {
			stat.Sets = append(stat.Sets, dto.MapToFormattedSet(s, ex))
		}
		result = append(result, stat)
	}

	total, err := uc.exercisesRepo.CountByUserIDAndExTypeID(userID, exerciseTypeID)
	if err != nil {
		return nil, err
	}

	return &dto.ExercisesStats{
		Items: result,
		Total: total,
	}, nil
}
