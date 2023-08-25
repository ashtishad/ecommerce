run:
	export SERVER_ADDRESS=localhost \
	export SERVER_PORT=8000 \
	export DB_USER=postgres \
	export DB_PASSWD=potgres \
	export DB_ADDR=localhost \
	export DB_PORT=5432 \
	export DB_NAME=users \
	&& go run main.go
