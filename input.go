package aoc

import (
	"io"

	"github.com/tajtiattila/aoc/input"
)

func MustInts(day int) []int {
	return input.MustInts(day)
}

func MustString(day int) string {
	return input.MustString(day)
}

func Reader(day int) io.Reader {
	return input.Reader(day)
}
