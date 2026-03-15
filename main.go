package main

import (
	"os"

	"github.com/maxzhang/istoria/cmd"
)

func main() {
	// Fix macOS PATH for gopsutil (ioreg, etc.)
	os.Setenv("PATH", "/usr/sbin:/usr/bin:"+os.Getenv("PATH"))

	cmd.Execute()
}
