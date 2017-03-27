package ioproc

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func IsJSONString(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil

}

func IsJSON(s string) bool {
	// var js map[string]interface{}
	// return json.Unmarshal([]byte(s), &js) == nil
	var val map[string]interface{}
	if err := json.Unmarshal([]byte(s), &val); err != nil {
		return false
	}
	return true

}
func DeleteGzip(path string, filenameWithoutFormat string) {
	if FileExistGzip(path, filenameWithoutFormat) {
		err := os.Remove(filepath.Join(path, filenameWithoutFormat+".gz"))
		if err != nil {
			println(err.Error())
		}

	}
}
func ReadGzip(path string, filenameWithoutFormat string) string {
	//file, err := os.Open(filepath.Join(path, filenameWithoutFormat+".gz"))
	file, err := os.OpenFile(filepath.Join(path, filenameWithoutFormat+".gz"), syscall.O_RDONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	gz, err := gzip.NewReader(file)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	defer gz.Close()

	//scanner := bufio.NewScanner(gz)
	scanner := bufio.NewReader(gz)
	data, err := ioutil.ReadAll(scanner)

	return string(data)
	// //bt := make([]byte, 1024) // Read by 1 KiB.
	// bt := make([]byte, 102400)
	// //var bt []byte
	// f, _ := os.Open(filepath.Join(path, filenameWithoutFormat+".gz"))
	// //f, _ := os.OpenFile(filepath.Join(path, filenameWithoutFormat+".gz"), syscall.O_RDONLY, 0644)
	// gz, _ := gzip.NewReader(f)
	// d, _ := gz.Read(bt)
	// fileContents := bt[:d]
	// gz.Close()
	// f.Close()
	// return string(fileContents)

}
func WriteGzipStr(path string, filenameWithoutFormat string, data string) {
	WriteGzip(path, filenameWithoutFormat, []byte(data))
}
func WriteGzip(path string, filenameWithoutFormat string, data []byte) {
	f, _ := os.Create(filepath.Join(path, filenameWithoutFormat+".gz"))
	gz := gzip.NewWriter(f)
	gz.Write(data)
	gz.Close()
	f.Close()
}

// fileExist return true is exist
func FileExist(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}
func FileExistGzip(path string, filenameWithoutFormat string) bool {
	if _, err := os.Stat(filepath.Join(path, filenameWithoutFormat+".gz")); err == nil {
		return true
	}
	return false
}

// directoryexist return true is exist
func Directoryexist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func Createdirectory(path string, folderName string) {
	if !Directoryexist(filepath.Join(path, folderName)) {
		os.MkdirAll(filepath.Join(path, folderName), 0777)
	}
}
func ReadFile(filename string) []byte {
	file, _ := ioutil.ReadFile(filename)
	return file
}
func ReadFileStr(filename string) string {
	file, _ := ioutil.ReadFile(filename)
	return string(file)
}
func WriteFile(filename string, data []byte) {
	ioutil.WriteFile(filename, data, 0644)
}

func WriteFileStr(filename string, data string) {
	ioutil.WriteFile(filename, []byte(data), 0644)
}

func AppendFileStr(filename string, data string, newLine bool) bool {
	if !FileExist(filename) {
		WriteFileStr(filename, "")
	}
	//os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	f, err := os.OpenFile(filename, os.O_APPEND, 0600)
	if err == nil {
		if newLine {
			w := csv.NewWriter(f)
			w.UseCRLF = true
			w.Write([]string{data})
			w.Flush()
		} else {
			f.WriteString(data)
		}
		f.Close()
		return true
	}
	defer f.Close()
	return false
}
func ReadLines(filename string) []string {
	var lines []string
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return lines
	}
	buf := bytes.NewBuffer(file)
	for {
		line, err := buf.ReadString('\n')
		if len(line) == 0 {
			if err != nil {
				if err == io.EOF {
					break
				}
				return lines
			}
		}
		lines = append(lines, line)
		if err != nil && err != io.EOF {
			return lines
		}
	}
	return lines
}

func ReadLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return scanner.Text()
	}
	return ""
}
