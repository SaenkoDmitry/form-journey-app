package timers

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/timer"
)

type StopUseCase struct {
	timerStore *timer.Store
}

func NewStopUseCase(timerStore *timer.Store) *StopUseCase {
	return &StopUseCase{
		timerStore: timerStore,
	}
}

func (uc *StopUseCase) Name() string {
	return "Остановка таймера"
}

func (uc *StopUseCase) Execute(chatID int64, timerID string) {
	uc.timerStore.StopTimer(chatID, timerID)
}
