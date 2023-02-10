//go:build unix
// +build unix

package fileinfo

import "strings"

func IsHidden(filename string) bool {
	return strings.HasPrefix(filename, ".")
}
