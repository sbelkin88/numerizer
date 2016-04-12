package numerizer

import "testing"

func TestParse(t *testing.T) {
	var tests = []struct {
		in   string
		want int
	}{
		{"zero dollars and thirty four cents", 34},
		{"zero dollars", 000},
		{"one dollar", 100},
		{"three dollars", 300},
		{"eleven dollars", 1100},
		{"seventeen dollars", 1700},
		{"sixty dollars", 6000},
		{"one hundred dollars", 10000},
		{"two hundred dollars", 20000},
		{"three thousand dollars", 300000},
		{"four million dollars", 400000000},
		{"five billion dollars", 500000000000},
		{"two hundred four dollars", 20400},
		{"two hundred four dollars and eighteen cents", 20418},
		{"two hundred and four dollars", 20400},
		{"seventeen hundred dollars", 170000},
		{"three thousand four dollars", 300400},
		{"three thousand and four dollars", 300400},
		{"three thousand sixteen dollars", 301600},
		{"three thousand and sixteen dollars", 301600},
		{"three thousand thirty dollars", 303000},
		{"three thousand and thirty dollars", 303000},
		{"three thousand thirty three dollars", 303300},
		{"three thousand thirty three dollars and twenty cents", 303320},
		{"three thousand four hundred dollars", 340000},
		{"three thousand and four hundred dollars", 340000},
		{"three thousand five hundred one dollars", 350100},
		{"three thousand five hundred and one dollars", 350100},
		{"three thousand six hundred twelve dollars", 361200},
		{"three thousand six hundred twelve dollars and sixty one cents", 361261},
		{"three thousand six hundred and twelve dollars", 361200},
		{"three thousand six hundred eighty dollars", 368000},
		{"three thousand six hundred and eighty dollars", 368000},
		{"three thousand six hundred eighty four dollars", 368400},
		{"three thousand six hundred and eighty four dollars", 368400},
		{"ten thousand dollars", 1000000},
		{"ten thousand dollars and two cents", 1000002},
		{"twenty thousand dollars", 2000000},
		{"three hundred thousand dollars", 30000000},
		{"forty five dollars", 4500},
		{"forty five hundred dollars", 450000},
		{"forty five thousand dollars", 4500000},
		{"forty five thousand dollars and one cent", 4500001},
		{"nine hundred ninety nine thousand nine hundred ninety nine dollars", 99999900},
		{"nine hundred ninety nine million nine hundred ninety nine thousand nine hundred ninety nine dollars", 99999999900},
		{"nine hundred ninety nine million nine hundred ninety nine thousand nine hundred ninety nine dollars and ninety nine cents", 99999999999},
		{"forty-five dollars", 4500},
		{"four thousand, four hundred dollars", 440000},
		{"four thousand, four hundred thirty-two dollars", 443200},
	}

	for i, tt := range tests {
		have, err := Parse(tt.in)
		if err != nil {
			t.Fatalf("Parse(%q) %v.  Test %d", tt.in, err, i)
		} else if have != tt.want {
			t.Errorf("Parse(%q)\nhave %d\nwant %d", tt.in, have, tt.want)
		}
	}
}
