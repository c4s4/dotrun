package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

// Version as printed -version option
var Version = "UNKNOWN"

const (
	// Help is the command line help
	Help = `Usage: dotrun [-version] [-env .env] [-shell] [-only] command args...
-version    Print version and exit
-env file   Alternative dotenv file (may be repeated to load multiple files)
-shell      Use to a shell to run command
-only       Delete all environment variables except those defined in env files
command     The command to run
args        The command arguments`
)

// ParseCommandLine parses command line and returns:
// - options: passed on command line
// Returns:
// - boolean telling if we should print version
// - boolean telling if we should print help
// - boolean telling if we run command in a shell
// - list of environment file to load
// - command to run
// - command arguments
// - error if any
func ParseCommandLine(options []string) (version bool, help bool, shell bool,
	only bool, envFiles []string, command string, args []string, err error) {
	nextOption := true
	nextEnvFile := false
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
				} else if arg == "-shell" {
					shell = true
				} else if arg == "-only" {
					only = true
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
		return false, false, false, false, nil, "", nil, fmt.Errorf("You must pass command to run on command line")
	}
	return version, help, shell, only, envFiles, command, args, nil
}

// Execute runs command with given arguments and return exit value.
func Execute(shell bool, command string, args ...string) int {
	if shell {
		args = append([]string{command}, args...)
		command = strings.Join(args, " ")
		if runtime.GOOS != "windows" {
			args = append([]string{"-c"}, command)
			command = "sh"
		} else {
			args = append([]string{"/c"}, command)
			command = "cmd"
		}
	}
	cmd := exec.Command(command, args...) // #nosec G204
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

// LoadEnv loads environment in given file
func LoadEnv(filename string) error {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return err
	}
	defer file.Close() // #nosec G307
	reader := bufio.NewReader(file)
	for {
		bytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(bytes))
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		index := strings.Index(line, "=")
		if index < 0 {
			return fmt.Errorf("bad environment line: '%s'", line)
		}
		name := strings.TrimSpace(line[:index])
		value := strings.TrimSpace(line[index+1:])
		if err = os.Setenv(name, value); err != nil {
			return fmt.Errorf("setting environment variable '%s': %v", name, err)
		}
	}
	return nil
}

func main() {
	version, help, shell, only, envFiles, command, args, err := ParseCommandLine(os.Args[1:])
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
	if only {
		os.Clearenv()
	}
	for _, file := range envFiles {
		err := LoadEnv(file)
		if err != nil {
			println(fmt.Sprintf("ERROR loading dotenv file '%s': %v", file, err))
			os.Exit(2)
		}
	}
	os.Exit(Execute(shell, command, args...))
}
