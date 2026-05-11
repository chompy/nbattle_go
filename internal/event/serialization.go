package event

import (
	"encoding/binary"
	"io"
)

func serialize(data ...any) ([]byte, error) {
	var out []byte
	var err error
	for _, v := range data {
		switch v := v.(type) {
		case int, float32, float64, int32, uint32, int16, uint16, int8, uint8, uint64, Type:
			var toAppend any = v
			if _, ok := v.(int); ok {
				toAppend = int64(v.(int))
			}
			out, err = binary.Append(out, binary.BigEndian, toAppend)
			if err != nil {
				return nil, err
			}
		case bool:
			byteValue := byte(1)
			if !v {
				byteValue = 0
			}
			out = append(out, byteValue)
		case string:
			strLen := uint16(len(v))
			out, err = binary.Append(out, binary.BigEndian, strLen)
			out = append(out, []byte(v)...)
		default:
			return nil, ErrUnserializableValueType
		}

	}
	return out, nil
}

type deserializer struct {
	data []byte
	pos  int
}

func newDeserializer(data []byte) *deserializer {
	return &deserializer{data, 0}
}

func (d *deserializer) ReadBytes(length int) ([]byte, error) {
	d.pos += length
	if d.pos > len(d.data) {
		return nil, io.EOF
	}
	return d.data[d.pos-length : d.pos], nil
}

func (d *deserializer) ReadByte() (byte, error) {
	data, err := d.ReadBytes(1)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

func (d *deserializer) ReadBool() (bool, error) {
	data, err := d.ReadBytes(1)
	if err != nil {
		return false, err
	}
	return data[0] != 0, nil
}

func (d *deserializer) ReadInt() (int, error) {
	out := int64(0)
	data, err := d.ReadBytes(binary.Size(out))
	if err != nil {
		return 0, err
	}
	if _, err := binary.Decode(data, binary.BigEndian, &out); err != nil {
		return 0, err
	}
	return int(out), nil
}

func (d *deserializer) ReadUint64() (uint64, error) {
	out := uint64(0)
	data, err := d.ReadBytes(binary.Size(out))
	if err != nil {
		return 0, err
	}
	if _, err := binary.Decode(data, binary.BigEndian, &out); err != nil {
		return 0, err
	}
	return out, nil
}
