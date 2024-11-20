# lenslocked
This follows a webdev course by Jon Calhoun found at https://courses.calhoun.io/courses/cor_wdv2

# Running the application locally

ensure `.env` file is configured and upto date. See `.env.template` in the base directory

`go run cmd/server/server.go` from the base directory

alternatively run `modd` from the base directory for a live-reload server on file system changes

# Connecting to postgres

ensure postgres container is running \
`docker compose up` \
connect with `docker compose exec -it db psql -U <user> -d lenslocked`

# Running migrations

set env variable \
`export GOOSE_POSTGRES_CFG="host=localhost port=5432 user=<user> password=<password> dbname=<db_name> sslmode=disable"`
<br>
ensure postgres container is running \
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

# Tailwind
Installation \
`npm install -D tailwindcss
npx tailwindcss init`

Build \
From the directory containing tailwind.config.js \
input and output css files are personal preference \
`npx tailwindcss -i ./styles.css -o ../assets/styles.css --watch`

# Production
`docker-compose -f docker-compose.yml -f docker-compose.production.yml up --build`