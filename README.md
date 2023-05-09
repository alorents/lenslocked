# lenslocked




# Connecting to postgres

ensure postgres container is running
`docker compose up`
connect with `docker compose exec -it db psql -U <user> -d lenslocked`