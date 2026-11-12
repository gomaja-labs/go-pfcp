package ie

// IsGrouped reports whether the type is a grouped IE defined by
// TS 29.244 V20.1.0 Table 8.1.2-1 (all assigned grouped types).
// Vendor-specific IEs require explicit interpretation with ValueAsGrouped.
func (i *IE) IsGrouped() bool {
	if i == nil {
		return false
	}
	return isGroupedType(i.Type)
}

// Add appends non-nil children to a grouped IE. A raw payload is parsed first.
// If that payload is malformed, it returns an error and leaves the IE unchanged.
// Appending does not overwrite a backing array shared with a copied IE.
func (i *IE) Add(ies ...*IE) error {
	if i == nil {
		return ErrIENotFound
	}

	if !i.IsGrouped() && i.ChildIEs == nil {
		return ErrInvalidType
	}
	children, err := i.ValueAsGrouped()
	if err != nil {
		return err
	}
	children = children[:len(children):len(children)]
	for _, child := range ies {
		if child != nil {
			children = append(children, child)
		}
	}
	i.ChildIEs, i.Payload = children, nil
	return nil
}

// Remove removes all children of typ from a grouped IE. A raw payload is parsed
// first; a malformed payload returns an error and leaves the IE unchanged.
func (i *IE) Remove(typ IEType) error {
	if i == nil {
		return ErrIENotFound
	}

	if !i.IsGrouped() && i.ChildIEs == nil {
		return ErrInvalidType
	}
	children, err := i.ValueAsGrouped()
	if err != nil {
		return err
	}
	kept := make([]*IE, 0, len(children))
	for _, child := range children {
		if child != nil && child.Type != typ {
			kept = append(kept, child)
		}
	}
	i.ChildIEs, i.Payload = kept, nil
	return nil
}

// FindByType returns the first child with the requested type.
func (i *IE) FindByType(typ IEType) (*IE, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	if !i.IsGrouped() && i.ChildIEs == nil {
		return nil, ErrInvalidType
	}
	children, err := i.ValueAsGrouped()
	if err != nil {
		return nil, err
	}
	for _, child := range children {
		if child != nil && child.Type == typ {
			return child, nil
		}
	}
	return nil, ErrIENotFound
}
