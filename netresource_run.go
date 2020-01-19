// +build run

package main

import (
	"github.com/zetamatta/go-windows-netresource"
)

func main() {
	machines := []string{}

	err := netresource.EnumFileServer(func(node *netresource.NetResource) bool {
		machines = append(machines, node.RemoteName())
		return true
	})
	if err != nil {
		println(err.Error())
	}

	for _, name := range machines {
		println("machine:", name)
		if fs, err := netresource.NewFileServer(name); err == nil {
			fs.Enum(func(node *netresource.NetResource) bool {
				println("  ", node.RemoteName())
				return true
			})
		}
	}

	if err != nil {
		println(err.Error())
	}
}

// https://msdn.microsoft.com/ja-jp/library/cc447030.aspx
// http://eternalwindows.jp/security/share/share06.html
