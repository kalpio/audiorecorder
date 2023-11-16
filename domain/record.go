package domain

import "io"

type Record struct {
	Data      []byte
	readIndex int64
}

func (r *Record) Read(p []byte) (n int, err error) {
	if r.readIndex >= int64(len(r.Data)) {
		return 0, io.EOF
	}

	n = copy(p, r.Data[r.readIndex:])
	r.readIndex += int64(n)
	return n, nil
}

func (r *Record) Write(p []byte) (n int, err error) {
	r.Data = append(r.Data, p...)
	return len(p), nil
}
