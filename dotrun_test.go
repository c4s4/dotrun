package main

import (
	"reflect"
	"testing"
)

func TestParseCommandLineError(t *testing.T) {
	version, help, shell, only, envFiles, command, args, err := ParseCommandLine(nil)
	if err == nil {
		t.Fatalf("Bad command line parsing: %v, %v, %v, %v, %v, %v, %v, %v", version, help, shell, only, envFiles, command, args, err)
	}
}

func TestParseCommandLineHelp(t *testing.T) {
	version, help, shell, only, envFiles, command, args, err := ParseCommandLine([]string{"-help"})
	if !help {
		t.Fatalf("Bad command line parsing: %v, %v, %v, %v, %v, %v, %v, %v", version, help, shell, only, envFiles, command, args, err)
	}
}

func TestParseCommandLineEnvFiles(t *testing.T) {
	version, help, shell, only, envFiles, command, args, err := ParseCommandLine([]string{"-env", "env1", "-env", "env2", "command"})
	if !reflect.DeepEqual([]string{"env1", "env2"}, envFiles) {
		t.Fatalf("Bad command line parsing: %v, %v, %v, %v, %v, %v, %v, %v", version, help, shell, only, envFiles, command, args, err)
	}
}

func TestParseCommandLineEnvFilesMissingCommand(t *testing.T) {
	version, help, shell, only, envFiles, command, args, err := ParseCommandLine([]string{"-env", "env1", "-env", "env2"})
	if err == nil {
		t.Fatalf("Bad command line parsing: %v, %v, %v, %v, %v, %v, %v, %v", version, help, shell, only, envFiles, command, args, err)
	}
}

func TestParseCommandLineNominal(t *testing.T) {
	version, help, shell, only, envFiles, command, args, err := ParseCommandLine([]string{"-env", "env1", "-env", "env2", "command", "arg1", "arg2"})
	if help ||
		!reflect.DeepEqual([]string{"env1", "env2"}, envFiles) ||
		command != "command" ||
		!reflect.DeepEqual([]string{"arg1", "arg2"}, args) ||
		err != nil {
		t.Fatalf("Bad command line parsing: %v, %v, %v, %v, %v, %v, %v, %v", version, help, shell, only, envFiles, command, args, err)
	}
}
