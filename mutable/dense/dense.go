/*
Package "dense" provides an implementation of "Matrix" which stores elements in a slide.
*/
package dense

import (
	"errors"
	"fmt"

	"github.com/mitsuse/matrix-go"
)

type matrixImpl struct {
	rows     int
	columns  int
	elements []float64
}

func New(rows, columns int) func(elements ...float64) (matrix.Matrix, error) {
	rowsShouldBePositiveNumber(rows)
	columnShouldBePositiveNumber(rows)

	constructor := func(elements ...float64) (matrix.Matrix, error) {
		size := rows * columns

		if len(elements) != size {
			template := "The number of %q should equal to %q * %q."
			message := fmt.Sprintf(template, "elements", "rows", "columns")

			return nil, errors.New(message)
		}

		m := &matrixImpl{
			rows:     rows,
			columns:  columns,
			elements: make([]float64, size),
		}
		copy(m.elements, elements)

		return m, nil
	}

	return constructor
}
