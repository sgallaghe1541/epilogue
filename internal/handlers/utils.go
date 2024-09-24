package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/sgallaghe1541/epilogue/internal/timeentry"
	"github.com/sgallaghe1541/epilogue/internal/utils"
)

func HandleWeekDay(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals := r.URL.Query()
			dateString := vals.Get("workdate")
			if dateString == "" {
				utils.ServerError(w, r, logger, fmt.Errorf("no date found in request url: %s", r.URL.Path))
				return
			}

			date, err := time.Parse("2006-01-02", dateString)
			if err != nil {
				utils.ServerError(w, r, logger, err)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			timeentry.DayofWeek(date).Render(context.Background(), w)
		})
}
