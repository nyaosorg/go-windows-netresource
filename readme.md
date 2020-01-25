go-windows-netresource
======================

net use
-------

```go
// +build run

package main

import (
    "os"
    "os/exec"

    "github.com/zetamatta/go-windows-netresource"
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
```

Find file servers and their shared folders
-------------------------------------------


```
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
```
