package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	libsass "github.com/wellington/go-libsass"
)

var quitCommandReg = regexp.MustCompile(`^quit`)
var runCommandRegexp = regexp.MustCompile(`^command\|`)
var illegalPathContents = regexp.MustCompile(`[\n\r]|\\\\`)

// command|D:/Projects/esbuild-angular|D:/Projects/esbuild-angular/src/styles.scss
// Compile sass files.
func SassCompiler(workingDir string, scssPath string) string {
  scssPath = illegalPathContents.ReplaceAllString(scssPath, "")
  fileContent, err := os.ReadFile( path.Clean(scssPath) )
  if err != nil {
    panic(err)
  }

  styleContent := string(fileContent)
  buf := bytes.NewBufferString(styleContent)
  var compiled bytes.Buffer

  comp, err := libsass.New(&compiled, buf)
  if err != nil {
    log.Fatal(err)
  }

  includePaths := []string{workingDir}
	optionError := comp.Option(libsass.IncludePaths(includePaths))
  if optionError != nil {
    log.Fatal(optionError)
  }

  if err := comp.Run(); err != nil {
    log.Fatal(err)
  }
  return compiled.String()
}

func IoReader() {
  var reader = bufio.NewReader(os.Stdin)
	message, _ := reader.ReadString('\n')

  if (quitCommandReg.Match([]byte(message))) {
    os.Exit(0)
  } else if runCommandRegexp.Match([]byte(message)) {
    fmt.Println(message)
    args := strings.Split(message, "|")
    css := SassCompiler(args[1], args[2])
    fmt.Println(css)

    reader = nil
    IoReader()
  } else {
    reader = nil
    IoReader()
  }
}

func main() {
  IoReader()
}
