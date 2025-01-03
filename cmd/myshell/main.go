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
var backslash = `\\`
var dollar = `$`
var newline = `\n`
var space = ` `
var tilde = `~`
var redirect = `>`
var redirectOne = `1>`
var redirectTwo = `2>`
// 
var commands map[string]func(string, string)

func init() {
  commands = make(map[string]func(string, string)) 
  commands["exit"] = exitCommand
  commands["echo"] = echoCommand
  commands["type"] = typeCommand
  commands["pwd"] = pwdCommand
  commands["cd"] = cdCommand
  commands["~"] = homeCommand
  // commands[">"] = redirectCommand
}

func checkCommand(command string, args []string) {
  // fmt.Printf("command: %s\n", command)
  // fmt.Printf("args: %s\n", args)

    var prev string
    // check if the command is a system command 
    for index, arg := range args {
      if arg == redirect || arg == redirectOne || arg == redirectTwo {
        // fmt.Printf("%s - %s: we want to redirect here\n", arg, args)
        if index+1 < len(args) {
          prev = args[index+1] 
          // fmt.Printf("arg: %s\n", arg)
          // fmt.Printf("prev: %s\n", prev)
          // fmt.Printf("next: %s\n", next)
          args = args[:index]
          break
        } else {
          fmt.Printf("Error, no file to redirect to\n")
          return
        }
      }
    }
    if cmd, ok := commands[command]; ok {
      cmd(strings.Join(args, " "), prev) // execute the command
      return
    } else {

      var cmd *exec.Cmd
    // fmt.Printf("command: %s\n", command)
      _, err := exec.LookPath(command)
      if err != nil {
        if prev != "" {
          var file *os.File
          file, err = os.Create(prev)
          if err != nil {
            fmt.Printf("Error creating file: %s\n", err)
            return
          }
          defer file.Close()
          fmt.Fprintf(file, "%s: command not found\n", command)
        } else {
          fmt.Printf("%s: command not found\n", command)
          return
        }
      } else {
        cmd = exec.Command(command, args...)
      }
      if prev != "" {
        var file *os.File
        file, err = os.Create(prev)
        if err != nil {
          fmt.Printf("Error creating file: %s\n", err)
          return
        }
        defer file.Close()
        cmd.Stdout = file
      } else {
        cmd.Stdout = os.Stdout
      }
      cmd.Stderr = os.Stderr
      err = cmd.Run()
        if err != nil {
          // fmt.Printf("%s: command not found (2)\n", command)
          return
        }
    }
}




func exitCommand(args string, redirect string) {
  if len(args) == 0 {
    os.Exit(0)
  }
  number, err := strconv.Atoi(args)
  if err != nil {
    fmt.Printf("%v - %s: not a valid number\n", number, err)
    return
  }
  os.Exit(number)
}

func echoCommand(args string, redirectBool string) {
  // fmt.Printf("%s\n", strings.Join(strings.Fields(args), " "))
  //  fmt.Printf("args: %s\n", trimFieldByQuotes(args))
  //  fmt.Printf("args: %s\n", args)
  // test
  if redirectBool != "" {
    file, err := os.Create(redirectBool)
      if err != nil {
          fmt.Printf("Error creating file: %v\n", err)
          return
      }
      defer file.Close()

      // Write the args to the file
      _, err = file.WriteString(args + "\n")
      if err != nil {
          fmt.Printf("Error writing to file: %v\n", err)
      }
      return
    }
  for _, arg := range args {
    fmt.Printf("%s", string(arg))
  }
  fmt.Println()
}


func typeCommand(args string, redirect string) {
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

func pwdCommand(_ string, redirect string) {
  dir, err := os.Getwd()
  if err != nil {
    fmt.Printf("Error getting current directory: %s\n", err)
    return
  }
  fmt.Printf("%s\n", dir)
}

func homeCommand(_ string, redirect string) {
  homeDir, err := os.UserHomeDir()
  if err != nil {
    fmt.Printf("Error retrieving home dir : %s\n", err)
    return
  }
  fmt.Printf("%s\n", homeDir)
}

func cdCommand(args string, redirect string) {
  // abs path
  var cmd string
  if args == tilde {
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
  
  cmd, args := trimFieldByQuotes(input)[0], trimFieldByQuotes(input)[1:]
  for _, arg := range args {
    if arg == redirect || arg == redirectOne {
      // fmt.Printf("%s - %s: we want to redirect here\n", arg, args)
      checkCommand(cmd, args)
      return
    }
  }
  
  checkCommand(cmd, args)
  // test
}

func trimFieldByQuotes(s string) []string {
    a := []string{}
    sb := &strings.Builder{}
    quoted := false
    var quoteChar byte 
    // escaping := false
    for i := 0; i < len(s); i++ {
        // Handle quotes
        r := s[i]
        if r == singleQuotes[0] || r == doubleQuotes[0] {
            if quoted && r == quoteChar {
                // End of quoted field
                quoted = false
                quoteChar = 0
            } else if !quoted {
                // Start of quoted field
                quoted = true
                quoteChar = r
            } else {
                  sb.WriteByte(r)
                }
                
            continue
        }
        if r == backslash[0] {
          if quoted && quoteChar == doubleQuotes[0] {
            if i+1 < len(s) {
              nextChar := s[i+1]
              // fmt.Printf("debug 1\n")
              if nextChar == dollar[0] || nextChar == backslash[0] || nextChar == doubleQuotes[0] || nextChar == newline[0] {
                i++
          
                sb.WriteByte(nextChar)
                // fmt.Printf("debug 2\n")
                continue
              }
            }
          } else if !quoted {
          // Skip the backslash outside of quotes
            if i+1 < len(s) {
              i++
              sb.WriteByte(s[i])
              // fmt.Printf("debug 3\n")
            }
            continue
          }
        } 

        // Handle spaces outside quoted fields
        if !quoted && r == space[0] {
            if sb.Len() > 0 {
                a = append(a, sb.String())
                sb.Reset()
            }
              
            // fmt.Printf("debug 4 %v\n", string(r))
            continue
        }
        if !quoted && r == backslash[0] {
            // Skip the backslash outside of quotes
            if i+1 < len(s) {
              i++
              sb.WriteByte(r)
            }
            continue
        }
      sb.WriteByte(r)
    }

    // Add the last field if there's any
    if sb.Len() > 0 {
        a = append(a, sb.String())
    }

    return a
}
