package sorting

import "fmt"

type ErrFileNotFound struct {
	File string
}

func (e ErrFileNotFound) Error() string {
	return fmt.Sprintf("file not found: %s", e.File)
}

type ErrNotSorted struct {
	Line int
}

func (e ErrNotSorted) Error() string {
	return fmt.Sprintf("lines not sorted at line %d", e.Line)
}

type ErrInvalidColumn struct {
	Key int
}

func (e ErrInvalidColumn) Error() string {
	return fmt.Sprintf("invalid column key: %d", e.Key)
}
