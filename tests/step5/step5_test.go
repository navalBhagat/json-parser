package step1

import (
	"json-parser/pkg/parser"
	"json-parser/pkg/tokenizer"
	"log"
	"os"
	"testing"
)

func TestValidJsonFromFiles(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		filename string
	}{
		{"ValidJson1", "../../testdata/tests/step5/pass1.json"},
		{"ValidJson2", "../../testdata/tests/step5/pass2.json"},
		{"ValidJson3", "../../testdata/tests/step5/pass3.json"},
	}

	// Loop through each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.Open(tc.filename)
			if err != nil {
				t.Fatalf("Unable to read file: %v", err)
			}
			defer file.Close()

			parsed := parser.Parse(tokenizer.NewTokenizerFromReader(file))
			if !parsed {
				t.Errorf("Expected successful parsing but could not parse file: %s", tc.filename)
			}
		})
	}
}

func TestInvalidJsonFromFile(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
	}{
		{"InvalidJson1", "../../testdata/tests/step5/fail1.json"},
		{"InvalidJson2", "../../testdata/tests/step5/fail2.json"},
		{"InvalidJson3", "../../testdata/tests/step5/fail3.json"},
		{"InvalidJson4", "../../testdata/tests/step5/fail4.json"},
		{"InvalidJson5", "../../testdata/tests/step5/fail5.json"},
		{"InvalidJson6", "../../testdata/tests/step5/fail6.json"},
		{"InvalidJson7", "../../testdata/tests/step5/fail7.json"},
		{"InvalidJson8", "../../testdata/tests/step5/fail8.json"},
		{"InvalidJson9", "../../testdata/tests/step5/fail9.json"},
		{"InvalidJson10", "../../testdata/tests/step5/fail10.json"},
		{"InvalidJson11", "../../testdata/tests/step5/fail11.json"},
		{"InvalidJson12", "../../testdata/tests/step5/fail12.json"},
		{"InvalidJson13", "../../testdata/tests/step5/fail13.json"},
		{"InvalidJson14", "../../testdata/tests/step5/fail14.json"},
		{"InvalidJson15", "../../testdata/tests/step5/fail15.json"},
		{"InvalidJson16", "../../testdata/tests/step5/fail16.json"},
		{"InvalidJson17", "../../testdata/tests/step5/fail17.json"},
		{"InvalidJson18", "../../testdata/tests/step5/fail18.json"},
		{"InvalidJson19", "../../testdata/tests/step5/fail19.json"},
		{"InvalidJson20", "../../testdata/tests/step5/fail20.json"},
		{"InvalidJson21", "../../testdata/tests/step5/fail21.json"},
		{"InvalidJson22", "../../testdata/tests/step5/fail22.json"},
		{"InvalidJson23", "../../testdata/tests/step5/fail23.json"},
		{"InvalidJson24", "../../testdata/tests/step5/fail24.json"},
		{"InvalidJson25", "../../testdata/tests/step5/fail25.json"},
		{"InvalidJson26", "../../testdata/tests/step5/fail26.json"},
		{"InvalidJson27", "../../testdata/tests/step5/fail27.json"},
		{"InvalidJson28", "../../testdata/tests/step5/fail28.json"},
		{"InvalidJson29", "../../testdata/tests/step5/fail29.json"},
		{"InvalidJson30", "../../testdata/tests/step5/fail30.json"},
		{"InvalidJson31", "../../testdata/tests/step5/fail31.json"},
		{"InvalidJson32", "../../testdata/tests/step5/fail32.json"},
		{"InvalidJson33", "../../testdata/tests/step5/fail33.json"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.Open(tc.filename)
			if err != nil {
				log.Fatalf("Unable to read file: %v", err)
			}
			defer file.Close()
			parsed := parser.Parse(tokenizer.NewTokenizerFromReader(file))

			if parsed {
				t.Errorf("Expected unsuccessful parsing but was able to parse.")
			}
		})
	}
}
