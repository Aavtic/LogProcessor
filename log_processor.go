package main

import (
	"fmt"
	"bufio"
	"strings"
	"os"
	"sync"
)

type LogProcessort interface {
	read_log(scanner *bufio.Scanner)
	is_error(line string) bool
	write_log(log string)
	process_logs(scanner *bufio.Scanner)
}

type LogProcessor struct {
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
	for scanner.Scan() {
		line := scanner.Text()
		if proc.is_error(line) {
			proc.error_w <- line
		}
	}
	wg.Done()
}

func (proc *LogProcessor) write_log(line string, wg *sync.WaitGroup) {
	fmt.Println(line)
	f, err := os.OpenFile(proc.output_file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	if _, err = f.WriteString(line)
	err != nil {
    panic(err)
	}
	wg.Done()
}

func (proc *LogProcessor) process_logs() {
	error_ch := make(chan string, 2)
	proc.error_r = error_ch
	proc.error_w = error_ch

	wg := sync.WaitGroup{}
	
	for _, scanner := range proc.scanners {
		wg.Add(1)
		go proc.read_log(scanner, &wg)
	}

	wg.Add(1)
	go func () {
		for {
			line, ok := <- proc.error_r
			if !ok {
				fmt.Println("Closed")
			}
			proc.write_log(line, &wg)
		}
	}()

	wg.Wait()
	close(error_ch)
}

func get_scanner(path string) (*bufio.Scanner, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bufio.NewReader(file))

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

	log_processor.process_logs()

	return nil
}
