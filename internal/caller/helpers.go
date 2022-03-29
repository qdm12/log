package caller

func boolPtr(b bool) *bool { return &b }

func defaultBoolPtr(existing *bool, value bool) *bool {
	if existing != nil {
		return existing
	}
	return &value
}

func copyBoolPtr(b *bool) *bool {
	if b == nil {
		return nil
	}
	return boolPtr(*b)
}

func overrideBoolPtr(existing, newValue *bool) *bool {
	if newValue == nil {
		return existing
	}
	return copyBoolPtr(newValue)
}
