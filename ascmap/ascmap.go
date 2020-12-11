package ascmap

import (
	"bufio"
	"errors"
	"io"
)

type Ascmap struct {
	Dx, Dy int
	P      []byte

	Padx, Pady int

	Origin int
	Stride int
}

var ErrDim = errors.New("inconsistent dimensions")

func FromReaderPad(r io.Reader, padx, pady int) (*Ascmap, error) {
	m := &Ascmap{
		Padx: padx,
		Pady: pady,
	}

	pad := func(n int) {
		m.P = append(m.P, make([]byte, n)...)
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Bytes()
		if m.Dx == 0 {
			m.Dx = len(t)
			if m.Dx == 0 {
				return nil, ErrDim
			}

			m.Stride = m.Dx + 2*padx
			m.Origin = pady*m.Stride + padx

			pad(pady * m.Stride)
		} else {
			if len(t) != m.Dx {
				return nil, ErrDim
			}
		}

		pad(padx)
		m.P = append(m.P, t...)
		pad(padx)

		m.Dy++
	}

	pad(pady * m.Stride)

	return m, scanner.Err()
}

func FromReader(r io.Reader) (*Ascmap, error) {
	return FromReaderPad(r, 0, 0)
}

func (m *Ascmap) In(x, y int) bool {
	return x >= 0 && x < m.Dx && y >= 0 && y < m.Dy
}

func (m *Ascmap) Offset(x, y int) int {
	return m.Origin + x + y*m.Stride
}

func (m *Ascmap) Delta(dx, dy int) int {
	return dx + dy*m.Stride
}

func (m *Ascmap) At(x, y int) byte {
	if !m.In(x, y) {
		return 0
	}
	return m.P[m.Offset(x, y)]
}

func (m *Ascmap) Clone() *Ascmap {
	x := &Ascmap{}
	*x = *m

	x.P = make([]byte, len(m.P))
	copy(x.P, m.P)

	return x
}
