package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/kolya59/matrix/pkg/worker"
)

func main() {
	args := os.Args
	if len(args) != 1 {
		os.Exit(1)
	}

	rawAddr := args[0]

	addr, err := net.ResolveUnixAddr("unix", rawAddr)
	if err != nil {
		os.Exit(1)
	}

	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		_ = conn.Close()
	}()

	data, err := ioutil.ReadAll(conn)
	if err != nil {
		os.Exit(1)
	}

	var workerData worker.Worker
	if err := json.Unmarshal(data, workerData); err != nil {
		os.Exit(1)
	}

	result := 0
	for i := 0; i < workerData.Columns; i++ {
		result += workerData.Matrix[workerData.Line][i] * workerData.Vector[i]
	}

	writer := bufio.NewWriter(os.Stdout)
	if _, err = writer.WriteString(fmt.Sprintf("%d", result)); err != nil {
		os.Exit(1)
	}
}
