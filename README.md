# Shortpost

A minimal twitter clone with user accounts, sessions, and a timeline.

![Screenshot of the website with some sample posts](./screenshot.png)

## Building & Running

To compile the program simply run `go build -o shortpost.out ./shortpost`, this will
generate an executable with the same name.

In order for the executable to work you need to launch it with proper values
in the `PORT` and `POSTGRES_URL` environment variables. For example:

```
$ PORT=8080 \
  POSTGRES_URL="postgres://user:password@localhost:5432/shortpost:?sslmode=disable \
  ./shortpost.out 
```

The server will then listen on the configured port.

### With Docker

You can also easily launch the server along with a Postgres instance using the provided
`docker-compose.yml`. Simply run `docker-compose up -d` in the project root directory
to run the service as a daemon.

By default the webserver can be accessed on port `8080`, but this can be configured in the 
`docker-compose.yml` file.

## License

Licensed under the [MIT License](./LICENSE).
