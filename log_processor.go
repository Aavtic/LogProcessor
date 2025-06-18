package main

import (
	"bufio"
	"strings"
	"os"
	"sync"
	"log"
)

const (
	// 1 MB
	BUFFER_SIZE = 1024 * 1024
)

type LogProcessort interface {
	read_log(scanner *bufio.Scanner)
	is_error(line string) bool
	write_log(log string)
	process_logs(scanner *bufio.Scanner)
}

type LogProcessor struct {
	mut sync.RWMutex
	scanners []*bufio.Scanner
	error_w chan<-string
	error_r <-chan string
	output_file string
}

func (_*LogProcessor) is_error(line string) bool {
	if strings.HasPrefix(line, "ERROR:") {
		return true
	}
	return false
}

func (proc* LogProcessor) read_log(scanner *bufio.Scanner, wg *sync.WaitGroup) {
	proc.mut.Lock()
	defer wg.Done()
	defer proc.mut.Unlock()

	for scanner.Scan() {
		line := scanner.Text()
		if proc.is_error(line) {
			proc.error_w <- line
		}
	}
}

func (proc *LogProcessor) write_log(line string, writer *bufio.Writer) {
	n, err := writer.Write([]byte(line + "\n"))
	if err != nil {
		log.Fatalf("ERROR: written %d, Could not write to file, due to %s", n, err)
	}
}

func (proc *LogProcessor) process_logs()  error {
	error_ch := make(chan string, 200)
	proc.error_r = error_ch
	proc.error_w = error_ch

	read_wg := sync.WaitGroup{}
	
	for _, scanner := range proc.scanners {
		read_wg.Add(1)
		go proc.read_log(scanner, &read_wg)
	}

	write_wg := sync.WaitGroup{}

	write_file, err := os.OpenFile(proc.output_file, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0644)
	writer := bufio.NewWriterSize(write_file, BUFFER_SIZE)

	if err != nil { return err }

	write_wg.Add(1)
	go func () {
		defer write_wg.Done()

		for {
			line, ok := <- proc.error_r
			if !ok {
				log.Println("OVER")
				break
			}
			proc.write_log(line, writer)
		}
	}()

	read_wg.Wait()

	close(error_ch)

	write_wg.Wait()

	error := writer.Flush()
	if error != nil {
		log.Printf("ERROR: Could not write to file due to: %s \n", error)
	}
	error = write_file.Close()
	
	if error != nil {
		log.Printf("ERROR: Could not close file due to: %s \n", error)
	}

	return nil
}

func get_scanner(path string) (*bufio.Scanner, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	buffer := make([]uint8, BUFFER_SIZE)
	scanner.Buffer(buffer, BUFFER_SIZE)

	return scanner, nil
}

func new_log_processor() LogProcessor {
	return LogProcessor {
		scanners: []*bufio.Scanner{},
	}
}

func (log_processor *LogProcessor) scanner_add(scanner *bufio.Scanner) {
	log_processor.scanners = append(log_processor.scanners, scanner)
}

func ProcessLogs(inputFiles []string, outputFile string) error {
	log_processor := new_log_processor()
	log_processor.output_file = outputFile

	for _, file := range inputFiles {
		scanner, err := get_scanner(file)
		if err != nil { return err }
		log_processor.scanner_add(scanner)
	}

	return log_processor.process_logs()
}
