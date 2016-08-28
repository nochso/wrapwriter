package wrapwriter

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
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
	},
	"Hard-wrap long words": {
		{"thisiswaytoolong", "this\niswa\nytoo\nlong", 4},
		{"it iswaytoolong", "it i\nsway\ntool\nong", 4},
		{"\nthisiswaytoolong\n", "\nthis\niswa\nytoo\nlong\n", 4},
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
			"Width: %d\nInput: %s\nDiff (-got +want):\n%s",
			test.width,
			pretty.Sprint(test.in),
			pretty.Compare(actual, test.out),
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
