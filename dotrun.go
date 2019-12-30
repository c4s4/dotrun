package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
)

// Version as printed -version option
var Version = "UNKNOWN"

const (
	// Help is the command line help
	Help = `Usage: dotrun [-version] [-env .env] command args...
-version    Print version and exit
-env file   Alternative dotenv file (mais be repeated to load multiple files)
command     The command to run
args        The command arguments`
)

// ParseCommandLine parses command line and returns:
// - options: passed on command line
// Returns:
// - boolean telling if we should print help
// - list of environment file to load
// - command to run
// - command arguments
// - error if any
func ParseCommandLine(options []string) (bool, bool, []string, string, []string, error) {
	nextOption := true
	nextEnvFile := false
	var version bool
	var help bool
	var envFiles []string
	var command string
	var args []string
	for _, arg := range options {
		if nextOption {
			if nextEnvFile {
				envFiles = append(envFiles, arg)
				nextEnvFile = false
			} else {
				if arg == "-version" {
					version = true
				} else if arg == "-help" {
					help = true
				} else if arg == "-env" {
					nextEnvFile = true
				} else {
					command = arg
					nextOption = false
				}
			}
		} else {
			args = append(args, arg)
		}
	}
	if !version && !help && command == "" {
		return false, false, nil, "", nil, fmt.Errorf("You must pass command to run on command line")
	}
	return version, help, envFiles, command, args, nil
}

// Execute runs command with given arguments and return exit value.
func Execute(command string, args ...string) int {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	exit := 0
	if err != nil {
		message := err.Error()
		if !strings.HasPrefix(message, "exit status") {
			println(message)
		}
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exit = ws.ExitStatus()
		} else {
			exit = -4
		}
	}
	return exit
}

// ExpandPath expands given path wit user home directory
func ExpandPath(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(dir, path[2:])
	}
	return path
}

// LoadEnv loads environment in given file
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		bytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(bytes))
		if line[0] == '#' {
			continue
		}
		index := strings.Index(line, "=")
		if index < 0 {
			return fmt.Errorf("bad environment line: '%s'", line)
		}
		name := strings.TrimSpace(line[:index])
		value := strings.TrimSpace(line[index+1:])
		os.Setenv(name, value)
	}
	return nil
}

func main() {
	version, help, envFiles, command, args, err := ParseCommandLine(os.Args[1:])
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	if version {
		println(Version)
		os.Exit(0)
	}
	if help {
		println(Help)
		os.Exit(0)
	}
	if envFiles == nil {
		envFiles = []string{".env"}
	}
	for _, file := range envFiles {
		err := LoadEnv(file)
		if err != nil {
			println(fmt.Sprintf("ERROR loading dotenv file '%s': %v", file, err))
			os.Exit(2)
		}
	}
	os.Exit(Execute(command, args...))
}
