package netresource

import (
	"errors"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

var mpr = windows.NewLazySystemDLL("mpr.dll")
var procWNetGetConnectionW = mpr.NewProc("WNetGetConnectionW")

func WNetGetConnectionUTF16s(localName []uint16) (string, error) {
	var buffer [1024]uint16
	size := uintptr(len(buffer))
	rc, _, err := procWNetGetConnectionW.Call(
		uintptr(unsafe.Pointer(&localName[0])),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)))

	if uint32(rc) != 0 {
		return "", err
	}
	return windows.UTF16ToString(buffer[:]), nil
}

func WNetGetConnectionUTF16a(drive uint16) (string, error) {
	return WNetGetConnectionUTF16s([]uint16{drive, ':', 0})
}

type NetDrive struct {
	Letter rune
	Remote string
}

func GetNetDrives() ([]*NetDrive, error) {
	drives, err := GetDrives()
	if err != nil {
		return nil, err
	}
	result := []*NetDrive{}
	for _, d := range drives {
		if d.Type == windows.DRIVE_REMOTE {
			path, err := WNetGetConnectionUTF16a(uint16(d.Letter))
			if err == nil {
				node := &NetDrive{Letter: d.Letter, Remote: path}
				result = append(result, node)
			}
		}
	}
	return result, nil
}

// FindVacantDrive returns the last unused drive like 'Z'
func FindVacantDrive() (uint, error) {
	bits, err := windows.GetLogicalDrives()
	if err != nil {
		return 0, err
	}
	for d := uint('Z'); d >= 'A'; d-- {
		if (bits & (1 << (d - 'A'))) == 0 {
			return d, nil
		}
	}
	return 0, errors.New("vacant drive is not found")
}

// NetUse do same thing as `net use X: \\server\path...`
//    drive - `X:`
//    vol - `server\\path\...`
// returns
//    func(update,force) - function release the drive
//        update - true: updates connection as not a persistent one
//        force - true: disconnect even if open process exists.
//    error
func NetUse(drive, vol string) (func(bool, bool), error) {
	if err := WNetAddConnection2(vol, drive, "", ""); err != nil {
		return func(bool, bool) {}, err
	}
	return func(update, force bool) {
		WNetCancelConnection2(drive, update, force)
	}, nil
}

// UNCtoNetDrive replace UNCPath to path using netdrive.
//    uncpath - for example \\server\path\folder\name
// returns
//    newpath - X:\folder\name
func UNCtoNetDrive(uncpath string) (newpath string, closer func(bool, bool)) {
	vol := filepath.VolumeName(uncpath)
	d, err := FindVacantDrive()
	if err != nil {
		return "", func(bool, bool) {}
	}
	netdrive := string([]byte{byte(d), ':'})
	newpath = filepath.Join(netdrive, uncpath[len(vol):])
	if closer, err = NetUse(netdrive, vol); err != nil {
		return "", func(bool, bool) {}
	} else {
		return newpath, closer
	}
}

// https://msdn.microsoft.com/ja-jp/library/cc447030.aspx
// http://eternalwindows.jp/security/share/share06.html
