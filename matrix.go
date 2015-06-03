/*
Package "matrix" provides types and operations for matrix manipulation.
*/
package matrix

type Matrix interface {
	// Return the shape of matrix, which consists of the "rows" and the "columns".
	Shape() (rows, columns int)

	// Return the "rows" of matrix.
	Rows() (rows int)

	// Return the "columns" of matrix.
	Columns() (columns int)

	// Return a row of matrix speficied with index "row".
	Row(row int) Row

	// Return a column of matrix speficied with index "column".
	Column(column int) Column
}

type Row interface {
}

type Column interface {
}