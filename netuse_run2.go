//go:build run
// +build run

package main

import (
	"os"
	"os/exec"

	"github.com/nyaosorg/go-windows-netresource"
)

func main() {
	cancel, err := netresource.NetUse(`X:`, `\\localhost\C$`)
	if err != nil {
		println(err.Error())
		return
	}
	defer cancel(true, false)

	cmd := exec.Command("cmd.exe", "/c", "dir", `X:\`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}
