package delcar

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
	op = "Delete Car"
)

type Response struct {
	IsDel bool `json:"isDel"`
}

type CarDeleter interface {
	Delete(ctx context.Context, car car.CarDelete) (bool, error)
}

func New(log *slog.Logger, carDeleter CarDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: deleting", op))
		ctx := r.Context()
		
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			log.Error(fmt.Sprintf("%s: invalid argument", op))
            http.Error(w, "invalid argument", http.StatusBadRequest)
            return			
		}
		isDel, err := carDeleter.Delete(ctx, car.CarDelete{
			Id: id,
		})
		if err != nil {
			log.Error("%s: deleting error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		if err := json.NewEncoder(w).Encode(Response{
			IsDel: isDel,
		}); err != nil {
			log.Error("%s: encoding error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%s: deleted", op))
	}
}
