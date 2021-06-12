package util

import (
	"os"
	"strings"
)

// ErrorContains returns NoCaseContains(err.Error(), substr)
// Returns false if err is nil.
func ErrorContains(err error, substr string) bool {
	if err == nil {
		return false
	}
	return noCaseContains(err.Error(), substr)
}

// noCaseContains reports whether substr is within s case-insensitive.
func noCaseContains(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// Panic panics if err != nil
func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func PathIsNotExist(name string) (ok bool) {
	_, err := os.Lstat(name)
	if os.IsNotExist(err) {
		ok = true
		err = nil
	}
	Panic(err)
	return
}

// PathIsExist .
func PathIsExist(name string) bool {
	return !PathIsNotExist(name)
}
