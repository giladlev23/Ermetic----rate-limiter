# Ermetic

## Further implementations left out of the exercise's scope

- Tests
- Better OOP - base structs for future extending (i.e. base Client struct for different requests formats)
- Some fallback and retry mechanism (i.e. retries on the client side, buffering requests on server side)
- Other notes left at the code

# Server

## Command Line

- cd server
- go run cmd/main.go [-r <rate_limit>, -s <window_size>]

## Default params:

- rate_limit - 5 (requests)
- window_size - 5 (seconds)

# Client

## Command Line

- cd client
- go run cmd/main.go [-u <server_url>, -c <clients_count>, -w <wait_interval_range_milliseconds>]

## Default params:

- clients_count - 1
- server_url - http://localhost:8081/
- wait_interval_range_milliseconds - 1000
