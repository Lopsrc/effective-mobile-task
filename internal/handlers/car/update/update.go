package upcar

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
	op = "Update Car"
)


type Response struct {
}

type CarUpdater interface {
	Update(ctx context.Context, car car.CarUpdate) error
}

func New(log *slog.Logger, carUpdater CarUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: updating", op))
		ctx := r.Context()
		
		
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
            log.Error(fmt.Sprintf("%s: invalid argument", op))
            http.Error(w, "invalid argument", http.StatusBadRequest)
            return            
        }
		ownerID, err := strconv.Atoi(r.FormValue("ownerid"))
		if err != nil {
            log.Error(fmt.Sprintf("%s: invalid argument", op))
            http.Error(w, "invalid argument", http.StatusBadRequest)
            return            
        }
		regNum := r.FormValue("regnum")
		mark := r.FormValue("mark")
		model := r.FormValue("model")
		year := r.FormValue("year")
		if regNum == "" || mark == "" || model == "" || year == "" || id == 0 {
            log.Error(fmt.Sprintf("%s: invalid argument", op))
            http.Error(w, "invalid argument", http.StatusBadRequest)
            return            
        }
		if err := carUpdater.Update(ctx, car.CarUpdate{
			Id: id,
            RegNum: regNum,
            Mark: mark,
            Model: model,
            Year: year,
			OwnerID: ownerID,
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
