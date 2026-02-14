package sets

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
	"time"
)

type CompleteByIDUseCase struct {
	setsRepo sets.Repo
}

func NewCompleteByIDUseCase(
	setsRepo sets.Repo,
) *CompleteByIDUseCase {
	return &CompleteByIDUseCase{
		setsRepo: setsRepo,
	}
}

func (uc *CompleteByIDUseCase) Name() string {
	return "Завершить/отменить подход"
}

func (uc *CompleteByIDUseCase) Execute(setID int64) error {
	set, err := uc.setsRepo.Get(setID)
	if err != nil {
		return err
	}
	set.Completed = !set.Completed
	now := time.Now()
	set.CompletedAt = &now

	err = uc.setsRepo.Save(set)
	if err != nil {
		return err
	}

	return nil
}
