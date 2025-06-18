package main

import (
	"log"
	"testing"
)

func TestCache(t *testing.T) {
	var input_logs = []string {"logs/server1.log", "logs/server2.log"}

	error := ProcessLogs(input_logs, "output.log")

	if error != nil {
		log.Fatal(error)
	}

	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }
}
