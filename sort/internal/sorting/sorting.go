package sorting

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func SortLines(opts SortOptions) error {
	var lines []string
	for _, file := range opts.Files {
		fileLines, err := readFileLines(file)
		if err != nil {
			return err
		}
		lines = append(lines, fileLines...)
	}

	if len(lines) == 0 {
		return nil
	}

	if opts.Check {
		if err := checkSorted(lines, opts); err != nil {
			return err
		}
		return nil
	}

	sort.SliceStable(lines, func(i, j int) bool {
		keyI := getKeyColumn(lines[i], opts)
		keyJ := getKeyColumn(lines[j], opts)

		hasKeyI := keyI != ""
		hasKeyJ := keyJ != ""

		var less bool
		if hasKeyI && hasKeyJ {
			less = compareKeys(lines[i], lines[j], keyI, keyJ, opts)
		} else {
			less = lines[i] < lines[j]
		}

		if opts.Reverse {
			return !less
		}
		return less
	})

	return writeLines(lines, opts)
}

func getKeyColumn(line string, opts SortOptions) string {
	if opts.Key <= 0 {
		return line
	}

	var fields []string
	if opts.Separator != 0 {
		fields = strings.Split(line, string(opts.Separator))
	} else {
		fields = strings.Fields(line)
	}

	if opts.Key <= len(fields) {
		return fields[opts.Key-1]
	}

	// GNU sort: если поля нет — использовать всю строку
	return line
}

func compareKeys(a, b, keyA, keyB string, opts SortOptions) bool {
	keyATrim, keyBTrim := keyA, keyB
	leadingA, leadingB := 0, 0

	if opts.IgnoreBlanks {
		keyATrim = strings.TrimLeft(keyA, " \t")
		keyBTrim = strings.TrimLeft(keyB, " \t")
		leadingA = len(keyA) - len(keyATrim)
		leadingB = len(keyB) - len(keyBTrim)
	}
	if keyA == "" || keyB == "" {
		return a < b
	}

	if opts.Month {
		ma, mb := monthOrder(keyATrim), monthOrder(keyBTrim)

		if ma == 0 && mb == 0 {
		} else if ma == 0 {
			return true
		} else if mb == 0 {
			return false
		} else if ma != mb {
			return ma < mb
		}
	}

	if opts.Human {
		okA, numA := parseHuman(keyATrim)
		okB, numB := parseHuman(keyBTrim)

		if okA && okB {
			if numA != numB {
				return numA < numB
			}
		} else if okA && !okB {
			return false
		} else if !okA && okB {
			return true
		}
	}

	if opts.NumericSort {
		numA, errA := extractNumber(keyATrim)
		numB, errB := extractNumber(keyBTrim)

		if errA == nil && errB == nil {
			if numA != numB {
				return numA < numB
			}
		} else if errA == nil {
			return false
		} else if errB == nil {
			return true
		}
	}

	if keyATrim == keyBTrim && (opts.IgnoreBlanks || opts.Key > 0) && leadingA != leadingB {
		return leadingA > leadingB
	}

	if keyATrim != keyBTrim {
		return keyATrim < keyBTrim
	}

	return false
}

func checkSorted(lines []string, opts SortOptions) error {
	for i := 1; i < len(lines); i++ {
		prevKey := getKeyColumn(lines[i-1], opts)
		curKey := getKeyColumn(lines[i], opts)

		less := compareKeys(lines[i-1], lines[i], prevKey, curKey, opts)
		if opts.Reverse {
			less = !less
		}
		if !less {
			return ErrNotSorted{Line: i + 1}
		}
	}
	return nil
}

func extractNumber(s string) (float64, error) {
	if s == "" {
		return 0, strconv.ErrSyntax
	}
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	s = s[i:]
	if s == "" {
		return 0, strconv.ErrSyntax
	}

	var numPart strings.Builder
	for j, r := range s {
		if unicode.IsDigit(r) || r == '.' {
			numPart.WriteRune(r)
		} else {
			if j == 0 {
				break
			}
			break
		}
	}
	if numPart.Len() == 0 {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseFloat(numPart.String(), 64)
}

func parseHuman(s string) (bool, float64) {
	if s == "" {
		return false, 0
	}

	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	s2 := s[i:]
	if s2 == "" {
		return false, 0
	}

	var num float64
	var unit string

	for j, r := range s2 {
		if (r < '0' || r > '9') && r != '.' {
			numPart := s2[:j]
			if j < len(s2) {
				unit = strings.ToUpper(string(s2[j]))
			}
			n, err := strconv.ParseFloat(numPart, 64)
			if err != nil {
				return false, 0
			}
			num = n
			break
		}
	}

	if unit == "" {
		n, err := strconv.ParseFloat(s2, 64)
		if err != nil {
			return false, 0
		}
		return true, n
	}

	switch unit {
	case "K":
		num *= 1 << 10
	case "M":
		num *= 1 << 20
	case "G":
		num *= 1 << 30
	case "T":
		num *= 1 << 40
	case "P":
		num *= 1 << 50
	case "E":
		num *= 1 << 60
	}

	return true, num
}

func writeLines(lines []string, opts SortOptions) error {
	writer := bufio.NewWriterSize(os.Stdout, 4<<20)
	defer writer.Flush()

	var prev string
	write := func(s string) error {
		trimmed := strings.TrimRight(s, "\r\n")
		if opts.Unique && trimmed == prev {
			return nil
		}
		prev = trimmed
		if _, err := writer.WriteString(trimmed); err != nil {
			return err
		}
		if err := writer.WriteByte('\n'); err != nil {
			return err
		}
		return nil
	}

	for _, line := range lines {
		if err := write(line); err != nil {
			return err
		}
	}

	return nil
}

var monthTable = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4,
	"may": 5, "jun": 6, "jul": 7, "aug": 8,
	"sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

func monthOrder(s string) int {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	if i+3 > len(s) {
		return 0
	}

	token := strings.ToLower(s[i : i+3])
	if num, ok := monthTable[token]; ok {
		return num
	}
	return 0
}
