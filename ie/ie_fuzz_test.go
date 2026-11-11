package ie_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/internal/utils"
)

func callHelpers(i *ie.IE) {
	v := reflect.ValueOf(i)
	for m := 0; m < v.NumMethod(); m++ {
		method := v.Method(m)
		if method.Type().NumIn() == 0 {
			method.Call(nil)
		}
	}
}

func FuzzParse(f *testing.F) {
	for _, v := range ieVectors() {
		f.Add(v.serialized)
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		i, err := ie.Parse(b)
		if err != nil {
			return
		}
		callHelpers(i)
		out, err := i.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(out, b[:i.MarshalLen()]) {
			t.Fatal("successful IE decode changed wire bytes")
		}
	})
}

func FuzzParseMultiIEs(f *testing.F) {
	for _, v := range ieVectors() {
		f.Add(v.serialized)
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		ies, err := ie.ParseMultiIEs(b)
		if err != nil {
			return
		}
		var out []byte
		for _, i := range ies {
			callHelpers(i)
			v, err := i.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			out = append(out, v...)
		}
		if !bytes.Equal(out, b) {
			t.Fatal("successful IE list decode changed wire bytes")
		}
	})
}

func FuzzValueAs(f *testing.F) {
	f.Add(uint16(ie.SDFFilter), []byte{1, 0, 255, 255, 0})
	for _, v := range ieVectors() {
		f.Add(uint16(v.structured.Type), v.serialized[4:])
	}
	f.Fuzz(func(t *testing.T, typ uint16, b []byte) {
		i := ie.New(ie.IEType(typ), b)
		callHelpers(i)
		out, err := i.Marshal()
		if err != nil {
			return
		}
		decoded, err := ie.Parse(out)
		if err != nil {
			return
		} // raw grouped payloads may be malformed
		again, err := decoded.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(out, again) {
			t.Fatal("raw IE round trip changed bytes")
		}
	})
}

func FuzzFields(f *testing.F) {
	f.Add([]byte{1, 0, 255, 255, 0})
	for _, v := range ieVectors() {
		f.Add(v.serialized[4:])
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		for _, parse := range fieldParsers {
			values := reflect.ValueOf(parse).Call([]reflect.Value{reflect.ValueOf(b)})
			if !values[1].IsNil() {
				continue
			}
			fields := values[0].Interface().(interface {
				Marshal() ([]byte, error)
				MarshalTo([]byte) error
				MarshalLen() int
			})
			if _, err := fields.Marshal(); err != nil {
				continue
			}
			if n := fields.MarshalLen(); n > 0 {
				if err := fields.MarshalTo(make([]byte, n-1)); err == nil {
					t.Fatal("accepted short fields buffer")
				}
			}
		}
		utils.DecodeFQDN(b)
		utils.EncodeFQDN(string(b))
		_, _, _ = utils.DecodePLMN(b)
		_, _ = utils.EncodePLMN(string(b), string(b))
	})
}
