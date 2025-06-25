# 🔍 Go Log Processor

A fast and efficient log parser written in Go that scans multiple log files **simultaneously**, extracts lines that start with `ERROR`, and writes them to a specified output file. Designed with performance and concurrency in mind.

---

## 🧠 What it Does

- 🗂️ Parses **multiple log files** at once.
- 🚨 Filters lines that start with `ERROR`.
- 📄 Writes the results to a specified output file.
- 🧵 Uses **concurrent goroutines** for parallel processing.
- 🧠 Employs **RWMutex** to ensure thread-safe access to shared channels.
- 🧠 Reads and writes using buffers with a max size of **1MB** to keep memory usage under control – even for large log files.

---

## 🏁 How to Use

1. **Build the binary**:

   ```bash
   go build -o logprocessor main.go
