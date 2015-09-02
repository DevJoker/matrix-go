/*
Package "hash" provides an hash-based implementation of mutable sparse matrix.
*/
package hash

import (
	"github.com/mitsuse/matrix-go/internal/rewriters"
	"github.com/mitsuse/matrix-go/internal/types"
	"github.com/mitsuse/matrix-go/internal/validates"
)

type Matrix struct {
	initialized bool
	base        *types.Shape
	view        *types.Shape
	offset      *types.Index
	elements    map[int]float64
	rewriter    rewriters.Rewriter
}

type Element struct {
	Row    int
	Column int
	Value  float64
}

func New(rows, columns int) func(elements ...Element) *Matrix {
	validates.ShapeShouldBePositive(rows, columns)

	constructor := func(elements ...Element) *Matrix {
		shape := types.NewShape(rows, columns)
		offset := types.NewIndex(0, 0)

		m := &Matrix{
			initialized: true,
			base:        shape,
			view:        shape,
			offset:      offset,
			elements:    make(map[int]float64),
			rewriter:    rewriters.Reflect(),
		}

		for _, element := range elements {
			key := element.Row*columns + element.Column
			if _, exist := m.elements[key]; exist {
				panic(validates.INVALID_ELEMENTS_PANIC)
			}

			m.elements[key] = element.Value
		}

		return m
	}

	return constructor
}

func Zeros(rows, columns int) *Matrix {
	return New(rows, columns)()
}
