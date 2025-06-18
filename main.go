package main

import (
	"log"
)

func main() {
	var input_logs = []string {"logs/server1.log", "logs/server2.log"}

	error := ProcessLogs(input_logs, "output.log")

	if error != nil {
		log.Fatal(error)
	}

	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }
}
