package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-jwt/jwt/v4"
)

var (
	botToken  = os.Getenv("TELEGRAM_TOKEN")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

type TelegramUser struct {
	ID           int64  `json:"id"` // chat_id
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	PhotoURL     string `json:"photo_url"`
	AuthDate     int64  `json:"auth_date"`
	Hash         string `json:"hash"`
}

func (s *serviceImpl) TelegramLoginHandler(w http.ResponseWriter, r *http.Request) {

	var tgUser TelegramUser

	if err := json.NewDecoder(r.Body).Decode(&tgUser); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if !verifyTelegram(tgUser, botToken) {
		http.Error(w, "invalid telegram hash", http.StatusUnauthorized)
		return
	}

	chatID := tgUser.ID

	user, err := s.container.GetUserUC.Execute(chatID)
	if err != nil && errors.Is(err, users.NotFoundUserErr) {
		user, err = s.container.CreateUserUC.Execute(tgUser.ID, &tgbotapi.User{
			ID:        tgUser.ID,
			FirstName: tgUser.FirstName,
			LastName:  tgUser.LastName,
			UserName:  tgUser.Username,
		})
	}
	if user == nil || err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"chat_id":    tgUser.ID,
		"first_name": tgUser.FirstName,
		"last_name":  tgUser.LastName,
		"photo_url":  tgUser.PhotoURL,
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(jwtSecret)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": signed,
	})
}

func verifyTelegram(user TelegramUser, botToken string) bool {
	data := map[string]string{
		"id":         strconv.FormatInt(user.ID, 10),
		"first_name": user.FirstName,
		"auth_date":  strconv.FormatInt(user.AuthDate, 10),
	}
	if user.Username != "" {
		data["username"] = user.Username
	}
	if user.PhotoURL != "" {
		data["photo_url"] = user.PhotoURL
	}

	// сортируем ключи
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// формируем check string
	var parts []string
	for _, k := range keys {
		parts = append(parts, k+"="+data[k])
	}
	checkString := strings.Join(parts, "\n")

	secret := sha256.Sum256([]byte(botToken))

	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(checkString))
	expected := hex.EncodeToString(mac.Sum(nil))

	return expected == user.Hash
}

func (s *serviceImpl) TelegramRedirectHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	if origin == "" {
		http.Error(w, "missing origin", http.StatusBadRequest)
		return
	}

	// 🔐 ОБЯЗАТЕЛЬНО — whitelist origin
	if !s.isAllowedOrigin(origin) {
		http.Error(w, "invalid origin", http.StatusForbidden)
		return
	}

	botID := os.Getenv("TELEGRAM_BOT_ID")
	if botID == "" {
		http.Error(w, "bot id not configured", http.StatusInternalServerError)
		return
	}

	returnTo := origin + "/auth-telegram"

	telegramURL := fmt.Sprintf(
		"https://oauth.telegram.org/auth?bot_id=%s&origin=%s&return_to=%s",
		url.QueryEscape(botID),
		url.QueryEscape(origin),
		url.QueryEscape(returnTo),
	)

	http.Redirect(w, r, telegramURL, http.StatusFound)
}

func (s *serviceImpl) isAllowedOrigin(origin string) bool {
	if strings.HasSuffix(origin, ".lhr.life") {
		return true
	}
	allowed := []string{
		"http://localhost:3000",
		"https://form-journey.ru",
	}

	for _, o := range allowed {
		if o == origin {
			return true
		}
	}

	return false
}
