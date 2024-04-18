# effective-mobile-task

## About it
This is test task for Effective Mobile.

### Technologies

>Swagger: https://redis.io/docs/install/install-redis/
>
>Migrations golang-migrate/migrate: https://github.com/golang-migrate/migrate
>
>PostgreSQL: https://www.postgresql.org
>
>Router go-chi/chi: https://github.com/go-chi/chi

## Using

### Clone the repositiry
```
git clone https://github.com/Lopsrc/effective-mobile-task
```

### Preparation

Edit the local.env file. Specify your bot's token:
```
vim config/api.env
vim config/local.yaml
```
Install migrator https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

### Running

Run:
```
# migrate up.
make migrate-up
# run.
make run
```

Migrate down:
```
# migrate down
make migrate-down
```