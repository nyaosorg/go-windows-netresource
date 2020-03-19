package main

import (
	"fmt"
	"time"

	"github.com/zetamatta/go-windows-netresource"
)

func main() {
	indent := 1
	var callback func(*netresource.NetResource) bool

	callback = func(node *netresource.NetResource) bool {
		name := node.RemoteName()
		now := time.Now()
		fmt.Printf("%02d:%02d:%02d %*s%s Scope=%X Type=%X DisplayType=%X Usage=%X LocalName=\"%s\" Comment=\"%s\" Provide=\"%s\"\n",
			now.Hour(),
			now.Minute(),
			now.Second(),
			indent*2,
			"",
			name,
			node.Scope,
			node.Type,
			node.DisplayType,
			node.Usage,
			node.LocalName(),
			node.Comment(),
			node.Provider())
		if len(name) <= 0 || name[0] != '\\' {
			indent++
			node.Enum(callback)
			indent--
		}
		return true
	}
	err := netresource.WNetEnum(callback)
	if err != nil {
		println(err.Error())
	}
}
