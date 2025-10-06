package sorting

import (
	flag "github.com/spf13/pflag"
)

// SortOptions holds all command-line options that control the sort behavior.
type SortOptions struct {
	Key          int      // -k N: sort by the Nth field (1-based index)
	Separator    rune     // -t C: field separator character (default is tab '\t')
	NumericSort  bool     // -n: compare according to numerical value
	Reverse      bool     // -r: reverse the result of comparisons
	Unique       bool     // -u: output only the first of lines with equal keys
	Month        bool     // -M: compare months (JAN < FEB < ... < DEC)
	Human        bool     // -h: compare human-readable numbers (e.g., 2K, 1G)
	IgnoreBlanks bool     // -b: ignore trailing blanks
	Check        bool     // -c: check whether the input is sorted; do not sort
	Files        []string // input files; if empty, stdin ("-") is used
	BufferMb     int
}

// ParseFlags parses command-line flags using pflag (supports GNU-style combined short flags).
func ParseFlags() SortOptions {
	key := flag.IntP("key", "k", 0, "Sort by the Nth field (1-based index)")
	sep := flag.StringP("separator", "t", "\t", "Field separator character (default is tab '\\t')")
	numeric := flag.BoolP("numeric", "n", false, "Compare according to numerical value")
	reverse := flag.BoolP("reverse", "r", false, "Reverse the result of comparisons")
	unique := flag.BoolP("unique", "u", false, "Output only the first of lines with equal keys")
	month := flag.BoolP("month", "M", false, "Compare months (JAN < FEB < ... < DEC)")
	human := flag.BoolP("human", "h", false, "Compare human-readable numbers (e.g., 2K, 1G)")
	ignore := flag.BoolP("ignore-blanks", "b", false, "Ignore trailing blanks")
	check := flag.BoolP("check", "c", false, "Check whether the input is sorted; do not sort")

	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		files = []string{"-"} // use stdin if no files provided
	}

	var separator rune
	runes := []rune(*sep)
	if len(runes) > 0 {
		separator = runes[0]
	}

	return SortOptions{
		Key:          *key,
		NumericSort:  *numeric,
		Reverse:      *reverse,
		Unique:       *unique,
		Month:        *month,
		Human:        *human,
		IgnoreBlanks: *ignore,
		Check:        *check,
		Separator:    separator,
		Files:        files,
	}
}
