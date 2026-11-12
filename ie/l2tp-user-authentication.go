package ie

// L2TPAuthenticationFlags selects the optional authentication values.
type L2TPAuthenticationFlags uint8

const (
	// L2TPAuthenticationPAN selects PAN.
	L2TPAuthenticationPAN L2TPAuthenticationFlags = 1
	// L2TPAuthenticationPAC selects PAC.
	L2TPAuthenticationPAC L2TPAuthenticationFlags = 2
	// L2TPAuthenticationPAR selects PAR.
	L2TPAuthenticationPAR L2TPAuthenticationFlags = 4
	// L2TPAuthenticationPAI selects PAI.
	L2TPAuthenticationPAI L2TPAuthenticationFlags = 8
)

// L2TPUserAuthenticationFields is the structured value of L2TPUserAuthentication.
// TS 29.244 V20.1.0 clause 8.2.187.
// Proxy Authen ID is one octet as Figure 8.2.187-1 specifies: RFC 2661
// section 4.4.5 reserves the high octet of the AVP's two-octet value as zero.
// The IE constructor enforces the type-dependent required flags; the decoder
// reads exactly the received flags without inferring missing fields.
type L2TPUserAuthenticationFields struct {
	// Type selects the authentication method.
	Type uint16
	// Flags contains the presence and option bits.
	Flags L2TPAuthenticationFlags
	// Name is the name.
	Name []byte
	// Challenge is the challenge.
	Challenge []byte
	// Response is the response.
	Response []byte
	// ID is the authentication identifier.
	ID uint8
}

// NewL2TPUserAuthentication creates an IE, retaining invalid arguments as an encoding error.
// TS 29.244 V20.1.0 clause 8.2.187.
func NewL2TPUserAuthentication(fields L2TPUserAuthenticationFields) *IE {
	var required L2TPAuthenticationFlags
	switch fields.Type {
	case 1:
		required = 5
	case 2, 5:
		required = 15
	case 3:
		required = 13
	}
	if fields.Flags&required != required {
		return newInvalidIE(L2TPUserAuthentication, ErrInvalidValue)
	}
	b, err := fields.Marshal()
	return addedFieldsIE(L2TPUserAuthentication, b, err)
}

// L2TPUserAuthentication returns the known fields, ignoring explicitly allowed extensions.
// TS 29.244 V20.1.0 clauses 8.2.187 and 8.1.2.
func (i *IE) L2TPUserAuthentication() (*L2TPUserAuthenticationFields, error) {
	b, e := addedValuePayload(i, L2TPUserAuthentication, 3, false)
	if e != nil {
		return nil, e
	}
	f, e := ParseL2TPUserAuthenticationFields(b)
	if e != nil {
		return nil, &DecodeError{Type: L2TPUserAuthentication, Err: e}
	}
	return f, nil
}

// ParseL2TPUserAuthenticationFields decodes fields according to clause 8.2.187.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseL2TPUserAuthenticationFields(b []byte) (*L2TPUserAuthenticationFields, error) {
	f := &L2TPUserAuthenticationFields{}
	if e := f.UnmarshalBinary(b); e != nil {
		return nil, e
	}
	return f, nil
}

// UnmarshalBinary resets the value and decodes the known fields.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *L2TPUserAuthenticationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = L2TPUserAuthenticationFields{}
	r := &addedReader{b: b}
	f.Type = r.u16()
	f.Flags = L2TPAuthenticationFlags(r.u8() & 15)
	if f.Flags&1 != 0 {
		f.Name = append([]byte(nil), r.take(int(r.u8()))...)
	}
	if f.Flags&2 != 0 {
		f.Challenge = append([]byte(nil), r.take(int(r.u8()))...)
	}
	if f.Flags&4 != 0 {
		f.Response = append([]byte(nil), r.take(int(r.u8()))...)
	}
	if f.Flags&8 != 0 {
		f.ID = r.u8()
	}
	return r.err
}
func (f *L2TPUserAuthenticationFields) valid() bool {
	return f.Flags&^15 == 0 && len(f.Name) <= 255 && len(f.Challenge) <= 255 && len(f.Response) <= 255 && (f.Flags&1 != 0 || len(f.Name) == 0) && (f.Flags&2 != 0 || len(f.Challenge) == 0) && (f.Flags&4 != 0 || len(f.Response) == 0) && (f.Flags&8 != 0 || f.ID == 0)
}

// Marshal encodes known fields with spare bits set to zero.
func (f *L2TPUserAuthenticationFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := addedPut16(nil, f.Type)
	b = append(b, byte(f.Flags))
	if f.Flags&1 != 0 {
		b = append(b, byte(len(f.Name)))
		b = append(b, f.Name...)
	}
	if f.Flags&2 != 0 {
		b = append(b, byte(len(f.Challenge)))
		b = append(b, f.Challenge...)
	}
	if f.Flags&4 != 0 {
		b = append(b, byte(len(f.Response)))
		b = append(b, f.Response...)
	}
	if f.Flags&8 != 0 {
		b = append(b, f.ID)
	}
	return b, nil
}

// MarshalLen returns the encoded length, or zero for nil or invalid fields.
func (f *L2TPUserAuthenticationFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the fields, rejecting undersized buffers.
func (f *L2TPUserAuthenticationFields) MarshalTo(b []byte) error {
	wire, e := f.Marshal()
	return addedCopyTo(b, wire, e)
}
