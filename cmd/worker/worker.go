package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/kolya59/matrix/pkg/worker"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("invalid args: %v", args)
	}

	rawAddr := args[1]

	addr, err := net.ResolveUnixAddr("unix", rawAddr)
	if err != nil {
		log.Fatalf("failed to resolve unix addr: %v", err)
	}

	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	data, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Fatalf("failed to read from conn: %v", err)
	}

	var workerData worker.Worker
	if err := json.Unmarshal(data, &workerData); err != nil {
		log.Fatalf("failed to unmarshal data: %v", err)
	}

	result := 0
	for i := 0; i < workerData.Columns; i++ {
		result += workerData.Matrix[workerData.Line][i] * workerData.Vector[i]
	}

	fmt.Print(result)
}
