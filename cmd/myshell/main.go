package main

import (
	"bufio"
	"fmt"
	"os"
  "io"
  "log"
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
var commands map[string]func(string, string, bool)

func init() {
  commands = make(map[string]func(string, string, bool)) 
  commands["exit"] = exitCommand
  commands["echo"] = echoCommand
  commands["type"] = typeCommand
  commands["pwd"] = pwdCommand
  commands["cd"] = cdCommand
  commands["~"] = homeCommand
  // commands[">"] = redirectCommand
}

func checkCommand(command string, args []string) {
	var redirectFile string
	var isStderr bool
  isStderr = false
  originalStdout := os.Stdout
  originalStderr := os.Stderr

	// Handle redirection operators ">" and "2>"
	for index, arg := range args {
		if arg == redirect || arg == redirectOne || arg == redirectTwo {
			if index+1 < len(args) {
				redirectFile = args[index+1]
				if arg == redirectTwo {
					isStderr = true
				}
				args = args[:index]
				break
			} else {
				fmt.Fprintln(os.Stderr, "Error: No file specified for redirection")
				return
			}
		}
	}

  // Check if the command is a built-in checkCommand(
  if cmd, ok := commands[command]; ok {
    cmd(strings.Join(args, " "), redirectFile, isStderr)
    return
  }

	var file *os.File
	if redirectFile != "" {
		var err error
		file, err = os.Create(redirectFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating file:", err)
			return
		}
    defer file.Close()

		if isStderr {
			// Redirect stderr to both terminal and file
			os.Stderr = file 
		} else {
			// Redirect stdout to both terminal and file
			os.Stdout = file
		}
	}
	cmd := exec.Command(command, args...)
  cmd.Env = os.Environ()
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
	_ = cmd.Run()
  os.Stdout = originalStdout
  os.Stderr = originalStderr
}

func exitCommand(args string, redirect string, _ bool) {
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


func echoCommand(args string, redirectFile string, isStderr bool) {
     var output io.Writer
     if redirectFile != "" {
         // Open or create the file for writing
         file, err := os.Create(redirectFile)
         if err != nil {
             log.Fatalf("Error creating file: %v\n", err)
         }
         defer file.Close()

         if isStderr {
             // Write to stderr and to the file
            os.Stderr = file
             
         } else {
         // Write to the file
            os.Stdout = file
         }
     } else {
       output = os.Stdout
     }

    // Default behavior: Print to stdout
    fmt.Fprintln(output, args)
}



func typeCommand(args string, redirect string, _ bool) {
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

func pwdCommand(_ string, _ string, _ bool) {
  dir, err := os.Getwd()
  fmt.Printf("dir: %s\n", dir)
  if err != nil {
    fmt.Printf("Error getting current directory: %s\n", err)
    return
  }
  fmt.Printf("%s\n", dir)
}

func homeCommand(_ string, redirect string, _ bool) {
  homeDir, err := os.UserHomeDir()
  if err != nil {
    fmt.Printf("Error retrieving home dir : %s\n", err)
    return
  }
  fmt.Printf("%s\n", homeDir)
}

func cdCommand(args string, redirect string, _ bool) {
  // abs path
  var cmd string
  if args == tilde {
    cmd, _ = os.UserHomeDir()
    fmt.Printf("cmd (2): %s\n", cmd)
  } else {
    cmd = args
    fmt.Printf("cmd (3): %s\n", cmd)
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
    if err != nil {
      fmt.Printf("%s: invalid input\n", input)
    }
    trimmedInput := strings.TrimSpace(input)
    if len(trimmedInput) == 0 {
			continue
		}   // fmt.Printf("%s: command not found\n", trimmpedInput)
    

    handleCommands(trimmedInput)
  }
}

// comment again
func handleCommands(input string) {
  
  cmd, args := trimFieldByQuotes(input)[0], trimFieldByQuotes(input)[1:]
  // fmt.Printf("cmd (1): %s\n", args)
  // for _, arg := range args {
  //   if arg == redirect || arg == redirectOne || arg == redirectTwo {
  //     // fmt.Printf("%s - %s: we want to redirect here\n", arg, args)
  //     checkCommand(cmd, args)
  //     return
  //   }
  // }
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
