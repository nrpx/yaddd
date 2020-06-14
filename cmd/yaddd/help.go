package main

import (
	"io"
	"path"
)

func printHelp(w io.Writer, name string) {
	printVersion(w)

	_, _ = io.WriteString(w, "\nUsage: ")
	_, _ = io.WriteString(w, path.Base(name))
	_, _ = io.WriteString(w, ` COMMAND [OPTION...]
Commands:
  run         start YaDDD service
  version     show version number

Options:
  -conf=FILE  use specified configuration file
  -ip=IP      find A-record with specified IP address
  -debug      enable debug mode
`)
}
