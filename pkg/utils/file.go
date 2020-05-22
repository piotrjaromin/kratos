package utils

import "os"

const StdOut = "stdout"

func File(name string, create bool) (*os.File, error) {
	switch name {
	case StdOut:
		return os.Stdout, nil
	default:
		if create {
			return os.Create(name)
		}
		return os.Open(name)
	}
}
