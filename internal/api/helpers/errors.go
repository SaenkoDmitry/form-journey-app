package helpers

import (
	"errors"
	"net/http"

	"github.com/SaenkoDmitry/training-tg-bot/internal/api/errorslist"
)

func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	switch {
	case errors.Is(err, errorslist.ErrAccessDenied):
		http.Error(w, err.Error(), http.StatusForbidden)
	case errors.Is(err, errorslist.ErrInternalMsg):
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
