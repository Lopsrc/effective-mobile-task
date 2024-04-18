package upperson

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
	op = "Update Person"
)

type Response struct {
}

type PersonUpdater interface {
	Update(ctx context.Context, person person.PersonUpdate) error
}

func New(log *slog.Logger, carUpdater PersonUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: updating", op))
		ctx := r.Context()
		
		id, err := strconv.Atoi(r.FormValue("id"))
		if err!= nil {
            log.Error(fmt.Sprintf("%s: invalid request", op))
            http.Error(w, "invalid request", http.StatusBadRequest)
            return
        }
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		patronymic := r.FormValue("patronymic")
		if id == 0 || name == "" || surname == "" || patronymic == "" {
			log.Error(fmt.Sprintf("%s: invalid argument", op))
            http.Error(w, "invalid argument", http.StatusBadRequest)
            return
		}

		if err := carUpdater.Update(ctx, person.PersonUpdate{
			Id: id,
            Name: name,
            Surname: surname,
            Patronymic: patronymic,
		}); err != nil {
			log.Error("%s: updating error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(Response{}); err != nil {
			log.Error("%s: encoding error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%s: updated", op))
	}
}
