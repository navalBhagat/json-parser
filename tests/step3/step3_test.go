package step1

import (
	"json-parser/pkg/parser"
	"json-parser/pkg/tokenizer"
	"log"
	"os"
	"testing"
)

func TestValidJsonFromFile(t *testing.T) {
	filename := "../../testdata/tests/step3/valid.json"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}
	defer file.Close()
	parsed := parser.Parse(tokenizer.NewTokenizerFromReader(file))

	if !parsed {
		t.Errorf("Expected successful parsing but could not parse")
	}

}

func TestInvalidJsonFromFile(t *testing.T) {
	filename := "../../testdata/tests/step3/invalid.json"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}
	defer file.Close()
	parsed := parser.Parse(tokenizer.NewTokenizerFromReader(file))

	if parsed {
		t.Errorf("Expected unsuccessful parsing but was able to parse.")
	}
}
