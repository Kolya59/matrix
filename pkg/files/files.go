package files

import (
	"bytes"
	"errors"
	"strconv"
)

var (
	lineDelimiter = []byte{'\n'}
	delimiter     = []byte{' '}
)

func ParseFiles(matrixFile, vectorFile []byte) ([][]int, []int, int, int, error) {
	matrixSplit := bytes.Split(matrixFile, lineDelimiter)
	if len(matrixSplit) == 1 {
		return nil, nil, 0, 0, errors.New("invalid matrix")
	}

	header := bytes.Split(matrixSplit[0], delimiter)
	if len(header) != 2 {
		return nil, nil, 0, 0, errors.New("invalid header")
	}

	n, err := strconv.Atoi(string(header[0]))
	if err != nil {
		return nil, nil, 0, 0, errors.New("invalid header")
	}

	m, err := strconv.Atoi(string(header[1]))
	if err != nil {
		return nil, nil, 0, 0, errors.New("invalid header")
	}

	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		tmp := make([]int, m)
		tmpSplit := bytes.Split(matrixSplit[i+1], delimiter)
		for j := 0; j < m; j++ {
			tmpInt, err := strconv.Atoi(string(tmpSplit[j]))
			if err != nil {
				return nil, nil, 0, 0, errors.New("invalid matrix")
			}
			tmp[j] = tmpInt
		}
		matrix[i] = tmp
	}

	vectorSplit := bytes.Split(vectorFile, lineDelimiter)
	vector := make([]int, m)
	for i := 0; i < m; i++ {
		tmp, err := strconv.Atoi(string(vectorSplit[i]))
		if err != nil {
			return nil, nil, 0, 0, errors.New("invalid vector")
		}
		vector[i] = tmp
	}

	return matrix, vector, n, m, nil
}
