package ie

import "fmt"

func groupedChildren(i *IE) ([]*IE, error) {
	if i.ChildIEs != nil {
		return i.ValueAsGrouped()
	}
	return parseGroupedFields(i.Type, i.Payload)
}

// parseGroupedFields starts at the same nesting depth as parsing a parent IE.
// TS 29.244 V20.1.0 clause 8.1.2 permits an empty grouped value.
func parseGroupedFields(_ IEType, b []byte) ([]*IE, error) {
	if len(b) > 65535 {
		return nil, ErrInvalidLength
	}
	return decodeIEs(b, 1)
}

func marshalGroupedFields(parent *IE, b []byte) error {
	length := parent.MarshalLen() - 4
	if length > 65535 || len(b) < length {
		return ErrInvalidLength
	}
	offset := 0
	for _, child := range parent.ChildIEs {
		if child == nil {
			continue
		}
		if err := child.MarshalTo(b[offset:]); err != nil {
			return fmt.Errorf("child of IE %d: %w", parent.Type, err)
		}
		offset += child.MarshalLen()
	}
	return nil
}
