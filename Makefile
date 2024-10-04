build: 
	go build -o json-parser.exe cmd/json-parser/main.go

test: 
	go test ./tests/step1
	go test ./tests/step2
	go test ./tests/step3