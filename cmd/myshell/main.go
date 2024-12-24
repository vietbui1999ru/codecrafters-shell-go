package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
  "strconv"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

var commands map[string]func(string)

func init() {
  commands = make(map[string]func(string)) 
  commands["exit"] = exitCommand
}

func checkCommand(command string, args string) {
  if cmd, ok := commands[command]; ok {
    cmd(args) // execute the command
  } else {
    fmt.Printf("Command %s not found\n", command)
  }
}

func exitCommand(args string) {
  number, err := strconv.Atoi(args)
  if err != nil {
    fmt.Printf("%v - %s: not a valid number\n", number, err)
    return
  }
  os.Exit(number)
}
func main() {
  for {
    fmt.Fprint(os.Stdout, "$ ")

    // Wait for user input
    input, err := bufio.NewReader(os.Stdin).ReadString('\n')
    fmt.Printf("%s: command not found\n", strings.TrimSpace(input))
    
    if err != nil {
      fmt.Printf("%s: invalid input\n", input)
    }

    trimmpedInput := strings.TrimSpace(input)
    handleCommands(trimmpedInput)
  }
}

// comment again
func handleCommands(input string) {
  cmd, args, _ := strings.Cut(input, " ")
  checkCommand(cmd, args)
}

