package gscript

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"math/rand"
	"time"
)

func CalledBy() string {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "Unknown"
	}
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "N/A"
	}
	return fun.Name()
}

func LocalFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func LocalFileCreate(path string, bytes []byte) error {
	if LocalFileExists(path) {
		return errors.New("The file to create already exists so we won't overwite it")
	}
	err := ioutil.WriteFile(path, bytes, 0700)
	if err != nil {
		return err
	}
	return nil
}

// LocalFileAppendBytes adds bytes to the end of filename's path.
func LocalFileAppendBytes(filename string, bytes []byte) error {
	if LocalFileExists(filename) {
  	fileInfo, err := os.Stat(filename)
	  if err != nil {
		  return err
	  }
	  file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, fileInfo.Mode())
	  if err != nil {
		  return err
	  }
	  if _, err = file.Write(bytes); err != nil {
		  return err
	  }
	  file.Close()
	  return nil
	} else {
		return errors.New("The file dosn't exist so we should create it in the future")
	}
}

// LocalFileAppendString adds input as strings to the end of filename's path.
func LocalFileAppendString(input, filename string) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_APPEND, fileInfo.Mode())
	if err != nil {
		return err
	}
	if _, err = file.WriteString(input); err != nil {
		return err
	}
	file.Close()
	return nil
}

// Replace will replace all instances of match with replace in file.
func LocalFileReplace(file, match, replacement string) error {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return err
	}
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	lines := strings.Split(string(contents), "\n")
	for index, line := range lines {
		if strings.Contains(line, match) {
			lines[index] = replacement
		}
	}
	ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), fileInfo.Mode())
	return nil
}

// ReplaceMulti will replace all instances of possible matches with replacement in file.
func LocalFileReplaceMulti(file string, matches []string, replacement string) error {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return err
	}
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	lines := strings.Split(string(contents), "\n")
	for index, line := range lines {
		for _, match := range matches {
			if strings.Contains(line, match) {
				lines[index] = replacement
			}
		}
	}
	ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), fileInfo.Mode())
	return nil
}

// LocalReadFile takes a file path and returns the byte array of the file there
func LocalFileRead(path string) ([]byte, error) {
	if LocalFileExists(path) {
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return dat, nil
	}
	return nil, errors.New("The file to read does not exist")
}

// ExecuteCommand function
func ExecuteCommand(c string, args ...string) VMExecResponse {
	cmd := exec.Command(c, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	respObj := VMExecResponse{
		Stdout: strings.Split(stdout.String(), "\n"),
		Stderr: strings.Split(stderr.String(), "\n"),
		PID:    cmd.Process.Pid,
	}
	if err != nil {
		respObj.ErrorMsg = err.Error()
		respObj.Success = false
	} else {
		respObj.Success = true
	}
	return respObj
}

// HTTPGetFile takes a url and returns a byte slice of the file there
func HTTPGetFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	pageData, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return pageData, nil
}

// RandString returns a string the length of strlen
func RandString(strlen int) string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

// RandomInt returns an int inbetween min and max.
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}
