build:
	GOOS=darwin  GOARCH=amd64 go build -o ./bin/mac/csv_parser      csv_parser.go
	GOOS=linux   GOARCH=amd64 go build -o ./bin/linux/csv_parser    csv_parser.go
	GOOS=windows GOARCH=amd64 go build -o ./bin/windows/csv_parser  csv_parser.go

	GOOS=darwin  GOARCH=amd64 go build -o ./bin/mac/web_watcher     web_watcher.go
	GOOS=linux   GOARCH=amd64 go build -o ./bin/linux/web_watcher   web_watcher.go
	GOOS=windows GOARCH=amd64 go build -o ./bin/windows/web_watcher web_watcher.go
