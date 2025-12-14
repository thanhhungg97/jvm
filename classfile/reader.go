package classfile

import (
	"encoding/binary"
	"io"
)

// ClassReader wraps a byte slice for reading class file data
type ClassReader struct {
	data []byte
	pos  int
}

// NewClassReader creates a new ClassReader
func NewClassReader(data []byte) *ClassReader {
	return &ClassReader{data: data, pos: 0}
}

// ReadU1 reads an unsigned 8-bit integer
func (r *ClassReader) ReadU1() uint8 {
	val := r.data[r.pos]
	r.pos++
	return val
}

// ReadU2 reads an unsigned 16-bit integer (big-endian)
func (r *ClassReader) ReadU2() uint16 {
	val := binary.BigEndian.Uint16(r.data[r.pos:])
	r.pos += 2
	return val
}

// ReadU4 reads an unsigned 32-bit integer (big-endian)
func (r *ClassReader) ReadU4() uint32 {
	val := binary.BigEndian.Uint32(r.data[r.pos:])
	r.pos += 4
	return val
}

// ReadBytes reads n bytes
func (r *ClassReader) ReadBytes(n int) []byte {
	bytes := r.data[r.pos : r.pos+n]
	r.pos += n
	return bytes
}

// ReadU2s reads a slice of uint16 values
func (r *ClassReader) ReadU2s() []uint16 {
	count := r.ReadU2()
	result := make([]uint16, count)
	for i := range result {
		result[i] = r.ReadU2()
	}
	return result
}

// EOF returns true if end of data is reached
func (r *ClassReader) EOF() bool {
	return r.pos >= len(r.data)
}

// Position returns current position
func (r *ClassReader) Position() int {
	return r.pos
}

// Ensure ClassReader can be used as io.Reader
var _ io.Reader = (*ClassReader)(nil)

func (r *ClassReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
