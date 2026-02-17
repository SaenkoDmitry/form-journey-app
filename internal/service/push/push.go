package push

import (
	"encoding/json"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SherClockHolmes/webpush-go"
	"gorm.io/gorm"
	"net/http"
	"os"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type Payload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
	Tag   string `json:"tag"`
}

func (p *Service) SendWorkoutFinished(userID, workoutID int64) error {
	var workout models.WorkoutDay

	if err := p.db.Preload("WorkoutDayType").First(&workout, workoutID).Error; err != nil {
		return err
	}

	var subs []models.PushSubscription
	if err := p.db.Where("user_id = ?", userID).Find(&subs).Error; err != nil {
		return err
	}

	payload := &Payload{
		Title: "Отдых закончен 💪",
		Body:  workout.WorkoutDayType.Name,
		URL:   fmt.Sprintf("/sessions/%d", workout.ID),
		Tag:   fmt.Sprintf("workout-%d", workout.ID),
	}

	payloadJSON, _ := json.Marshal(payload)

	for _, sub := range subs {
		status, err := sendPush(&sub, payloadJSON)
		if err != nil {
			if status == http.StatusGone || status == http.StatusNotFound {
				p.db.Delete(&sub)
			}
		}
	}

	return nil
}

func sendPush(sub *models.PushSubscription, payload []byte) (int, error) {
	subscription := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			Auth:   sub.Auth,
			P256dh: sub.P256dh,
		},
	}

	options := &webpush.Options{
		TTL:             3600,
		Urgency:         webpush.UrgencyHigh,
		Topic:           constants.Origin, // ← ВАЖНО
		VAPIDPrivateKey: os.Getenv("VAPID_PRIVATE_KEY"),
		VAPIDPublicKey:  os.Getenv("VAPID_PUBLIC_KEY"),
		Subscriber:      constants.Domain,
	}

	resp, err := webpush.SendNotification(payload, subscription, options)
	if err != nil {
		if resp != nil {
			return resp.StatusCode, err
		}
		return 0, err
	}

	defer resp.Body.Close()
	fmt.Println("Push status:", resp.StatusCode)
	return resp.StatusCode, nil
}
