run:
	export SERVER_ADDRESS=localhost \
	export SERVER_PORT=8000 \
	export DB_USER=postgres \
	export DB_PASSWD=postgres \
	export DB_ADDR=localhost \
	export DB_PORT=5432 \
	export DB_NAME=ecommerce \
	export POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/ecommerce?sslmode=disable&timezone=UTC'  \
	&& go run main.go
