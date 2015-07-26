/*
Package "dense" provides an implementation of mutable dense matrix.
*/
package dense

import (
	"github.com/mitsuse/matrix-go/internal/rewriters"
	"github.com/mitsuse/matrix-go/internal/types"
	"github.com/mitsuse/matrix-go/internal/validates"
)

type matrixImpl struct {
	rows     int
	columns  int
	elements []float64
	rewriter rewriters.Rewriter
}

func New(rows, columns int) func(elements ...float64) types.Matrix {
	validates.ShapeShouldBePositive(rows, columns)

	constructor := func(elements ...float64) types.Matrix {
		size := rows * columns

		if len(elements) != size {
			panic(validates.INVALID_ELEMENTS_PANIC)
		}

		m := &matrixImpl{
			rows:     rows,
			columns:  columns,
			elements: make([]float64, size),
			rewriter: rewriters.Reflect(),
		}
		copy(m.elements, elements)

		return m
	}

	return constructor
}

func Zeros(rows, columns int) types.Matrix {
	return New(rows, columns)(make([]float64, rows*columns)...)
}

func (m *matrixImpl) Shape() (rows, columns int) {
	return m.rewriter.Rewrite(m.rows, m.columns)
}

func (m *matrixImpl) Rows() (rows int) {
	rows, _ = m.Shape()
	return rows
}

func (m *matrixImpl) Columns() (columns int) {
	_, columns = m.Shape()
	return columns
}

func (m *matrixImpl) All() types.Cursor {
	return newAllCursor(m)
}

func (m *matrixImpl) NonZeros() types.Cursor {
	return newNonZerosCursor(m)
}

func (m *matrixImpl) Diagonal() types.Cursor {
	return newDiagonalCursor(m)
}

func (m *matrixImpl) Get(row, column int) (element float64) {
	rows, columns := m.Shape()

	validates.IndexShouldBeInRange(rows, columns, row, column)

	row, column = m.rewriter.Rewrite(row, column)

	return m.elements[row*m.columns+column]
}

func (m *matrixImpl) Update(row, column int, element float64) types.Matrix {
	rows, columns := m.Shape()

	validates.IndexShouldBeInRange(rows, columns, row, column)

	row, column = m.rewriter.Rewrite(row, column)

	m.elements[row*m.columns+column] = element

	return m
}

func (m *matrixImpl) Equal(n types.Matrix) bool {
	validates.ShapeShouldBeSame(m, n)

	cursor := n.All()

	for cursor.HasNext() {
		element, row, column := cursor.Get()
		if m.Get(row, column) != element {
			return false
		}
	}

	return true
}

func (m *matrixImpl) Add(n types.Matrix) types.Matrix {
	validates.ShapeShouldBeSame(m, n)

	cursor := n.NonZeros()

	for cursor.HasNext() {
		element, row, column := cursor.Get()
		m.Update(row, column, m.Get(row, column)+element)
	}

	return m
}

func (m *matrixImpl) Subtract(n types.Matrix) types.Matrix {
	validates.ShapeShouldBeSame(m, n)

	cursor := n.NonZeros()

	for cursor.HasNext() {
		element, row, column := cursor.Get()
		m.Update(row, column, m.Get(row, column)-element)
	}

	return m
}

func (m *matrixImpl) Dot(n types.Matrix) types.Matrix {
	validates.ShapeShouldBeMultipliable(m, n)

	rows := m.Rows()
	columns := n.Columns()

	r := Zeros(rows, columns)

	cursor := n.NonZeros()

	for cursor.HasNext() {
		element, j, k := cursor.Get()

		for i := 0; i < rows; i++ {
			r.Update(i, k, r.Get(i, k)+m.Get(i, j)*element)
		}
	}

	return r
}

func (m *matrixImpl) Multiply(s float64) types.Matrix {
	for index, element := range m.elements {
		m.elements[index] = element * s
	}

	return m
}

func (m *matrixImpl) Transpose() types.Matrix {
	n := &matrixImpl{
		rows:     m.rows,
		columns:  m.columns,
		elements: m.elements,
		rewriter: m.rewriter.Transpose(),
	}

	return n
}
