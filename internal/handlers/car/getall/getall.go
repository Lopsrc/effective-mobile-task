package getallcars

import (
	"context"
	"effectiveM-test-task/internal/models/car"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	op = "Get All Cars"
)

type Response struct {
	Cars []car.Car
}

type CarAllGetter interface {
	GetAll(ctx context.Context) ([]car.Car, error)
}

func New(log *slog.Logger, carAllGetter CarAllGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: getting", op))
		ctx := r.Context()

		cars, err := carAllGetter.GetAll(ctx)
		if err != nil {
			log.Error("%s: getting error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		if err := json.NewEncoder(w).Encode(Response{
			Cars: cars,
		}); err != nil {
			log.Error("%s: encoding error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%s: got", op))
	}
}


