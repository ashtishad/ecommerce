run:
	export SERVER_ADDRESS=127.0.0.1 \
	export USER_API_PORT=8000 \
	export PRODUCT_API_PORT=8001 \
	export DB_USER=postgres \
	export DB_PASSWD=postgres \
	export DB_ADDR=127.0.0.1 \
	export DB_PORT=5432 \
	export DB_NAME=ecommerce \
	&& go run main.go
