build: 
	go build -o json-parser.exe cmd/json-parser/main.go

test: 
	go test ./tests/step1
	go test ./tests/step2
	go test ./tests/step3
	go test ./tests/step4
	go test ./tests/step5

run: 
	go run cmd/json-parser/main.go ${file}