//go:build windows
// +build windows

package fileinfo

import (
	"fmt"
	"log"
	"runtime"

	"golang.org/x/sys/windows"
)

func IsHidden(filename string) bool {
	pointer, err := windows.UTF16PtrFromString(filename)
	if err != nil {
		panic(fmt.Sprint(filename, err))
	}
	attributes, err := windows.GetFileAttributes(pointer)
	if err != nil {
		panic(fmt.Sprint(filename, err))
	}
	return attributes&windows.FILE_ATTRIBUTE_HIDDEN != 0
}
