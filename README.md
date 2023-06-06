# Ermetic

# Server

## Command Line

- cd server
- go run main.go [-r <rate_limit>, -s <window_size>]

## Default params:

- rate_limit - 5 (requests)
- window_size - 5 (seconds)

# Client

## Command Line

- cd client
- go run client/main.go [-u <server_url>, -c <clients_count>]

## Default params:

- clients_count - 1
- server_url - http://localhost:8081/
