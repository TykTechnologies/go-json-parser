package ast

type ByteSliceReference struct {
	Start uint32
	End   uint32
}

func (b ByteSliceReference) Length() uint32 {
	return b.End - b.Start
}

type Input struct {
	// RawBytes is the raw byte input
	RawBytes []byte
	// Length of RawBytes
	Length int
	// InputPosition is the current position in the RawBytes
	InputPosition int
}

// Reset empties the Input
func (i *Input) Reset() {
	i.RawBytes = i.RawBytes[:0]
	i.InputPosition = 0
}

// AppendInputBytes appends a byte slice to the current input and returns the ByteSliceReference
func (i *Input) AppendInputBytes(bytes []byte) (ref ByteSliceReference) {
	ref.Start = uint32(len(i.RawBytes))
	i.RawBytes = append(i.RawBytes, bytes...)
	i.Length = len(i.RawBytes)
	ref.End = uint32(len(i.RawBytes))
	return
}

// ResetInputBytes empties the input and sets it to bytes argument
func (i *Input) ResetInputBytes(bytes []byte) {
	i.Reset()
	i.AppendInputBytes(bytes)
	i.Length = len(i.RawBytes)
}

// ByteSlice returns the byte slice for a given byte ByteSliceReference
func (i *Input) ByteSlice(reference ByteSliceReference) ByteSlice {
	return i.RawBytes[reference.Start:reference.End]
}

// ByteSlice is an alias for []byte
type ByteSlice []byte
