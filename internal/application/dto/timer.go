package dto

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

type StartTimer struct {
	RemainingCh chan int // channel for getting remaining seconds
	NewTimerID  string   // created timer ID
	Exercise    models.Exercise
	IsStopped   func(chatID int64, newTimerID string) bool
	StopTimer   func(chatID int64, newTimerID string)
}
