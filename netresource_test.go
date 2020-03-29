package netresource

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"
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

func getOutputOfNetExeShare() ([]string, error) {
	cmd := exec.Command("net", "share")
	in, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer in.Close()
	cmd.Start()

	sc := bufio.NewScanner(in)
	for sc.Scan() {
		line := sc.Text()
		if len(line) > 0 && line[0] == '-' {
			nodes := []string{}
			for sc.Scan() {
				line = sc.Text()
				if len(line) > 0 && line[0] == ' ' {
					// skip continued line.
					continue
				}
				field := strings.Fields(line)
				if len(field) > 0 && !strings.ContainsRune(field[0], '$') {
					nodes = append(nodes, field[0])
				}
			}
			if len(nodes) >= 1 {
				nodes = nodes[:len(nodes)-1]
			}
			return nodes, nil
		}
	}
	return nil, errors.New("net share does not output as expected")
}

func TestNewFileServerEnum(t *testing.T) {
	share, err := getOutputOfNetExeShare()
	if err != nil {
		t.Fatal(err.Error())
	}
	computerName := `\\` + os.Getenv("COMPUTERNAME")
	expected := make(map[string]struct{})
	for _, s := range share {
		expected[computerName+`\`+s] = struct{}{}
	}
	fs, err := NewFileServer(os.Getenv("COMPUTERNAME"))
	if err != nil {
		t.Fatal(err.Error())
	}
	err = fs.Enum(func(nr *NetResource) bool {
		rn := nr.RemoteName()
		if _, ok := expected[rn]; ok {
			delete(expected, rn)
			// println(rn,": found")
		} else {
			t.Logf("%s is not found in expected list", rn)
			t.Fail()
		}
		return true
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(expected) > 0 {
		t.Log("not all expected nodes are listuped")
		for s, _ := range expected {
			t.Log(s)
		}
		t.FailNow()
	}
}
