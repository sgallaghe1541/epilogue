package app

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func Make(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				writeJSON(w, apiErr.StatusCode, apiErr)
			} else {
				errResp := map[string]any{
					"statusCode": http.StatusInternalServerError,
					"msg":        "internal server error",
				}
				writeJSON(w, http.StatusInternalServerError, errResp)
			}
			slog.Error("HTTP API error", "err", err.Error(), "path", r.URL.Path)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func htmlDatetoQueryDate(date string) string {
	datetime, _ := time.Parse("2006-01-02", date)
	return datetime.Format("01/02/2006")
}

func getRecentWEDate() string {
	now := time.Now()

	now = now.Add(time.Hour * -24 * 5)

	for now.Weekday() != 6 {
		now = now.Add(time.Hour * -24)
		fmt.Println(now)
	}

	return now.Format("2006-01-02")
}

func getAbsURL(u url.URL) string {
	b := u.JoinPath(u.Scheme, u.Host)
	return fmt.Sprint(b, "/")
}
