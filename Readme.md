Run from the root

- `go run ./cmd/web`
- `go test -v ./cmd/web/`

Start the DB

- `mysqld`

Connect with

- `mysql -D snippetbox -u web -p`
- 'pass'

Use DB

- `USE snippetbox`

Stop SQL

- `ps -a | grep mysqld (show PID)`
- `kill -TERM PID`
