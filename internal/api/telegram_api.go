package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

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
		"https://189cfed595c8de.lhr.life", // tunnel for dev
	}

	for _, o := range allowed {
		if o == origin {
			return true
		}
	}

	return false
}
