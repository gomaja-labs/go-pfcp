package ie

import (
	"encoding/binary"
	"io"
)

// takePayload bounds a field against the remaining payload, not its capacity.
// The full slice expression prevents an append from corrupting a following field.
func takePayload(b []byte, offset *int, length int) ([]byte, error) {
	if length < 0 || *offset < 0 || *offset > len(b) || length > len(b)-*offset {
		return nil, io.ErrUnexpectedEOF
	}
	end := *offset + length
	value := b[*offset:end:end]
	*offset = end
	return value, nil
}

// takeLength16 reads the length before checking its value against the payload.
func takeLength16(b []byte, offset *int) ([]byte, error) {
	length, err := takePayload(b, offset, 2)
	if err != nil {
		return nil, err
	}
	return takePayload(b, offset, int(binary.BigEndian.Uint16(length)))
}
