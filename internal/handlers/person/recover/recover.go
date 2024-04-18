package recperson

import (
	"context"
	"effectiveM-test-task/internal/models/person"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

const (
	op = "Recover Person"
)

type Response struct {
	IsRec bool `json:"isRec"`
}

type PersonRecover interface {
	Recover(ctx context.Context, person person.PersonRecover) (bool,error)
}

func New(log *slog.Logger, personRecover PersonRecover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: deleting", op))
		ctx := r.Context()
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			log.Error(fmt.Sprintf("%s: invalid argument", op))
            http.Error(w, "invalid argument", http.StatusBadRequest)
            return
		}
		isRec, err := personRecover.Recover(ctx, person.PersonRecover{
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
