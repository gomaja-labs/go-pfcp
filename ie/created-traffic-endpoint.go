// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// LocalFTEID returns the first Local F-TEID in structured form.
// TS 29.244 V20.1.0 Table 7.5.3.3-1 shares the F-TEID type between local roles.
func (f *CreatedTrafficEndpointFields) LocalFTEID() (*FTEIDFields, error) {
	return f.LocalFTEIDN(1)
}

// LocalFTEIDN returns the nth Local F-TEID in structured form, counting from one.
// It reports ErrIENotFound for a nil receiver or a missing occurrence.
func (f *CreatedTrafficEndpointFields) LocalFTEIDN(n int) (*FTEIDFields, error) {
	if f == nil || n < 1 {
		return nil, ErrIENotFound
	}
	count := 0
	for _, child := range f.FTEID {
		if child == nil {
			continue
		}
		count++
		if count == n {
			return child.FTEID()
		}
	}
	return nil, ErrIENotFound
}
