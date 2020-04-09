package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kolya59/matrix/pkg/channels"
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

	done := make(chan os.Signal)
	go func() {
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
		<-done
		close(done)
	}()

	resultVector, err := calculate(matrix, vector, n, m, currDir, done)
	if err != nil {
		log.Fatal("Failed to calculate: ", err)
	}

	resultVectorString := ""
	for i := 0; i < m; i++ {
		resultVectorString += fmt.Sprintf("%d\n", resultVector[i])
	}

	if err = ioutil.WriteFile(resultPath, []byte(resultVectorString), 0644); err != nil {
		log.Fatal("Failed to write result")
	}

	fmt.Printf("Result:\n%s", resultVectorString)
}

func calculate(matrix [][]int, vector []int, n, m int, currDir string, done chan os.Signal) ([]int, error) {
	resultVector := make([]int, n)

	chGroup := channels.NewChanGroups(done)
	wgCounter := 0

	errs := make(chan error)

	for i := 0; i < m; i++ {
		matrixLine := make([]int, n)
		for j := 0; j < n; j++ {
			matrixLine[j] = matrix[i][j]
		}
		wgCounter++
		go worker.CalculateMultiplyByWorker(matrix, vector, i, n, m, currDir, chGroup)
	}

	go func() {
		for {
			select {
			case <-chGroup.Done:
				errs <- nil
				return
			case err := <-chGroup.Errors:
				errs <- err
				return
			case res := <-chGroup.Results:
				wgCounter--
				resultVector[res.WorkerNumber] = res.Value
				if wgCounter == 0 {
					errs <- nil
					return
				}
			}
		}
	}()

	if err := <-errs; err != nil {
		return nil, err
	}

	return resultVector, nil
}
