package reccar

import (
	"context"
	"effectiveM-test-task/internal/models/car"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

const (
	op = "Recover Car"
)

type Response struct {
	IsRec bool `json:"isRec"`
}

type CarRecover interface {
	Recover(ctx context.Context, car car.CarRecover) (bool,error)
}

func New(log *slog.Logger, carRecover CarRecover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: deleting", op))
		ctx := r.Context()
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			log.Error(fmt.Sprintf("%s: invalid argument", op))
            http.Error(w, "invalid argument", http.StatusBadRequest)
            return
		}
		isRec, err := carRecover.Recover(ctx, car.CarRecover{
			Id: id,
		})
		if err != nil {
			log.Error("%s: deleting error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(Response{
			IsRec: isRec,
		}); err != nil {
			log.Error("%s: encoding error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%s: deleted", op))
	}
}


