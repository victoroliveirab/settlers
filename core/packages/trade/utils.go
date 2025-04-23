package trade

func equalResourceMaps(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for key, valA := range a {
		valB, ok := b[key]
		if !ok || valB != valA {
			return false
		}
	}
	return true
}
