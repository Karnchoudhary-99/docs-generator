package watch

// Changed reports whether two snapshots differ: a file added,
// removed, resized, or modified since the earlier snapshot.
func Changed(before, after map[string]fileState) bool {
	if len(before) != len(after) {
		return true
	}
	for p, a := range after {
		b, ok := before[p]
		if !ok || b != a {
			return true
		}
	}
	return false
}
