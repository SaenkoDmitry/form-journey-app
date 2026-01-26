package stats

import (
	"strings"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
)

type GetPeriodStatsUseCase struct {
	usersRepo    users.Repo
	workoutsRepo workouts.Repo
}

func NewGetPeriodStatsUseCase(usersRepo users.Repo, workoutsRepo workouts.Repo) *GetPeriodStatsUseCase {
	return &GetPeriodStatsUseCase{
		usersRepo:    usersRepo,
		workoutsRepo: workoutsRepo,
	}
}

func (uc *GetPeriodStatsUseCase) Name() string {
	return "Статистика за период"
}

func (uc *GetPeriodStatsUseCase) Execute(chatID int64, period string) (*dto.PeriodStats, error) {
	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}

	workoutObjs, err := uc.workoutsRepo.FindAll(user.ID)
	if err != nil {
		return nil, err
	}

	completedWorkouts := 0
	sumTime := time.Duration(0)
	cardioTime := 0

	for _, w := range workoutObjs {
		if !w.Completed {
			continue
		}
		switch period {
		case week:
			if time.Since(w.StartedAt).Hours() > 7*24 {
				continue
			}
		case month:
			if time.Since(w.StartedAt).Hours() > 30*24 {
				continue
			}
		default:
		}

		completedWorkouts++
		sumTime += w.EndedAt.Sub(*&w.StartedAt)
		for _, e := range w.Exercises {
			if len(e.Sets) == 0 {
				continue
			}
			for _, s := range e.Sets {
				if !s.Completed {
					continue
				}
				if s.GetRealMinutes() > 0 {
					cardioTime += s.GetRealMinutes()
				}
			}
		}
	}

	avgTime := time.Duration(0)
	if completedWorkouts != 0 {
		avgTime = sumTime / time.Duration(completedWorkouts)
	}

	return &dto.PeriodStats{
		CompletedWorkouts: completedWorkouts,
		CardioTime:        cardioTime,
		SumTime:           sumTime,
		AvgTime:           avgTime,
		IsWeek:            strings.EqualFold(period, week),
		IsMonth:           strings.EqualFold(period, month),
	}, nil
}

const (
	week  = "week"
	month = "month"
)
