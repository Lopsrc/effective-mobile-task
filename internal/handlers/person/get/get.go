package getperson

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
	op = "Get Person"
)

type Response struct {
	Person person.Person
}

type PersonGetter interface {
	Get(ctx context.Context, p person.PersonGet) (person person.Person, err error)
}

func New(log *slog.Logger, personGetter PersonGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: getting", op))
		ctx := r.Context()

		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			log.Error(fmt.Sprintf("%s: invalid argument", op))
            http.Error(w, "invalid argument", http.StatusBadRequest)
            return
		}
		p, err := personGetter.Get(ctx, person.PersonGet{
			Id: id,
		})
		if err != nil {
			log.Error("%s: getting error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		if err := json.NewEncoder(w).Encode(Response{
			Person: p,
		}); err != nil {
			log.Error("%s: encoding error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%s: got", op))
	}
}
