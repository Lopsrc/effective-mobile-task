package personstore

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"effectiveM-test-task/internal/models/person"
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

func (r *Repository) Save(ctx context.Context, person person.PersonSave) error {
	const op = "Store. SavePerson"
	r.log.Debug(op)

	query := "INSERT INTO person (name, surname, patronymic, del) VALUES ($1, $2, $3, $4)"
	_, err := r.client.Exec(ctx, query, person.Name, person.Surname, person.Patronymic, false)
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
	r.log.Debug(fmt.Sprintf("%s: saved", op))
	return nil
}

func (r *Repository) Get(ctx context.Context, p person.PersonGet) (person person.Person, err error) {
	const op = "Store. GetPerson"
	var isDel bool
	r.log.Debug(op)

	query := "SELECT name, surname, patronymic, del FROM person WHERE id = $1"
	if err := r.client.QueryRow(ctx, query, p.Id).Scan(&person.Name, &person.Surname, &person.Patronymic, &isDel); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.log.Error("%s: %v", op, newErr)
			return person, newErr
		}
		r.log.Error("%s: %v", op, err)
		return person, err
	}
	if isDel {
		r.log.Debug(fmt.Sprintf("%s: person is deleted", op))
		return person, errors.New("person is deleted")
	}
	r.log.Debug(fmt.Sprintf("%s: got", op))
	return
}

func (r *Repository) Update(ctx context.Context, person person.PersonUpdate) error {
	const op = "Store. UpdatePeson"
	query := "UPDATE person SET name = $1, surname = $2, patronymic = $3 WHERE id = $4"
	r.log.Debug(op)

	_, err := r.client.Exec(ctx, query, person.Name, person.Surname, person.Patronymic, person.Id)
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
	r.log.Debug(fmt.Sprintf("%s: updated", op))
	return nil
}

func (r *Repository) Delete(ctx context.Context, person person.PersonDelete) (bool, error) {
	const op = "Store. DeletePerson"
	query := "UPDATE person SET del = $1 WHERE id = $2"

	r.log.Debug(op)

	_, err := r.client.Exec(ctx, query, true, person.Id)
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

	r.log.Debug(fmt.Sprintf("%s: deleted", op))
	return true, nil
}

func (r *Repository) Recover(ctx context.Context, person person.PersonRecover) (bool, error) {
	const op = "Store. RecoverPerson"
	query := "UPDATE person SET del = $1 WHERE id = $2"

	r.log.Debug(op)

	_, err := r.client.Exec(ctx, query, false, person.Id)
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

	r.log.Debug(fmt.Sprintf("%s: recovered", op))
	return true, nil
}