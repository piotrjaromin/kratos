package utils

import "io"

type StringWriterReader struct {
	buffer    []byte
	readIndex int
}

func (sw *StringWriterReader) Write(p []byte) (n int, err error) {
	sw.buffer = append(sw.buffer, p...)
	return len(p), nil
}

func (sw *StringWriterReader) Read(p []byte) (n int, err error) {
	if sw.readIndex >= len(sw.buffer) {
		return 0, io.EOF
	}

	bytesLeft := len(sw.buffer) - sw.readIndex
	bytesToRead := len(p)

	if bytesLeft < bytesToRead {
		bytesToRead = bytesLeft
	}

	for index := 0; index < bytesToRead; index++ {
		p[index] = sw.buffer[sw.readIndex]
		sw.readIndex++
	}

	return bytesToRead, nil
}
