package utils

import "os"

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

func MustFile(file *os.File, err error) *os.File {
	if err != nil {
		panic(err)
	}
	return file
}
