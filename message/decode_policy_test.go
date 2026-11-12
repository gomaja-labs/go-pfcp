package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func policyIE(typ ie.IEType, payload []byte) []byte {
	b := make([]byte, 4+len(payload))
	binary.BigEndian.PutUint16(b, uint16(typ))
	binary.BigEndian.PutUint16(b[2:], uint16(len(payload)))
	copy(b[4:], payload)
	return b
}

func policyFrame(typ uint8, payload []byte) []byte {
	h := NewHeader(typ, 1, 2, payload)
	b, err := h.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}

func requirePolicyError(t *testing.T, err error, typ ie.IEType, path []ie.IEType, cause error) {
	t.Helper()
	var de *ie.DecodeError
	if !errors.As(err, &de) || de.Type != typ || !errors.Is(err, cause) {
		t.Fatalf("error = %v, want leaf %d and %v", err, typ, cause)
	}
	v := reflect.ValueOf(de).Elem().FieldByName("Path")
	if !v.IsValid() {
		t.Fatal("DecodeError lacks enclosing Path")
	}
	if !reflect.DeepEqual(v.Interface(), path) {
		t.Fatalf("path = %v, want %v", v.Interface(), path)
	}
	for e := errors.Unwrap(err); e != nil; e = errors.Unwrap(e) {
		var nested *ie.DecodeError
		if errors.As(e, &nested) {
			t.Fatal("nested DecodeError")
		}
	}
}

func TestDecodeErrorLeafPath(t *testing.T) {
	leaf := []byte{0, byte(ie.FTEID), 0, 4, 1}
	group := policyIE(ie.CreatePDR, policyIE(ie.PDI, leaf))
	_, err := ie.Parse(group)
	requirePolicyError(t, err, ie.FTEID, []ie.IEType{ie.CreatePDR, ie.PDI}, io.ErrUnexpectedEOF)
	_, err = ie.ParseMultiIEs(group)
	requirePolicyError(t, err, ie.FTEID, []ie.IEType{ie.CreatePDR, ie.PDI}, io.ErrUnexpectedEOF)
}

func TestMessageGroupedDecodePolicy(t *testing.T) {
	malformed := []byte{0, byte(ie.FTEID), 0, 4, 1}
	deep := policyIE(ie.FTEID, nil)
	for range 16 {
		deep = policyIE(ie.CreatePDR, deep)
	}
	tests := []struct {
		name    string
		typ     uint8
		prefix  []byte
		group   ie.IEType
		payload []byte
		cause   error
		leaf    ie.IEType
		path    []ie.IEType
	}{
		{"expected", MsgTypeSessionReportRequest, nil, ie.DownlinkDataReport, malformed, io.ErrUnexpectedEOF, ie.FTEID, []ie.IEType{ie.DownlinkDataReport}},
		{"unexpected", MsgTypeHeartbeatRequest, nil, ie.CreatePDR, malformed, io.ErrUnexpectedEOF, ie.FTEID, []ie.IEType{ie.CreatePDR}},
		{"repeated", MsgTypeSessionReportRequest, policyIE(ie.DownlinkDataReport, nil), ie.DownlinkDataReport, malformed, io.ErrUnexpectedEOF, ie.FTEID, []ie.IEType{ie.DownlinkDataReport}},
		{"nested", MsgTypeHeartbeatRequest, nil, ie.CreatePDR, policyIE(ie.PDI, malformed), io.ErrUnexpectedEOF, ie.FTEID, []ie.IEType{ie.CreatePDR, ie.PDI}},
		{"short_child_header", MsgTypeHeartbeatRequest, nil, ie.CreatePDR, []byte{1, 2, 3}, io.ErrUnexpectedEOF, ie.CreatePDR, nil},
		{"depth", MsgTypeHeartbeatRequest, nil, ie.CreatePDR, deep, ie.ErrNestingTooDeep, ie.CreatePDR, repeatPolicyType(ie.CreatePDR, 16)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := policyIE(tt.group, tt.payload)
			wire := policyFrame(tt.typ, append(append([]byte{}, tt.prefix...), encoded...))
			m, err := Parse(wire)
			if err != nil {
				t.Fatalf("message decode: %v", err)
			}
			out := make([]byte, m.MarshalLen())
			if err = m.MarshalTo(out); err != nil || !bytes.Equal(wire, out) {
				t.Fatalf("replay changed: %v", err)
			}
			ms, err := ParseDatagram(wire)
			if err != nil || len(ms) != 1 {
				t.Fatalf("datagram decode: %v", err)
			}
			var opaque *ie.IE
			switch v := m.(type) {
			case *HeartbeatRequest:
				opaque = v.OtherIEs[0]
			case *SessionReportRequest:
				if len(tt.prefix) > 0 {
					opaque = v.OtherIEs[0]
					if v.DownlinkDataReport.ChildIEs == nil {
						t.Fatal("first occurrence replaced")
					}
				} else {
					opaque = v.DownlinkDataReport
				}
			}
			if opaque.ChildIEs != nil || !bytes.Equal(opaque.Payload, tt.payload) {
				t.Fatal("failed group was not retained opaquely")
			}
			_, err = opaque.ValueAsGrouped()
			requirePolicyError(t, err, tt.leaf, tt.path, tt.cause)
			accessor := reflect.ValueOf(opaque).MethodByName(tt.group.String())
			values := accessor.Call(nil)
			requirePolicyError(t, values[1].Interface().(error), tt.leaf, tt.path, tt.cause)
			_, err = ie.Parse(encoded)
			if err == nil {
				t.Fatal("explicit IE parser became tolerant")
			}
		})
	}
}

