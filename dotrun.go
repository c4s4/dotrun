package main

import (
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
	help = `Usage: dotrun [-file .env] command args...
-file file   Alternative dotenv file
command      The command to run
args         The command arguments`
)

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

func expandPath(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(dir, path[2:])
	}
	return path
}

func main() {
	if os.Args[1] == "-help" {
		println(help)
		os.Exit(0)
	}
	if len(os.Args) < 2 {
		println("ERROR you must pass command to run on command line")
		println(help)
		os.Exit(-1)
	}
	var file = ".env"
	var command string
	var args []string
	if os.Args[1] == "-file" {
		if len(os.Args) < 3 {
			println("ERROR you must specify dotenv file after '-file' option")
			os.Exit(-1)
		}
		file = os.Args[2]
		command = os.Args[3]
		args = os.Args[4:]
	} else {
		command = os.Args[1]
		args = os.Args[2:]
	}
	file = expandPath(file)
	err := godotenv.Load(file)
	if err != nil {
		println(fmt.Sprintf("ERROR loading dotenv file '%s': %v", file, err))
	}
	exit := Execute(command, args...)
	if exit != 0 {
		os.Exit(exit)
	}
}
