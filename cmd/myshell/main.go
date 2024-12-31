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
var backslash = `\`
var dollar = `$`
var newline = `\n`
var space = ` `
var tilde = `~`

var commands map[string]func([]string)

func init() {
  commands = make(map[string]func([]string)) 
  commands["exit"] = exitCommand
  commands["echo"] = echoCommand
  commands["type"] = typeCommand
  commands["pwd"] = pwdCommand
  commands["cd"] = cdCommand
  commands["~"] = homeCommand
}

func checkCommand(command string, args []string) {
  if cmd, ok := commands[command]; ok {
    cmd(args) // execute the command
    return
  } else {

    // check if the command is a system command 

    _, err := exec.LookPath(command)
    if err != nil {
      fmt.Printf("%s: command not found\n", command)
      return
    }

    for _, arg := range args {
      // unicode print
      // fmt.Printf("arg: %s\n", arg)
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

func trimFieldByQuotes(s []string) []string {
    a := []string{}
    sb := &strings.Builder{}
    quoted := false
    var quoteChar string
    // escaping := false
    for i := 0; i < len(s); i++ {
        // Handle quotes
        r := s[i]
        if r == string(singleQuotes[0]) || r == string(doubleQuotes[0]) {
            if quoted && r == string(quoteChar) {
                // End of quoted field
                quoted = false
                quoteChar = ""
            } else if !quoted {
                // Start of quoted field
                quoted = true
                quoteChar = r
            } else {
                  sb.WriteString(r)
                }
                
            continue
        }
        if r == string(backslash[0]) {
          if quoted && quoteChar == string(doubleQuotes[0]) {
            if i+1 < len(s) {
              nextChar := s[i+1]
              // fmt.Printf("debug 1\n")
              if nextChar == string(dollar[0]) || nextChar == string(backslash[0]) || nextChar == string(doubleQuotes[0]) || nextChar == string(newline[0]) {
                i++
          
                sb.WriteString(nextChar)
                // fmt.Printf("debug 2\n")
                continue
              }
            }
          } else if !quoted {
          // Skip the backslash outside of quotes
            if i+1 < len(s) {
              i++
              sb.WriteString(s[i])
              // fmt.Printf("debug 3\n")
            }
            continue
          }
        } 

        // Handle spaces outside quoted fields
        if !quoted && r == string(space[0]) {
            if sb.Len() > 0 {
                a = append(a, sb.String())
                sb.Reset()
            }
              
            // fmt.Printf("debug 4 %v\n", string(r))
            continue
        }
        if !quoted && r == string(backslash[0]) {
            // Skip the backslash outside of quotes
            if i+1 < len(s) {
              i++
              sb.WriteString(r)
            }
            continue
        }
      sb.WriteString(r)
    }

    // Add the last field if there's any
    if sb.Len() > 0 {
        a = append(a, sb.String())
    }

    return a
}




func exitCommand(args []string) {
  number, err := strconv.Atoi(args[0])
  if err != nil {
    fmt.Printf("%v - %s: not a valid number\n", number, err)
    return
  }
  os.Exit(number)
}

func echoCommand(args []string) {
  // fmt.Printf("%s\n", strings.Join(strings.Fields(args), " "))
  fmt.Printf("args: %s\n", trimFieldByQuotes(args))
  for _, arg := range trimFieldByQuotes(args) {
    fmt.Printf("%s ", arg)
  }
  fmt.Println()
}


func typeCommand(args []string) {
  if _, ok := commands[args[0]]; ok {
      fmt.Printf("%s is a shell builtin\n", args)
      return
  } else {

    paths := os.Getenv("PATH")
    pathList := strings.Split(paths, ":")
    for _, path := range pathList {
      if _, err := os.Stat(filepath.Join(path, args[0])); err == nil {
        fmt.Printf("%s is %s/%s\n", args, path, args)
        return
      }
    }

  }
  fmt.Printf("%s: not found\n", args)
}

func pwdCommand(_ []string) {
  dir, err := os.Getwd()
  if err != nil {
    fmt.Printf("Error getting current directory: %s\n", err)
    return
  }
  fmt.Printf("%s\n", dir)
}

func homeCommand(_ []string) {
  homeDir, err := os.UserHomeDir()
  if err != nil {
    fmt.Printf("Error retrieving home dir : %s", err)
    return
  }
  fmt.Printf("%s\n", homeDir)
}

func cdCommand(args []string) {
  // abs path
  var cmd string
  if args[0] == tilde {
    cmd, _ = os.UserHomeDir()
  } else {
    cmd = args[0]
  }
  if err := os.Chdir(cmd); err != nil {
  fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", args)
  }
}

func main() {
  for {
    fmt.Fprint(os.Stdout, "$ ",)

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
  
  parsedInput := trimFieldByQuotes(strings.Fields(input))
  // fmt.Printf("parsedInput: %s\n", parsedInput)
  cmd := parsedInput[0]
  args := parsedInput[1:]
  checkCommand(cmd, args)
  // test
}

