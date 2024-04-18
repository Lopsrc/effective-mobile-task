package carstore

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"effectiveM-test-task/internal/models/car"
	// "effectiveM-test-task/internal/storage"
	postgresqlcli "effectiveM-test-task/pkg/client/postgresql"

	"github.com/jackc/pgconn"
)

type Repository struct {
	client postgresqlcli.Client
	log    *slog.Logger
}

func NewRepository(client postgresqlcli.Client, log *slog.Logger) *Repository {
	return &Repository{
		client: client,
		log:    log,
	}
}

func (r *Repository) Save(ctx context.Context, car car.CarSave) error {
	const op = "Store. CarSave"
	query := "INSERT INTO car (reg_number, owner_id, del) VALUES ($1, $2, $3)"

	r.log.Debug(op)
	for i := 0; i < len(car.RegNums); i++ {
		_, err := r.client.Exec(ctx, query, car.RegNums[i], 1, false)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				pgErr = err.(*pgconn.PgError)
				newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
				r.log.Error("%s: %v", op, newErr)
				return newErr
			}
			r.log.Error("%s: %v", op, err)
			return err
		}
	}
	
	r.log.Debug(fmt.Sprintf("%s: saved car", op))
	return nil
}

func (r *Repository) Get(ctx context.Context, c car.CarGet) (car car.Car, err error) {
	const op = "Store. CarGet"
	var isDel bool
	query := "SELECT mark, model, year, del FROM car WHERE reg_number = $1"

	r.log.Debug(op)

	if err := r.client.QueryRow(ctx, query, c.RegNum).Scan( &car.Mark, &car.Model, &car.Year, &isDel); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.log.Error("%s: %v", op, newErr)
			return car, newErr
		}
		r.log.Error("%s: %v", op, err)
		return car, err
	}
	if isDel{
		r.log.Debug(fmt.Sprintf("%s: deleted car", op))
        return car, errors.New("car not found")
	}

	r.log.Debug(fmt.Sprintf("%s: got car", op))

	return car, nil
}

func (r *Repository) Update(ctx context.Context, car car.CarUpdate) error {
	const op = "Store. CarUpdate"

	query := "UPDATE car SET reg_number = $1, mark = $2, model = $3, year = $4, owner_id = $5 WHERE id = $6"

	r.log.Debug(op)

	_, err := r.client.Exec(ctx, query, car.RegNum, car.Mark, car.Model, car.Year, car.OwnerID ,car.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.log.Error("%s: %v", op, newErr)
			return newErr
		}
		r.log.Error("%s: %v", op, err)
		return err
	}

	r.log.Debug(fmt.Sprintf("%s: updated car", op))

	return nil
}

func (r *Repository) Delete(ctx context.Context, car car.CarDelete) (bool, error) {
	const op = "Store. CarDelete"
	query := "UPDATE car SET del = $1 WHERE id = $2"

	r.log.Debug(op)

	_, err := r.client.Exec(ctx, query, true, car.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.log.Error("%s: %v", op, newErr)
			return false, newErr
		}
		r.log.Error("%s: %v", op, err)
		return false, err
	}
	
	r.log.Debug(fmt.Sprintf("%s: deleted car", op))

	return true, nil
}

func (r *Repository) Recover(ctx context.Context, car car.CarRecover) (bool, error) {
	const op = "Store. CarRecover"
	query := "UPDATE car SET del = $1 WHERE id = $2"
	r.log.Debug(op)

	_, err := r.client.Exec(ctx, query, false, car.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.log.Error("%s: %v", op, newErr)
			return false, newErr
		}
		r.log.Error("%s: %v", op, err)
		return false, err
	}

    r.log.Debug(fmt.Sprintf("%s: recovered car", op))
	return true, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]car.Car, error) { 
	const op = "Store. GetAll"
	query := "SELECT reg_number, mark, model, year FROM car WHERE del = $1"

	r.log.Debug(op)
	
	rows, err := r.client.Query(ctx, query, false)
	if err!= nil {
		var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            pgErr = err.(*pgconn.PgError)
            newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
            r.log.Error("%s: %v", op, newErr)
            return nil, newErr
        }
        r.log.Error("%s: %v", op, err)
        return nil, err
    }
	defer rows.Close()
	var cars []car.Car
	for rows.Next() {
		var car car.Car
        if err := rows.Scan(&car.RegNum, &car.Mark, &car.Model, &car.Year); err!= nil {
            r.log.Error("%s: %v", op, err)
            return nil, err
        }
        cars = append(cars, car)
    }

	r.log.Debug(fmt.Sprintf("%s: got all cars", op))
	return cars, nil
}
