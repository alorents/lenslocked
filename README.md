# lenslocked

# Running the application locally

ensure `.env` file is configured and upto date. See `.env.template` in the base directory

`go run main.go` from the base directory

# Connecting to postgres

ensure postgres container is running
`docker compose up`
connect with `docker compose exec -it db psql -U <user> -d lenslocked`

# Running migrations

set env variable
`export GOOSE_POSTGRES_CFG="host=localhost port=5432 user=<user> password=<password> dbname=<db_name> sslmode=disable"`
<br>
ensure postgres container is running
`docker compose up`
<br>
run migrations with goose
<br>
update `goose postgres $GOOSE_POSTGRES_CFG up`
<br>
downgrade run migrations with goose `goose postgres $GOOSE_POSTGRES_CFG down`
<br>
status run migrations with goose `goose postgres $GOOSE_POSTGRES_CFG status`
<br>
