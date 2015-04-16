package numerizer

import "testing"

func TestParse(t *testing.T) {
	var tests = []struct {
		in   string
		want int64
	}{
		{"zero", 0},
		{"one", 1},
		{"three", 3},
		{"eleven", 11},
		{"seventeen", 17},
		{"sixty", 60},
		{"one hundred", 100},
		{"two hundred", 200},
		{"three thousand", 3000},
		{"four million", 4000000},
		{"five billion", 5000000000},
		{"six trillion", 6000000000000},
		{"two hundred four", 204},
		{"two hundred and four", 204},
		{"seventeen hundred", 1700},
		{"three thousand four", 3004},
		{"three thousand and four", 3004},
		{"three thousand sixteen", 3016},
		{"three thousand and sixteen", 3016},
		{"three thousand thirty", 3030},
		{"three thousand and thirty", 3030},
		{"three thousand thirty three", 3033},
		{"three thousand four hundred", 3400},
		{"three thousand and four hundred", 3400},
		{"three thousand five hundred one", 3501},
		{"three thousand five hundred and one", 3501},
		{"three thousand six hundred twelve", 3612},
		{"three thousand six hundred and twelve", 3612},
		{"three thousand six hundred eighty", 3680},
		{"three thousand six hundred and eighty", 3680},
		{"three thousand six hundred eighty four", 3684},
		{"three thousand six hundred and eighty four", 3684},
		{"ten thousand", 10000},
		{"twenty thousand", 20000},
		{"three hundred thousand", 300000},
		{"forty five", 45},
		{"forty five hundred", 4500},
		{"forty five thousand", 45000},
		{"nine hundred ninety nine thousand nine hundred ninety nine", 999999},
		{"nine hundred ninety nine million nine hundred ninety nine thousand nine hundred ninety nine", 999999999},
		{"forty-five", 45},
		{"four thousand, four hundred", 4400},
		{"four thousand, four hundred thirty-two", 4432},
	}

	for _, tt := range tests {
		have, err := Parse(tt.in)
		if err != nil {
			t.Fatalf("Parse(%q) %v", tt.in, err)
		} else if have != tt.want {
			t.Errorf("Parse(%q)\nhave %d\nwant %d", tt.in, have, tt.want)
		}
	}
}

func TestParseError(t *testing.T) {
	var tests = []string{
		"",
		"a",
		"and",
		"hundred",
		"thousand",
		"zero four",
		"four zero",
		"five three",
		"four eight",
		"eight seventy",
		"eleven six",
		"twelve seventeen",
		"twelve thirty",
		"forty eighteen",
		"forty twenty",
		"five hundred hundred",
		"six thousand hundred",
		"seven thousand thousand",
		"forty five thousand and",
		"two and three",
	}

	for _, tt := range tests {
		have, err := Parse(tt)
		if err == nil {
			t.Errorf("Parse(%q)\nhave %d\nwant parse error", tt, have)
		}
	}
}
