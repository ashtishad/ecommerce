### Auth Challenge

##### Run `./start.sh` to download the dependencies and run the application.
* Make the Script Executable: You must give the script execute permissions before you can run it. Use the following command:
  `chmod +x start.sh`


To run the application, you have to define the environment variables, default values of the variables are defined inside `start.sh`

- SERVER_ADDRESS    `[IP Address of the machine]` : `localhost`
- SERVER_PORT       `[Port of the machine]` : `8000`
- DB_USER           `[Database username]` : `root`
- DB_PASSWD         `[Database password]`: `root`
- DB_ADDR           `[IP address of the database]` : `localhost`
- DB_PORT           `[Port of the database]` : `3306`
- DB_NAME           `[Name of the database]` : `user_auth`

#### MySQL Database
Make the changes to your `start.sh` file for modifying default db configurations.
`docker-compose.yml` file. This contains the database migrations scripts. You just need to bring the container up.

   To start the docker container, run the `docker-compose up`.
