package step1

import (
	"json-parser/pkg/parser"
	"json-parser/pkg/tokenizer"
	"log"
	"os"
	"testing"
)

func TestValidJsonFromFile(t *testing.T) {
	filename := "../../testdata/tests/step2/valid.json"
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

func TestValidJson2FromFile(t *testing.T) {
	filename := "../../testdata/tests/step2/valid2.json"
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
	filename := "../../testdata/tests/step2/invalid.json"
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

func TestInvalidJson2FromFile(t *testing.T) {
	filename := "../../testdata/tests/step2/invalid2.json"
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
