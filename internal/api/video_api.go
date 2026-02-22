package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
)

func (s *serviceImpl) LinkVideo(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "missing url", 400)
		return
	}

	token, _ := generateVideoToken(url, claims.UserID)

	json.NewEncoder(w).Encode(map[string]string{
		"url": "/api/video/stream?token=" + token,
	})
}

func (s *serviceImpl) StreamVideo(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", 401)
		return
	}

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		http.Error(w, "invalid token", 401)
		return
	}

	dataToken, _ := base64.URLEncoding.DecodeString(parts[0])
	signature, _ := base64.URLEncoding.DecodeString(parts[1])

	mac := hmac.New(sha256.New, secret)
	mac.Write(dataToken)
	expected := mac.Sum(nil)

	if !hmac.Equal(signature, expected) {
		http.Error(w, "invalid signature", 401)
		return
	}

	var payload VideoToken
	json.Unmarshal(dataToken, &payload)

	if time.Now().Unix() > payload.Expires {
		http.Error(w, "expired", 401)
		return
	}
	publicURL := payload.URL
	if publicURL == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}

	// ---------------------------
	// 1️⃣ Получаем download href
	// ---------------------------
	apiURL := "https://cloud-api.yandex.net/v1/disk/public/resources/download?public_key=" +
		url.QueryEscape(publicURL)

	apiReq, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	apiResp, err := client.Do(apiReq)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer apiResp.Body.Close()

	if apiResp.StatusCode != http.StatusOK {
		http.Error(w, "failed to get download link", apiResp.StatusCode)
		return
	}

	var data struct {
		Href string `json:"href"`
	}

	if err := json.NewDecoder(apiResp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if data.Href == "" {
		http.Error(w, "empty download link", 500)
		return
	}

	// ---------------------------
	// 2️⃣ Запрашиваем само видео
	// ---------------------------
	videoReq, err := http.NewRequest("GET", data.Href, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// ВАЖНО: поддержка перемотки
	if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
		videoReq.Header.Set("Range", rangeHeader)
	}

	videoReq.Header.Set("User-Agent", "Mozilla/5.0")
	videoReq.Header.Set("Referer", "https://disk.yandex.ru/")

	videoResp, err := client.Do(videoReq)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer videoResp.Body.Close()

	// ---------------------------
	// 3️⃣ Копируем заголовки
	// ---------------------------

	// Content-Type
	contentType := videoResp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "video/mp4"
	}
	w.Header().Set("Content-Type", contentType)

	// Range поддержка
	if videoResp.Header.Get("Content-Range") != "" {
		w.Header().Set("Content-Range", videoResp.Header.Get("Content-Range"))
	}

	if videoResp.Header.Get("Content-Length") != "" {
		w.Header().Set("Content-Length", videoResp.Header.Get("Content-Length"))
	}

	w.Header().Set("Accept-Ranges", "bytes")

	// Очень важно для iOS
	w.Header().Set("Cache-Control", "no-cache")

	w.WriteHeader(videoResp.StatusCode)

	// ---------------------------
	// 4️⃣ Стримим тело
	// ---------------------------
	io.Copy(w, videoResp.Body)
}

var secret = []byte("SUPER_SECRET_KEY")

type VideoToken struct {
	URL     string
	Expires int64
	UserID  int64
}

func generateVideoToken(url string, userID int64) (string, error) {
	payload := VideoToken{
		URL:     url,
		Expires: time.Now().Add(5 * time.Minute).Unix(),
		UserID:  userID,
	}

	data, _ := json.Marshal(payload)

	mac := hmac.New(sha256.New, secret)
	mac.Write(data)
	signature := mac.Sum(nil)

	token := base64.URLEncoding.EncodeToString(data) + "." +
		base64.URLEncoding.EncodeToString(signature)

	return token, nil
}
