package saveperson

import (
	"context"
	"effectiveM-test-task/internal/models/person"
	"encoding/json"
	// "errors"
	"fmt"
	// "io"
	"log/slog"
	"net/http"
	// "strings"
)

const (
	op = "Save Person"
)

type Request struct {
	// Person person.PersonSave
	Name string `json:"name"`
    Surname string `json:"surname"`
    Patronymic string `json:"patronymic"`
}

type Response struct {
	// Person person.PersonSave
	
}

type PersonSaver interface {
	Save(ctx context.Context, person person.PersonSave) error
}

func New(log *slog.Logger, personSaver PersonSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: saving", op))
		ctx := r.Context()
		name 		:= r.FormValue("name")
		surname 	:= r.FormValue("surname")
		patronymic 	:= r.FormValue("patronymic")
		log.Debug(fmt.Sprintf("Request params = name:'%s', surname:'%s', patronymic:'%s'", name, surname, patronymic))
		if name == "" || surname == "" || patronymic == ""{
			log.Error(fmt.Sprintf("%s: invalid request", op))
            http.Error(w, "invalid request", http.StatusBadRequest)
            return
        }
		
		if err := personSaver.Save(ctx, person.PersonSave{
			Name: name,
            Surname: surname,
            Patronymic: surname,
		}); err != nil {
			log.Error("%s: saving error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		if err := json.NewEncoder(w).Encode(Response{}); err != nil {
			log.Error("%s: encoding error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%s: saved", op))
	}
}
