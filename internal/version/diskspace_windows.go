//go:build windows
// +build windows

package version

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"
)

// getAvailableDiskSpace returns the available disk space in bytes for the given path
func (i *VersionInstaller) getAvailableDiskSpace(path string) (int64, error) {
	// Get the volume root path
	volumePath := filepath.VolumeName(path)
	if volumePath == "" {
		volumePath = filepath.Dir(path)
	}
	volumePath += string(filepath.Separator)

	// Windows API call to GetDiskFreeSpaceExW
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")

	var freeBytesAvailable int64
	var totalBytes int64
	var totalFreeBytes int64

	volumePathPtr, err := syscall.UTF16PtrFromString(volumePath)
	if err != nil {
		return 0, fmt.Errorf("failed to convert path to UTF16: %w", err)
	}

	ret, _, err := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(volumePathPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)),
	)

	if ret == 0 {
		return 0, fmt.Errorf("GetDiskFreeSpaceEx failed: %w", err)
	}

	return freeBytesAvailable, nil
}
