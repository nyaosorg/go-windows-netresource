package netresource

import (
	"os"
	"testing"
)

func TestNetresourceEnumFileServer(t *testing.T) {
	computerName := `\\` + os.Getenv("COMPUTERNAME")
	found := false

	err := EnumFileServer(func(node *NetResource) bool {
		remoteName := node.RemoteName()
		// println("DEBUG: remoteName==", remoteName)
		if remoteName == computerName {
			found = true
			return false
		}
		return true
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if !found {
		t.Fatalf("self computer name '%s' not found.", computerName)
	}
}