func repeatPolicyType(typ ie.IEType, n int) []ie.IEType {
	p := make([]ie.IEType, n)
	for i := range p {
		p[i] = typ
	}
	return p
}

func TestMessageTopLevelFramingStaysStrict(t *testing.T) {
	for _, payload := range [][]byte{{1, 2, 3}, {0, byte(ie.CreatePDR), 0, 5, 1}, append(policyIE(ie.CreatePDR, []byte{1}), 1)} {
		for _, typ := range []uint8{MsgTypeHeartbeatRequest, MsgTypeSessionReportRequest} {
			b := policyFrame(typ, payload)
			if _, err := Parse(b); err == nil {
				t.Fatal("accepted malformed top-level frame")
			}
			if _, err := ParseDatagram(b); err == nil {
				t.Fatal("datagram accepted malformed top-level frame")
			}
		}
	}
}

func TestHeaderErrorPolicy(t *testing.T) {
	parsers := []struct {
		name  string
		parse func([]byte) error
	}{{"header", func(b []byte) error { _, e := ParseHeader(b); return e }}, {"message", func(b []byte) error { _, e := Parse(b); return e }}, {"datagram", func(b []byte) error { _, e := ParseDatagram(b); return e }}}
	for _, typ := range []uint8{MsgTypeHeartbeatRequest, MsgTypeSessionReportRequest} {
		b := policyFrame(typ, nil)
		for _, p := range parsers {
			t.Run(p.name+MessageTypeName(typ), func(t *testing.T) {
				for n := range len(b) {
					short := append([]byte{}, b[:n]...)
					if n >= 4 {
						binary.BigEndian.PutUint16(short[2:], 0)
					}
					if e := p.parse(short); !errors.Is(e, io.ErrUnexpectedEOF) {
						t.Fatalf("short header %d: %v", n, e)
					}
				}
				for _, delta := range []int{-1, 1} {
					v := append([]byte{}, b...)
					binary.BigEndian.PutUint16(v[2:], uint16(len(b)-4+delta))
					if e := p.parse(v); !errors.Is(e, ErrInvalidLength) {
						t.Fatalf("length delta %d: %v", delta, e)
					}
				}
				if e := p.parse(append(append([]byte{}, b...), 0)); !errors.Is(e, ErrInvalidLength) {
					t.Fatalf("long body: %v", e)
				}
			})
		}
	}
}

func TestAllMessageMethodsNilSafe(t *testing.T) {
	for _, m := range catalogueMessages() {
		typ := reflect.TypeOf(m)
		t.Run(typ.Elem().Name(), func(t *testing.T) {
			for n := range typ.NumMethod() {
				method := typ.Method(n)
				t.Run(method.Name, func(t *testing.T) {
					defer func() {
						if p := recover(); p != nil {
							t.Fatalf("nil receiver panicked: %v", p)
						}
					}()
					args := []reflect.Value{reflect.Zero(typ)}
					for i := 1; i < method.Type.NumIn(); i++ {
						args = append(args, reflect.Zero(method.Type.In(i)))
					}
					values := method.Func.Call(args)
					for _, v := range values {
						if v.Type().Implements(reflect.TypeFor[error]()) {
							if v.IsNil() || !errors.Is(v.Interface().(error), ErrNilMessage) {
								t.Errorf("error = %v, want ErrNilMessage", v)
							}
						} else if !v.IsZero() {
							t.Errorf("nonzero result: %v", v)
						}
					}
				})
			}
		})
	}
}
