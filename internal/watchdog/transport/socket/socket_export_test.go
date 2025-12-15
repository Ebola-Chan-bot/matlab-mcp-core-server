// Copyright 2025 The MathWorks, Inc.

package socket

func NewSocket(
	directory Directory,
) (Socket, error) {
	return newSocket(
		directory,
	)
}
