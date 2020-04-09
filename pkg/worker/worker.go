package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"

	"github.com/kolya59/matrix/pkg/channels"
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

func CalculateMultiplyByWorker(matrix [][]int, vector []int, line, n, m int, currDir string, chGroup *channels.ChanGroups) {
	addr := fmt.Sprintf("%s/socket/%d", currDir, line)
	go handleWorker(matrix, vector, line, n, m, addr, chGroup)
	cmd := exec.Command(fmt.Sprintf(`%s/%s`, currDir, workerPath), addr)

	resultBytes, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			chGroup.Errors <- fmt.Errorf("error from worker %d: %s", line, exitErr.Stderr)
			return
		}
		chGroup.Errors <- fmt.Errorf("failed to read worker %d result: %s", line, err)
		return
	}

	result, err := strconv.Atoi(string(resultBytes))
	if err != nil {
		chGroup.Errors <- fmt.Errorf("failed to convert worker %d result: %s", line, err)
		return
	}

	chGroup.Results <- channels.Result{
		WorkerNumber: line,
		Value:        result,
	}
}

func handleWorker(matrix [][]int, vector []int, line, n, m int, rawAddr string, chGroup *channels.ChanGroups) {
	addr, err := net.ResolveUnixAddr("unix", rawAddr)
	if err != nil {
		chGroup.Errors <- fmt.Errorf("failed to resolve UNIX addr for worker %d: %v", line, err)
		return
	}

	clean := func() {
		if err := os.RemoveAll(rawAddr); err != nil {
			chGroup.Errors <- fmt.Errorf("failed to remove sock for worker %d: %v", line, err)
		}
	}

	clean()
	defer clean()

	ls, err := net.ListenUnix("unix", addr)
	if err != nil {
		chGroup.Errors <- fmt.Errorf("failed to start listen for worker %d: %v", line, err)
		return
	}
	defer func() {
		if err := ls.Close(); err != nil {
			chGroup.Errors <- fmt.Errorf("failed to close conn for worker %d: %v", line, err)
		}
	}()

	conn, err := ls.Accept()
	if err != nil {
		chGroup.Errors <- fmt.Errorf("failed to accept conn for worker %d: %v", line, err)
		return
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
		chGroup.Errors <- fmt.Errorf("failed to marshal data for worker %d: %v", line, err)
		return
	}

	if _, err := conn.Write(data); err != nil {
		chGroup.Errors <- fmt.Errorf("failed to write data for worker %d: %v", line, err)
		return
	}

	return
}
