package cs2inspect

import (
	"encoding/binary"
	"math"
)

const (
	wireVarint = 0
	wireLen    = 2
	wire32Bit  = 5
)

// ProtoWriter is a hand-written protobuf binary writer.
// Fields with default/zero values are omitted (proto3 semantics).
type ProtoWriter struct {
	buf []byte
}

// Bytes returns the encoded bytes.
func (w *ProtoWriter) Bytes() []byte {
	return w.buf
}

// writeVarint writes a base-128 varint.
func (w *ProtoWriter) writeVarint(v uint64) {
	for {
		b := byte(v & 0x7F)
		v >>= 7
		if v != 0 {
			b |= 0x80
		}
		w.buf = append(w.buf, b)
		if v == 0 {
			break
		}
	}
}

func (w *ProtoWriter) writeTag(fieldNum int, wireType int) {
	w.writeVarint(uint64((fieldNum << 3) | wireType))
}

// WriteUint32 writes a uint32 field. Zero values are omitted.
func (w *ProtoWriter) WriteUint32(fieldNum int, value uint32) {
	if value == 0 {
		return
	}
	w.writeTag(fieldNum, wireVarint)
	w.writeVarint(uint64(value))
}

// WriteUint64 writes a uint64 field. Zero values are omitted.
func (w *ProtoWriter) WriteUint64(fieldNum int, value uint64) {
	if value == 0 {
		return
	}
	w.writeTag(fieldNum, wireVarint)
	w.writeVarint(value)
}

// WriteInt32 writes an int32 field. Zero values are omitted.
// Negative values are treated as unsigned 64-bit two's complement (proto3).
func (w *ProtoWriter) WriteInt32(fieldNum int, value int32) {
	if value == 0 {
		return
	}
	w.writeTag(fieldNum, wireVarint)
	// Negative int32 is sign-extended to 64 bits and written as a 10-byte varint
	w.writeVarint(uint64(int64(value)))
}

// WriteString writes a string field. Empty strings are omitted.
func (w *ProtoWriter) WriteString(fieldNum int, value string) {
	if value == "" {
		return
	}
	encoded := []byte(value)
	w.writeTag(fieldNum, wireLen)
	w.writeVarint(uint64(len(encoded)))
	w.buf = append(w.buf, encoded...)
}

// WriteFloat32Fixed writes a float32 as wire type 5 (fixed 32-bit, little-endian).
// Used for sticker float fields (wear, scale, rotation, etc.).
func (w *ProtoWriter) WriteFloat32Fixed(fieldNum int, value float32) {
	w.writeTag(fieldNum, wire32Bit)
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], math.Float32bits(value))
	w.buf = append(w.buf, b[:]...)
}

// WriteRawBytes writes raw bytes as a length-delimited field (wire type 2).
// Empty slices are omitted.
func (w *ProtoWriter) WriteRawBytes(fieldNum int, data []byte) {
	if len(data) == 0 {
		return
	}
	w.writeTag(fieldNum, wireLen)
	w.writeVarint(uint64(len(data)))
	w.buf = append(w.buf, data...)
}

// WriteEmbedded writes a nested message as a length-delimited field.
func (w *ProtoWriter) WriteEmbedded(fieldNum int, nested *ProtoWriter) {
	w.WriteRawBytes(fieldNum, nested.Bytes())
}
