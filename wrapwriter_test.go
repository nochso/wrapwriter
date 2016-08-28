package wrapwriter

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

type testInput struct {
	in, out string
	width   int
}

var tests = map[string][]testInput{
	"Soft-wrap on word boundaries": {
		{"1 23 4", "1\n23\n4", 2},
		{"wrap me!", "wra\np\nme!", 3},
		{"wrap me!", "wrap\nme!", 4},
		{"wrap me!", "wrap\nme!", 5},
		{"wrap me!", "wrap\nme!", 7},
		{"wrap me!", "wrap me!", 8},
		{"ääää ääää", "ääää\nääää", 4},
	},
	"Hard-wrap long words": {
		{"thisiswaytoolong", "this\niswa\nytoo\nlong", 4},
		{"it iswaytoolong", "it i\nsway\ntool\nong", 4},
		{"\nthisiswaytoolong\n", "\nthis\niswa\nytoo\nlong\n", 4},
		{"ääää", "ää\nää", 2},
	},
	"Keep linefeeds": {
		{"\n\n", "\n\n", 4},
		{"keep\nme", "keep\nme", 4},
		{"keep\nme", "keep\nme", 5},
		{"keep\nme\n", "keep\nme\n", 5},
		{"\nkeep\nme\n", "\nkeep\nme\n", 5},
	},
	"No extra whitespace": {
		{"", "", 2},
	},
	"Collapse lone whitespace": {
		{" ", "", 1},
		{" ", "", 5},
		{"  ", "", 1},
		{"  ", "", 5},
		{" xx", "x\nx", 1},
		{" xx", "xx", 2},
	},
}

func TestWrap(t *testing.T) {
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			for _, test := range data {
				testWrap(test, t)
			}
		})
	}
}

func testWrap(test testInput, t *testing.T) {
	actual, err := Wrap(test.in, test.width)
	if err != nil {
		t.Fatal(err)
	}
	if actual != test.out {
		t.Errorf(
			"Width: %d\nDiff (-got +want):\n%s",
			test.width,
			diff(test.out, actual),
		)
	}
}

func TestWrap_widthMustBePositive(t *testing.T) {
	_, err := Wrap("", 0)
	if err == nil {
		t.Error("expecting error when using width=0, got nil")
	}
	_, err = Wrap("", -1)
	if err == nil {
		t.Error("expecting error when using width=1, got nil")
	}
}

var update = flag.Bool("update", false, "update golden files")

func TestWrap_golden(t *testing.T) {
	var inputFiles []string
	filepath.Walk("test-fixtures/golden", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() && !strings.HasSuffix(path, ".golden") {
			inputFiles = append(inputFiles, path)
		}
		return nil
	})
	for _, inputFile := range inputFiles {
		input, err := ioutil.ReadFile(inputFile)
		if err != nil {
			t.Fatal(err)
		}
		inputStr := string(input)
		firstLinePos := bytes.IndexByte(input, '\n')
		if firstLinePos == -1 {
			t.Fatalf("input file contains no newline: %s", inputFile)
		}
		width, err := strconv.Atoi(string(input[0:firstLinePos]))
		if err != nil {
			t.Fatalf("parsing width: %s: %s", inputFile, err)
		}
		actualStr, err := Wrap(inputStr, width)
		if err != nil {
			t.Error(err)
		}
		actual := []byte(actualStr)
		golden := inputFile + ".golden"
		if *update {
			ioutil.WriteFile(golden, actual, 0644)
		}
		expected, _ := ioutil.ReadFile(golden)
		if !bytes.Equal(actual, expected) {
			t.Errorf(
				"Width: %d\nDiff (-got +want):\n%s",
				width,
				diff(string(expected), actualStr),
			)
		}
	}
}

func diff(expected, actual string) string {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(expected),
		B:        difflib.SplitLines(actual),
		FromFile: "Expected",
		ToFile:   "Actual",
		Context:  1,
	}
	out, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		log.Fatal(err)
	}
	return out
}
