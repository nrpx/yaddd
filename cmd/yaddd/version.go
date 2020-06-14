package main

import "io"

func printVersion(w io.Writer) {
	_, _ = io.WriteString(w, "Yandex Connect DynDNS Daemon v.")
	_, _ = io.WriteString(w, version)

	_, _ = io.WriteString(w, "\nBuilt on ")
	_, _ = io.WriteString(w, buildTime)
	_, _ = io.WriteString(w, "\n")
}
