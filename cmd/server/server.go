package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kolya59/matrix/pkg/files"
	"github.com/kolya59/matrix/pkg/worker"
)

const (
	matrixPath = "static/matrix.txt"
	vectorPath = "static/vector.txt"
	resultPath = "static/result.txt"
)

func main() {
	matrixFile, err := ioutil.ReadFile(matrixPath)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}

	vectorFile, err := ioutil.ReadFile(vectorPath)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}

	matrix, vector, n, m, err := files.ParseFiles(matrixFile, vectorFile)
	if err != nil {
		log.Fatal("Failed to parse files: ", err)
	}

	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current dir: ", err)
	}

	resultVector, err := calculate(matrix, vector, n, m, currDir)
	if err != nil {
		log.Fatal("Failed to calculate: ", err)
	}

	resultVectorString := ""
	for i := 0; i < m; i++ {
		resultVectorString += fmt.Sprintf("%d", resultVector[i])
	}

	if err = ioutil.WriteFile(resultPath, []byte(resultVectorString), 0644); err != nil {
		log.Fatal("Failed to write result")
	}

	fmt.Printf("Result:\n%s", resultVectorString)
}

func calculate(matrix [][]int, vector []int, n, m int, currDir string) ([]int, error) {
	resultVector := make([]int, n)

	for i := 0; i < m; i++ {
		matrixLine := make([]int, n)
		for j := 0; j < n; j++ {
			matrixLine[j] = matrix[i][j]
		}
		result, err := worker.CalculateMultiplyByWorker(matrix, vector, i, n, m, currDir)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate by worker %d: %v", i, err)
		}
		resultVector[i] = result
	}

	return resultVector, nil
}
