package sorting

import (
	"bufio"
	"io"
	"os"
)

func readFileLines(fileName string) ([]string, error) {
	if fileName == "-" {
		return readLines(os.Stdin)
	}

	f, err := os.Open(fileName)
	if err != nil {
		return nil, ErrFileNotFound{File: fileName}
	}
	defer f.Close()

	return readLines(f)
}

func readLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
