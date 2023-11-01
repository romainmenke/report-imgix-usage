package prompt

import "fmt"

type matrix interface {
	cumulative() matrix
	toString() [][]string
	size(x int, y int) matrix
	insert(x int, y int, value interface{}) matrix
}

type intMatrix [][]int

func (v intMatrix) cumulative() matrix {
	output := make([][]int, len(v))

	for y, yv := range v {
		output[y] = make([]int, len(yv))
		var cumulative int
		for x, xv := range yv {
			cumulative = cumulative + xv
			output[y][x] = cumulative
		}
	}

	return intMatrix(output)
}

func (v intMatrix) toString() [][]string {
	output := make([][]string, len(v))

	for y, yv := range v {
		output[y] = make([]string, len(yv))
		for x, xv := range yv {
			output[y][x] = fmt.Sprint(xv)
		}
	}

	return output
}

func (v intMatrix) size(x int, y int) matrix {
	growY := y - len(v)
	if growY > 0 {
		v = append(v, make([][]int, growY)...)
	}

	for i, yv := range v {
		growX := x - len(yv)
		if growX > 0 {
			v[i] = append(v[i], make([]int, growX)...)
		}
	}

	return v
}

func (v intMatrix) insert(x int, y int, value interface{}) matrix {
	if intval, ok := value.(int); ok {
		v[y][x] = intval
	}

	return v
}

type floatMatrix [][]float64

func (v floatMatrix) cumulative() matrix {
	output := make([][]float64, len(v))

	for y, yv := range v {
		output[y] = make([]float64, len(yv))
		var cumulative float64
		for x, xv := range yv {
			cumulative = cumulative + xv
			output[y][x] = cumulative
		}
	}

	return floatMatrix(output)
}

func (v floatMatrix) toString() [][]string {
	output := make([][]string, len(v))

	for y, yv := range v {
		output[y] = make([]string, len(yv))
		for x, xv := range yv {
			output[y][x] = fmt.Sprintf("%.2f", xv)
		}
	}

	return output
}

func (v floatMatrix) size(x int, y int) matrix {
	growY := y - len(v)
	if growY > 0 {
		v = append(v, make([][]float64, growY)...)
	}

	for i, yv := range v {
		growX := x - len(yv)
		if growX > 0 {
			v[i] = append(v[i], make([]float64, growX)...)
		}
	}

	return v
}

func (v floatMatrix) insert(x int, y int, value interface{}) matrix {
	if intval, ok := value.(float64); ok {
		v[y][x] = intval
	}

	return v
}
