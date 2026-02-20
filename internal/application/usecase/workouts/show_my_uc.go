package workouts

import (
	"errors"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
)

type FindMyUseCase struct {
	workoutsRepo workouts.Repo
	usersRepo    users.Repo
}

func NewFindMyUseCase(repo workouts.Repo, usersRepo users.Repo) *FindMyUseCase {
	return &FindMyUseCase{workoutsRepo: repo, usersRepo: usersRepo}
}

var (
	NotFoundAllErr = errors.New("no workouts for user")
)

func (uc *FindMyUseCase) Name() string {
	return "Показать мои тренировки"
}

func (uc *FindMyUseCase) ExecuteByChatID(chatID int64, offset, limit int) (*dto.ShowMyWorkoutsResult, error) {
	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}
	return uc.Execute(user.ID, offset, limit)
}

func (uc *FindMyUseCase) Execute(userID int64, offset, limit int) (*dto.ShowMyWorkoutsResult, error) {
	total, err := uc.workoutsRepo.Count(userID)
	if err != nil {
		return nil, err
	}

	workoutObjs, err := uc.workoutsRepo.Find(userID, offset, limit)
	if err != nil {
		return nil, err
	}

	items := make([]dto.WorkoutItem, 0, len(workoutObjs))
	for _, w := range workoutObjs {
		duration := ""
		if w.Completed && w.EndedAt != nil {
			duration = utils.BetweenTimes(w.StartedAt, w.EndedAt)
		}

		status := "🟡"
		if w.Completed {
			status = "✅"
			if w.EndedAt != nil {
				status += fmt.Sprintf(" ~ %s", duration)
			}
		}

		items = append(items, dto.WorkoutItem{
			ID:        w.ID,
			Name:      w.WorkoutDayType.Name,
			StartedAt: "📆️ " + utils.FormatDateTimeWithDayOfWeek(w.StartedAt),
			Duration:  duration,
			Completed: w.Completed,
			Status:    status,
		})
	}

	return &dto.ShowMyWorkoutsResult{
		Items: items,
		Pagination: dto.Pagination{
			Offset: offset,
			Limit:  limit,
			Total:  int(total),
		},
	}, nil
}
