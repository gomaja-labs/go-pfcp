package ie_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestKeyIDAndSharedSecretOWAMPOrTWAMP(t *testing.T) {
	for _, key := range []string{"", "key-1", "clé", strings.Repeat("k", 80)} {
		for _, secret := range []string{"", "password", strings.Repeat("s", 255)} {
			f := ie.NewKeyIDAndSharedSecretOWAMPOrTWAMPFields(key, secret)
			i := ie.NewKeyIDAndSharedSecretOWAMPOrTWAMP(key, secret)
			want := append([]byte{byte(len(key))}, []byte(key)...)
			want = append(want, byte(len(secret)))
			want = append(want, []byte(secret)...)
			checkNewIEWire(t, i, want)
			checkNewFields(t, f, func(b []byte) (newFields, error) { return ie.ParseKeyIDAndSharedSecretOWAMPOrTWAMPFields(b) })
			got, err := i.KeyIDAndSharedSecretOWAMPOrTWAMP()
			if err != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("fields", err)
			}
		}
	}
	for _, key := range []string{strings.Repeat("k", 81), string([]byte{255})} {
		assertNewIEInvalid(t, ie.NewKeyIDAndSharedSecretOWAMPOrTWAMP(key, "secret"), ie.KeyIDAndSharedSecretOWAMPOrTWAMP)
	}
	for _, secret := range []string{strings.Repeat("s", 256), "line\nfeed", "carriage\rreturn", "nonasciié", string([]byte{255})} {
		assertNewIEInvalid(t, ie.NewKeyIDAndSharedSecretOWAMPOrTWAMP("key", secret), ie.KeyIDAndSharedSecretOWAMPOrTWAMP)
	}
	checkNewFieldsNil(t, (*ie.KeyIDAndSharedSecretOWAMPOrTWAMPFields)(nil))
	checkNewIEFailures(t, ie.KeyIDAndSharedSecretOWAMPOrTWAMP, func(i *ie.IE) (any, error) { return i.KeyIDAndSharedSecretOWAMPOrTWAMP() }, 2)
}
