// +build run

package main

import (
	"github.com/zetamatta/go-windows-netresource"
)

func main() {
	err := netresource.WNetAddConnection2(`\\localhost\C$`, "O:", "", "")
	if err != nil {
		println(err.Error())
	}
	err = netresource.WNetCancelConnection2("O:", false, false)
	if err != nil {
		println(err.Error())
	}
}
