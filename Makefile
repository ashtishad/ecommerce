run:
	export SERVER_ADDRESS=127.0.0.1 \
	export SERVER_PORT=8000 \
	export DB_USER=postgres \
	export DB_PASSWD=postgres \
	export DB_ADDR=127.0.0.1 \
	export DB_PORT=5432 \
	export DB_NAME=ecommerce \
	export POSTGRESQL_URL='postgres://postgres:postgres@127.0.0.1:5432/ecommerce?sslmode=disable&timezone=UTC'  \
	&& go run main.go
