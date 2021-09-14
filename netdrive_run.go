//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/nyaosorg/go-windows-netresource"
)

func main() {
	netdrive, err := netresource.GetNetDrives()
	if err != nil {
		panic(err.Error())
	}
	for _, d := range netdrive {
		fmt.Printf("net use %c: \"%s\"\n", d.Letter, d.Remote)
	}

	d, err := netresource.FindVacantDrive()
	if err != nil {
		println(err)
		return
	}
	fmt.Printf("last drive=%c\n", d)
}
