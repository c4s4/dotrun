package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
)

const (
	help = `Usage: dotrun [-env .env] command args...
-env file   Alternative dotenv file (mais be repeated to load multiple files)
command     The command to run
args        The command arguments`
)

// Help tells if we should print help and exit
var Help bool

// EnvFiles is a list of environment files passed on command line
var EnvFiles Strings

// Strings is the type for list of strings
type Strings []string

// String representation for a list of strings
func (s *Strings) String() string {
	return "[" + strings.Join(*s, ", ") + "]"
}

// Set append a string to the list
func (s *Strings) Set(path string) error {
	*s = append(*s, ExpandPath(path))
	return nil
}

// ParseCommandLine parses command line
func init() {
	flag.Bool("-help", false, "Print help and exit")
	flag.Var(&EnvFiles, "env", "Environment file to load")
	flag.Parse()
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

func main() {
	if Help {
		println(help)
		os.Exit(0)
	}
	if EnvFiles == nil {
		EnvFiles = []string{".env"}
	}
	if len(flag.Args()) < 1 {
		println("ERROR you must pass command to run on command line")
		os.Exit(1)
	}
	command := flag.Args()[0]
	args := flag.Args()[1:]
	for _, file := range EnvFiles {
		err := godotenv.Overload(file)
		if err != nil {
			println(fmt.Sprintf("ERROR loading dotenv file '%s': %v", file, err))
		}
	}
	os.Exit(Execute(command, args...))
}
