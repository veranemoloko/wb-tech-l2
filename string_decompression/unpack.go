package main

import (
	"errors"
	"strings"
)

const (
	StateNone = iota // Initial state, no rune processed yet
	StateRune
	StateDigit
	StateEscape
)

const MaxRepeat = 1_000_000 // Maximum allowed repeat count to avoid excessive memory allocation

var (
	ErrInvalidStartDigit = errors.New("invalid string: starts with digit")
	ErrDigitAfterNumber  = errors.New("invalid string: digit after number")
	ErrInvalidEscape     = errors.New("invalid string: ends with escape character")
	ErrRepeatTooLarge    = errors.New("invalid string: repeat number too large")
)

// unpackString decompresses a string with repetition numbers and escape sequences.
func unpackString(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var b strings.Builder
	b.Grow(len(s))

	state := StateNone
	var lastRune rune
	var num int

	for _, r := range s {

		switch state {

		case StateNone:
			if r == '\\' {
				state = StateEscape
				continue
			}
			if isDigit(r) {
				return "", ErrInvalidStartDigit
			}
			lastRune = r
			state = StateRune

		case StateRune:
			if r == '\\' {
				b.WriteRune(lastRune)
				state = StateEscape
				continue
			}
			if isDigit(r) {
				num = int(r - '0')
				if num > MaxRepeat {
					return "", ErrRepeatTooLarge
				}
				state = StateDigit
				continue
			}
			b.WriteRune(lastRune)
			lastRune = r
			state = StateRune

		case StateDigit:
			if isDigit(r) { // Continue building multi-digit number
				num = num*10 + int(r-'0')
				if num > MaxRepeat {
					return "", ErrRepeatTooLarge
				}
				continue
			}

			// Finish the number: repeat lastRune 'num' times
			for i := 0; i < num; i++ {
				b.WriteRune(lastRune)
			}
			num = 0

			if r == '\\' {
				state = StateEscape
				continue
			}

			lastRune = r
			state = StateRune

		case StateEscape:
			lastRune = r
			state = StateRune
		}
	}

	// Handle the final state after the loop
	switch state {
	case StateRune:
		b.WriteRune(lastRune)
	case StateDigit:
		for i := 0; i < num; i++ {
			b.WriteRune(lastRune)
		}
	case StateEscape: // Cannot end on escape
		return "", ErrInvalidEscape
	}

	return b.String(), nil
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
