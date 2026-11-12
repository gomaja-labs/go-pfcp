package ie

import (
	"encoding/binary"
	"io"
)

// SDFMultiplexedMedia contains one (S)RTP media identification instance.
// Flags bit 0..4 are SSRCPI, PTPI, MIDPI, RHIDMPI and RPTPI respectively.
// TS 29.244 V20.1.0 clause 8.2.5, Figure 8.2.5-2.
type SDFMultiplexedMedia struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// SynchronizationSource is the synchronization source.
	SynchronizationSource uint32
	// PayloadType is the payload type.
	PayloadType uint8
	// MediaIdentificationTag is the media identification tag.
	MediaIdentificationTag string
	// RTPHeaderExtensionID is the RTP header extension ID.
	RTPHeaderExtensionID uint8
	// RTCPPacketType is the RTCP packet type.
	RTCPPacketType uint8
}

func (f *SDFMultiplexedMedia) validate() error {
	if f.Flags&1 == 0 && f.SynchronizationSource != 0 || f.Flags&2 == 0 && f.PayloadType != 0 || f.Flags&4 == 0 && f.MediaIdentificationTag != "" || f.Flags&8 == 0 && f.RTPHeaderExtensionID != 0 || f.Flags&16 == 0 && f.RTCPPacketType != 0 {
		return ErrInvalidValue
	}
	if f.Flags&^uint8(31) != 0 || f.Flags&2 != 0 && f.PayloadType > 127 {
		return ErrInvalidValue
	}
	if len(f.MediaIdentificationTag) > 255 {
		return ErrInvalidLength
	}
	if f.Flags&8 != 0 && f.RTPHeaderExtensionID == 0 {
		return ErrInvalidValue
	}
	if f.Flags&4 != 0 && f.Flags&16 == 0 && f.Flags&8 == 0 {
		return ErrInvalidValue
	}
	if f.Flags&16 != 0 && (f.Flags&2 != 0 || f.Flags&5 == 0) {
		return ErrInvalidValue
	}
	if f.Flags&18 == 0 && f.Flags&5 == 0 {
		return ErrInvalidValue
	}
	return nil
}

// MarshalLen returns the encoded nested-instance length.
func (f *SDFMultiplexedMedia) MarshalLen() int {
	if f == nil {
		return 0
	}
	n := 1
	if f.Flags&1 != 0 {
		n += 4
	}
	if f.Flags&2 != 0 {
		n++
	}
	if f.Flags&4 != 0 {
		n += 1 + len(f.MediaIdentificationTag)
	}
	if f.Flags&8 != 0 {
		n++
	}
	if f.Flags&16 != 0 {
		n++
	}
	return n
}

// Marshal serializes a media identification instance.
func (f *SDFMultiplexedMedia) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	b := make([]byte, f.MarshalLen())
	if e := f.MarshalTo(b); e != nil {
		return nil, e
	}
	return b, nil
}

// MarshalTo writes a length-delimited instance's contents.
func (f *SDFMultiplexedMedia) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if e := f.validate(); e != nil {
		return e
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.Flags
	o := 1
	if f.Flags&1 != 0 {
		binary.BigEndian.PutUint32(b[o:o+4], f.SynchronizationSource)
		o += 4
	}
	if f.Flags&2 != 0 {
		b[o] = f.PayloadType
		o++
	}
	if f.Flags&4 != 0 {
		b[o] = uint8(len(f.MediaIdentificationTag))
		o++
		copy(b[o:], f.MediaIdentificationTag)
		o += len(f.MediaIdentificationTag)
	}
	if f.Flags&8 != 0 {
		b[o] = f.RTPHeaderExtensionID
		o++
	}
	if f.Flags&16 != 0 {
		b[o] = f.RTCPPacketType
	}
	return nil
}

// UnmarshalBinary parses one media identification instance and ignores extensions.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *SDFMultiplexedMedia) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = SDFMultiplexedMedia{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0] & 31
	o := 1
	if f.Flags&1 != 0 {
		if len(b) < o+4 {
			return io.ErrUnexpectedEOF
		}
		f.SynchronizationSource = binary.BigEndian.Uint32(b[o : o+4])
		o += 4
	}
	if f.Flags&2 != 0 {
		if len(b) < o+1 {
			return io.ErrUnexpectedEOF
		}
		f.PayloadType = b[o] & 127
		o++
	}
	if f.Flags&4 != 0 {
		if len(b) < o+1 {
			return io.ErrUnexpectedEOF
		}
		n := int(b[o])
		o++
		if len(b) < o+n {
			return io.ErrUnexpectedEOF
		}
		f.MediaIdentificationTag = string(b[o : o+n])
		o += n
	}
	if f.Flags&8 != 0 {
		if len(b) < o+1 {
			return io.ErrUnexpectedEOF
		}
		f.RTPHeaderExtensionID = b[o]
		o++
	}
	if f.Flags&16 != 0 {
		if len(b) < o+1 {
			return io.ErrUnexpectedEOF
		}
		f.RTCPPacketType = b[o]
	}
	return f.validate()
}
