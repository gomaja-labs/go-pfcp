package ie

// QERIndicationFlags encodes QER Indications; TS 29.244 V20.1.0 clause 8.2.216.
type QERIndicationFlags uint8

const (
	// QERIndicationsIQFISN is the IQFISN value defined in clause 8.2.216.
	QERIndicationsIQFISN QERIndicationFlags = 1
	// QERIndicationsEDBMI is the EDBMI value defined in clause 8.2.216.
	QERIndicationsEDBMI QERIndicationFlags = 2
	// QERIndicationsEML4S is the EML4S value defined in clause 8.2.216.
	QERIndicationsEML4S QERIndicationFlags = 4
	// QERIndicationsPDUSM is the PDUSM value defined in clause 8.2.216.
	QERIndicationsPDUSM QERIndicationFlags = 8
	// QERIndicationsDBSMI is the DBSMI value defined in clause 8.2.216.
	QERIndicationsDBSMI QERIndicationFlags = 16
	// QERIndicationsTNBMI is the TNBMI value defined in clause 8.2.216.
	QERIndicationsTNBMI QERIndicationFlags = 32
	// QERIndicationsSCOMI is the SCOMI value defined in clause 8.2.216.
	QERIndicationsSCOMI QERIndicationFlags = 64
)

// NewQERIndications creates a QER Indications IE.
// TS 29.244 V20.1.0 clause 8.2.216.
func NewQERIndications(value QERIndicationFlags) *IE {
	if value > 127 {
		return newInvalidIE(QERIndications, ErrInvalidValue)
	}
	return New(QERIndications, []byte{byte(value)})
}

// QERIndications returns the QER Indications value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.216 and 8.1.2.
func (i *IE) QERIndications() (QERIndicationFlags, error) {
	b, err := addedValuePayload(i, QERIndications, 1, false)
	if err != nil {
		return 0, err
	}
	return QERIndicationFlags(b[0] & 127), nil
}
