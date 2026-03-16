package cs2inspect

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	wireVarintR = 0
	wire64Bit   = 1
	wireLenR    = 2
	wire32BitR  = 5

	maxFields = 100
)

// ProtoField holds the raw value of a decoded protobuf field.
type ProtoField struct {
	FieldNum int
	WireType int
	// For varint fields: the decoded uint64 value.
	Varint uint64
	// For length-delimited and 32/64-bit fields: the raw bytes.
	Bytes []byte
}

// Float32LE interprets the field's 4 bytes as a little-endian float32.
func (f *ProtoField) Float32LE() float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(f.Bytes))
}

// ProtoReader is a hand-written protobuf binary reader.
// Supports:
//   - Wire type 0: varint (uint32, uint64, int32)
//   - Wire type 2: length-delimited (string, bytes, nested messages)
//   - Wire type 5: 32-bit fixed (float32)
type ProtoReader struct {
	data []byte
	pos  int
}

// NewProtoReader creates a new ProtoReader for the given data.
func NewProtoReader(data []byte) *ProtoReader {
	return &ProtoReader{data: data}
}

func (r *ProtoReader) remaining() int {
	return len(r.data) - r.pos
}

func (r *ProtoReader) readByte() (byte, error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("unexpected end of protobuf data")
	}
	b := r.data[r.pos]
	r.pos++
	return b, nil
}

func (r *ProtoReader) readBytes(n int) ([]byte, error) {
	if r.pos+n > len(r.data) {
		return nil, fmt.Errorf("need %d bytes but only %d remain", n, len(r.data)-r.pos)
	}
	chunk := make([]byte, n)
	copy(chunk, r.data[r.pos:r.pos+n])
	r.pos += n
	return chunk, nil
}

// readVarint reads a base-128 varint and returns it as uint64.
func (r *ProtoReader) readVarint() (uint64, error) {
	var result uint64
	var shift uint

	for {
		b, err := r.readByte()
		if err != nil {
			return 0, err
		}
		result |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
		if shift > 63 {
			return 0, fmt.Errorf("varint too long")
		}
	}

	return result, nil
}

func (r *ProtoReader) readTag() (fieldNum int, wireType int, err error) {
	tag, err := r.readVarint()
	if err != nil {
		return 0, 0, err
	}
	return int(tag >> 3), int(tag & 7), nil
}

func (r *ProtoReader) readLengthDelimited() ([]byte, error) {
	length, err := r.readVarint()
	if err != nil {
		return nil, err
	}
	return r.readBytes(int(length))
}

// ReadAllFields reads all fields from the data and returns them.
// Returns an error if the field count exceeds maxFields or an unknown wire type is encountered.
func (r *ProtoReader) ReadAllFields() ([]ProtoField, error) {
	var fields []ProtoField
	fieldCount := 0

	for r.remaining() > 0 {
		fieldCount++
		if fieldCount > maxFields {
			return nil, fmt.Errorf("protobuf field count exceeds limit of %d", maxFields)
		}

		fieldNum, wireType, err := r.readTag()
		if err != nil {
			return nil, err
		}

		f := ProtoField{FieldNum: fieldNum, WireType: wireType}

		switch wireType {
		case wireVarintR:
			v, err := r.readVarint()
			if err != nil {
				return nil, err
			}
			f.Varint = v

		case wire64Bit:
			b, err := r.readBytes(8)
			if err != nil {
				return nil, err
			}
			f.Bytes = b

		case wireLenR:
			b, err := r.readLengthDelimited()
			if err != nil {
				return nil, err
			}
			f.Bytes = b

		case wire32BitR:
			b, err := r.readBytes(4)
			if err != nil {
				return nil, err
			}
			f.Bytes = b

		default:
			return nil, fmt.Errorf("unknown wire type %d for field %d", wireType, fieldNum)
		}

		fields = append(fields, f)
	}

	return fields, nil
}
