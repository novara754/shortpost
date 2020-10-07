# Shortpost

A minimal twitter clone with user accounts, sessions, and a timeline.

![Screenshot of the website with some sample posts](./screenshot.png)

## Building & Running

To compile the program simply run `go build ./shortpost`, this will
generate an executable with the same name.

In order for the executable to work you need to create a file named `config.json`
with the same struct as the `config.json.sample`. Fill in the URL
for your instance of postgres and you're good to go!

The server will run on port `8080`.

## License

Licensed under the [MIT License](./LICENSE).
