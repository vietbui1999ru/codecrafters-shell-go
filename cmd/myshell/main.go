package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

var commands map[string]func(string)

func init() {
  commands = make(map[string]func(string)) 
  commands["exit"] = exitCommand
  commands["echo"] = echoCommand
  commands["type"] = typeCommand
}

func checkCommand(command string, args string) {
  if cmd, ok := commands[command]; ok {
    cmd(args) // execute the command
    return
  } else {
    cmd := exec.Command(command, args)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
      fmt.Printf("%s: command not found\n", command)
      return
    }
  }
  fmt.Printf("%s: command not found\n", command)
}

func exitCommand(args string) {
  number, err := strconv.Atoi(args)
  if err != nil {
    fmt.Printf("%v - %s: not a valid number\n", number, err)
    return
  }
  os.Exit(number)
}

func echoCommand(args string) {
  fmt.Printf("%s\n", args)
}

func typeCommand(args string) {
  if _, ok := commands[args]; ok {
      fmt.Printf("%s is a shell builtin\n", args)
      return
  } else {

    paths := os.Getenv("PATH")
    pathList := strings.Split(paths, ":")
    for _, path := range pathList {
      if _, err := os.Stat(filepath.Join(path, args)); err == nil {
        fmt.Printf("%s is %s/%s\n", args, path, args)
        return
      }
    }

  }
  fmt.Printf("%s: not found\n", args)
}

func main() {
  for {
    fmt.Fprint(os.Stdout, "$ ")

    // Wait for user input
    input, err := bufio.NewReader(os.Stdin).ReadString('\n')
    trimmpedInput := strings.TrimSpace(input)
    // fmt.Printf("%s: command not found\n", trimmpedInput)
    
    if err != nil {
      fmt.Printf("%s: invalid input\n", input)
    }

    handleCommands(trimmpedInput)
  }
}

// comment again
func handleCommands(input string) {
  cmd, args, _ := strings.Cut(input, " ")
  checkCommand(cmd, args)
}

