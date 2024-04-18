package getcar

import (
	"context"
	"effectiveM-test-task/internal/models/car"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	op = "Get Car"
)

type Response struct {
	Car car.Car
}

type CarGetter interface {
	Get(ctx context.Context, car car.CarGet) (car.Car, error)
}

func New(log *slog.Logger, carGetter CarGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: getting", op))
		ctx := r.Context()
		
		regNum := r.FormValue("regnum")
		if regNum == "" {
            log.Error(fmt.Sprintf("%s: invalid argument", op))
            http.Error(w, "invalid argument", http.StatusBadRequest)
            return
        }

		car, err := carGetter.Get(ctx, car.CarGet{
			RegNum: regNum,
		})
		if err != nil {
			log.Error("%s: getting error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		if err := json.NewEncoder(w).Encode(Response{
			Car: car,
		}); err != nil {
			log.Error("%s: encoding error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%s: got", op))
	}
}
