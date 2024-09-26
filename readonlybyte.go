package gocache

type ReadOnlyByte struct {
	bytes []byte
}

func (r ReadOnlyByte) Len() int {
	return len(r.bytes)
}

func (r *ReadOnlyByte) Slice() []byte {
	return cloneBytes(r.bytes)
}

func (r ReadOnlyByte) String() string {
	return string(r.bytes)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
