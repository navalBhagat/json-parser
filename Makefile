build: 
	go build -o json-parser.exe cmd/json-parser/main.go

test: 
	go test ./tests/step1