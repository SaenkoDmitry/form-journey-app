package timers

import (
	"errors"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/timer"
	"time"
)

type StartUseCase struct {
	timerStore    *timer.Store
	exercisesRepo exercises.Repo
}

func NewStartUseCase(timerStore *timer.Store, exercisesRepo exercises.Repo) *StartUseCase {
	return &StartUseCase{
		timerStore:    timerStore,
		exercisesRepo: exercisesRepo,
	}
}

func (uc *StartUseCase) Name() string {
	return "Включение таймера"
}

var (
	TimerNotSupported = errors.New(messages.RestNotSupported)
)

func (uc *StartUseCase) Execute(chatID, exerciseID int64, seconds int) (*dto.StartTimer, error) {
	if seconds == 0 {
		return nil, TimerNotSupported
	}

	newTimerID := uc.timerStore.NewTimer(chatID)

	remainingCh := make(chan int)

	go func() {

		for remaining := seconds; remaining > 0; remaining-- {
			time.Sleep(1 * time.Second)
			if !uc.timerStore.HasTimer(chatID, newTimerID) {
				fmt.Println("stopped timer by user:", newTimerID)
				break
			}
			remainingCh <- remaining
		}
		close(remainingCh)
	}()

	exercise, _ := uc.exercisesRepo.Get(exerciseID)

	return &dto.StartTimer{
		RemainingCh: remainingCh,
		NewTimerID:  newTimerID,
		Exercise:    exercise,
		IsStopped: func(chatID int64, newTimerID string) bool {
			return !uc.timerStore.HasTimer(chatID, newTimerID)
		},
		StopTimer: func(chatID int64, newTimerID string) {
			fmt.Println("stopped timer by ending period:", newTimerID)
			uc.timerStore.StopTimer(chatID, newTimerID)
		},
	}, nil
}
