package workouts

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
)

type FindByUserIDUseCase struct {
	workoutsRepo workouts.Repo
	usersRepo    users.Repo
}

func NewFindByUserUseCase(workoutsRepo workouts.Repo, usersRepo users.Repo) *FindByUserIDUseCase {
	return &FindByUserIDUseCase{workoutsRepo: workoutsRepo, usersRepo: usersRepo}
}

func (uc *FindByUserIDUseCase) Name() string {
	return "Найти тренировки пользователя"
}

var (
	EmptyWorkoutsErr = errors.New("empty workouts found for user")
)

func (uc *FindByUserIDUseCase) Execute(userID int64) (*dto.ShowWorkoutByUserID, error) {
	user, err := uc.usersRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	workoutObjs, err := uc.workoutsRepo.FindAll(user.ID)
	if err != nil {
		return nil, err
	}

	if len(workoutObjs) == 0 {
		return nil, EmptyWorkoutsErr
	}

	return &dto.ShowWorkoutByUserID{
		User:     user,
		Workouts: workoutObjs,
	}, nil
}
