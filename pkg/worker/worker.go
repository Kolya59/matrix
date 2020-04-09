package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
)

const (
	workerPath = `bin/worker`
)

type Worker struct {
	Matrix  [][]int `json:"matrix"`
	Vector  []int   `json:"vector"`
	Columns int     `json:"columns"`
	Lines   int     `json:"lines"`
	Line    int     `json:"line"`
}

func CalculateMultiplyByWorker(matrix [][]int, vector []int, line, n, m int, currDir string) (int, error) {
	addr := fmt.Sprintf("%s/socket/%d", currDir, line)
	go func() {
		if err := handleWorker(matrix, vector, line, n, m, addr); err != nil {
			log.Fatal("Failed to handle worker: ", err)
		}
	}()
	cmd := exec.Command(fmt.Sprintf(`%s/%s`, currDir, workerPath), addr)

	resultBytes, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return 0, fmt.Errorf("error from worker: %s", exitErr.Stderr)
		}
		return 0, errors.New("failed to read proc result")
	}

	result, err := strconv.Atoi(string(resultBytes))
	if err != nil {
		return 0, errors.New("failed to convert result")
	}

	return result, nil
}

func handleWorker(matrix [][]int, vector []int, line, n, m int, rawAddr string) error {
	addr, err := net.ResolveUnixAddr("unix", rawAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve UNIX addr: %v", err)
	}

	clean := func() {
		if err := os.RemoveAll(rawAddr); err != nil {
			log.Println("Error remove all: ", err)
		}
	}

	clean()
	defer clean()

	ls, err := net.ListenUnix("unix", addr)
	if err != nil {
		return fmt.Errorf("failed to start listen: %v", err)
	}
	defer func() {
		if err := ls.Close(); err != nil {
			log.Println("Failed to close connection: ", err)
		}
	}()

	conn, err := ls.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept conn: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Failed to close conn: ", err)
		}
	}()

	workerStruct := Worker{
		Matrix:  matrix,
		Vector:  vector,
		Columns: n,
		Lines:   m,
		Line:    line,
	}

	data, err := json.Marshal(workerStruct)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %v", err)
	}

	return nil
}
