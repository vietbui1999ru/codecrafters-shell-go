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

    // check if the command is a system command 

    _, err := exec.LookPath(command)
    if err != nil {
      fmt.Printf("%s: command not found\n", command)
      return
    }

    for _, arg := range trimFieldByQuotes(args) {
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

func trimFieldByQuotes(s string) []string {
    // s := `Foo bar random "letters lol" stuff`
    a := []string{}
    sb := &strings.Builder{}
    quoted := false
    var quoteChar rune
    for i, r := range s {
        if r == rune(singleQuotes[0]) || r == rune(doubleQuotes[0]) {

          // if field is quoted(true) and the quoteChar is the same as the current char then
          if quoted && r == quoteChar {             
            quoted = !quoted // end of quote -> false
            quoteChar = 0 // reset quoteChar

          // if field is not quoted and the quoteChar is not the same as the current char then
          } else if !quoted {
            quoted = true // start of quote = true
            quoteChar = r // set quoteChar to be the single/double quote

          // if field is quoted and the quoteChar is not the same as the current char then  
          } else {
            sb.WriteRune(r) //mismatch quote inside quote field, treat as normal char
          }
            // sb.WriteRune(r) // keep '"' otherwise comment this line

        // if field is not quoted and the current char is a space
        } else if !quoted && r == ' ' {

          // if the string builder has a length greater than 0
          if sb.Len() > 0 {
            a = append(a, sb.String())
            sb.Reset()
          }

       // if field is not quoted and the current char is not a space
       } else {
          // if is not quoted and the current char is a backslash then add the next char to the string builder 
          if !quoted && r == rune(backslash[0]) {
            if i+1 < len(s) {
              // sb.WriteRune(r)
              // write the next char to the string builder, ignore the backslash
              sb.WriteRune(rune(s[i+1]))
            }
            continue
          } else {

          sb.WriteRune(r)
        }
          // fmt.Printf("r: %v\n", r)
       }
        // sb.WriteRune(r)
      }
    if sb.Len() > 0 {
        a = append(a, sb.String())
    }

    return a
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
  // fmt.Printf("%s\n", strings.Join(strings.Fields(args), " "))
  for _, arg := range trimFieldByQuotes(args) {
    fmt.Printf("%s ", arg)
  }
  fmt.Println()
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
  cmd, args, _ := strings.Cut(input, " ")
  checkCommand(cmd, args)
}

