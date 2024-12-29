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

var singleQuotes = "'"
var doubleQuotes = `"`

var commands map[string]func(string)

func init() {
  commands = make(map[string]func(string)) 
  commands["exit"] = exitCommand
  commands["echo"] = echoCommand
  commands["type"] = typeCommand
  commands["pwd"] = pwdCommand
  commands["cd"] = cdCommand
  commands["~"] = homeCommand
}

func checkCommand(command string, args string) {
  if cmd, ok := commands[command]; ok {
    cmd(args) // execute the command
    return
  } else {
    fmt.Printf("args %s\n", args)
    for _, arg := range trimFieldByQuotes(args) {
      // unicode print
      fmt.Printf("arg: %s\n", arg)
      cmd := exec.Command(command, arg)
      cmd.Stdout = os.Stdout
      cmd.Stderr = os.Stderr
      err := cmd.Run()
      if err != nil {
        fmt.Printf("%s: command not found\n", command)
        return
      }
    }
  }
}

func trimFieldByQuotes(s string) []string {
  var fields []string
  var curField []rune
  inQuotes := false
  quoteType := rune(0) // single or double quotes

  for _, char := range s {
    switch {
      case char == '\'' || char == '"':
        if inQuotes  && char == quoteType {
          inQuotes = false
          fields = append(fields, string(curField))
          quoteType = 0 
        } else if !inQuotes {
          inQuotes = true
          quoteType = char
        } else {
          curField = append(curField, char)
        }
      default:
        if inQuotes || char != ' ' {
          curField = append(curField, char)
      }
    }
  }
  if len(curField) > 0 {
    fields = append(fields, string(curField))
  }
  return fields
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
  if strings.HasPrefix(args, singleQuotes) && strings.HasSuffix(args, singleQuotes) {
    fmt.Printf("%s\n", trimCoupledQuotes(args))
  } else if strings.HasPrefix(args, doubleQuotes) && strings.HasSuffix(args, doubleQuotes) {
    fmt.Printf("%s\n", trimCoupledQuotes(args))
  } else {
  fmt.Printf("%s\n", strings.Join(strings.Fields(args), " "))
  }
}

func trimCoupledQuotes(s string) string {
  if strings.HasSuffix(s, singleQuotes) && strings.HasPrefix(s, singleQuotes) {
    s = strings.TrimSuffix(s, singleQuotes)
    s = strings.TrimPrefix(s, singleQuotes)
  }
  if strings.HasSuffix(s, doubleQuotes) && strings.HasPrefix(s, doubleQuotes) {
    s = strings.TrimSuffix(s, doubleQuotes)
    s = strings.TrimPrefix(s, doubleQuotes)
  }

  // fmt.Println("we callin here?")
  return s
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

func pwdCommand(_ string) {
  dir, err := os.Getwd()
  if err != nil {
    fmt.Printf("Error getting current directory: %s\n", err)
    return
  }
  fmt.Printf("%s\n", dir)
}

func homeCommand(_ string) {
  homeDir, err := os.UserHomeDir()
  if err != nil {
    fmt.Printf("Error retrieving home dir : %s", err)
    return
  }
  fmt.Printf("%s\n", homeDir)
}

func cdCommand(args string) {
  // abs path
  var cmd string
  if args == "~" {
    cmd, _ = os.UserHomeDir()
  } else {
    cmd = args
  }
  if err := os.Chdir(cmd); err != nil {
  fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", args)
  }
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
  cmd, args, err := strings.Cut(input, " ")
  if !err {
    fmt.Printf("%s: invalid input. Seperation is illegal\n", input)
    return
  }
  checkCommand(cmd, args)
}

